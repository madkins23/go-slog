package flash

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"sync"
)

const lenLog = 1024
const lenBasic = 4
const lenPrefix = 512
const lenSuffix = 32

var _ slog.Handler = &Handler{}

type Handler struct {
	options        *slog.HandlerOptions
	writer         io.Writer
	mutex          *sync.Mutex
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
		mutex:   &sync.Mutex{},
		prefix:  make([]byte, 0, lenPrefix),
		suffix:  make([]byte, 0, lenSuffix),
	}
	return hdlr
}

// -----------------------------------------------------------------------------

func (h *Handler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.options.Level.Level()
}

var logPool = newArrayPool[byte](lenLog)

func (h *Handler) Handle(_ context.Context, record slog.Record) error {
	// The x[:0] should reset len(x) to zero but leave cap(x) and
	// the underlying array space intact for reuse.
	buffer := logPool.get()[:0]
	defer func() { logPool.put(buffer) }()

	c := newComposer(buffer, false, h.options.ReplaceAttr, h.groups)
	defer reuseComposer(c)
	c.addBytes('{')

	// Adding attributes to the composer one at a time instead of
	// adding them to an array of attributes and
	// adding the list to the composer all at once.
	// See BenchmarkBasicManual and BenchmarkBasicMultiple in speed_test.go.
	if !record.Time.IsZero() {
		if err := c.addAttribute(slog.Time(slog.TimeKey, record.Time)); err != nil {
			return fmt.Errorf("add time: %w", err)
		}
	}
	if err := c.addAttribute(slog.String(slog.LevelKey, record.Level.String())); err != nil {
		return fmt.Errorf("add level: %w", err)
	}
	if err := c.addAttribute(slog.String(slog.MessageKey, record.Message)); err != nil {
		return fmt.Errorf("add message: %w", err)
	}
	if h.options.AddSource && record.PC != 0 {
		// Using local variable and loadSource instead of newSource and reuseSource.
		// See BenchmarkSourceLoad and BenchmarkSourceNewReuse in speed_test.go.
		var src source
		loadSource(record.PC, &src)
		if err := c.addAttribute(slog.Any(slog.SourceKey, &src)); err != nil {
			return fmt.Errorf("add source: %w", err)
		}
	}

	if len(h.prefix) > 0 {
		c.addBytes(',', ' ')
		c.addByteArray(h.prefix)
		if bytes.HasSuffix(h.prefix, []byte{'{'}) {
			// Inside a group, reset composer (started = false) to avoid comma.
			c.reset()
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

	h.mutex.Lock()
	defer h.mutex.Unlock()
	if _, err := h.writer.Write(c.getBytes()); err != nil {
		return fmt.Errorf("write log Line: %w", err)
	}

	return nil
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	hdlr := &Handler{
		options: h.options,
		writer:  h.writer,
		mutex:   h.mutex,
	}
	var prefixStarted bool
	if len(h.prefix) > 0 {
		hdlr.prefix = h.prefix
		if !bytes.HasSuffix(hdlr.prefix, []byte{'{'}) {
			prefixStarted = true
		}
	} else {
		hdlr.prefix = make([]byte, 0, lenPrefix)
	}
	if len(h.suffix) > 0 {
		hdlr.suffix = h.suffix
	} else {
		hdlr.suffix = make([]byte, 0, lenSuffix)
	}

	c := newComposer(hdlr.prefix, prefixStarted, h.options.ReplaceAttr, h.groups)
	defer reuseComposer(c)

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
			mutex:   h.mutex,
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
	} else {
		hdlr.prefix = make([]byte, 0, lenPrefix)
	}
	if len(h.suffix) > 0 {
		hdlr.suffix = h.suffix
	} else {
		hdlr.suffix = make([]byte, 0, lenSuffix)
	}
	hdlr.prefix = append(hdlr.prefix, []byte(fmt.Sprintf("%s\"%s\": {", prefixStart, name))...)
	// Should really be prepending the right brace into the suffix,
	// but suffix only contains right braces, so it doesn't really matter.
	hdlr.suffix = append(h.suffix, '}')
	return hdlr
}

// -----------------------------------------------------------------------------