package creator

import (
	"io"
	"log/slog"

	"github.com/phsym/zeroslog"

	"github.com/madkins23/go-slog/infra"
)

// SlogPhsymZerolog returns a Creator object for the phsym/zeroslog handler.
func SlogPhsymZerolog() infra.Creator {
	return infra.NewCreator("phsym/zeroslog", SlogPhsymZerologHandlerFn)
}

func SlogPhsymZerologHandlerFn(w io.Writer, options *slog.HandlerOptions) *slog.Logger {
	return slog.New(zeroslog.NewJsonHandler(w, &zeroslog.HandlerOptions{
		Level:     options.Level,
		AddSource: options.AddSource,
	}))
}
