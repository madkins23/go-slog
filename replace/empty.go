package replace

import (
	"log/slog"
)

var _ AttrFn = MessageToMsg

// RemoveEmptyKey removes attributes with empty strings as key.
func RemoveEmptyKey(_ []string, a slog.Attr) slog.Attr {
	if a.Key == "" {
		return EmptyAttr
	}
	return a
}
