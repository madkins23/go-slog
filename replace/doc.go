// Package replace defines functions that can be used as slog.HandlerOptions.ReplaceAttr values.
// For example:
// * Remove attributes with empty key strings.
// * Change level attributes named "lvl" to be named slog.LevelKey.
// * Change message attributes named "message" to be named slog.MessageKey.
// * Remove the "time" basic attribute.
package replace
