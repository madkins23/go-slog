package replace

import "log/slog"

type AttrFn func(groups []string, a slog.Attr) slog.Attr

var EmptyAttr = slog.Attr{}
