package main

import (
	_ "embed"
	"flag"
	"fmt"
	"html/template"
	"log/slog"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/phsym/console-slog"

	ginslog "github.com/madkins23/go-slog/gin"
	"github.com/madkins23/go-slog/infra"
)

// -----------------------------------------------------------------------------

type pageType string

const (
	pageRoot    = "root"
	pageBench   = "bench"
	pageHandler = "handler"
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
	router.GET("/bench/:bench", pageFunction(pageBench))
	router.GET("/handler/:handler", pageFunction(pageHandler))
	router.GET("/style.css", textFunction(css))

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
			tmpl, err = tmpl.Parse(tmplBench)
		case pageHandler:
			tmpl, err = tmpl.Parse(tmplHandler)
		case pageError:
			tmpl, err = tmpl.Parse(tmplError)
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

type pageData struct {
	Data    *infra.BenchData
	Bench   infra.BenchTag
	Handler infra.HandlerTag
}

func pageFunction(page pageType) func(c *gin.Context) {
	return func(c *gin.Context) {
		pageData := pageData{Data: data}
		if name := c.Param("bench"); name != "" {
			pageData.Bench = infra.BenchTag(name)
		} else if name := c.Param("handler"); name != "" {
			pageData.Handler = infra.HandlerTag(name)
		}
		if err := templates[page].Execute(c.Writer, pageData); err != nil {
			slog.Error("Error in page function", "err", err)
		}
	}
}

func textFunction(text string) func(c *gin.Context) {
	return func(c *gin.Context) {
		slog.Debug("textFunction()", "text", text)
		_, _ = c.Writer.Write([]byte(text))
	}
}
