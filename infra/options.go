package infra

import (
	"log/slog"
)

// SimpleOptions returns a default, simple, slog.HandlerOptions.
func SimpleOptions() *slog.HandlerOptions {
	return &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
}

// LevelOptions returns a slog.HandlerOptions with the specified level.
func LevelOptions(level slog.Leveler) *slog.HandlerOptions {
	return &slog.HandlerOptions{
		Level: level,
	}
}

// SourceOptions returns a slog.HandlerOptions with the specified level
// and the AddSource field set to true.
func SourceOptions() *slog.HandlerOptions {
	return &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	}
}

// ReplaceAttrOptions returns a slog.HandlerOptions with the specified ReplaceAttr function.
func ReplaceAttrOptions(fn AttrFn) *slog.HandlerOptions {
	return &slog.HandlerOptions{
		Level:       slog.LevelInfo,
		ReplaceAttr: fn,
	}
}
