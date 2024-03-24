package slogjson

import (
	"io"
	"log/slog"

	"github.com/madkins23/go-slog/infra"
)

// Creator returns a Creator object for the [slog/JSONHandler] handler.
//
// [slog/JSONHandler]: https://pkg.go.dev/log/slog#JSONHandler
func Creator() infra.Creator {
	return infra.NewCreator("slog/JSONHandler", handlerFn, nil,
		`^slog/JSONHandler^ is the JSON handler provided with the ^slog^ library.
		It is fast and as a part of the Go distribution it is used
		along with published documentation as a model for ^slog.Handler^ behavior.`,
		map[string]string{
			"slog/JSONHandler": "https://pkg.go.dev/log/slog#JSONHandler",
		})
}

func handlerFn(w io.Writer, options *slog.HandlerOptions) slog.Handler {
	return slog.NewJSONHandler(w, options)
}
