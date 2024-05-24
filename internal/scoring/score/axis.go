package score

import (
	"html/template"

	"github.com/madkins23/go-slog/internal/data"
)

type Value float64

type Axis interface {
	Setup(bench *data.Benchmarks, warns *data.Warnings) error
	AxisTitle() string
	Name() string
	HandlerScore(handler data.HandlerTag) Value
	Documentation() template.HTML
	ExhibitCount() uint
	Exhibits() []Exhibit
}
