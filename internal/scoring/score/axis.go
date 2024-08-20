package score

import (
	"html/template"
	"log/slog"
	"math"

	"github.com/madkins23/go-slog/internal/data"
)

// -----------------------------------------------------------------------------

//go:generate go run github.com/dmarkham/enumer -type=Type
type Type uint8

const (
	ByData Type = iota
	Original
	ByTest
)

var colNames = map[Type]string{
	ByData:   "by Data",
	ByTest:   "by Test",
	Original: "Original",
}

func (t Type) ColHeader() string {
	return colNames[t]
}

type Axis interface {
	Setup(bench *data.Benchmarks, warns *data.Warnings) error
	Name() string
	HasTest(test data.TestTag) bool
	ScoreFor(handler data.HandlerTag) Value
	ScoreForTest(handler data.HandlerTag, test data.TestTag) Value
	ScoreForType(handler data.HandlerTag, which Type) Value
	Summary() template.HTML
	Exhibits() []Exhibit
	Documentation() template.HTML
	Type() string
}

type Value float64

func (v Value) Round() Value {
	const rounder = 1_000_000_000.0
	return Value(math.Round(float64(v)*rounder) / rounder)
}

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
