package json

import (
	"encoding/json"
	"log/slog"
)

// Parse JSON byte array and return map[string]any.
// In the event of an error resulting map will contain the error object as "err".
func Parse(asJSON []byte) map[string]any {
	var parsed map[string]any
	if err := json.Unmarshal(asJSON, &parsed); err != nil {
		slog.Error("unable to parse expected JSON", "json", string(asJSON), "err", err)
		parsed = map[string]any{
			"err": err,
		}
	}
	return parsed
}
