package warning

import (
	"html/template"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSomeWarnings(t *testing.T) {
	assert.Greater(t, len(allWarnings), 0)
}

func TestAllWarnings(t *testing.T) {
	assert.Len(t, allWarnings,
		testCounts[LevelRequired]+
			testCounts[LevelImplied]+
			testCounts[LevelSuggested]+
			testCounts[LevelAdmin])
}

func TestWarnings(t *testing.T) {
	assert.Len(t, Required(), testCounts[LevelRequired])
	assert.Len(t, Implied(), testCounts[LevelImplied])
	assert.Len(t, Suggested(), testCounts[LevelSuggested])
	assert.Len(t, Administrative(), testCounts[LevelAdmin])
}

func TestDescription(t *testing.T) {
	assert.Equal(t,
		template.HTML("<p>Handlers are supposed to avoid logging empty attributes.</p>\n\n<ul>\n<li><a href=\"https://pkg.go.dev/log/slog@master#Handler\" target=\"_blank\">&quot;- If an Attr's key and value are both the zero value, ignore the Attr.&quot;</a></li>\n</ul>\n"),
		EmptyAttributes.Description())
}
