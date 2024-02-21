package creator

import (
	"log/slog"

	"go.uber.org/zap/zapcore"
)

// convertLevelToZap maps slog Levels to zap Levels.
//
// Note: copied from https://github.com/uber-go/zap/blob/d27427d23f81dba1f048d6034d5f286572049e1e/exp/zapslog/handler.go
//
// Note: there is some room between slog levels while zap levels are continuous, so we can't 1:1 map them.
// See also https://go.googlesource.com/proposal/+/master/design/56345-structured-logging.md?pli=1#levels
func convertLevelToZap(l slog.Level) zapcore.Level {
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
