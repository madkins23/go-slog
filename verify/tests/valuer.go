package tests

import "log/slog"

// -----------------------------------------------------------------------------
// Instance of slog.LogValuer used for testing.

var _ slog.LogValuer = &hiddenValuer{}

type hiddenValuer struct {
	v any
}

func (r *hiddenValuer) LogValue() slog.Value {
	return slog.AnyValue(r.v)
}
