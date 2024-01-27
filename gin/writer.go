package gin

import (
	"fmt"
	"io"
	"log/slog"
	"regexp"
	"strings"
)

const (
	GinGroup  = "gin"
	NoTraffic = ""
	Traffic   = "*"
)

var TrafficMessage = "Gin Traffic"

// NewWriter returns an io.Writer object with the specified zerolog.Level.
// There are two gin output streams: gin.DefaultWriter and gin.DefaultErrorWriter.
// These streams are used by gin internal Code outside the request middleware loop.
// Create a separate Writer object with a different zerolog.Level for each stream
// or create a single object for both streams (untested but should work).
//
// If the group argument is not the empty string then Gin traffic messages of the form:
//
//	200 |    2.512908ms |             ::1 | GET      "/handler?tag=samber_zap"
//
// will be parsed into further fields which will be output within a group of the specified name,
// else the pre-formatted message lines will be output as is (which looks fine in a text logger).
// A group name of "*" results in the group contents spliced in at the top level.
func NewWriter(level slog.Leveler, group string) io.Writer {
	w := &writer{level: level}
	if group != "" {
		w.group = group
	}
	return w
}

// writer object returned by NewWriter function.
type writer struct {
	level slog.Leveler
	group string
}

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
// The data will be parsed, if possible, and converted into a slog record.
// For the moment we're assuming that there is a single Write() call for each log record.
// TODO: Fix Code to handle partial/multiple Write() calls per log record.
func (w *writer) Write(p []byte) (int, error) {
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
			msg = msg[len(matches[0]):]
		} else {
			break
		}
	}

	// Format of many of the lines at this point:
	// 2024/01/26 - 13:21:32 | 200 |  5.529751605s |             ::1 | GET      "/chart.svg?tag=With_Attrs_Attributes&item=MemBytes"
	if matches := ptnTimePrefix.FindStringSubmatch(msg); len(matches) == 2 {
		msg = matches[1]
	}

	var args []any
	var err error
	if w.group != "" {
		// Attempt to parse Gin traffic data out of the message.
		//  200 |  5.529751605s |             ::1 | GET      "/chart.svg?tag=With_Attrs_Attributes&item=MemBytes"
		if args, err = Parse(msg); err != nil {
			// The error isn't really an error, it just couldn't parse.
			slog.Debug("gin traffic parse", "err", err)
			// Just fall through now and use the pre-existing msg.
		} else {
			level = slog.LevelInfo
			msg = TrafficMessage
			if w.group != "*" {
				args = []any{slog.Group(w.group, args...)}
			}
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
