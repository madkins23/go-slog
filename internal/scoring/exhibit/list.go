package exhibit

import (
	_ "embed"
	"fmt"
	"html/template"

	"github.com/madkins23/go-slog/internal/scoring/score"
)

var (
	listName string = "List"
	listTmpl *template.Template

	//go:embed list.gohtml
	listSrc string
)

func setupList() error {
	var err error
	if listTmpl, err = template.New(listName).Parse(listSrc); err != nil {
		return fmt.Errorf("parse template '%s': %w", listName, err)
	}
	return nil
}

var _ score.Exhibit = &List{}

type List struct {
	score.ExhibitCore
	caption string
	items   []string
}

func NewList(caption string, items []string) *List {
	return &List{
		ExhibitCore: score.NewExhibitCore(listTmpl),
		caption:     caption,
		items:       items,
	}
}

func (t *List) HasCaption() bool {
	return len(t.caption) > 0
}

func (t *List) Caption() template.HTML {
	return template.HTML(t.caption)
}

func (t *List) Items() []string {
	return t.items
}
