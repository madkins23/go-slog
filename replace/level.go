package replace

import (
	"log/slog"
	"strings"
)

var _ AttrFn = LvlToLevel

// LvlToLevel replaces attribute keys matching "lvl" with the correct slog.LevelKey.
func LvlToLevel(groups []string, a slog.Attr) slog.Attr {
	if strings.ToLower(a.Key) == "lvl" && len(groups) == 0 {
		a.Key = slog.LevelKey
	}
	return a
}

// LevelCase changes the values of "level" attributes to uppercase.
func LevelCase(groups []string, a slog.Attr) slog.Attr {
	if strings.ToLower(a.Key) == "level" && len(groups) == 0 {
		return slog.String(a.Key, strings.ToUpper(a.String()))
	}

	return a
}
