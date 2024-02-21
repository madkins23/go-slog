package warning

import (
	"fmt"
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

func MD2HTML(md string) template.HTML {
	fmt.Printf(">-> %s\n", md)
	// Can't pre-build the parser in init(), it fails the second time it's used.
	mdParser := parser.NewWithExtensions(
		parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock)
	return template.HTML(markdown.Render(
		mdParser.Parse([]byte(strings.Replace(md, "^", "`", -1))),
		renderer))
}
