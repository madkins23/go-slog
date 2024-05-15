package exhibit

import (
	_ "embed"
	"fmt"
	"html/template"

	"github.com/madkins23/go-slog/internal/scoring/score"
)

var (
	tableName string = "Table"
	tableTmpl *template.Template

	//go:embed table.gohtml
	tableSrc string
)

func setupTable() error {
	var err error
	if tableTmpl, err = template.New(tableName).Parse(tableSrc); err != nil {
		return fmt.Errorf("parse template '%s': %w", tableName, err)
	}
	return nil
}

// TODO: This should probably work.
var _ score.Exhibit = &Table{}

type Table struct {
	score.ExhibitCore
	caption string
	columns []string
	rows    [][]string
}

func NewTable(caption string, columns []string, rows [][]string) *Table {
	return &Table{
		ExhibitCore: score.NewExhibitCore(tableTmpl),
		caption:     caption,
		columns:     columns,
		rows:        rows,
	}
}

func (t *Table) HasCaption() bool {
	return len(t.caption) > 0
}

func (t *Table) Caption() template.HTML {
	return template.HTML(t.caption)
}

func (t *Table) Columns() []string {
	return t.columns
}

func (t *Table) Rows() [][]string {
	return t.rows
}
