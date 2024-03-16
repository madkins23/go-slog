/*
server parses benchmark test and verification test output and displays it via web pages.

Usage:

	go run cmd/server/server.go [flags]

The flags are:

	-bench string
	    Load benchmark data from path (optional)
	-language value
	    One or more language tags to be tried, defaults to US English.
	-useWarnings
	    Show warning instead of known errors
	-verify string
	    Load verification data from path (optional)

See scripts/bench, scripts/verify and scripts/server for usage examples.
*/
package main

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"log/slog"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/phsym/console-slog"
	"github.com/vicanso/go-charts/v2"
	"github.com/wcharczuk/go-chart/v2"
	"github.com/wcharczuk/go-chart/v2/drawing"
	"golang.org/x/text/message"

	"github.com/madkins23/gin-utils/pkg/handler"
	"github.com/madkins23/gin-utils/pkg/shutdown"

	ginslog "github.com/madkins23/go-slog/gin"
	"github.com/madkins23/go-slog/internal/data"
	"github.com/madkins23/go-slog/internal/language"
	"github.com/madkins23/go-slog/internal/warning"
)

// Server reads output from `test -bench` and serves tables and charts via HTTP.
// See scripts/server for usage example.

type pageType string

const port = 8080

const (
	pageHome     = "pageHome"
	pageTest     = "pageBench"
	pageHandler  = "pageHandler"
	pageScores   = "pageScores"
	pageWarnings = "pageWarnings"
	pageGuts     = "pageGuts"
	pageError    = "pageError"

	partFooter   = "partFooter"
	partHeader   = "partHeader"
	partSource   = "partSource"
	partWarnings = "partWarnings"
)

var (
	//go:embed stuff/style.css
	css string

	//go:embed stuff/home.svg
	home []byte

	//go:embed stuff/scripts.js
	scripts string
)

func main() {
	// Necessary for -bench=<file> and -verify=<file> arguments
	// defined in internal/bench and internal/verify packages.
	flag.Parse()

	gin.DefaultWriter = ginslog.NewWriter(&ginslog.Options{})
	gin.DefaultErrorWriter = ginslog.NewWriter(&ginslog.Options{Level: slog.LevelError})
	logger := slog.New(console.NewHandler(os.Stderr, &console.HandlerOptions{
		Level:      slog.LevelInfo,
		TimeFormat: time.TimeOnly,
	}))
	slog.SetDefault(logger)

	if err := setup(); err != nil {
		slog.Error("Setup error", "err", err)
		return
	}

	graceful := &shutdown.Graceful{}
	graceful.Initialize()
	defer graceful.Close()

	router := gin.Default()
	homePageFn := pageFunction(pageHome)
	router.GET("/go-slog/", homePageFn)
	router.GET("/go-slog/index.html", homePageFn)
	router.GET("/go-slog/test/:tag", pageFunction(pageTest))
	router.GET("/go-slog/handler/:tag", pageFunction(pageHandler))
	router.GET("/go-slog/scores.html", pageFunction(pageScores))
	router.GET("/go-slog/warnings.html", pageFunction(pageWarnings))
	router.GET("/go-slog/guts.html", pageFunction(pageGuts))
	router.GET("/go-slog/error.html", pageFunction(pageError))
	router.GET("/go-slog/chart/scores.svg", scoreFunction)
	router.GET("/go-slog/chart/:tag/:item", chartFunction)
	router.GET("/go-slog/home.svg", svgFunction(home))
	router.GET("/go-slog/scripts.js", textFunction(scripts))
	router.GET("/go-slog/style.css", textFunction(css))
	router.GET("/go-slog/exit", handler.Exit)

	if err := router.SetTrustedProxies(nil); err != nil {
		slog.Error("Don't trust proxies", "err", err)
		os.Exit(1)
	}

	// Listen and serve on 0.0.0.0:8080 (for windows "localhost:8080"). {
	slog.Info("Web Server @ http://localhost:8080/go-slog")

	if err := graceful.Serve(router, port); err != nil {
		slog.Error("Running gin server", "err", err)
	}
}

// -----------------------------------------------------------------------------

var (
	bench     = data.NewBenchmarks()
	warns     = data.NewWarningData()
	score     = data.NewScoreKeeper()
	pages     = []pageType{pageHome, pageTest, pageHandler, pageScores, pageWarnings, pageGuts, pageError}
	templates map[pageType]*template.Template

	//go:embed pages/home.gohtml
	tmplPageHome string

	//go:embed pages/test.gohtml
	tmplPageTest string

	//go:embed pages/handler.gohtml
	tmplPageHandler string

	//go:embed pages/scores.gohtml
	tmplPageScores string

	//go:embed pages/warnings.gohtml
	tmplPageWarnings string

	//go:embed pages/guts.gohtml
	tmplPageGuts string

	//go:embed pages/error.gohtml
	tmplPageError string

	//go:embed parts/footer.gohtml
	tmplPartFooter string

	//go:embed parts/header.gohtml
	tmplPartHeader string

	//go:embed parts/source.gohtml
	tmplPartSource string

	//go:embed parts/warnings.gohtml
	tmplPartWarnings string
)

// setup server data structures and templates.
func setup() error {
	if err := language.Setup(); err != nil {
		return fmt.Errorf("language setup: %w", err)
	}

	if err := data.Setup(bench, warns, score); err != nil {
		return fmt.Errorf("data setup: %w", err)
	}

	templates = make(map[pageType]*template.Template)
	for _, page := range pages {
		var err error
		tmpl := template.New(string(page))
		tmpl.Funcs(functions())
		switch page {
		case pageHome:
			tmpl, err = tmpl.Parse(tmplPageHome)
			if err == nil {
				_, err = tmpl.New(partHeader).Parse(tmplPartHeader)
			}
			if err == nil {
				_, err = tmpl.New(partFooter).Parse(tmplPartFooter)
			}
		case pageTest:
			_, err = tmpl.Parse(tmplPageTest)
			if err == nil {
				_, err = tmpl.New(partHeader).Parse(tmplPartHeader)
			}
			if err == nil {
				_, err = tmpl.New(partFooter).Parse(tmplPartFooter)
			}
			if err == nil {
				_, err = tmpl.New(partWarnings).Parse(tmplPartWarnings)
			}
		case pageHandler:
			tmpl, err = tmpl.Parse(tmplPageHandler)
			if err == nil {
				_, err = tmpl.New(partHeader).Parse(tmplPartHeader)
			}
			if err == nil {
				_, err = tmpl.New(partFooter).Parse(tmplPartFooter)
			}
			if err == nil {
				_, err = tmpl.New(partWarnings).Parse(tmplPartWarnings)
			}
		case pageScores:
			tmpl, err = tmpl.Parse(tmplPageScores)
			if err == nil {
				_, err = tmpl.New(partHeader).Parse(tmplPartHeader)
			}
			if err == nil {
				_, err = tmpl.New(partFooter).Parse(tmplPartFooter)
			}
		case pageWarnings:
			tmpl, err = tmpl.Parse(tmplPageWarnings)
			if err == nil {
				_, err = tmpl.New(partHeader).Parse(tmplPartHeader)
			}
			if err == nil {
				_, err = tmpl.New(partFooter).Parse(tmplPartFooter)
			}
		case pageGuts:
			tmpl, err = tmpl.Parse(tmplPageGuts)
			if err == nil {
				_, err = tmpl.New(partHeader).Parse(tmplPartHeader)
			}
			if err == nil {
				_, err = tmpl.New(partFooter).Parse(tmplPartFooter)
			}
			if err == nil {
				_, err = tmpl.New(partSource).Parse(tmplPartSource)
			}
		case pageError:
			tmpl, err = tmpl.Parse(tmplPageError)
			if err == nil {
				_, err = tmpl.New(partHeader).Parse(tmplPartHeader)
			}
			if err == nil {
				_, err = tmpl.New(partFooter).Parse(tmplPartFooter)
			}
		default:
			return fmt.Errorf("unknown page name: %s", page)
		}
		if err != nil {
			return fmt.Errorf("parse template %s: %w", page, err)
		}
		templates[page] = tmpl
	}

	return nil
}

// -----------------------------------------------------------------------------

var (
	chartCache      = make(map[string][]byte)
	chartCacheMutex sync.Mutex
)

// chartFunction generates an SVG chart for the current object tags.
func chartFunction(c *gin.Context) {
	itemArg := strings.TrimSuffix(c.Param("item"), ".svg")
	item, err := data.BenchItemsString(itemArg)
	if err != nil {
		slog.Error("Bad URL parameter", "param", itemArg, "err", err)
		// TODO: what to do here?
		return
	}
	tag := c.Param("tag")
	cacheKey := tag + ":" + item.String()
	chartCacheMutex.Lock()
	ch, found := chartCache[cacheKey]
	chartCacheMutex.Unlock()
	if !found {
		var labels []string
		var values []float64
		if records := bench.HandlerRecords(data.TestTag(tag)); records != nil {
			labels, values = chartTest(records, item)
		} else if records := bench.TestRecords(data.HandlerTag(tag)); records != nil {
			labels, values = chartHandler(records, item)
		} else {
			slog.Error("Neither handler nor benchmark records found", "fn", "chartFunction")
			c.HTML(http.StatusBadRequest, "pageFunction", gin.H{
				"ErrorTitle":   "Template failed execution",
				"ErrorMessage": "No records for " + tag})
			return
		}
		const verticalPadding = 100
		const barMargin = 5
		const barWidth = 15
		painter, err := charts.HorizontalBarRender(
			[][]float64{values},
			charts.SVGTypeOption(),
			charts.TitleTextOptionFunc(item.LongName()),
			charts.YAxisDataOptionFunc(labels),
			charts.WidthOptionFunc(400),
			charts.HeightOptionFunc(verticalPadding+(barMargin+barWidth)*len(values)),
			charts.PaddingOptionFunc(charts.Box{
				Top:    10,
				Right:  20,
				Bottom: 10,
				Left:   10,
			}),
			func(opt *charts.ChartOption) {
				opt.BarWidth = barWidth
			},
		)
		if err != nil {
			panic(err)
		}
		ch, err = painter.Bytes()
		if err != nil {
			panic(err)
		}
		chartCacheMutex.Lock()
		chartCache[cacheKey] = ch
		chartCacheMutex.Unlock()
	}
	c.Data(http.StatusOK, "image/svg+xml", ch)
}

// charTest returns labels and values for a Test chart.
func chartTest(records data.HandlerRecords, item data.BenchItems) (labels []string, values []float64) {
	labels = make([]string, 0, len(records))
	values = make([]float64, 0, len(records))
	for _, tag := range bench.HandlerTags() {
		if record, found := records[tag]; found {
			labels = append(labels, bench.HandlerName(tag))
			values = append(values, record.ItemValue(item))
		}
	}
	reverse(labels)
	reverse(values)
	return
}

// chartHandler returns labels and values for a Handler chart.
func chartHandler(records data.TestRecords, item data.BenchItems) (labels []string, values []float64) {
	labels = make([]string, 0, len(records))
	values = make([]float64, 0, len(records))
	for _, tag := range bench.TestTags() {
		if record, found := records[tag]; found {
			labels = append(labels, bench.TestName(tag))
			values = append(values, record.ItemValue(item))
		}
	}
	reverse(labels)
	reverse(values)
	return
}

// scoreFunction generates an SVG chart for the score visualization.
func scoreFunction(c *gin.Context) {
	cacheKey := "score"
	chartCacheMutex.Lock()
	ch, found := chartCache[cacheKey]
	chartCacheMutex.Unlock()
	if !found {
		graph := scoreChartFunction()
		var buf bytes.Buffer
		if err := graph.Render(chart.SVG, &buf); err != nil {
			log.Println(err.Error())
		} else {
			ch = buf.Bytes()
			chartCacheMutex.Lock()
			chartCache[cacheKey] = ch
			chartCacheMutex.Unlock()
		}
	}
	c.Data(http.StatusOK, chart.ContentTypeSVG, ch)
}

const shiftAnnotations = 2

func scoreChartFunction() chart.Chart {
	type handlerCoords struct{ x, y float64 }
	handlers := make(map[data.HandlerTag]*handlerCoords)
	for _, hdlr := range bench.HandlerTags() {
		handlers[hdlr] = &handlerCoords{y: score.HandlerBenchScores(hdlr).Overall}
	}
	for _, hdlr := range warns.HandlerTags() {
		if coords, found := handlers[hdlr]; found {
			// Only add value if there is already a benchmark score.
			coords.x = score.HandlerWarningScore(hdlr)
		}
	}
	aValues := make([]chart.Value2, 0, len(handlers)+1)
	xValues := make([]float64, 0, len(handlers)+1)
	yValues := make([]float64, 0, len(handlers)+1)
	for hdlr, coords := range handlers {
		if coords.y > 0.00001 {
			// This is an attempt to not have annotations sit on top of each other.
			// It won't handle cases where more than two dots overlap.
			y := coords.y
			for tag, a := range aValues {
				yDist := a.YValue - y
				if math.Abs(yDist) < 2*shiftAnnotations && math.Abs(a.XValue-coords.x) < 10 {
					if yDist >= 0 {
						aValues[tag].YValue += shiftAnnotations
						y -= shiftAnnotations
					} else {
						aValues[tag].YValue -= shiftAnnotations
						y += shiftAnnotations
					}
				}
			}
			aValues = append(aValues, chart.Value2{
				Label:  warns.HandlerName(hdlr),
				XValue: coords.x + 1,
				YValue: y,
			})
			xValues = append(xValues, coords.x)
			yValues = append(yValues, coords.y)
		}
	}
	ticks := make([]chart.Tick, 11)
	for i := 0; i < 11; i++ {
		ticks[i].Value = float64(i * 10)
		ticks[i].Label = strconv.FormatFloat(ticks[i].Value, 'f', 1, 64)
	}
	return chart.Chart{
		Height: 600,
		Width:  750,
		XAxis: chart.XAxis{
			Name:  "Warning Score",
			Range: &chart.ContinuousRange{Min: 0, Max: 100.0, Domain: 100.0},
			Ticks: ticks,
		},
		YAxis: chart.YAxis{
			Name: "Benchmark Score",
			//AxisType: chart.YAxisSecondary, // cuts off axis labels on left
			Range: &chart.ContinuousRange{Min: 0, Max: 100.0, Domain: 100.0},
			Ticks: ticks,
		},
		Series: []chart.Series{
			chart.ContinuousSeries{
				Style: chart.Style{
					DotWidth:         5,
					DotColorProvider: scoreChartColorFunction,
					StrokeWidth:      chart.Disabled,
				},
				XValues: xValues,
				YValues: yValues,
			},
			chart.AnnotationSeries{
				Style: chart.Style{
					StrokeWidth: chart.Disabled,
					DotWidth:    chart.Disabled,
				},
				Annotations: aValues,
			},
		},
	}
}

const bigByte = float64(0xFF)

func scoreChartColorFunction(xr, yr chart.Range, _ int, x, y float64) drawing.Color {
	ratio := scoreChartRatio(xr.GetMin(), yr.GetMin(), xr.GetMax(), yr.GetMax())
	diagonal := scoreChartDistance(xr.GetMin(), yr.GetMin(), xr.GetMax(), yr.GetMax(), ratio)
	distLow := scoreChartDistance(x, y, xr.GetMin(), yr.GetMin(), ratio)
	distHigh := scoreChartDistance(xr.GetMax(), yr.GetMax(), x, y, ratio)
	return drawing.Color{
		R: scoreChartColor(distLow, diagonal),
		G: scoreChartColor(distHigh, diagonal),
		B: 0x00,
		A: 0xff,
	}
}

func scoreChartRatio(xMin, yMin, xMax, yMax float64) float64 {
	return (xMax - xMin) / (yMax - yMin)
}

func scoreChartColor(distance, diagonal float64) uint8 {
	if distance > diagonal {
		return 0
	}
	return byte(bigByte - bigByte*distance/diagonal)
}

func scoreChartDistance(xMin, yMin, xMax, yMax, ratio float64) float64 {
	return math.Sqrt(math.Pow((xMax-xMin)/ratio, 2) + math.Pow(yMax-yMin, 2))
}

// -----------------------------------------------------------------------------

// templateData is all the data that will be available during template execution.
type templateData struct {
	*data.Benchmarks
	*data.Warnings
	data.Scores
	Handler   data.HandlerTag
	Test      data.TestTag
	Levels    []warning.Level
	Printer   *message.Printer
	Page      string
	Timestamp string
	Errors    []string
}

// FixUint converts a uint64 into a string using the language printer.
// This will apply the proper numeric separators.
func (pd *templateData) FixUint(number uint64) string {
	return pd.Printer.Sprintf("%d", number)
}

// FixFloat converts a float64 into a string using the language printer.
// This will apply the proper decimal and numeric separators.
func (pd *templateData) FixFloat(number float64) string {
	return pd.Printer.Sprintf("%0.2f", number)
}

// pageFunction returns a Gin handler function for generating an HTML page for the server.
// During execution of the handler function,
// URL parameter values will be read and appropriate object tags configured as template data.
// The appropriate template will be executed using the template data to generate the HTML page.
func pageFunction(page pageType) gin.HandlerFunc {
	return func(c *gin.Context) {
		tmplData := &templateData{
			Benchmarks: bench,
			Warnings:   warns,
			Scores:     score,
			Levels:     warning.LevelOrder,
			Printer:    language.Printer(),
			Page:       string(page),
			Timestamp:  time.Now().Format(time.DateTime + " MST"),
		}
		switch page {
		case pageTest, pageHandler:
			if tag := c.Param("tag"); tag == "" {
				slog.Error("No URL parameter", "param", "tag")
			} else {
				tag := strings.TrimSuffix(tag, ".html")
				if page == pageTest {
					tmplData.Test = data.TestTag(tag)
				} else if page == pageHandler {
					tmplData.Handler = data.HandlerTag(tag)
				}
			}
		case pageError:
			tmplData.Errors = c.Errors.Errors()
		}
		var buf bytes.Buffer
		if err := templates[page].Execute(&buf, tmplData); err != nil {
			slog.Error("Error executing template", "err", err, "page", page)
			_ = c.Error(err)
			if page != pageError {
				pageFunction(pageError)(c)
			}
		} else if _, err := io.Copy(c.Writer, &buf); err != nil {
			slog.Error("Error writing page", "err", err, "page", page)
			_ = c.Error(err)
		}
	}
}

// functions returns required template functions.
func functions() map[string]any {
	return map[string]any{
		"mod": func(a, b int) int {
			return a % b
		},
		"dict": func(v ...interface{}) map[string]any {
			dict := map[string]interface{}{}
			vLen := len(v)
			for i := 0; i < vLen; i += 2 {
				if key, ok := v[i].(string); ok && i+1 < vLen {
					dict[key] = v[i+1]
				}
			}
			return dict
		},
		"unescape": func(s string) template.HTML {
			return template.HTML(s)
		},
	}
}

// reverse an array.
func reverse[T any](array []T) {
	i := 0
	j := len(array) - 1
	for i < j {
		array[i], array[j] = array[j], array[i]
		i++
		j--
	}
}

// svgFunction returns a Gin handler function for generating an SVG image.
// During execution of the handler function the SVG image will be generated from the specified string.
func svgFunction(svg []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		slog.Debug("svgFunction()", "svg", string(svg))
		c.Data(http.StatusOK, "image/svg+xml", svg)
	}
}

// textFunction returns a Gin handler function to return a block of ASCII text.
// During execution of the handler function the text will be generated from the specified string.
func textFunction(text string) gin.HandlerFunc {
	return func(c *gin.Context) {
		slog.Debug("textFunction()", "text", text)
		if _, err := c.Writer.Write([]byte(text)); err != nil {
			c.HTML(http.StatusInternalServerError, "textFunction", gin.H{
				"ErrorTitle":   "Failed to write string",
				"ErrorMessage": err.Error()})
		}
	}
}
