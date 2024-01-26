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
	ptnCode  = regexp.MustCompile(`\s(\d+)\s*$`)
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
		slog.Warn("Unable to parse code", "from", parts[0])
	} else if num, err := strconv.ParseInt(matches[1], 10, 64); err != nil {
		slog.Warn("Unable to parse int", "from", parts[0], "func", "getLogData")
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
