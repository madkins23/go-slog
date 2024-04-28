package infra

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"math"

	"github.com/madkins23/go-slog/internal/json"
)

func ExampleCreator() {
	creator := NewCreator(
		"example",
		func(w io.Writer, options *slog.HandlerOptions) slog.Handler {
			return slog.NewJSONHandler(w, options)
		},
		nil, // loggerFn optional if handlerFn provided
		// summary and links optional
		"", nil)
	var buffer bytes.Buffer
	logger := creator.NewLogger(&buffer, nil)
	logger.Info("message", "pi", math.Pi)
	logMap, err := json.Parse(buffer.Bytes())
	if err == nil {
		fmt.Printf("msg:%s pi:%7.5f\n", logMap[slog.MessageKey], logMap["pi"])
	} else {
		fmt.Printf("!!! %s\n", err.Error())
	}
	// Output: msg:message pi:3.14159
}
