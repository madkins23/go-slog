package score

import (
	"html/template"

	"github.com/madkins23/go-slog/internal/data"
)

type Value float64

type Axis interface {
	Initialize(bench *data.Benchmarks, warns *data.Warnings) error
	ColumnHeader() string
	HandlerScore(handler data.HandlerTag) Value
	Documentation() template.HTML
	// DocTables() []???
}
