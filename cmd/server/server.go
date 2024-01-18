package main

import (
	_ "embed"
	"flag"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/phsym/console-slog"
	charts "github.com/vicanso/go-charts/v2"

	ginslog "github.com/madkins23/go-slog/gin"
	"github.com/madkins23/go-slog/internal/bench"
)

// -----------------------------------------------------------------------------

type pageType string

const (
	pageRoot    = "root"
	pageBench   = "bench"
	pageHandler = "handler"
	pageChoices = "choices"
	pageError   = "error"
)

var (
	//go:embed embed/style.css
	css string
)

func main() {
	flag.Parse() // Necessary for -json=<file> argument defined in infra package.

	gin.DefaultWriter = ginslog.NewWriter(slog.LevelInfo)
	gin.DefaultErrorWriter = ginslog.NewWriter(slog.LevelError)
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
	router.GET("/bench", pageFunction(pageBench))
	router.GET("/handler", pageFunction(pageHandler))
	router.GET("/chart.svg", chartFunction)
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
	pages     = []pageType{pageRoot, pageBench, pageHandler}
	templates map[pageType]*template.Template

	//go:embed embed/root.tmpl
	tmplRoot string

	//go:embed embed/bench.tmpl
	tmplBench string

	//go:embed embed/handler.tmpl
	tmplHandler string

	//go:embed embed/choices.tmpl
	tmplChoices string

	//go:embed embed/error.tmpl
	tmplError string
)

func setup() error {
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
		case pageBench:
			_, err = tmpl.Parse(tmplBench)
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

//go:generate go run github.com/dmarkham/enumer -type=dataFields
type dataFields uint

const (
	dfNanos dataFields = iota
	dfMemAllocs
	dfMemBytes
	dfMemMB
	dfLen
)

// -----------------------------------------------------------------------------

var (
	chartCache = make(map[string][]byte)
)

func chartFunction(c *gin.Context) {
	tag := c.Query("tag")
	//tag := c.Query("field")
	ch, found := chartCache[tag]
	if !found {
		if records := data.HandlerRecords(bench.TestTag(tag)); records != nil {
			_, _, _ = chartBench(bench.TestTag(tag), records)
		} else if records := data.TestRecords(bench.HandlerTag(tag)); records != nil {
			_, _, _ = chartHandler(bench.HandlerTag(tag), records)
		} else {
			slog.Error("Neither handler nor benchmark records found", "fn", "chartFunction")
			c.HTML(http.StatusBadRequest, "pageFunction", gin.H{
				"ErrorTitle":   "Template failed execution",
				"ErrorMessage": "No records for " + tag})
			return
		}
		painter, err := charts.HorizontalBarRender(
			// TODO: here
			[][]float64{
				{
					10,
					30,
					50,
					70,
					90,
					110,
					130,
				},
			},
			charts.SVGTypeOption(),
			// TODO: here
			charts.TitleTextOptionFunc("World Population"),
			charts.PaddingOptionFunc(charts.Box{
				Top:    20,
				Right:  40,
				Bottom: 20,
				Left:   20,
			}),
			// TODO: here
			charts.YAxisDataOptionFunc([]string{
				"UN",
				"Brazil",
				"Indonesia",
				"USA",
				"India",
				"China",
				"World",
			}),
		)
		if err != nil {
			panic(err)
		}
		ch, err = painter.Bytes()
		if err != nil {
			panic(err)
		}
		chartCache[tag] = ch
	}
	c.Data(http.StatusOK, "image/svg+xml", ch)
}

func chartBench(test bench.TestTag, records bench.HandlerRecords) (
	title string, labels []string, values [][]float64) {

	title = data.TestName(test)

	order := make([]string, 0, len(records))
	for benchTag := range records {
		order = append(order, string(benchTag))
	}

	sort.Strings(order)
	labels = make([]string, len(records))

	values = make([][]float64, len(records))
	for i := 0; i < len(records); i++ {
		values[i] = make([]float64, 4)
	}
	for i, tag := range order {
		labels[i] = data.HandlerName(bench.HandlerTag(tag))
		record := records[bench.HandlerTag(tag)]
		values[i][0] = record.NanosPerOp
		values[i][1] = float64(record.MemAllocsPerOp)
		values[i][2] = float64(record.MemBytesPerOp)
		values[i][3] = float64(record.MemMbPerSec)
	}

	return
}

func chartHandler(handler bench.HandlerTag, records bench.TestRecords) (
	title string, labels []string, values [][]float64) {

	title = data.HandlerName(handler)

	order := make([]string, 0, len(records))
	for benchTag := range records {
		order = append(order, string(benchTag))
	}
	sort.Strings(order)
	labels = make([]string, len(records))
	for i, tag := range order {
		labels[i] = data.TestName(bench.TestTag(tag))
	}

	values = make([][]float64, 4)
	for i := 0; i < 4; i++ {
		values[i] = make([]float64, 0, len(records))
	}
	for _, tag := range order {
		record := records[bench.TestTag(tag)]
		values[0] = append(values[0], record.NanosPerOp)
		values[1] = append(values[1], float64(record.MemAllocsPerOp))
		values[2] = append(values[2], float64(record.MemBytesPerOp))
		values[3] = append(values[3], float64(record.MemMbPerSec))
	}

	return
}

type pageData struct {
	Data    *bench.Data
	Bench   bench.TestTag
	Handler bench.HandlerTag
}

func pageFunction(page pageType) gin.HandlerFunc {
	return func(c *gin.Context) {
		pageData := pageData{Data: data}
		if page == pageBench || page == pageHandler {
			if name := c.Query("tag"); name == "" {
				slog.Error("No 'tag' URL argument")
			} else if page == pageBench {
				pageData.Bench = bench.TestTag(name)
			} else if page == pageHandler {
				pageData.Handler = bench.HandlerTag(name)
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
