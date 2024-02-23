package slog

import (
	"io"
	"log/slog"

	"github.com/madkins23/go-slog/infra"
)

// Slog returns a Creator object for the slog/json.
func Slog() infra.Creator {
	return infra.NewCreator("slog/json", SlogHandlerFn, nil)
}

func SlogHandlerFn(w io.Writer, options *slog.HandlerOptions) slog.Handler {
	return slog.NewJSONHandler(w, options)
}
