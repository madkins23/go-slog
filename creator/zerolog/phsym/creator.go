package phsym

import (
	"io"
	"log/slog"

	"github.com/phsym/zeroslog"

	"github.com/madkins23/go-slog/infra"
)

// SlogPhsymZerolog returns a Creator object for the phsym/zerolog handler.
func SlogPhsymZerolog() infra.Creator {
	return infra.NewCreator("phsym/zerolog", SlogPhsymZerologHandlerFn, nil)
}

func SlogPhsymZerologHandlerFn(w io.Writer, options *slog.HandlerOptions) slog.Handler {
	return zeroslog.NewJsonHandler(w, &zeroslog.HandlerOptions{
		Level:     options.Level,
		AddSource: options.AddSource,
	})
}
