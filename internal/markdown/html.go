package markdown

import (
	"html/template"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

var renderer *html.Renderer

func init() {
	// Prebuild the HTML renderer.
	renderer = html.NewRenderer(html.RendererOptions{Flags: html.HrefTargetBlank})
}

// TemplateHTML renders Markdown into template.HTML suitable for use in cmd/server.
// Embedded Markdown in Go code blocks is typically quoted using back quotes ("`")
// which are used within Markdown for fixed font markup.
//
// The fixCarets flag will change every caret character ("^") in the provided string
// to a back quote ("`") prior to parsing the Markdown text,
// allowing carets to be used for fixed font markup instead of back quotes.
// Set fixCarets to true for use with back quoted strings with carets,
// false for use with embedded Markdown files containing back quotes.
func TemplateHTML(md string, fixCarets bool) template.HTML {
	// Can't pre-build the parser in init(), it fails the second time it's used.
	mdParser := parser.NewWithExtensions(
		parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock)
	if fixCarets {
		md = strings.Replace(md, "^", "`", -1)
	}
	return template.HTML(markdown.Render(mdParser.Parse([]byte(md)), renderer))
}
