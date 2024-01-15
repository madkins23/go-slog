package infra

import "log/slog"

// AttrFn defines a type for ReplaceAttr functions.
// The slog.HandlersOptions struct defines this inline without defining a type.
type AttrFn func(groups []string, a slog.Attr) slog.Attr

// emptyAttr defines a package-internal empty attribute as a convenience.
var emptyAttr = slog.Attr{}

// EmptyAttr returns an empty attribute object as a convenience.
func EmptyAttr() slog.Attr {
	return emptyAttr
}
