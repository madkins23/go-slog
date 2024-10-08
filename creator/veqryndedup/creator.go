package veqryndedup

import (
	"io"
	"log/slog"

	slogdedup "github.com/veqryn/slog-dedup"

	"github.com/madkins23/go-slog/infra"
)

const BaseName = "veqryn/dedup"

//go:generate go run github.com/dmarkham/enumer -type=Mode
type Mode uint8

const (
	// None mode
	None Mode = iota
	// Append mode collects all values in a group.
	Append
	// Ignore mode uses the first value found.
	Ignore
	// Increment mode increments the field name for each value.
	Increment
	// Overwrite mode uses the last value found.
	Overwrite
)

func Name(mode Mode) string {
	return BaseName + "/" + mode.String()
}

// Creator returns a Creator object for the [veqryn/dedup] handler.
//
// [veqryn/dedup]: https://github.com/veqryn/slog-dedup
func Creator(mode Mode) infra.Creator {
	return infra.NewCreator(Name(mode), handler(mode), nil,
		`^veqryn/dedup^ provides a variety of ^slog.Handler^ options
		for deduplicating the keys: overwriting, ignoring, appending, and incrementing.`,
		map[string]string{
			"veqryn/dedup": "https://github.com/veqryn/slog-dedup",
		})
}

func handler(mode Mode) infra.CreateHandlerFn {
	switch mode {
	case Ignore:
		return ignoreHandler
	case Overwrite:
		return overwriteHandler
	case Append:
		return appendHandler
	case Increment:
		return incrementHandler
	default:
		slog.Error("Unknown creator mode", "mode", mode)
		return noHandler
	}
}

func appendHandler(w io.Writer, options *slog.HandlerOptions) slog.Handler {
	return slogdedup.NewAppendHandler(slog.NewJSONHandler(w, options), nil)
}

func noHandler(w io.Writer, options *slog.HandlerOptions) slog.Handler {
	return slog.NewJSONHandler(w, options)
}

func ignoreHandler(w io.Writer, options *slog.HandlerOptions) slog.Handler {
	return slogdedup.NewIgnoreHandler(slog.NewJSONHandler(w, options), nil)
}

func incrementHandler(w io.Writer, options *slog.HandlerOptions) slog.Handler {
	return slogdedup.NewIncrementHandler(slog.NewJSONHandler(w, options), nil)
}

func overwriteHandler(w io.Writer, options *slog.HandlerOptions) slog.Handler {
	return slogdedup.NewOverwriteHandler(slog.NewJSONHandler(w, options), nil)
}
