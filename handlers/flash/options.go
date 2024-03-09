package flash

import "log/slog"

// Options defines options for a flash.Handler.
// The core of Options is slog.HandlerOptions.
// Additional options particular to flash may be added over time.
type Options struct {
	*slog.HandlerOptions
	// Room for expansion.
}

type Option func(*Options)

// MakeOptions creates a flash.Options object from a slog.HandlerOptions object.
// This is the recommended way to acquire a new flash.Options object.
func MakeOptions(sho *slog.HandlerOptions, opts ...Option) *Options {
	options := &Options{HandlerOptions: sho}
	for _, opt := range opts {
		opt(options)
	}
	return options
}
