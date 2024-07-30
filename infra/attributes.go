package infra

import "log/slog"

// AttrFn defines a type for ReplaceAttr functions.
// The slog.HandlerOptions struct defines this inline without defining a type.
type AttrFn func(groups []string, a slog.Attr) slog.Attr

// EmptyAttr returns an empty attribute object as a convenience.
// An empty attribute has both an empty key and a nil value.
func EmptyAttr() slog.Attr {
	return slog.Attr{}
}
