package madkinsreplattr

import (
	"io"
	"log/slog"

	"github.com/madkins23/go-slog/handlers/flash"
	"github.com/madkins23/go-slog/infra"
	"github.com/madkins23/go-slog/replace"
)

const Name = "madkins/replattr"

// Creator returns a Creator object for the [madkins/replattr] handler.
// This is the madkins/flash handler configured with both
// flash.Extras customization (to mimic a badly behaved handler) and
// a set of ReplaceAttr functions to remove that customization.
// The goal is to benchmark the overhead in using ReplaceAttr functions.
func Creator() infra.Creator {
	return infra.NewCreator(Name, handlerFn, nil,
		`^madkins/replattr^ is the [^madkins/flash^ handler](/go-slog/handler/MadkinsFlash.html)
		setup to test ^slog.HandlerOptions.ReplaceAttr^ performance.
		The ^madkins/flash^ handler is first configured to generate various warnings using ^flash.Extras^ options,
		then several ^ReplaceAttr^ functions are used to correct the aberrant behavior.
		This is intended to measure the effect of ^ReplaceAttr^ usage on performance by comparison with ^madkins/flash^`,
		map[string]string{
			"madkins/flash":           "https://pkg.go.dev/github.com/madkins23/go-slog/handlers/flash",
			"flash.Extras":            "https://pkg.go.dev/github.com/madkins23/go-slog/handlers/flash#Extras",
			"madkinsreplattr.Creator": "https://github.com/madkins23/go-slog/blob/main/creator/madkinsreplattr/creator.go",
		})
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
