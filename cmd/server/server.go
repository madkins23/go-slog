package main

import (
	_ "embed"
	"flag"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/phsym/console-slog"
	"github.com/vicanso/go-charts/v2"
	"golang.org/x/text/message"

	ginslog "github.com/madkins23/go-slog/gin"
	"github.com/madkins23/go-slog/internal/bench"
	"github.com/madkins23/go-slog/internal/language"
)

// Server reads output from `test -bench` and serves tables and charts via HTTP.
// See scripts/server for usage example.

type pageType string

const (
	pageRoot    = "root"
	pageTest    = "bench"
	pageHandler = "handler"
	pageChoices = "choices"
	pageError   = "error"
)

var (
	//go:embed embed/style.css
	css string

	//go:embed embed/home.svg
	home []byte
)

func main() {
	flag.Parse() // Necessary for -json=<file> argument defined in infra package.

	gin.DefaultWriter = ginslog.NewWriter(slog.LevelInfo, ginslog.NoTraffic)
	gin.DefaultErrorWriter = ginslog.NewWriter(slog.LevelError, ginslog.NoTraffic)
	logger := slog.New(console.NewHandler(os.Stderr, &console.HandlerOptions{
		Level:      slog.LevelInfo,
		TimeFormat: time.TimeOnly,
	}))
	slog.SetDefault(logger)

	if err := setup(); err != nil {
		slog.Error("Error during setup", "err", err)
		return
	}

	router := gin.Default()
	router.GET("/", pageFunction(pageRoot))
	router.GET("/test", pageFunction(pageTest))
	router.GET("/handler", pageFunction(pageHandler))
	router.GET("/chart.svg", chartFunction)
	router.GET("/home.svg", svgFunction(home))
	router.GET("/style.css", textFunction(css))

	if err := router.SetTrustedProxies(nil); err != nil {
		slog.Error("Don't trust proxies", "err", err)
		os.Exit(1)
	}

	// Listen and serve on 0.0.0.0:8080 (for windows "localhost:8080"). {
	slog.Info("Web Server @ http://localhost:8080")
	if err := router.Run(); err != nil {
		slog.Error("Error during ListenAndServe()", "err", err)
	}
}

// -----------------------------------------------------------------------------

var (
	data      = &bench.Data{}
	pages     = []pageType{pageRoot, pageTest, pageHandler}
	templates map[pageType]*template.Template

	//go:embed embed/root.tmpl
	tmplRoot string

	//go:embed embed/test.tmpl
	tmplTest string

	//go:embed embed/handler.tmpl
	tmplHandler string

	//go:embed embed/choices.tmpl
	tmplChoices string

	//go:embed embed/error.tmpl
	tmplError string
)

func setup() error {
	if err := language.Setup(); err != nil {
		return fmt.Errorf("language setup: %w", err)
	}

	if err := data.LoadDataJSON(); err != nil {
		return fmt.Errorf("load benchmark JSON: %w", err)
	}

	templates = make(map[pageType]*template.Template)
	for _, page := range pages {
		var err error
		tmpl := template.New(string(page))
		switch page {
		case pageRoot:
			tmpl, err = tmpl.Parse(tmplRoot)
		case pageTest:
			_, err = tmpl.Parse(tmplTest)
			if err == nil {
				_, err = tmpl.New(pageChoices).Parse(tmplChoices)
			}
		case pageHandler:
			tmpl, err = tmpl.Parse(tmplHandler)
			if err == nil {
				_, err = tmpl.New(pageChoices).Parse(tmplChoices)
			}
		case pageError:
			tmpl, err = tmpl.Parse(tmplError)
		case pageChoices:
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

func chartFunction(c *gin.Context) {
	itemArg := c.Query("item")
	item, err := bench.TestItemsString(itemArg)
	if err != nil {
		slog.Error("Bad item URL argument", "arg", itemArg, "err", err)
		// TODO: what to do here?
		return
	}
	tag := c.Query("tag")
	cacheKey := tag + ":" + item.String()
	chartCacheMutex.Lock()
	ch, found := chartCache[cacheKey]
	chartCacheMutex.Unlock()
	if !found {
		var labels []string
		var values []float64
		if records := data.HandlerRecords(bench.TestTag(tag)); records != nil {
			labels, values = chartTest(records, item)
		} else if records := data.TestRecords(bench.HandlerTag(tag)); records != nil {
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

func chartTest(records bench.HandlerRecords, item bench.TestItems) (labels []string, values []float64) {
	labels = make([]string, 0, len(records))
	values = make([]float64, 0, len(records))
	for _, tag := range data.HandlerTags() {
		if record, found := records[tag]; found {
			labels = append(labels, data.HandlerName(tag))
			values = append(values, record.ItemValue(item))
		}
	}
	return
}

func chartHandler(records bench.TestRecords, item bench.TestItems) (labels []string, values []float64) {
	labels = make([]string, 0, len(records))
	values = make([]float64, 0, len(records))
	for _, tag := range data.TestTags() {
		if record, found := records[tag]; found {
			labels = append(labels, data.TestName(tag))
			values = append(values, record.ItemValue(item))
		}
	}
	return
}

// -----------------------------------------------------------------------------

type pageData struct {
	Data    *bench.Data
	Test    bench.TestTag
	Handler bench.HandlerTag
	Printer *message.Printer
}

func (pd *pageData) FixUint(number uint64) string {
	return pd.Printer.Sprintf("%d", number)
}

func (pd *pageData) FixFloat(number float64) string {
	return pd.Printer.Sprintf("%0.2f", number)
}

func pageFunction(page pageType) gin.HandlerFunc {
	return func(c *gin.Context) {
		pageData := &pageData{Data: data, Printer: language.Printer()}
		if page == pageTest || page == pageHandler {
			if tag := c.Query("tag"); tag == "" {
				slog.Error("No URL argument", "arg", "tag")
			} else if page == pageTest {
				pageData.Test = bench.TestTag(tag)
			} else if page == pageHandler {
				pageData.Handler = bench.HandlerTag(tag)
			}
		}
		if err := templates[page].Execute(c.Writer, pageData); err != nil {
			slog.Error("Error in page function", "err", err)
			c.HTML(http.StatusInternalServerError, "pageFunction", gin.H{
				"ErrorTitle":   "Template failed execution",
				"ErrorMessage": err.Error()})
		}
	}
}

func svgFunction(svg []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		slog.Debug("svgFunction()", "svg", string(svg))
		c.Data(http.StatusOK, "image/svg+xml", svg)
	}
}

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
