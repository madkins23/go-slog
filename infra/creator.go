package infra

import (
	"io"
	"log/slog"
)

// A CreatorFn is a function that can create new slog.Handler objects.
type CreatorFn func(w io.Writer, options *slog.HandlerOptions) slog.Handler

// A Creator object encapsulates the creation of new slog.Handler objects.
// This includes both the name of the handler and a CreatorFn.
type Creator struct {
	name string
	fn   CreatorFn
}

// NewCreator returns a new Creator object for the specified CreatorFn.
func NewCreator(name string, fn CreatorFn) Creator {
	return Creator{
		name: name,
		fn:   fn,
	}
}

// NewHandle returns a new slog.Handler for the specified writer and options.
// The actual creation is done by invoking the embedded CreatorFn.
func (c *Creator) NewHandle(w io.Writer, options *slog.HandlerOptions) slog.Handler {
	return c.fn(w, options)
}

// Name returns the name of the handler.
func (c *Creator) Name() string {
	return c.name
}
