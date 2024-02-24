package samber_zap

import (
	"io"
	"log/slog"

	slogzap "github.com/samber/slog-zap/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/madkins23/go-slog/creator/util_zap"
	"github.com/madkins23/go-slog/infra"
)

// SlogSamberZap returns a Creator object for the samber/zap handler.
func SlogSamberZap() infra.Creator {
	return infra.NewCreator("samber/zap", SlogSamberZapHandlerFn, nil)
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
			zap.NewAtomicLevelAt(util_zap.ConvertLevelToZap(level.Level())))),
		Converter:   nil,
		AddSource:   options.AddSource,
		ReplaceAttr: options.ReplaceAttr,
	}.NewZapHandler()
}
