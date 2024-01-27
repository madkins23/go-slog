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
	Code    Field = "Code"
	Client  Field = "Client"
	Elapsed Field = "Elapsed"
	Method  Field = "Method"
	Url     Field = "Url"
)

var (
	ptnCode  = regexp.MustCompile(`^\s*(\d+)\s*$`)
	ptnSplit = regexp.MustCompile(`\s+`)
)

// TODO: document

// Parse a Gin traffic record (specified as message) to return a set of slog.Attr items.
// There should be a single such item for each Field constant defined above.
func Parse(message string) ([]any, error) {
	// Example line:
	//  200 |    9.522199ms |             ::1 | GET      "/chart.svg?tag=samber_zap&item=MemAllocs"
	parts := strings.Split(message, "|")
	if len(parts) != 4 {
		return nil, fmt.Errorf("wrong number of parts: %d", len(parts))
	}
	var result []any
	if matches := ptnCode.FindStringSubmatch(parts[0]); len(matches) != 2 {
		// TODO: if it doesn't parse it's not an error, just return the message as is
		return nil, fmt.Errorf("parse Code from '%s'", parts[0])
	} else if num, err := strconv.ParseInt(matches[1], 10, 64); err != nil {
		// TODO: if it doesn't parse it's not an error, just return the message as is
		return nil, fmt.Errorf("parse int from '%s': %w", matches[1], err)
	} else {
		result = append(result, slog.Int64(string(Code), num))
	}
	result = append(result, slog.String(string(Elapsed), strings.Trim(parts[1], " ")))
	result = append(result, slog.String(string(Client), strings.Trim(parts[2], " ")))
	// Example parts[3]:
	//  GET      "/chart.svg?tag=samber_zap&item=MemAllocs" System=gin
	parts = ptnSplit.Split(strings.Trim(parts[3], " "), -1)
	if len(parts) == 2 {
		result = append(result, slog.String(string(Method), parts[0]))
		result = append(result, slog.String(string(Url), strings.Trim(parts[1], "\"")))
	}
	return result, nil
}
