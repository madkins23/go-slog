package gin

import (
	"fmt"
	"io"
	"log/slog"
	"regexp"
	"strings"
)

const (
	// DefaultTrafficGroup is the default group for traffic data if parsing is required and no group is provided.
	DefaultTrafficGroup = "gin"

	// DefaultTrafficMessage to be logged when the original log message is
	// a Gin traffic data line which has been parsed under other field names.
	DefaultTrafficMessage = "Gin Traffic"
)

// ----------------------------------------------------------------------------

// Options for NewWriter.
type Options struct {
	// Level sets the default log level for everything except Gin Traffic messages.
	Level slog.Leveler

	// Traffic collects options related to parsing Gin traffic data records.
	Traffic Traffic
}

// Traffic collects options related to parsing Gin traffic data records.
// Records have the form:
//
//	200 |    2.512908ms |             ::1 | GET      "/handler?tag=samber_zap"
type Traffic struct {
	// Parse Gin traffic data records when true.
	// Nothing else in this struct matters if this is false.
	Parse bool

	// Embed parsed traffic data at the top level of the log message when true.
	// Options data will be mixed in with other log data.
	// Only used if Parse is true.
	Embed bool

	// Level is the slog.Level for traffic data lines (defaults to slog.LevelInfo).
	// Overrides the level attached to the writer object for relevant messages.
	// Only used if Parse is true.
	Level slog.Leveler

	// Message to be used when traffic data has been parsed from the original message.
	// When no message is provided the default value of DefaultTrafficMessage is used.
	// Only used if Parse is true.
	Message string

	// Group provides a group name under which traffic data will be gathered.
	// When no group name is provided the value of DefaultTrafficGroup is used.
	// Only used if Parse is true and Embed is false.
	Group string
}

// ----------------------------------------------------------------------------

// NewWriter returns an io.Writer object with the specified zerolog.Level.
// There are two gin output streams: gin.DefaultWriter and gin.DefaultErrorWriter.
// These streams are used by gin internal Code outside the request middleware loop.
// Create a separate Writer object with a different zerolog.Level for each stream
// or create a single object for both streams (untested but should work).
//
// The options argument holds settings for the underlying writer object.
// See the documentation for Options.
func NewWriter(options *Options) io.Writer {
	w := &writer{
		Options: *options,
	}
	if w.Level == nil {
		w.Level = slog.LevelInfo
	}
	if w.Traffic.Parse {
		// Fix some default values here, so it doesn't have to be done in Write().
		if w.Traffic.Level == nil {
			w.Traffic.Level = slog.LevelInfo
		}
		if w.Traffic.Group == "" && !w.Traffic.Embed {
			w.Traffic.Group = DefaultTrafficGroup
		}
		if w.Traffic.Message == "" {
			w.Traffic.Message = DefaultTrafficMessage
		}
	}
	return w
}

// ----------------------------------------------------------------------------

// writer object returned by NewWriter function.
type writer struct {
	Options
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
// TODO: Fix code to handle partial/multiple Write() calls per log record.
func (w *writer) Write(p []byte) (int, error) {
	level := w.Level
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
		// Remove time prefix.
		msg = matches[1]
	}

	var args []any
	var err error
	if w.Traffic.Parse {
		// Attempt to parse Gin traffic data out of the message.
		//  200 |  5.529751605s |             ::1 | GET      "/chart.svg?tag=With_Attrs_Attributes&item=MemBytes"
		if args, err = Parse(msg); err != nil {
			// The error isn't really an error, it just couldn't parse.
			slog.Debug("gin traffic parse", "err", err)
			// Use the pre-existing log message in variable msg.
		} else {
			// The traffic data in the log message parsed.
			level = slog.LevelInfo
			msg = w.Traffic.Message
			if !w.Traffic.Embed {
				args = []any{slog.Group(w.Traffic.Group, args...)}
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
		// Shouldn't ever happen so no test code coverage here.
		return 0, fmt.Errorf("unknown log level %s", w.Level)
	}

	return len(p), nil
}
