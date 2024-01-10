package replace

import (
	"log/slog"
	"strings"
)

var _ AttrFn = MessageToMsg

// MessageToMsg replaces attribute keys matching "message" with the correct value of slog.MessageKey
func MessageToMsg(groups []string, a slog.Attr) slog.Attr {
	if strings.ToLower(a.Key) == "message" && len(groups) == 0 {
		a.Key = slog.MessageKey
	}
	return a
}
