package utillogrus

import (
	"log/slog"

	"github.com/sirupsen/logrus"
)

// ConvertSlogLevel2Logrus maps slog Levels to logrus Levels.
func ConvertSlogLevel2Logrus(l slog.Level) logrus.Level {
	switch {
	case l >= slog.LevelError:
		return logrus.ErrorLevel
	case l >= slog.LevelWarn:
		return logrus.WarnLevel
	case l >= slog.LevelInfo:
		return logrus.InfoLevel
	default:
		return logrus.DebugLevel
	}
}
