package score

import (
	"fmt"
	"html/template"

	"github.com/madkins23/go-slog/internal/scoring/setup"
)

var _ setup.Item = &Template{}

type Template struct {
	name string
	text string
	tmpl *template.Template
}

func (t *Template) Name() string {
	return t.name
}

func (t *Template) Setup() error {
	var err error
	if t.tmpl, err = template.New(t.name).Parse(t.text); err != nil {
		return fmt.Errorf("parse template: %w", err)
	}
	return nil
}
