package sloggy

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"runtime"
	"sync"
)

var _ slog.Handler = &Handler{}

// Handler provides a fairly straightforward, feature-complete slog.Handler implementation.
type Handler struct {
	options        *slog.HandlerOptions
	writer         io.Writer
	mutex          *sync.Mutex
	prefix, suffix bytes.Buffer
	groups         []string
}

// NewHandler returns a new sloggy handler with the specified output writer and slog.HandlerOptions.
// If the options argument is nil it will be set to a level of slog.LevelInfo and nothing else.
func NewHandler(writer io.Writer, options *slog.HandlerOptions) *Handler {
	hdlr := &Handler{
		options: fixOptions(options),
		writer:  writer,
		mutex:   &sync.Mutex{},
	}
	return hdlr
}

// -----------------------------------------------------------------------------

func (h *Handler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.options.Level.Level()
}

func (h *Handler) Handle(_ context.Context, record slog.Record) error {
	c := newComposer(h.writer, false, h.options.ReplaceAttr, h.groups)
	if err := c.begin(); err != nil {
		return fmt.Errorf("begin: %w", err)
	}

	basic := make([]slog.Attr, 0, 4)
	if !record.Time.IsZero() {
		basic = append(basic, slog.Time(slog.TimeKey, record.Time))
	}
	basic = append(basic, slog.String(slog.LevelKey, record.Level.String()))
	basic = append(basic, slog.String(slog.MessageKey, record.Message))
	if h.options.AddSource && record.PC != 0 {
		fs := runtime.CallersFrames([]uintptr{record.PC})
		f, _ := fs.Next()
		source := map[string]any{
			"function": f.Function,
			"file":     f.File,
			"line":     f.Line,
		}
		basic = append(basic, slog.Any(slog.SourceKey, source))
	}
	if err := c.addAttributes(basic); err != nil {
		return fmt.Errorf("add basic attributes: %w", err)
	}

	if h.prefix.Len() > 0 {
		if _, err := c.Write(commaSpace); err != nil {
			return fmt.Errorf("comma space: %w", err)
		}
		if _, err := c.Write(h.prefix.Bytes()); err != nil {
			return fmt.Errorf("write prefix: %w", err)
		}
		if bytes.HasSuffix(h.prefix.Bytes(), []byte{'{'}) {
			c.setStarted(false)
		}
	}

	var err error
	record.Attrs(func(attr slog.Attr) bool {
		if err = c.addAttribute(attr); err != nil {
			return false
		}
		return true // keep going
	})
	if err != nil {
		return fmt.Errorf("add attributes: %w", err)
	}

	if _, err := c.Write(h.suffix.Bytes()); err != nil {
		return fmt.Errorf("write suffix: %w", err)
	}

	if err := c.end(); err != nil {
		return fmt.Errorf("end: %w", err)
	}

	h.mutex.Lock()
	defer h.mutex.Unlock()

	if _, err := c.Write(newLine); err != nil {
		return fmt.Errorf("newLine: %w", err)
	}

	return nil
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	hdlr := &Handler{
		options: h.options,
		writer:  h.writer,
		mutex:   h.mutex,
		groups:  h.groups,
		prefix:  bytes.Buffer{},
		suffix:  bytes.Buffer{},
	}
	var prefixStarted bool
	if h.prefix.Len() > 0 {
		hdlr.prefix.Write(h.prefix.Bytes())
		if !bytes.HasSuffix(hdlr.prefix.Bytes(), braceLeft) {
			prefixStarted = true
		}
	}
	if h.suffix.Len() > 0 {
		hdlr.suffix.Write(h.suffix.Bytes())
	}
	c := newComposer(&hdlr.prefix, prefixStarted, h.options.ReplaceAttr, h.groups)
	if err := c.addAttributes(attrs); err != nil {
		slog.Error("adding with attributes", "err", err)
	}

	return hdlr
}

func (h *Handler) WithGroup(name string) slog.Handler {
	if name == "" {
		// Groups with empty names are to be inlined.
		return h
	}
	hdlr := &group{
		Handler: &Handler{
			options: h.options,
			writer:  h.writer,
			mutex:   h.mutex,
			prefix:  bytes.Buffer{},
			suffix:  bytes.Buffer{},
			groups:  append(h.groups, name),
		},
		name:   name,
		parent: h,
	}
	prefixStart := emptyString
	if h.prefix.Len() > 0 {
		hdlr.prefix.Write(h.prefix.Bytes())
		if !bytes.HasSuffix(h.prefix.Bytes(), braceLeft) {
			prefixStart = commaSpace
		}
	}
	if h.suffix.Len() > 0 {
		hdlr.suffix.Write(h.suffix.Bytes())
	}

	if _, err := fmt.Fprintf(&hdlr.prefix, "%s\"%s\": {", prefixStart, name); err != nil {
		slog.Error("open group", "err", err)
	}
	if _, err := hdlr.suffix.Write(braceRight); err != nil {
		slog.Error("group right brace", "err", err)
	}

	return hdlr
}

// -----------------------------------------------------------------------------

// fixOptions makes certain that a slog.HandlerOptions object has been properly created and
// configured with default values.
func fixOptions(options *slog.HandlerOptions) *slog.HandlerOptions {
	if options == nil {
		options = &slog.HandlerOptions{}
	}
	if options.Level == nil {
		options.Level = slog.LevelInfo
	}
	return options
}
