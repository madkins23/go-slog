package tests

import "log/slog"

// -----------------------------------------------------------------------------
// Constant data used for tests.

const (
	message = "This is a message"
)

var logLevels = map[string]slog.Level{
	"DEBUG": slog.LevelDebug,
	"INFO":  slog.LevelInfo,
	"WARN":  slog.LevelWarn,
	"ERROR": slog.LevelError,
}
