package data

import (
	"strings"

	"github.com/fatih/camelcase"
)

// -----------------------------------------------------------------------------

// TestTag is a unique name for a Benchmark or Verification test.
// The type is an alias for string so that types can't be confused.
type TestTag string

var testTagNames = make(map[TestTag]string)

// Name returns a name string calculated from the TestTag string and cached for reuse.
func (tt TestTag) Name() string {
	name, found := testTagNames[tt]
	if !found {
		var builder strings.Builder
		tagString := string(tt)
		if parts := strings.Split(string(tt), ":"); len(parts) == 2 {
			tagString = parts[1]
		}
		for _, part := range camelcase.Split(tagString) {
			if builder.Len() > 0 {
				builder.WriteString(" ")
			}
			builder.WriteString(part)
		}
		name = builder.String()
		testTagNames[tt] = name
	}
	return name

}

// -----------------------------------------------------------------------------

// HandlerTag is a unique name for a slog handler.
// The type is an alias for string so that types can't be confused.
type HandlerTag string

var handlerTagNames = make(map[HandlerTag]string)

// Name returns a name string calculated from the HandlerTag string and cached for reuse.
func (ht HandlerTag) Name() string {
	name, found := handlerTagNames[ht]
	if !found {
		var builder strings.Builder
		tagString := string(ht)
		for _, part := range camelcase.Split(tagString) {
			if builder.Len() > 0 {
				builder.WriteString(" ")
			}
			builder.WriteString(part)
		}
		name = builder.String()
		handlerTagNames[ht] = name
	}
	return name
}
