package samberzap

import (
	"io"
	"log/slog"

	slogzap "github.com/samber/slog-zap/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/madkins23/go-slog/creator/utilzap"
	"github.com/madkins23/go-slog/infra"
)

// Creator returns a Creator object for the [samber/slog-zap] handler
// that wraps the [uber-go/zap] logger.
//
// [samber/slog-zap]: https://pkg.go.dev/github.com/samber/slog-zap
// [uber-go/zap]: https://pkg.go.dev/go.uber.org/zap
func Creator() infra.Creator {
	return infra.NewCreator("samber/slog-zap", handlerFn, nil,
		`^samber/slog-zap^ is a wrapper around the pre-existing ^uber-go/zap^ logging library.`,
		map[string]string{
			"samber/slog-zap": "https://pkg.go.dev/github.com/samber/slog-zap",
			"uber-go/zap":     "https://pkg.go.dev/go.uber.org/zap",
		})
}

func handlerFn(w io.Writer, options *slog.HandlerOptions) slog.Handler {
	level := options.Level
	if level == nil {
		level = slog.LevelInfo
	}
	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.TimeKey = "time"
	productionCfg.EncodeTime = zapcore.RFC3339TimeEncoder
	return slogzap.Option{
		Level: level,
		Logger: zap.New(zapcore.NewCore(
			zapcore.NewJSONEncoder(productionCfg),
			zapcore.AddSync(w),
			zap.NewAtomicLevelAt(utilzap.ConvertLevelToZap(level.Level())))),
		Converter:   nil,
		AddSource:   options.AddSource,
		ReplaceAttr: options.ReplaceAttr,
	}.NewZapHandler()
}
