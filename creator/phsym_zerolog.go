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

var _ infra.CreatorFn = SlogPhsymZerologHandlerFn

// SlogPhsymZerologHandlerFn returns a new slog.Handler for phsym/zeroslog.
func SlogPhsymZerologHandlerFn(w io.Writer, options *slog.HandlerOptions) slog.Handler {
	return zeroslog.NewJsonHandler(w, &zeroslog.HandlerOptions{
		Level:     options.Level,
		AddSource: options.AddSource,
	})
}
