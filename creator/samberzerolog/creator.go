package samberzerolog

import (
	"io"
	"log/slog"

	"github.com/rs/zerolog"
	samber "github.com/samber/slog-zerolog/v2"

	"github.com/madkins23/go-slog/infra"
)

// Creator returns a Creator object for the [samber/slog-zerolog] handler
// that wraps the [rs/zerolog] logger.
//
// [samber/slog-zerolog]: https://github.com/samber/slog-zerolog
// [rs/zerolog]: https://pkg.go.dev/github.com/rs/zerolog
func Creator() infra.Creator {
	return infra.NewCreator("samber/slog-zerolog", handlerFn, nil,
		`^samber/slog-zerolog^ is a wrapper around the pre-existing ^rs/zerolog^ logging library.`,
		map[string]string{
			"samber/slog-zerolog": "https://pkg.go.dev/github.com/samber/slog-zerolog",
			"rs/zerolog":          "https://github.com/rs/zerolog",
		})
}

func handlerFn(w io.Writer, options *slog.HandlerOptions) slog.Handler {
	zeroLogger := zerolog.New(w)
	return samber.Option{
		Logger:      &zeroLogger,
		Level:       options.Level,
		AddSource:   options.AddSource,
		ReplaceAttr: options.ReplaceAttr,
	}.NewZerologHandler()
}
