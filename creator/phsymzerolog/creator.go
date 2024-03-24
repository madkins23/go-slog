package phsymzerolog

import (
	"io"
	"log/slog"

	"github.com/phsym/zeroslog"

	"github.com/madkins23/go-slog/infra"
)

// Creator returns a Creator object for the [phsym/zeroslog] handler
// that wraps the [rs/zerolog] logger.
//
// [phsym/zeroslog]: https://github.com/phsym/zeroslog
// [rs/zerolog]: https://pkg.go.dev/github.com/rs/zerolog
func Creator() infra.Creator {
	return infra.NewCreator("phsym/zeroslog", handlerFn, nil,
		`^phsym/zeroslog^ is a wrapper around the pre-existing ^rs/zerolog^ logging library.
		The ^phsym/zerolog^ handler, like the underlying ^rs/zerolog^ logging library,
		is in the fastest ranks of ^slog.Handler^ implementations tested.`,
		map[string]string{
			"phsym/zeroslog": "https://pkg.go.dev/github.com/phsym/zeroslog",
			"rs/zerolog":     "https://pkg.go.dev/github.com/rs/zerolog",
		})
}

func handlerFn(w io.Writer, options *slog.HandlerOptions) slog.Handler {
	return zeroslog.NewJsonHandler(w, &zeroslog.HandlerOptions{
		Level:     options.Level,
		AddSource: options.AddSource,
	})
}
