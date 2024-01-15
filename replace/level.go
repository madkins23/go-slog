package replace

import (
	"log/slog"
	"strings"

	"github.com/madkins23/go-slog/infra"
)

// -----------------------------------------------------------------------------
// ReplaceAttr functions related to the "level" field.

var _ infra.AttrFn = LvlToLevel

// LvlToLevel replaces attribute keys matching "lvl" with the correct slog.LevelKey.
func LvlToLevel(groups []string, a slog.Attr) slog.Attr {
	if strings.ToLower(a.Key) == "lvl" && len(groups) == 0 {
		a.Key = slog.LevelKey
	}
	return a
}

// -----------------------------------------------------------------------------

var _ infra.AttrFn = LevelLowerCase

// LevelLowerCase changes the values of "level" attributes to lowercase.
func LevelLowerCase(groups []string, a slog.Attr) slog.Attr {
	if strings.ToLower(a.Key) == "level" && len(groups) == 0 {
		return slog.String(a.Key, strings.ToLower(a.Value.String()))
	}

	return a
}

var _ infra.AttrFn = LevelUpperCase

// LevelUpperCase changes the values of "level" attributes to uppercase.
// Based on the existing behavior of log/slog this is the correct output.
func LevelUpperCase(groups []string, a slog.Attr) slog.Attr {
	if strings.ToLower(a.Key) == "level" && len(groups) == 0 {
		return slog.String(a.Key, strings.ToUpper(a.Value.String()))
	}

	return a
}
