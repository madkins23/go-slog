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
	_ "embed"
	"flag"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/phsym/console-slog"
	"github.com/vicanso/go-charts/v2"
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
	pageWarnings = "pageWarnings"
	pageGuts     = "pageGuts"
	pageError    = "pageError"

	partFooter   = "partFooter"
	partHeader   = "partHeader"
	partSource   = "partSource"
	partWarnings = "partWarnings"
)

var (
	//go:embed parts/style.css
	css string

	//go:embed parts/home.svg
	home []byte

	//go:embed parts/scripts.js
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
	router.GET("/go-slog/warnings.html", pageFunction(pageWarnings))
	router.GET("/go-slog/guts.html", pageFunction(pageGuts))
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
	pages     = []pageType{pageHome, pageTest, pageHandler, pageWarnings, pageGuts}
	templates map[pageType]*template.Template

	//go:embed pages/home.tmpl
	tmplPageHome string

	//go:embed pages/test.tmpl
	tmplPageTest string

	//go:embed pages/handler.tmpl
	tmplPageHandler string

	//go:embed pages/warnings.tmpl
	tmplPageWarnings string

	//go:embed pages/guts.tmpl
	tmplPageGuts string

	//go:embed pages/error.tmpl
	tmplPageError string

	//go:embed parts/footer.tmpl
	tmplPartFooter string

	//go:embed parts/header.tmpl
	tmplPartHeader string

	//go:embed parts/source.tmpl
	tmplPartSource string

	//go:embed parts/warnings.tmpl
	tmplPartWarnings string
)

// setup server data structures and templates.
func setup() error {
	if err := language.Setup(); err != nil {
		return fmt.Errorf("language setup: %w", err)
	}

	if err := data.Setup(bench, warns); err != nil {
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
		case pageWarnings:
			tmpl, err = tmpl.Parse(tmplPageWarnings)
			if err == nil {
				_, err = tmpl.New(partHeader).Parse(tmplPartHeader)
			}
			if err == nil {
				_, err = tmpl.New(partFooter).Parse(tmplPartFooter)
			}
		case pageError:
			tmpl, err = tmpl.Parse(tmplPageError)
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

// -----------------------------------------------------------------------------

// templateData is all the data that will be available during template execution.
type templateData struct {
	*data.Benchmarks
	*data.Warnings
	Handler   data.HandlerTag
	Test      data.TestTag
	Levels    []warning.Level
	Printer   *message.Printer
	Page      string
	Timestamp string
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
			Levels:     warning.LevelOrder,
			Printer:    language.Printer(),
			Page:       string(page),
			Timestamp:  time.Now().Format(time.DateTime + " MST"),
		}
		if page == pageTest || page == pageHandler {
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
		}
		if err := templates[page].Execute(c.Writer, tmplData); err != nil {
			slog.Error("Error in page function", "err", err)
			c.HTML(http.StatusInternalServerError, "pageFunction", gin.H{
				"ErrorTitle":   "Template failed execution",
				"ErrorMessage": err.Error()})
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
