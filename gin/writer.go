package gin

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
	Parser
}

// NewWriter returns a Writer object with the specified zerolog.Level.
// There are two gin output streams: gin.DefaultWriter and gin.DefaultErrorWriter.
// These streams are used by gin internal code outside the request middleware loop.
// Create a separate Writer object with a different zerolog.Level for each stream
// or create a single object for both streams (untested but should work).
// If the parse flag is non-nil then Gin traffic lines:
//
//	200 |    2.512908ms |             ::1 | GET      "/handler?tag=samber_zap" system=gin
//
// will be parsed into further fields, otherwise that data will be the log message.
// In addition, it serves as a lookup to convert the output fields to other field names.
// This facility does not extend to renaming into groups.
// Use Empty() to specify default parsing behavior and field names.
func NewWriter(level slog.Leveler, parseTraffic FieldString) Writer {
	w := &writer{level: level}
	if parseTraffic != nil {
		w.Parser = NewParser(parseTraffic)
	}
	return w
}

// Make sure the writer struct implements ginzero.Writer.
var _ Writer = &writer{}

// writer object returned by NewWriter function.
type writer struct {
	// Empty zerolog level for this object.
	// Can be overridden by error levels (specified in logLevels variable)
	// in square brackets at the beginning of a log record line.
	level slog.Leveler
	Parser
	Group string
}

const groupName = "gin"

var (
	logLevels = map[string]slog.Leveler{
		"DEBUG":   slog.LevelDebug,
		"ERROR":   slog.LevelError,
		"INFO":    slog.LevelInfo,
		"WARNING": slog.LevelWarn,
	}
	ptnGIN, _      = regexp.Compile(`^\s*\[GIN]\s*`)
	ptnGINdebug, _ = regexp.Compile(`^\s*\[GIN-debug]\s*`)
	ptnLogLevel, _ = regexp.Compile(`^\s*\[(DEBUG|ERROR|INFO|WARNING|.*)]\s*`)
	ptnTimePrefix  = regexp.MustCompile(`^\s*\d+/\d+/\d+\s*-\s*\d+:\d+:\d+\s*\|\s*(.+)$`)
)

// Write a block of data to the (supposedly) stream object.
// For the moment we're assuming that there is a single Write() call for each log record.
// TODO: Fix code to handle multiple Write() calls per log record.
func (w *writer) Write(p []byte) (n int, err error) {
	level := w.level
	msg := strings.TrimRight(string(p), "\n")
	for x := 0; x < 10; x++ { // Don't use infinite for loop for safety
		// Pull off prefix sequences that represent log information.
		if match := ptnGIN.FindString(msg); match != "" {
			msg = msg[len(match):]
		} else if match := ptnGINdebug.FindString(msg); match != "" {
			level = slog.LevelDebug
			msg = msg[len(match):]
		} else if matches := ptnLogLevel.FindStringSubmatch(msg); len(matches) > 1 {
			var ok bool
			if level, ok = logLevels[matches[1]]; !ok {
				return 0, fmt.Errorf("no level %s", matches[1])
			}
			msg = matches[0]
		} else {
			break
		}
	}

	// Format of many of the lines:
	// 2024/01/20 - 07:18:37 | 200 |    1.226682ms |             ::1 | GET      "/"
	if matches := ptnTimePrefix.FindStringSubmatch(msg); len(matches) == 2 {
		msg = matches[1]
	}

	var args []any
	if w.Parser != nil {
		// Empty the traffic record out of the message:
		if w.Group == "" {
			w.Group = groupName
		}
		if args, err := w.Parse(msg); err != nil {
			slog.Warn("Unable to parse Gin traffic", "err", err)
		} else if w.Group != "*" {
			args = []any{slog.Group(w.Group, args)}
		}
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
