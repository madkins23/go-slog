package score

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/gomarkdown/markdown/html"
)

type Exhibit struct {
	name string
	tmpl *template.Template
}

func NewExhibit(name string, tmpl *template.Template) Exhibit {
	return Exhibit{
		name: name,
		tmpl: tmpl,
	}
}

func (e *Exhibit) Name() string {
	return e.name
}

func (e *Exhibit) HTML(data any) template.HTML {
	var buf bytes.Buffer
	if err := e.tmpl.Execute(&buf, data); err != nil {
		var errorText bytes.Buffer
		html.EscapeHTML(&errorText, []byte(err.Error()))
		// TODO: Make something more general, read markdown from a file or something?
		return template.HTML(fmt.Sprintf("<h3>Error:</h3><p>%s</p>", errorText.String()))
	} else {
		return template.HTML(buf.String())
	}
}
