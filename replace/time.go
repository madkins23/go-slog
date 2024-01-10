package replace

import "log/slog"

// RemoveTime removes the top-level time attribute.
// It is intended to be used as a ReplaceAttr function, to make example output deterministic.
// This code was lifted from slogtest/internal: https://pkg.go.dev/log/slog/internal/slogtest#RemoveTime
func RemoveTime(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.TimeKey && len(groups) == 0 {
		return emptyAttr
	}
	return a
}
