package test

import "log/slog"

// -----------------------------------------------------------------------------

var _ slog.LogValuer = &hiddenValuer{}

type hiddenValuer struct {
	v any
}

func (r *hiddenValuer) LogValue() slog.Value {
	return slog.AnyValue(r.v)
}
