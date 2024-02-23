package creator

import (
	"io"
	"log/slog"

	"go.mrchanchal.com/zaphandler"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/madkins23/go-slog/infra"
)

// SlogChanchalZapHandler returns a Creator object for the samber/slog-zerolog handler.
func SlogChanchalZapHandler() infra.Creator {
	return infra.NewCreator("chanchal/ZapHandler", SlogChanchalZapHandlerFn, nil)
}

func SlogChanchalZapHandlerFn(w io.Writer, options *slog.HandlerOptions) slog.Handler {
	level := options.Level
	if level == nil {
		level = slog.LevelInfo
	}
	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.TimeKey = "time"
	productionCfg.EncodeTime = zapcore.RFC3339TimeEncoder
	return zaphandler.New(
		zap.New(
			zapcore.NewCore(
				zapcore.NewJSONEncoder(productionCfg),
				zapcore.AddSync(w),
				zap.NewAtomicLevelAt(convertLevelToZap(level.Level())))))
}