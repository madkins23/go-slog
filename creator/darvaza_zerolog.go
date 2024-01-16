package creator

import (
	"io"
	"log/slog"

	"github.com/rs/zerolog"
	samber "github.com/samber/slog-zerolog/v2"

	"github.com/madkins23/go-slog/infra"
)

// SlogDarvazaZerolog returns a Creator object for the darvaza/zerolog handler.
func SlogDarvazaZerolog() infra.Creator {
	return &darvazaZeroCreator{CreatorData: infra.NewCreatorData("darvaza/zerolog")}
}

type darvazaZeroCreator struct {
	infra.CreatorData
}

func (c *darvazaZeroCreator) NewLogger(w io.Writer, options *slog.HandlerOptions) *slog.Logger {
	zeroLogger := zerolog.New(w)
	return slog.New(samber.Option{
		Logger:      &zeroLogger,
		Level:       options.Level,
		AddSource:   options.AddSource,
		ReplaceAttr: options.ReplaceAttr,
	}.NewZerologHandler())
}
