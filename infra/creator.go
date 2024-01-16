package infra

import (
	"io"
	"log/slog"
)

// A CreatorFn is a function that can create new slog.Handler objects.
type CreatorFn func(w io.Writer, options *slog.HandlerOptions) *slog.Logger

// A Creator object encapsulates the creation of new slog.Handler objects.
// This includes both the name of the handler and a CreatorFn.
type Creator struct {
	name string
	fn   CreatorFn
}

// NewCreator returns a new Creator object for the specified name and CreatorFn.
func NewCreator(name string, fn CreatorFn) Creator {
	return Creator{
		name: name,
		fn:   fn,
	}
}

// NewLogger returns a new slog.Logger object.
// The actual creation is done by invoking the embedded CreatorFn.
func (c *Creator) NewLogger(w io.Writer, options *slog.HandlerOptions) *slog.Logger {
	return c.fn(w, options)
}

// Name returns the name of the slog package.
func (c *Creator) Name() string {
	return c.name
}
