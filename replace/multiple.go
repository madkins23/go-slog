package replace

import (
	"log/slog"

	"github.com/madkins23/go-slog/infra"
)

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
