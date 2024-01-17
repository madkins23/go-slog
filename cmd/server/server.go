package main

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/phsym/console-slog"
	"github.com/wcharczuk/go-chart/v2"

	ginslog "github.com/madkins23/go-slog/gin"
	"github.com/madkins23/go-slog/infra"
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
	data      = &infra.BenchData{}
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
	if err := data.LoadBenchJSON(); err != nil {
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

var (
	charts = make(map[string][]byte)
)

func chartFunction(c *gin.Context) {
	tag := c.Query("tag")
	ch, found := charts[tag]
	if !found {
		// Generate and save chart.
		graph := chart.Chart{
			Series: []chart.Series{
				chart.ContinuousSeries{
					XValues: []float64{1.0, 2.0, 3.0, 4.0},
					YValues: []float64{1.0, 2.0, 3.0, 4.0},
				},
			},
		}
		b := &bytes.Buffer{}
		if err := graph.Render(chart.SVG, b); err != nil {
			slog.Error("Render graph", "err", err)
		}
		ch = b.Bytes()
		charts[tag] = ch
	}
	c.Data(http.StatusOK, "image/svg+xml", ch)
}

type pageData struct {
	Data    *infra.BenchData
	Bench   infra.BenchTag
	Handler infra.HandlerTag
}

func pageFunction(page pageType) gin.HandlerFunc {
	return func(c *gin.Context) {
		pageData := pageData{Data: data}
		if page == pageBench || page == pageHandler {
			if name := c.Query("tag"); name == "" {
				slog.Error("No 'tag' URL argument")
			} else if page == pageBench {
				pageData.Bench = infra.BenchTag(name)
			} else if page == pageHandler {
				pageData.Handler = infra.HandlerTag(name)
			}
		}
		if err := templates[page].Execute(c.Writer, pageData); err != nil {
			slog.Error("Error in page function", "err", err)
			c.HTML(http.StatusBadRequest, "pageFunction", gin.H{
				"ErrorTitle":   "Template failed execution",
				"ErrorMessage": err.Error()})
		}
	}
}

func textFunction(text string) gin.HandlerFunc {
	return func(c *gin.Context) {
		slog.Debug("textFunction()", "text", text)
		if _, err := c.Writer.Write([]byte(text)); err != nil {
			c.HTML(http.StatusBadRequest, "textFunction", gin.H{
				"ErrorTitle":   "Failed to write string",
				"ErrorMessage": err.Error()})
		}
	}
}
