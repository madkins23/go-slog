package score

import "log/slog"

// -----------------------------------------------------------------------------

//go:generate go run github.com/dmarkham/enumer -type=Type
type Type uint8

const (
	Default Type = iota
	ByData
	Original
	ByTest
)

var colNames = map[Type]string{
	Default:  "Score",
	ByData:   "by Data",
	ByTest:   "by Test",
	Original: "Original",
}

func (t Type) ColHeader() string {
	if hdr, found := colNames[t]; found {
		return hdr
	}
	return "Unknown:" + t.String()
}

// -----------------------------------------------------------------------------

func List(typeName ...string) []Type {
	result := make([]Type, 0, len(typeName))
	for _, name := range typeName {
		if st, err := TypeString(name); err != nil {
			slog.Error("convert name to score type", "name", name, "err", err)
		} else {
			result = append(result, st)
		}
	}
	return result
}
