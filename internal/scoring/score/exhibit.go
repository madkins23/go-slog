package score

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/gomarkdown/markdown/html"
)

type Exhibit interface {
	HTML(data any) template.HTML
}

type ExhibitCore struct {
	tmpl *template.Template
}

func NewExhibitCore(tmpl *template.Template) ExhibitCore {
	return ExhibitCore{tmpl: tmpl}
}

func (ec *ExhibitCore) HTML(data any) template.HTML {
	var buf bytes.Buffer
	if err := ec.tmpl.Execute(&buf, data); err != nil {
		var errorText bytes.Buffer
		html.EscapeHTML(&errorText, []byte(err.Error()))
		// TODO: Make something more general, read markdown from a file or something?
		return template.HTML(fmt.Sprintf("<h3>Error:</h3><p>%s</p>", errorText.String()))
	} else {
		return template.HTML(buf.String())
	}
}
