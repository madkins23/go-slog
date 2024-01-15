package ginzero

import (
	"fmt"
	"io"
	"log/slog"
	"regexp"
	"strings"
)

// Writer interface for replacing gin standard output and/or error streams.
type Writer interface {
	io.Writer
}

// NewWriter returns a Writer object with the specified zerolog.Level.
// There are two gin output streams: gin.DefaultWriter and gin.DefaultErrorWriter.
// These streams are used by gin internal code outside the request middleware loop.
// Create a separate Writer object with a different zerolog.Level for each stream
// or create a single object for both streams (untested but should work).
func NewWriter(level slog.Leveler) Writer {
	return &writer{level: level}
}

// Make sure the writer struct implements ginzero.Writer.
var _ = Writer(&writer{})

// writer object returned by NewWriter function.
type writer struct {
	// Default zerolog level for this object.
	// Can be overridden by error levels (specified in logLevels variable)
	// in square brackets at the beginning of a log record line.
	level slog.Leveler
}

var (
	logLevels = map[string]slog.Leveler{
		"DEBUG":   slog.LevelDebug,
		"ERROR":   slog.LevelError,
		"INFO":    slog.LevelInfo,
		"WARNING": slog.LevelWarn,
	}
	ptnGIN, _      = regexp.Compile(`^\s*\[GIN\]\s*`)
	ptnGINdebug, _ = regexp.Compile(`^\s*\[GIN-debug\]\s*`)
	ptnLogLevel, _ = regexp.Compile(`^\s*\[(DEBUG|ERROR|INFO|WARNING|.*)\]\s*`)
)

// Write a block of data to the (supposedly) stream object.
// For the moment we're assuming that there is a single Write() call for each log record.
// TODO: Fix code to handle multiple Write() calls per log record.
func (w *writer) Write(p []byte) (n int, err error) {
	level := w.level
	msg := strings.TrimRight(string(p), "\n")
	var sys string

	for x := 0; x < 10; x++ { // Don't use infinite for loop for safety
		// Pull off prefix sequences that represent log information.
		if match := ptnGIN.FindString(msg); match != "" {
			msg = msg[len(match):]
			sys = "gin"
		} else if match := ptnGINdebug.FindString(msg); match != "" {
			level = slog.LevelDebug
			msg = msg[len(match):]
			sys = "gin"
		} else if matches := ptnLogLevel.FindStringSubmatch(msg); len(matches) > 1 {
			var ok bool
			if level, ok = logLevels[matches[1]]; !ok {
				return 0, fmt.Errorf("no level %s", matches[1])
			}
			msg = msg[len(matches[0]):]
		} else {
			break
		}
	}

	var args []any
	if sys != "" {
		args = append(args, "sys", sys)
	}

	switch level {
	case slog.LevelDebug:
		slog.Debug(msg, args...)
	case slog.LevelError:
		slog.Error(msg, args...)
	case slog.LevelInfo:
		slog.Info(msg, args...)
	case slog.LevelWarn:
		slog.Warn(msg, args...)
	default:
		return 0, fmt.Errorf("unknown log level %s", w.level)
	}

	return len(p), nil
}
