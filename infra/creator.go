package infra

import (
	"errors"
	"io"
	"log/slog"
)

type CreatorFn func(w io.Writer, options *slog.HandlerOptions) slog.Handler

// A Creator object encapsulates the creation of new slog.Handler objects.
type Creator struct {
	name string
	fn   CreatorFn
}

// Creator primarily creates slog.Handler objects.
//
// TODO: Implementations need to be thread safe?

// NewCreator returns a new Creator object for the specified CreatorFn.
func NewCreator(name string, fn CreatorFn) Creator {
	return Creator{
		name: name,
		fn:   fn,
	}
}

// NewHandle returns a new slog.Handler for the specified writer and options.
func (c *Creator) NewHandle(w io.Writer, options *slog.HandlerOptions) slog.Handler {
	return c.fn(w, options)
}

// The global Creator is private and only accessible via functions in this package.
// There is no locking around access to the global Creator.
var globalCreator *Creator

var ErrAlreadySet = errors.New("global Creator already set")

// InitGlobalCreator should be called once at the beginning of the application,
// before any global Creator interaction, to set the global Creator.
// Prior to any InitGlobalCreator call the global Creator is nil.
// Returns an error if the global Creator has already been initialized.
// Setting a nil pointer will return the same error if appropriate but otherwise do nothing.
// There is no locking around access to the global Creator.
func InitGlobalCreator(creator *Creator) error {
	if globalCreator != nil {
		return ErrAlreadySet
	}
	if creator != nil {
		globalCreator = creator
	}
	return nil
}

// GlobalCreator returns the global Creator.
// If InitGlobalCreator has not been called this result will be nil.
// There is no locking around access to the global Creator.
func GlobalCreator() *Creator {
	return globalCreator
}
