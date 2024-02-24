package slog_json

import (
	"io"
	"log/slog"

	"github.com/madkins23/go-slog/infra"
)

// SlogJson returns a Creator object for the slog/json.
func SlogJson() infra.Creator {
	return infra.NewCreator("slog/json", SlogJsonHandlerFn, nil)
}

func SlogJsonHandlerFn(w io.Writer, options *slog.HandlerOptions) slog.Handler {
	return slog.NewJSONHandler(w, options)
}
