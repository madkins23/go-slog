package samberlogrus

import (
	"io"
	"log/slog"

	slogrus "github.com/samber/slog-logrus/v2"
	"github.com/sirupsen/logrus"

	"github.com/madkins23/go-slog/creator/utillogrus"
	"github.com/madkins23/go-slog/infra"
)

// Creator returns a Creator object for the [samber/slog-logrus] handler
// that wraps the [sirupsen/logrus] logger.
//
// [samber/slog-logrus]: https://pkg.go.dev/github.com/samber/slog-logrus
// [sirupsen/logrus]: https://github.com/sirupsen/logrus
func Creator() infra.Creator {
	return infra.NewCreator("samber/slog-logrus", handlerFn, nil,
		`^samber/slog-logrus^ is a wrapper around the pre-existing ^sirupsen/logrus^ logging library.`,
		map[string]string{
			"samber/slog-logrus": "https://pkg.go.dev/github.com/samber/slog-logrus",
			"sirupsen/logrus":    "https://pkg.go.dev/github.com/sirupsen/logrus",
		})
}

func handlerFn(w io.Writer, options *slog.HandlerOptions) slog.Handler {
	level := options.Level
	if level == nil {
		level = slog.LevelInfo
	}
	log := logrus.New()
	log.SetLevel(utillogrus.ConvertSlogLevel2Logrus(level.Level()))
	log.SetOutput(w)
	log.SetFormatter(&logrus.JSONFormatter{})
	return slogrus.Option{
		Level:       level,
		Logger:      log,
		AddSource:   options.AddSource,
		ReplaceAttr: options.ReplaceAttr,
	}.NewLogrusHandler()
}
