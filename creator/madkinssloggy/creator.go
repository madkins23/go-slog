package madkinssloggy

import (
	"io"
	"log/slog"

	"github.com/madkins23/go-slog/handlers/sloggy"
	"github.com/madkins23/go-slog/infra"
)

// Creator returns a Creator object for the [madkins/sloggy] handler.
// This is an experimental handler development.
func Creator() infra.Creator {
	return infra.NewCreator("madkins/sloggy", handlerFn, nil,
		`^madkins/sloggy^ is a new ^slog.Handler^ built from the ground up to generate JSON log records.
		While the performance isn't top-tier, this handler adheres faithfully to documented ^slog.Handler^
		and observed ^slog.JSONHandler^ behavior.
		The [^madkins/flash^ handler](/go-slog/handler/MadkinsFlash.html), derived from this one, has better performance.`,
		map[string]string{
			"madkins/sloggy": "https://pkg.go.dev/github.com/madkins23/go-slog/handlers/sloggy",
			"madkins/flash":  "https://pkg.go.dev/github.com/madkins23/go-slog/handlers/flash",
		})
}

func handlerFn(w io.Writer, options *slog.HandlerOptions) slog.Handler {
	return sloggy.NewHandler(w, options)
}
