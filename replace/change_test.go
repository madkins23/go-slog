package replace

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/madkins23/go-slog/handlers/flash"
	"github.com/madkins23/go-slog/handlers/sloggy/test"
	"github.com/madkins23/go-slog/internal/json"
)

const (
	alpha = "Home, Home on the Range!"
	bravo = "home, home on the range!"
)

// BenchmarkCompareChangeCase benchmarks comparing strings with two case conversions.
// This seems to take a factor of ten longer than the strings.EqualFold version.
// It also results in 2 allocs/op with 48 bytes.
func BenchmarkCompareChangeCase(b *testing.B) {
	var count uint
	for i := 0; i < b.N; i++ {
		if strings.ToUpper(alpha) == strings.ToUpper(bravo) {
			count++
		}
	}
}

// BenchmarkCompareEqualFold benchmarks comparing strings using strings.EqualFold.
// This appears to be the winner at 10% of the dual case conversion version.
// There are no memory allocations.
func BenchmarkCompareEqualFold(b *testing.B) {
	var count uint
	for i := 0; i < b.N; i++ {
		if strings.EqualFold(alpha, bravo) {
			count++
		}
	}
}

// -----------------------------------------------------------------------------

// TestChanges tests all current ReplaceAttr functions against a madkins/flash handler
// with all the flash.Extra fields set but time formatting.
func TestChanges(t *testing.T) {
	for _, lvl := range []slog.Level{slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError} {
		var buf bytes.Buffer
		log := slog.New(flash.NewHandler(&buf,
			&slog.HandlerOptions{
				AddSource: true,
				Level:     lvl,
				ReplaceAttr: Multiple(
					ChangeKey("What", slog.MessageKey, false, TopCheck),
					ChangeKey("When", slog.TimeKey, false, TopCheck),
					ChangeKey("Whence", slog.SourceKey, false, TopCheck),
					ChangeKey("Why", slog.LevelKey, false, TopCheck),
					ChangeCase(slog.LevelKey, CaseUpper, false, TopCheck),
					RemoveKey(slog.TimeKey, false, TopCheck),
					RemoveKey(slog.SourceKey, false, TopCheck),
				),
			},
			&flash.Extras{
				LevelNames: map[slog.Level]string{
					slog.LevelDebug: "Debug",
					slog.LevelInfo:  "Info",
					slog.LevelWarn:  "Warn",
					slog.LevelError: "Error",
				},
				LevelKey:   "Why",
				MessageKey: "What",
				SourceKey:  "Whence",
				TimeKey:    "When",
			}))
		switch lvl {
		case slog.LevelDebug:
			log.Debug(test.Message)
		case slog.LevelInfo:
			log.Info(test.Message)
		case slog.LevelWarn:
			log.Warn(test.Message)
		case slog.LevelError:
			log.Error(test.Message)
		}
		logMap := json.Parse(buf.String())
		_, found := logMap["When"]
		assert.False(t, found)
		_, found = logMap[slog.TimeKey]
		assert.False(t, found)
		_, found = logMap["Whence"]
		assert.False(t, found)
		_, found = logMap[slog.SourceKey]
		assert.False(t, found)
		expected := map[string]any{
			"level": strings.ToUpper(lvl.String()),
			"msg":   test.Message,
		}
		assert.Equal(t, expected, logMap)
	}
}
