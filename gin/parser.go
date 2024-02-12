package gin

import (
	"fmt"
	"log/slog"
	"regexp"
	"strconv"
	"strings"
)

// Field codifies the fields that may be parsed from a Gin traffic record.
type Field string

const (
	Code    Field = "code"
	Client  Field = "client"
	Elapsed Field = "elapsed"
	Method  Field = "method"
	Url     Field = "url"
)

var (
	ptnCode  = regexp.MustCompile(`^\s*(\d+)\s*$`)
	ptnSplit = regexp.MustCompile(`\s+`)
)

// Parse a Gin traffic record (specified as message) to return an array of slog.Attr items.
// There should be a single such item for each Field constant defined above.
func Parse(message string) ([]any, error) {
	// Example line:
	//  200 |    9.522199ms |             ::1 | GET      "/chart.svg?tag=samber_zap&item=MemAllocs"
	parts := strings.Split(message, "|")
	if len(parts) != 4 {
		return nil, fmt.Errorf("wrong number of parts: %d", len(parts))
	}
	var result []any
	// Parse HTTP code from first part.
	if matches := ptnCode.FindStringSubmatch(parts[0]); len(matches) != 2 {
		return nil, fmt.Errorf("parse Code from '%s'", parts[0])
	} else if num, err := strconv.ParseInt(matches[1], 10, 64); err != nil {
		// This should never happen so no test code coverage here.
		return nil, fmt.Errorf("parse int from '%s': %w", matches[1], err)
	} else {
		result = append(result, slog.Int64(string(Code), num))
	}
	// Parse elapsed time from second part.
	result = append(result, slog.String(string(Elapsed), strings.Trim(parts[1], " ")))
	// Parse client IP address from third part.
	result = append(result, slog.String(string(Client), strings.Trim(parts[2], " ")))
	// Parse HTTP method and URL from fourth part.
	// Example of parts[3]:
	//  GET      "/chart.svg?tag=samber_zap&item=MemAllocs" System=gin
	parts = ptnSplit.Split(strings.Trim(parts[3], " "), -1)
	if len(parts) == 2 {
		result = append(result, slog.String(string(Method), parts[0]))
		result = append(result, slog.String(string(Url), strings.Trim(parts[1], "\"")))
	}
	return result, nil
}
