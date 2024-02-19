package data

import "strings"

// -----------------------------------------------------------------------------

// TestTag is a unique name for a Benchmark or Verification test.
// The type is an alias for string so that types can't be confused.
type TestTag string

func (tt TestTag) Name() string {
	parts := strings.Split(string(tt), "_")
	for i, part := range parts {
		parts[i] = strings.ToUpper(part[:1]) + strings.ToLower(part[1:])
	}
	return strings.Join(parts, " ")
}

// -----------------------------------------------------------------------------

// HandlerTag is a unique name for a slog handler.
// The type is an alias for string so that types can't be confused.
type HandlerTag string

func (ht HandlerTag) Name() string {
	parts := strings.Split(string(ht), "_")
	for i, part := range parts {
		parts[i] = strings.ToUpper(part[:1]) + strings.ToLower(part[1:])
	}
	return strings.Join(parts, " ")
}

func FixBenchHandlerTag(hdlrBytes []byte) HandlerTag {
	if string(hdlrBytes) == "Benchmark_slog" {
		// Fix this so the handler name doesn't get edited down to nothing.
		hdlrBytes = []byte("Benchmark_slog_slog_JSONHandler")
	}
	tagString := strings.TrimPrefix(string(hdlrBytes), "Benchmark_slog_")
	return HandlerTag(tagString)
}
