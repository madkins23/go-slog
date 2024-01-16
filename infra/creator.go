package infra

import (
	"io"
	"log/slog"
)

// A Creator object encapsulates the creation of new slog.Handler objects.
// This includes both the name of the handler and a CreatorFn.
type Creator interface {
	CreatorData

	// NewLogger returns a new slog.Logger for the specified writer and options.
	NewLogger(w io.Writer, options *slog.HandlerOptions) *slog.Logger
}

// -----------------------------------------------------------------------------

// A CreatorData object encapsulates any ancillary Creator data.
type CreatorData interface {
	// Name returns the name of the slog.Logger to be created.
	Name() string
}

// creatorData instantiates the CreatorData interface.
type creatorData struct {
	name string
}

// NewCreatorData returns a new Creator object for the specified CreatorFn.
func NewCreatorData(name string) CreatorData {
	return &creatorData{
		name: name,
	}
}

// Name returns the name of the slog.Logger to be created.
func (c *creatorData) Name() string {
	return c.name
}
