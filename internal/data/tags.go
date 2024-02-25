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

func FixBenchHandlerTag(hdlrBytes []byte) HandlerTag {
	if string(hdlrBytes) == "Benchmark_slog_json" {
		// Fix this so the handler name doesn't get edited down to nothing.
		hdlrBytes = []byte("Benchmark_slog_slog_json")
	}
	tagString := strings.TrimPrefix(string(hdlrBytes), "Benchmark_slog_")
	return HandlerTag(tagString)
}
