package replace

import (
	"log/slog"
)

// RemoveEmptyKey removes attributes with empty strings as key.
// This is done by setting the attribute to "empty" (slog.Attr{} or JSON "": null).
// If the handler improperly shows empty keys in JSON records
// then the attribute will still be logged but will have a null JSON value.
//
// Deprecated: This function can be replaced by
//
//	replace.RemoveKey("", CaseLower, false, TopCheck)
//
// which is how it is now implemented.
func RemoveEmptyKey(groups []string, a slog.Attr) slog.Attr {
	return RemoveKey("", false, nil)(groups, a)
}
