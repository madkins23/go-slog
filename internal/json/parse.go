package json

import (
	"encoding/json"
	"fmt"
)

// Expect parses a JSON string and returns a map[string]any.
// In the event of a JSON unmarshal error
// the resulting map will contain only the error object as "error".
func Expect(asJSON string) map[string]any {
	if result, err := Parse([]byte(asJSON)); err == nil {
		return result
	} else {
		return map[string]any{"error": err.Error()}
	}
}

// Parse JSON byte array and return map[string]any.
func Parse(asJSON []byte) (map[string]any, error) {
	var parsed map[string]any
	if err := json.Unmarshal(asJSON, &parsed); err != nil {
		return nil, fmt.Errorf("unmarshal json: %w", err)
	}
	return parsed, nil
}
