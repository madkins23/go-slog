package replace

import (
	"log/slog"
)

var _ AttrFn = MessageToMsg
var emptyAttr = slog.Attr{}

// RemoveEmptyKey removes attributes with empty strings as key.
func RemoveEmptyKey(_ []string, a slog.Attr) slog.Attr {
	if a.Key == "" {
		return emptyAttr
	}
	return a
}
