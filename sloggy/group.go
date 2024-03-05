package sloggy

import (
	"context"
	"log/slog"

	"github.com/madkins23/go-slog/infra"
)

type group struct {
	*Handler
	name   string
	parent slog.Handler
}

func (g *group) Handle(ctx context.Context, record slog.Record) error {
	count := 0
	var dead = false
	deadGroup := ctx.Value("deadGroup")
	record.Attrs(func(attr slog.Attr) bool {
		if !attr.Equal(infra.EmptyAttr()) {
			count++
			if deadGroup != "" && attr.Key == deadGroup {
				if attr.Value.Kind() == slog.KindGroup {
					dead = true
				}
			}
			return false
		}
		return true
	})
	if count < 1 || // There is nothing in the record attributes, the group is empty or...
		(count < 2 && dead) { // ...there is only a single attribute which is the dead group that called.
		// Since the normal handler prefix already includes group name and open brace
		// it isn't possible to use the current handler prefix.
		// The parent handler prefix, however, should work just fine.
		return g.parent.Handle(context.WithValue(ctx, "deadGroup", g.name), record)
	}
	return g.Handler.Handle(ctx, record)
}
