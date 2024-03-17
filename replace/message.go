package replace

import (
	"log/slog"
)

// MessageToMsg replaces attribute keys matching "message" with the correct value of slog.MessageKey
//
// Deprecated: This function was never tested and can be replaced by
//
//	replace.ChangeKey("message", slog.MessageKey, false, replace.TopCheck)
//
// which is how it is now implemented.
func MessageToMsg(groups []string, a slog.Attr) slog.Attr {
	return ChangeKey("message", slog.MessageKey, false, TopCheck)(groups, a)
}
