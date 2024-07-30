/*
server parses benchmark test and verification test output and displays it via web pages.

# Usage

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

The scripts/server script will run cmd/server,
taking input from temporary files created by scripts/verify and scripts/bench.

The -language flag is used to enable proper formatting of displayed numbers.

# Output

	GOROOT=/snap/go/current #gosetup
	GOPATH=/home/madkins23/go #gosetup
	/snap/go/current/bin/go build -o /home/madkins23/.cache/JetBrains/GoLand2023.3/tmp/GoLand/___run_server /home/madkins23/work/go/src/github.com/madkins23/go-slog/cmd/server/server.go #gosetup
	/home/madkins23/.cache/JetBrains/GoLand2023.3/tmp/GoLand/___run_server -bench=/tmp/go-slog/bench.txt -verify=/tmp/go-slog/verify.txt
	12:31:18 WRN Creating an Engine instance with the Logger and Recovery middleware already attached.
	12:31:18 WRN Running in "debug" mode. Switch to "release" mode in production.
	 - using env:   export GIN_MODE=release
	 - using code:  gin.SetMode(gin.ReleaseMode)
	12:31:18 INF Web Server @ http://localhost:8080/go-slog
	12:31:19 INF HasTest() test=Verify$KeyCase found=true
	12:31:19 INF 200 |    5.012399ms |             ::1 | GET      "/go-slog/test/Verify$KeyCase.html"
	12:31:19 INF 200 |      46.763µs |             ::1 | GET      "/go-slog/home.svg"
	12:31:19 INF 200 |     169.933µs |             ::1 | GET      "/go-slog/style.css"

# Notes

The use of gin-gonic/gin is probably unnecessary aside from
demonstrating the use of the [gin] package defined in this repository.

[gin]: https://pkg.go.dev/github.com/madkins23/go-slog/gin
*/
package main

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/phsym/console-slog"
	"golang.org/x/text/message"

	"github.com/madkins23/gin-utils/pkg/handler"
	"github.com/madkins23/gin-utils/pkg/shutdown"

	ginslog "github.com/madkins23/go-slog/gin"
	"github.com/madkins23/go-slog/infra/warning"
	"github.com/madkins23/go-slog/internal/data"
	"github.com/madkins23/go-slog/internal/language"
	"github.com/madkins23/go-slog/internal/scoring"
	"github.com/madkins23/go-slog/internal/scoring/keeper"
	"github.com/madkins23/go-slog/internal/scoring/score"
)

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
	router.GET("/go-slog/scores/:keeper/summary.html", pageFunction(pageScores))
	router.GET("/go-slog/scores/:keeper/:size/chart.svg", scoreFunction)
	router.GET("/go-slog/warnings.html", pageFunction(pageWarnings))
	router.GET("/go-slog/guts.html", pageFunction(pageGuts))
	router.GET("/go-slog/error.html", pageFunction(pageError))
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
	warns     = data.NewWarnings()
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

	if err := data.Setup(bench, warns); err != nil {
		return fmt.Errorf("data setup: %w", err)
	}

	if err := scoring.Setup(bench, warns); err != nil {
		return fmt.Errorf("score keepers: %w", err)
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

// templateData is all the data that will be available during template execution.
type templateData struct {
	*data.Benchmarks
	*data.Warnings
	*score.Keeper
	Handler   data.HandlerTag
	Test      data.TestTag
	Keepers   []score.KeeperTag
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

// FixValue converts a score.Value into a string using the language printer.
// This will apply the proper decimal and numeric separators.
func (pd *templateData) FixValue(number score.Value) string {
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
			Keepers:    score.Keepers(),
			Levels:     warning.LevelOrder,
			Printer:    language.Printer(),
			Page:       string(page),
			Timestamp:  time.Now().Format(time.DateTime + " MST"),
		}
		tmplData.Keeper = score.GetKeeper(score.KeeperTag(c.Param("keeper")))
		switch page {
		case pageScores:
			if tmplData.Keeper == nil {
				slog.Warn("No Keeper")
				tmplData.Keeper = score.GetKeeper(keeper.DefaultName)
				if tmplData.Keeper == nil {
					slog.Error("No Keeper")
				}
			}
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
