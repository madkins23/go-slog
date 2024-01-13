package creator

import (
	"io"
	"log/slog"

	"github.com/madkins23/go-slog/infra"
)

// Slog returns a Creator object for the log/slog.JSONHandler.
func Slog() infra.Creator {
	return infra.NewCreator("log/slog.JSONHandler", slogHandlerFn)
}

var _ infra.CreatorFn = slogHandlerFn

// slogHandlerFn returns a new slog.Handler for log/slog.NewJSONHandler.
func slogHandlerFn(w io.Writer, options *slog.HandlerOptions) slog.Handler {
	return slog.NewJSONHandler(w, options)
}
