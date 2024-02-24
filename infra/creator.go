package infra

import (
	"io"
	"log/slog"
	"os"
)

// CreateLoggerFn is a function that can create new slog.Logger objects.
type CreateLoggerFn func(w io.Writer, options *slog.HandlerOptions) *slog.Logger

// CreateHandlerFn is a function that can create new slog.Handler objects.
type CreateHandlerFn func(w io.Writer, options *slog.HandlerOptions) slog.Handler

// A Creator object encapsulates the creation of new slog.Handler objects.
// This includes both the name of the handler and a CreateLoggerFn.
type Creator struct {
	name      string
	handlerFn CreateHandlerFn
	loggerFn  CreateLoggerFn
	docURL    string
}

// NewCreator returns a new Creator object for the specified name and CreateLoggerFn.
func NewCreator(name string, handlerFn CreateHandlerFn, loggerFn CreateLoggerFn, docURL string) Creator {
	if handlerFn == nil && loggerFn == nil {
		slog.Error("Creator must have either handlerFn or loggerFn")
		os.Exit(1)
	}
	return Creator{
		name:      name,
		handlerFn: handlerFn,
		loggerFn:  loggerFn,
		docURL:    docURL,
	}
}

// NewLogger returns a new slog.Logger object.
// The actual creation is done by invoking the embedded CreateLoggerFn,
// if it is non-nil, or the embedded CreateHandlerFn.
func (c *Creator) NewLogger(w io.Writer, options *slog.HandlerOptions) *slog.Logger {
	if c.loggerFn != nil {
		return c.loggerFn(w, options)
	} else {
		return slog.New(c.handlerFn(w, options))
	}
}

// NewHandler returns a new slog.Handler object.
// The actual creation is done by invoking the embedded CreateHandlerFn.
func (c *Creator) NewHandler(w io.Writer, options *slog.HandlerOptions) slog.Handler {
	if c.handlerFn != nil {
		return c.handlerFn(w, options)
	} else {
		return nil
	}
}

func (c *Creator) CanMakeHandler() bool {
	return c.handlerFn != nil
}

// Name returns the name of the slog package.
func (c *Creator) Name() string {
	return c.name
}

// DocURL returns the documentation URL of the slog package.
func (c *Creator) DocURL() string {
	return c.docURL
}
