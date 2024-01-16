package creator

import (
	"io"
	"log/slog"

	"github.com/rs/zerolog"
	samber "github.com/samber/slog-zerolog/v2"

	"github.com/madkins23/go-slog/infra"
)

// SlogSamberZerolog returns a Creator object for the samber/slog-zerolog handler.
func SlogSamberZerolog() infra.Creator {
	return &samberZeroCreator{CreatorData: infra.NewCreatorData("samber/slog-zerolog")}
}

type samberZeroCreator struct {
	infra.CreatorData
}

func (c *samberZeroCreator) NewLogger(w io.Writer, options *slog.HandlerOptions) *slog.Logger {
	zeroLogger := zerolog.New(w)
	return slog.New(samber.Option{
		Logger:      &zeroLogger,
		Level:       options.Level,
		AddSource:   options.AddSource,
		ReplaceAttr: options.ReplaceAttr,
	}.NewZerologHandler())
}
