package tests

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"runtime/debug"
	"strings"
	"testing"

	"github.com/madkins23/go-slog/internal/test"
)

// -----------------------------------------------------------------------------
// Utility methods.

// -----------------------------------------------------------------------------

var _ HandlerFn = withAllAttributes

func withAllAttributes(handler slog.Handler) slog.Handler {
	return handler.WithAttrs(withAttributes)
}

var _ HandlerFn = withGroup

func withGroup(handler slog.Handler) slog.Handler {
	return handler.WithGroup("withGroup")
}

// -----------------------------------------------------------------------------

func recoverAndFailOnPanic(b *testing.B) {
	r := recover()
	failOnPanic(b, r)
}

func failOnPanic(b *testing.B, r interface{}) {
	if r != nil {
		b.Errorf("test panicked: %v\n%s", r, debug.Stack())
		b.FailNow()
	}
}

const (
	message = "This is a message"
)

// -----------------------------------------------------------------------------

// fixLogMap destructively alters the specified log map to make it match actual results.
// The fixed parseLogMap is also returned as the function's result as a convenience.
func fixLogMap(logMap map[string]any) map[string]any {
	delete(logMap, slog.TimeKey)
	if lvl, found := logMap[slog.LevelKey]; found {
		if level, ok := lvl.(string); ok {
			logMap[slog.LevelKey] = strings.ToUpper(level)
		}
	}
	if msg, found := logMap["message"]; found {
		logMap[slog.MessageKey] = msg
		delete(logMap, "message")
	}
	//if src, found := logMap[slog.SourceKey]; found {
	//	if down, ok := src.(map[string]any); ok {
	//		// Can't match the contained data over time,
	//		// just return an alphabetized array of the keys.
	//		keys := make([]string, 0, len(down))
	//		for key := range down {
	//			keys = append(keys, key)
	//		}
	//		sort.Strings(keys)
	//		logMap[slog.SourceKey] = keys
	//	}
	//}
	return logMap
}

// parseLogMap unmarshals JSON in the output capture buffer into a map[string]any.
// The buffer is sent to test logging output if the -debug=<level> flag is >= 1.
func parseLogMap(b []byte) (map[string]any, error) {
	test.Debugf(1, ">>> JSON: %s", b)
	var results map[string]any
	if err := json.Unmarshal(b, &results); err != nil {
		return results, fmt.Errorf("unmarshal bytes: %w", err)
	}
	return results, nil
}
