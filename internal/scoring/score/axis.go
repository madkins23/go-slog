package score

import (
	"html/template"
	"log/slog"
	"math"

	"github.com/madkins23/go-slog/internal/data"
)

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
}

// -----------------------------------------------------------------------------

type Value float64

func (v Value) Round() Value {
	const rounder = 1_000_000_000.0
	return Value(math.Round(float64(v)*rounder) / rounder)
}

// -----------------------------------------------------------------------------

type AxisCore struct {
	exhibits    []Exhibit
	summaryHTML template.HTML
}

func (ac *AxisCore) SetSummary(summary template.HTML) {
	ac.summaryHTML = summary
}

func (ac *AxisCore) Summary() template.HTML {
	return ac.summaryHTML
}

func (ac *AxisCore) AddExhibit(exhibit Exhibit) {
	if ac.exhibits == nil {
		ac.exhibits = make([]Exhibit, 0, 2)
	}
	ac.exhibits = append(ac.exhibits, exhibit)
}

func (ac *AxisCore) Exhibits() []Exhibit {
	// TODO: Should be OK if this just returns nil but maybe not.
	return ac.exhibits
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
