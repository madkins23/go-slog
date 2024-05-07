package svcrunnerjsonlog

import (
	"io"
	"log/slog"

	"go.seankhliao.com/svcrunner/v3/jsonlog"

	"github.com/madkins23/go-slog/infra"
)

// Creator returns a Creator object for the [svcrunner/jsonlog] handler.
func Creator() infra.Creator {
	return infra.NewCreator("svcrunner/jsonlog", handlerFn, nil,
		`^svcrunner/jsonlog^.`,
		map[string]string{
			"svcrunner/jsonlog": "https://github.com/seankhliao/svcrunner/tree/main/jsonlog",
		})
}

func handlerFn(w io.Writer, options *slog.HandlerOptions) slog.Handler {
	if options.Level == nil {
		options.Level = slog.LevelInfo
	}
	return jsonlog.New(options.Level.Level(), w)
}
