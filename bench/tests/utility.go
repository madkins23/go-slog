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

// withAllAttributes derives a handler with the specified attributes.
// Intended for use as a HandlerFn when creating Benchmark objects.
func withAllAttributes(handler slog.Handler) slog.Handler {
	return handler.WithAttrs(withAttributes)
}

var _ HandlerFn = withGroup

// withAllAttributes derives a handler with the specified group open.
// Intended for use as a HandlerFn when creating Benchmark objects.
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
