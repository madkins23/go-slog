package flash

import (
	"bytes"
	"log/slog"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/madkins23/go-slog/internal/json"
	"github.com/madkins23/go-slog/internal/test"
)

var levelNames = map[slog.Level]string{
	slog.LevelDebug: "Debug",
	slog.LevelInfo:  "Info",
	slog.LevelWarn:  "Warn",
	slog.LevelError: "Error",
}

var ptnJustDate = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

// TestChanges_FlashExtras verifies that the madkins/flash handler
// configured with flash.Extras options produces the expected output.
// This test covers LevelNames, LevelKey, MessageKey, SourceKey, TimeKey, and TimeFormat.
func TestFlashExtras(t *testing.T) {
	for _, lvl := range []slog.Level{slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError} {
		var buf bytes.Buffer
		log := slog.New(NewHandler(&buf,
			&slog.HandlerOptions{
				AddSource: true,
				Level:     lvl,
			},
			&Extras{
				TimeFormat: time.DateOnly,
				LevelNames: levelNames,
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
		logMap, err := json.Parse(buf.Bytes())
		assert.NoError(t, err)
		_, found := logMap[slog.TimeKey]
		assert.False(t, found)
		whenVal, found := logMap["When"]
		assert.True(t, found)
		when, ok := whenVal.(string)
		assert.True(t, ok)
		assert.Regexp(t, ptnJustDate, when)
		delete(logMap, "When")
		_, found = logMap[slog.SourceKey]
		assert.False(t, found)
		_, found = logMap["Whence"]
		assert.True(t, found)
		delete(logMap, "Whence")
		expected := map[string]any{
			"Why":  levelNames[lvl],
			"What": test.Message,
		}
		assert.Equal(t, expected, logMap)
	}
}
