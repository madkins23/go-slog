package test

import "log/slog"

// -----------------------------------------------------------------------------

type hiddenValuer struct {
	v any
}

func (r *hiddenValuer) LogValue() slog.Value {
	return slog.AnyValue(r.v)
}
