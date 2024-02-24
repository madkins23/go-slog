package samber_logrus

import (
	"io"
	"log/slog"

	slogrus "github.com/samber/slog-logrus/v2"
	"github.com/sirupsen/logrus"

	"github.com/madkins23/go-slog/creator/util_logrus"
	"github.com/madkins23/go-slog/infra"
)

// SlogSamberLogrus returns a Creator object for the samber/logrus handler.
func SlogSamberLogrus() infra.Creator {
	return infra.NewCreator("samber/logrus", SlogSamberLogrusHandlerFn, nil)
}

func SlogSamberLogrusHandlerFn(w io.Writer, options *slog.HandlerOptions) slog.Handler {
	level := options.Level
	if level == nil {
		level = slog.LevelInfo
	}
	log := logrus.New()
	log.SetLevel(util_logrus.ConvertSlogLevel2Logrus(level.Level()))
	log.SetOutput(w)
	log.SetFormatter(&logrus.JSONFormatter{})
	return slogrus.Option{
		Level:       level,
		Logger:      log,
		AddSource:   options.AddSource,
		ReplaceAttr: options.ReplaceAttr,
	}.NewLogrusHandler()
}
