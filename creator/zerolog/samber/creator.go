package samber

import (
	"io"
	"log/slog"

	"github.com/rs/zerolog"
	samber "github.com/samber/slog-zerolog/v2"

	"github.com/madkins23/go-slog/infra"
)

// SlogSamberZerolog returns a Creator object for the samber/zerolog handler.
func SlogSamberZerolog() infra.Creator {
	return infra.NewCreator("samber/zerolog", SlogSamberZerologHandlerFn, nil)
}

func SlogSamberZerologHandlerFn(w io.Writer, options *slog.HandlerOptions) slog.Handler {
	zeroLogger := zerolog.New(w)
	return samber.Option{
		Logger:      &zeroLogger,
		Level:       options.Level,
		AddSource:   options.AddSource,
		ReplaceAttr: options.ReplaceAttr,
	}.NewZerologHandler()
}
