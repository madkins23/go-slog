package score

import (
	"html/template"
	"math"

	"github.com/madkins23/go-slog/internal/data"
)

type Axis interface {
	Setup(bench *data.Benchmarks, warns *data.Warnings) error
	AxisTitle() string
	Name() string
	HandlerScore(handler data.HandlerTag) Value
	Documentation() template.HTML
	ExhibitCount() uint
	Exhibits() []Exhibit
}

type Value float64

func (v Value) Round() Value {
	const rounder = 1_000_000_000.0
	return Value(math.Round(float64(v)*rounder) / rounder)
}
