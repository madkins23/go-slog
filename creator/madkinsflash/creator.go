package madkinsflash

import (
	"io"
	"log/slog"

	"github.com/madkins23/go-slog/handlers/flash"
	"github.com/madkins23/go-slog/infra"
)

// Creator returns a Creator object for the [madkins/flash] handler.
// This is an experimental handler development based on madkins/sloggy.
func Creator() infra.Creator {
	return infra.NewCreator("madkins/flash", handlerFn, nil)
}

func handlerFn(w io.Writer, options *slog.HandlerOptions) slog.Handler {
	return flash.NewHandler(w, options, nil)
}
