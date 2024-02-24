package data

import (
	"strings"

	"github.com/fatih/camelcase"
)

// -----------------------------------------------------------------------------

// TestTag is a unique name for a Benchmark or Verification test.
// The type is an alias for string so that types can't be confused.
type TestTag string

func (tt TestTag) Name() string {
	var builder strings.Builder
	for _, part1 := range strings.Split(string(tt), "_") {
		for _, part2 := range camelcase.Split(part1) {
			if builder.Len() > 0 {
				builder.WriteString(" ")
			}
			builder.WriteString(strings.ToUpper(part2[:1]) + strings.ToLower(part2[1:]))
		}
	}
	return builder.String()
}

// -----------------------------------------------------------------------------

// HandlerTag is a unique name for a slog handler.
// The type is an alias for string so that types can't be confused.
type HandlerTag string

func (ht HandlerTag) Name() string {
	var builder strings.Builder
	for _, part1 := range strings.Split(string(ht), "_") {
		for _, part2 := range camelcase.Split(part1) {
			if builder.Len() > 0 {
				builder.WriteString(" ")
			}
			builder.WriteString(strings.ToUpper(part2[:1]) + strings.ToUpper(part2[1:]))
		}
	}
	return builder.String()
}

func FixBenchHandlerTag(hdlrBytes []byte) HandlerTag {
	if string(hdlrBytes) == "Benchmark_slog_json" {
		// Fix this so the handler name doesn't get edited down to nothing.
		hdlrBytes = []byte("Benchmark_slog_slog_json")
	}
	tagString := strings.TrimPrefix(string(hdlrBytes), "Benchmark_slog_")
	return HandlerTag(tagString)
}
