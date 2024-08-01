package madkinsflash

import (
	"io"
	"log/slog"

	"github.com/madkins23/go-slog/handlers/flash"
	"github.com/madkins23/go-slog/infra"
)

const Name = "madkins/flash"

// Creator returns a Creator object for the [madkins/flash] handler.
// This is an experimental handler development based on madkins/sloggy.
func Creator() infra.Creator {
	return infra.NewCreator(Name, handlerFn, nil,
		`^madkins/flash^ is a clone of
		the [^madkins/sloggy^ handler](/go-slog/handler/MadkinsSloggy.html)
		with numerous performance improvements.
		In addition, ^madkins/flash^ can be configured with ^flash.Extras^ options.`,
		map[string]string{
			"madkins/flash":  "https://pkg.go.dev/github.com/madkins23/go-slog/handlers/flash",
			"madkins/sloggy": "https://pkg.go.dev/github.com/madkins23/go-slog/handlers/sloggy",
			"flash.Extras":   "https://pkg.go.dev/github.com/madkins23/go-slog/handlers/flash#Extras",
		})
}

func handlerFn(w io.Writer, options *slog.HandlerOptions) slog.Handler {
	return flash.NewHandler(w, options, nil)
}
