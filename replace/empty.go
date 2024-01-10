package replace

import (
	"log/slog"
)

var _ AttrFn = MessageToMsg

// RemoveEmptyKey removes attributes with empty strings as key.
// This is done by setting the attribute to "empty" (slog.Attr{} or JSON "": null).
// If the handler improperly shows empty keys in JSON records
// then the attribute will still be logged but will have a null JSON value.
func RemoveEmptyKey(_ []string, a slog.Attr) slog.Attr {
	if a.Key == "" {
		return emptyAttr
	}
	return a
}
