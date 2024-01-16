package creator

import (
	"io"
	"log/slog"

	"github.com/madkins23/go-slog/infra"
)

// Slog returns a Creator object for the log/slog.JSONHandler.
func Slog() infra.Creator {
	return infra.NewCreator("log/slog.JSONHandler", SlogHandlerFn)
}

func SlogHandlerFn(w io.Writer, options *slog.HandlerOptions) *slog.Logger {
	return slog.New(slog.NewJSONHandler(w, options))
}
