package madkinssloggy

import (
	"io"
	"log/slog"

	"github.com/madkins23/go-slog/handlers/sloggy"
	"github.com/madkins23/go-slog/infra"
)

// Creator returns a Creator object for the [madkins/sloggy] handler.
// This is an experimental handler development.
func Creator() infra.Creator {
	return infra.NewCreator("madkins/sloggy", handlerFn, nil)
}

func handlerFn(w io.Writer, options *slog.HandlerOptions) slog.Handler {
	return sloggy.NewHandler(w, options)
}
