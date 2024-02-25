package phsymzerolog

import (
	"io"
	"log/slog"

	"github.com/phsym/zeroslog"

	"github.com/madkins23/go-slog/infra"
)

// Creator returns a Creator object for the phsym/zeroslog handler.
func Creator() infra.Creator {
	return infra.NewCreator("phsym/zeroslog", handlerFn, nil, "https://github.com/phsym/zeroslog")
}

func handlerFn(w io.Writer, options *slog.HandlerOptions) slog.Handler {
	return zeroslog.NewJsonHandler(w, &zeroslog.HandlerOptions{
		Level:     options.Level,
		AddSource: options.AddSource,
	})
}
