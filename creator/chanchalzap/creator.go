package chanchalzap

import (
	"io"
	"log/slog"

	"go.mrchanchal.com/zaphandler"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/madkins23/go-slog/creator/utilzap"
	"github.com/madkins23/go-slog/infra"
)

// Creator returns a Creator object for the [chanchal/zaphandler] handler
// that wraps the [uber-go/zap] logger.
//
// [chanchal/zaphandler]: https://github.com/chanchal1987/zaphandler
// [uber-go/zap]: https://pkg.go.dev/go.uber.org/zap
func Creator() infra.Creator {
	return infra.NewCreator("chanchal/zaphandler", handlerFn, nil)
}

func handlerFn(w io.Writer, options *slog.HandlerOptions) slog.Handler {
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
				zap.NewAtomicLevelAt(utilzap.ConvertLevelToZap(level.Level())))))
}