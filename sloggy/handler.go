package sloggy

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
)

var _ slog.Handler = &Handler{}

type Handler struct {
	options        *slog.HandlerOptions
	writer         io.Writer
	prefix, suffix bytes.Buffer
}

func NewHandler(writer io.Writer, options *slog.HandlerOptions) *Handler {
	if options == nil {
		options = &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}
	}
	hdlr := &Handler{
		options: options,
		writer:  writer,
	}
	return hdlr
}

// -----------------------------------------------------------------------------

func (h *Handler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.options.Level.Level()
}

func (h *Handler) Handle(_ context.Context, record slog.Record) error {
	c := newComposer(h.writer, false)
	if err := c.begin(); err != nil {
		return fmt.Errorf("begin: %w", err)
	}

	basic := make([]slog.Attr, 0, 3)
	if !record.Time.IsZero() {
		basic = append(basic, slog.Time(slog.TimeKey, record.Time))
	}
	basic = append(basic, slog.String(slog.LevelKey, h.options.Level.Level().String()))
	basic = append(basic, slog.String(slog.MessageKey, record.Message))
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

	return nil
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	hdlr := &Handler{
		options: h.options,
		writer:  h.writer,
		prefix:  bytes.Buffer{},
		suffix:  bytes.Buffer{},
	}
	var prefixStarted bool
	if h.prefix.Len() > 0 {
		hdlr.prefix.Write(h.prefix.Bytes())
		prefixStarted = true
	}
	if h.suffix.Len() > 0 {
		hdlr.prefix.Write(h.suffix.Bytes())
	}
	if err := newComposer(&hdlr.prefix, prefixStarted).addAttributes(attrs); err != nil {
		slog.Error("adding with attributes", "err", err)
	}

	return hdlr
}

func (h *Handler) WithGroup(name string) slog.Handler {
	hdlr := &Handler{
		options: h.options,
		writer:  h.writer,
		prefix:  bytes.Buffer{},
		suffix:  bytes.Buffer{},
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
