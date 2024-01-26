package gin

import (
	"fmt"
	"log/slog"
	"regexp"
	"strconv"
	"strings"
)

type Field string

const (
	code    Field = "code"
	client  Field = "client"
	elapsed Field = "elapsed"
	method  Field = "method"
	system  Field = "system"
	url     Field = "url"
)

var (
	ptnCode  = regexp.MustCompile(`^\s*(\d+)\s*$`)
	ptnSplit = regexp.MustCompile(`\s+`)
)

func Parse(message string) ([]any, error) {
	// Example line:
	//  200 |    9.522199ms |             ::1 | GET      "/chart.svg?tag=samber_zap&item=MemAllocs" system=gin
	parts := strings.Split(message, "|")
	if len(parts) != 4 {
		return nil, fmt.Errorf("wrong number of parts: %d", len(parts))
	}
	var result []any
	if matches := ptnCode.FindStringSubmatch(parts[0]); len(matches) != 2 {
		// TODO: if it doesn't parse it's not an error, just return the message as is
		return nil, fmt.Errorf("parse code from '%s'", parts[0])
	} else if num, err := strconv.ParseInt(matches[1], 10, 64); err != nil {
		// TODO: if it doesn't parse it's not an error, just return the message as is
		return nil, fmt.Errorf("parse int from '%s': %w", matches[1], err)
	} else {
		result = append(result, slog.Int64(string(code), num))
	}
	result = append(result, slog.String(string(elapsed), strings.Trim(parts[1], " ")))
	result = append(result, slog.String(string(client), strings.Trim(parts[2], " ")))
	// Example parts[3]:
	//  GET      "/chart.svg?tag=samber_zap&item=MemAllocs" system=gin
	parts = ptnSplit.Split(strings.Trim(parts[3], " "), -1)
	if len(parts) == 3 {
		result = append(result, slog.String(string(method), parts[0]))
		result = append(result, slog.String(string(url), strings.Trim(parts[1], "\"")))
	}
	return result, nil
}
