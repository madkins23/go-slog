package slogjson

import (
	"io"
	"log/slog"

	"github.com/madkins23/go-slog/infra"
)

// Creator returns a Creator object for the slog/json.
func Creator() infra.Creator {
	return infra.NewCreator("slog/json", handlerFn, nil)
}

func handlerFn(w io.Writer, options *slog.HandlerOptions) slog.Handler {
	return slog.NewJSONHandler(w, options)
}
