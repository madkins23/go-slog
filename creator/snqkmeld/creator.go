package snqkmeld

import (
	"io"
	"log/slog"

	"github.com/madkins23/go-slog/infra"

	"snqk.dev/slog/meld"
)

// Creator returns a Creator object for the [snqk/meld] handler.
//
// [snqk/meld]: https://github.com/snqk/slog-meld
func Creator() infra.Creator {
	return infra.NewCreator("snqk/meld", handlerFn, nil,
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
