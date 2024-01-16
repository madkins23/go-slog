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
	return &samberZapCreator{CreatorData: infra.NewCreatorData("samber/slog-zap")}
}

type samberZapCreator struct {
	infra.CreatorData
}

func (c *samberZapCreator) NewLogger(w io.Writer, options *slog.HandlerOptions) *slog.Logger {
	level := options.Level
	if level == nil {
		level = slog.LevelInfo
	}
	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.TimeKey = "time"
	productionCfg.EncodeTime = zapcore.RFC3339TimeEncoder
	return slog.New(slogzap.Option{
		Level: level,
		Logger: zap.New(zapcore.NewCore(
			zapcore.NewJSONEncoder(productionCfg),
			zapcore.AddSync(w),
			zap.NewAtomicLevelAt(convertSlogLevel(level.Level())))),
		Converter:   nil,
		AddSource:   options.AddSource,
		ReplaceAttr: options.ReplaceAttr,
	}.NewZapHandler())
}

// convertSlogLevel maps slog Levels to zap Levels.
//
// Note: copied from https://github.com/uber-go/zap/blob/d27427d23f81dba1f048d6034d5f286572049e1e/exp/zapslog/handler.go
//
// Note: there is some room between slog levels while zap levels are continuous, so we can't 1:1 map them.
// See also https://go.googlesource.com/proposal/+/master/design/56345-structured-logging.md?pli=1#levels
func convertSlogLevel(l slog.Level) zapcore.Level {
	switch {
	case l >= slog.LevelError:
		return zapcore.ErrorLevel
	case l >= slog.LevelWarn:
		return zapcore.WarnLevel
	case l >= slog.LevelInfo:
		return zapcore.InfoLevel
	default:
		return zapcore.DebugLevel
	}
}
