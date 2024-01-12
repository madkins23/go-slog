package infra

import (
	"errors"
	"io"
	"log/slog"
)

type CreatorFn func(w io.Writer, options *slog.HandlerOptions) slog.Handler

// Creator primarily creates slog.Handler objects.
//
// Note: Implementations need to be thread safe.
//
// GOFU: Currently the Creator is just a CreatorFn.
// This may change (e.g. CreatorObjX below) in the future.
// Some nice (or questionable) Go fu is used to hide the function-ness
// of the Creator behind a facade of object-ness.
//
// TODO: Use it (replace with a struct or interface) or lose it.
type Creator CreatorFn

// NewHandle returns a new slog.Handler for the specified writer and options.
func (fn Creator) NewHandle(w io.Writer, options *slog.HandlerOptions) slog.Handler {
	// GOFU: Currently Creator is actually a func so invoke itself on the provided arguments.
	return fn(w, options)
}

// NewCreator returns a new Creator object for the specified CreatorFn.
//
// GOFU: Currently this is done by simply casting the CreatorFn to a Creator.
// If CreatorObjX replaces CreatorFn this will instantiate a CreatorObjX and
// store the CreatorFn value on a field.
func NewCreator(creatorFn CreatorFn) Creator {
	return Creator(creatorFn)
}

// The global Creator is private and only accessible via functions in this package.
// There is no locking around access to the global Creator.
var globalCreator Creator

var ErrAlreadySet = errors.New("global Creator already set")

// InitGlobalCreator should be called once at the beginning of the application,
// before any global Creator interaction, to set the global Creator.
// Prior to any InitGlobalCreator call the global Creator is nil.
// Returns an error if the global Creator has already been initialized.
// Setting a nil pointer will return the same error if appropriate but otherwise do nothing.
// There is no locking around access to the global Creator.
func InitGlobalCreator(creator Creator) error {
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
func GlobalCreator() Creator {
	return globalCreator
}

// -----------------------------------------------------------------------------

// Note: This would be a replacement for tests.CreateHandlerFn.
// Seems only necessary if there would be more fields than that one.
// TODO: Use it (change CreatorObjX to Creator) or lose it.

// A CreatorObjX object encapsulates the creation of new slog.Handler objects.
type CreatorObjX interface {
	NewHandle(w io.Writer, options *slog.HandlerOptions) slog.Handler
}
