package flash

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"runtime"
)

const lenLog = 1024
const lenPrefix = 512
const lenSuffix = 32

var logPool = newArrayPool[byte](lenLog)

var _ slog.Handler = &Handler{}

type Handler struct {
	options        *slog.HandlerOptions
	writer         io.Writer
	prefix, suffix []byte
	groups         []string
}

func NewHandler(writer io.Writer, options *slog.HandlerOptions) *Handler {
	if options == nil {
		options = &slog.HandlerOptions{}
	}
	if options.Level == nil {
		options.Level = slog.LevelInfo
	}
	hdlr := &Handler{
		options: options,
		writer:  writer,
		prefix:  make([]byte, 0, lenPrefix),
		suffix:  make([]byte, 0, lenSuffix),
	}
	return hdlr
}

// -----------------------------------------------------------------------------

func (h *Handler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.options.Level.Level()
}

func (h *Handler) Handle(_ context.Context, record slog.Record) error {
	buffer := logPool.get()
	defer func() {
		logPool.put(buffer)
	}()

	c := newComposer(buffer, false, h.options.ReplaceAttr, h.groups)
	c.addBytes('{')

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

	if len(h.prefix) > 0 {
		c.addBytes(',', ' ')
		c.addByteArray(h.prefix)
		if bytes.HasSuffix(h.prefix, []byte{'{'}) {
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
	if len(h.suffix) > 0 {
		c.addByteArray(h.suffix)
	}
	c.addBytes('}', '\n')
	if _, err := h.writer.Write(c.getBytes()); err != nil {
		return fmt.Errorf("write log line: %w", err)
	}

	return nil
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	hdlr := &Handler{
		options: h.options,
		writer:  h.writer,
		prefix:  make([]byte, 0, lenPrefix),
		suffix:  make([]byte, 0, lenSuffix),
	}
	var prefixStarted bool
	if len(h.prefix) > 0 {
		hdlr.prefix = h.prefix
		if !bytes.HasSuffix(hdlr.prefix, []byte{'{'}) {
			prefixStarted = true
		}
	}
	if len(h.suffix) > 0 {
		hdlr.suffix = h.suffix
	}
	c := newComposer(h.prefix, prefixStarted, h.options.ReplaceAttr, h.groups)
	if err := c.addAttributes(attrs); err != nil {
		slog.Error("adding with attributes", "err", err)
	}
	hdlr.prefix = c.getBytes()

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
			prefix:  make([]byte, 0, lenPrefix),
			suffix:  make([]byte, 0, lenSuffix),
			groups:  append(h.groups, name),
		},
		name:   name,
		parent: h,
	}
	prefixStart := ""
	if len(h.prefix) > 0 {
		hdlr.prefix = h.prefix
		if !bytes.HasSuffix(h.prefix, []byte{'{'}) {
			prefixStart = ", "
		}
	}
	if len(h.suffix) > 0 {
		hdlr.suffix = h.suffix
	}
	hdlr.prefix = append(hdlr.prefix, []byte(fmt.Sprintf("%s\"%s\": {", prefixStart, name))...)
	// Should really be prepending the right brace into the suffix,
	// but suffix only contains right braces, so it doesn't really matter.
	hdlr.suffix = append(h.suffix, '}')
	return hdlr
}

// -----------------------------------------------------------------------------
