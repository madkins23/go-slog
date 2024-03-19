package replace

import (
	"log/slog"
)

// -----------------------------------------------------------------------------
// ReplaceAttr functions related to the "level" field.

// LvlToLevel replaces attribute keys matching "lvl" with the correct slog.LevelKey.
//
// Deprecated: This function was never tested and can be replaced by
//
//	replace.ChangeKey("lvl", slog.LevelKey, false, replace.TopCheck)
//
// which is how it is now implemented.
func LvlToLevel(groups []string, a slog.Attr) slog.Attr {
	return ChangeKey("lvl", slog.LevelKey, false, TopCheck)(groups, a)
}

// -----------------------------------------------------------------------------

// LevelLowerCase changes the values of "level" attributes to lowercase.
//
// Deprecated: This function can be replaced by
//
//	replace.ChangeCase("level", CaseLower, false, TopCheck)
//
// which is how it is now implemented.
func LevelLowerCase(groups []string, a slog.Attr) slog.Attr {
	return ChangeCase("level", CaseLower, false, TopCheck)(groups, a)
}

// LevelUpperCase changes the values of "level" attributes to uppercase.
// Based on the existing behavior of log/slog this is the correct output.
//
// Deprecated: This function can be replaced by
//
//	replace.ChangeCase("level", CaseUpper, false, TopCheck)
//
// which is how it is now implemented.
func LevelUpperCase(groups []string, a slog.Attr) slog.Attr {
	return ChangeCase("level", CaseUpper, false, TopCheck)(groups, a)
}
