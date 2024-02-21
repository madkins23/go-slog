package creator

import (
	"io"
	"log/slog"

	slogzap "github.com/samber/slog-zap/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/madkins23/go-slog/infra"
)

// SlogSamberZap returns a Creator object for the samber/slog-zerolog handler.
func SlogSamberZap() infra.Creator {
	return infra.NewCreator("samber/slog-zap", SlogSamberZapHandlerFn, nil)
}

func SlogSamberZapHandlerFn(w io.Writer, options *slog.HandlerOptions) slog.Handler {
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
			zap.NewAtomicLevelAt(convertLevelToZap(level.Level())))),
		Converter:   nil,
		AddSource:   options.AddSource,
		ReplaceAttr: options.ReplaceAttr,
	}.NewZapHandler()
}
