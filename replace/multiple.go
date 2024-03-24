package replace

import (
	"log/slog"

	"github.com/madkins23/go-slog/infra"
)

// Multiple generates a ReplaceAttr function out of multiple such functions.
// They will be executed in sequence.
// If one of the subsidiary functions returns the empty attribute that is returned
// and the loop ceases.
func Multiple(fns ...infra.AttrFn) infra.AttrFn {
	return func(groups []string, a slog.Attr) slog.Attr {
		for _, fn := range fns {
			a = fn(groups, a)
			if a.Equal(infra.EmptyAttr()) {
				break
			}
		}
		return a
	}
}
