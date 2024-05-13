package score

import (
	"html/template"

	"github.com/madkins23/go-slog/internal/data"
)

type Value float64

func (v Value) Float64() float64 {
	return float64(v)
}

type Axis interface {
	Initialize(bench *data.Benchmarks, warns *data.Warnings) error
	ColumnHeader() string
	HandlerScore(handler data.HandlerTag) Value
	Documentation() template.HTML
	// DocTables() []???
}
