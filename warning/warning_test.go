package warning

import (
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
