package madkinsflash

import (
	"io"
	"log/slog"

	"github.com/madkins23/go-slog/handlers/flash"
	"github.com/madkins23/go-slog/infra"
	"github.com/madkins23/go-slog/replace"
)

// Creator returns a Creator object for the [madkins/replattr] handler.
// This is the madkins/flash handler configured with both
// flash.Extras customization (to mimic a badly behaved handler) and
// a set of ReplaceAttr functions to remove that customization.
// The goal is to benchmark the overhead in using ReplaceAttr functions.
func Creator() infra.Creator {
	return infra.NewCreator("madkins/replattr", handlerFn, nil)
}

func handlerFn(w io.Writer, options *slog.HandlerOptions) slog.Handler {
	fns := []infra.AttrFn{
		replace.ChangeKey("message", slog.MessageKey, false, replace.TopCheck),
		replace.ChangeKey("lvl", slog.LevelKey, false, replace.TopCheck),
		replace.ChangeCase(slog.LevelKey, replace.CaseUpper, false, replace.TopCheck),
	}
	if options.ReplaceAttr != nil {
		fns = append(fns, options.ReplaceAttr)
	}
	return flash.NewHandler(w,
		&slog.HandlerOptions{
			AddSource:   options.AddSource,
			Level:       options.Level,
			ReplaceAttr: replace.Multiple(fns...),
		},
		&flash.Extras{
			LevelNames: map[slog.Level]string{
				slog.LevelDebug: "Debug",
				slog.LevelInfo:  "Info",
				slog.LevelWarn:  "Warn",
				slog.LevelError: "Error",
			},
			LevelKey:   "lvl",
			MessageKey: "message",
		})
}
