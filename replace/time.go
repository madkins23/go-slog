package replace

import (
	"log/slog"
)

// RemoveTime removes the top-level time attribute.
// It is intended to be used as a ReplaceAttr function, to make example output deterministic.
// This code was lifted from slogtest/internal: https://pkg.go.dev/log/slog/internal/slogtest#RemoveTime
//
// Deprecated: This function was never tested and can be replaced by
//
//	replace.RemoveKey(slog.TimeKey, false, replace.TopCheck)
//
// which is how it is now implemented.
func RemoveTime(groups []string, a slog.Attr) slog.Attr {
	return RemoveKey(slog.TimeKey, false, TopCheck)(groups, a)
}
