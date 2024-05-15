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

type Table struct {
	score.Exhibit
	columns []string
	rows    [][]string
}

func setupTable() error {
	var err error
	if tableTmpl, err = template.New(tableName).Parse(tableSrc); err != nil {
		return fmt.Errorf("parse template '%s': %w", tableName, err)
	}
	return nil
}

func NewTable(name string) *Table {
	return &Table{
		Exhibit: score.NewExhibit(name, tableTmpl),
		columns: []string{"alpha", "bravo", "charlie"},
		rows: [][]string{
			{"one", "13", "booger"},
			{"two", "17", "goober"},
			{"three", "23", "snoofus"},
		},
	}
}

func (t *Table) Columns() []string {
	return t.columns
}

func (t *Table) Rows() [][]string {
	return t.rows
}
