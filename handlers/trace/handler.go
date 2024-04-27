package trace

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
)

var _ slog.Handler = &Handler{}

// Handler implements a simple slog.Handler that logs method calls to the handler
// that are part of the slog.Handler interface definition to STDOUT and then loses them.
// There is no log output to any io.Writer or other destination.
// Instantiate an object of this type and use it to make a slog.Logger.
// Log messages to that logger will dump lines of text representing Handler calls.
type Handler struct {
	indent string
}

// NewHandler returns a new trace.Handler.
// Provide a non-empty string for beginning of line indents
// as well as WithAttrs and WithGroup indents.
func NewHandler(indent string) *Handler {
	return &Handler{
		indent: indent,
	}
}

// -----------------------------------------------------------------------------
// Methods that implement the slog.Handler interface.

func (h Handler) Enabled(_ context.Context, level slog.Level) bool {
	h.show("Enabled", level)
	return true
}

func (h Handler) Handle(_ context.Context, record slog.Record) error {
	h.show("Handle", recordImage(record))
	return nil
}

func (h Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	h.show("WithAttrs", attrs)
	return Handler{indent: h.indent + "  "}
}

func (h Handler) WithGroup(name string) slog.Handler {
	h.show("WithGroup", name)
	return Handler{indent: h.indent + "  "}
}

// -----------------------------------------------------------------------------

func (h Handler) show(name string, arg any) {
	fmt.Printf("%s%s(%v)\n", h.indent, name, arg)
}

// -----------------------------------------------------------------------------

func recordImage(record slog.Record) string {
	var image strings.Builder
	image.WriteString(record.Level.String())
	image.WriteString(" \"")
	image.WriteString(record.Message)
	image.WriteString("\"")
	if record.NumAttrs() > 0 {
		separator := " {"
		record.Attrs(func(attr slog.Attr) bool {
			image.WriteString(separator)
			separator = ", "
			image.WriteString(attr.String())
			return true
		})
		image.WriteString("}")
	}
	return image.String()
}
