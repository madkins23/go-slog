package snqkmeld

import (
	"io"
	"log/slog"

	"snqk.dev/slog/meld"

	"github.com/madkins23/go-slog/infra"
)

const Name = "snqk/meld"

// Creator returns a Creator object for the [snqk/meld] handler.
//
// [snqk/meld]: https://github.com/snqk/slog-meld
func Creator() infra.Creator {
	return infra.NewCreator(Name, handlerFn, nil,
		`^snqk/slog-meld^ provides a simple slog.Handler
		designed to recursively merge and de-duplicate log attributes,
		ensuring clean, concise, and informative log entries.`,
		map[string]string{
			"snqk/meld": "https://github.com/snqk/slog-meld",
		})
}

func handlerFn(w io.Writer, options *slog.HandlerOptions) slog.Handler {
	return meld.NewHandler(slog.NewJSONHandler(w, options))
}
