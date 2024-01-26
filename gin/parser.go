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

type FieldString map[Field]string

// Empty returns an empty FieldString.
func Empty() FieldString {
	return make(FieldString)
}

var allFieldStrings = FieldString{
	code:    string(code),
	client:  string(client),
	elapsed: string(elapsed),
	method:  string(method),
	system:  string(system),
	url:     string(url),
}

var (
	ptnCode  = regexp.MustCompile(`\s(\d+)\s*$`)
	ptnSplit = regexp.MustCompile(`\s+`)
)

type Parser interface {
	Parse(message string) ([]any, error)
}

func NewParser(fields FieldString) Parser {
	lookup := make(FieldString)
	for field, str := range allFieldStrings {
		if val, found := fields[field]; found {
			lookup[field] = val
		} else {
			lookup[field] = str
		}
	}
	return &parser{
		lookup: lookup,
	}
}

var _ Parser = &parser{}

type parser struct {
	lookup FieldString
}

func (p *parser) Parse(message string) ([]any, error) {
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
		result = append(result, slog.Int64(p.lookup[code], num))
	}
	result = append(result, slog.String(p.lookup[elapsed], strings.Trim(parts[1], " ")))
	result = append(result, slog.String(p.lookup[client], strings.Trim(parts[2], " ")))
	// Example parts[3]:
	//  GET      "/chart.svg?tag=samber_zap&item=MemAllocs" system=gin
	parts = ptnSplit.Split(strings.Trim(parts[3], " "), -1)
	if len(parts) == 3 {
		result = append(result, slog.String(p.lookup[method], parts[0]))
		result = append(result, slog.String(p.lookup[url], strings.Trim(parts[1], "\"")))
	}
	return result, nil
}
