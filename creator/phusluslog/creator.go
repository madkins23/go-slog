package phusluslog

import (
	"io"
	"log/slog"

	"github.com/phuslu/log"

	"github.com/madkins23/go-slog/infra"
)

// Creator returns a Creator object for the [phuslu/slog] handler.
func Creator() infra.Creator {
	return infra.NewCreator("phuslu/slog", handlerFn, nil,
		`^phuslu/slog^ is a wrapper around the pre-existing ^phuslu/log^ logging library.`,
		map[string]string{
			"phuslu/log": "https://github.com/phuslu/log",
		})
}

func handlerFn(w io.Writer, options *slog.HandlerOptions) slog.Handler {
	return log.SlogNewJSONHandler(w, options)
}
