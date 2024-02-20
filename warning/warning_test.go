package warning

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	NumRequired  = 8
	NumImplied   = 5
	NumSuggested = 5
	NumAdmin     = 5
	NumBenchmark = 2
)

func TestAllWarnings(t *testing.T) {
	assert.Len(t, allWarnings, NumRequired+NumImplied+NumSuggested+NumAdmin)
}

func TestWarnings(t *testing.T) {
	assert.Len(t, Required(), NumRequired)
	assert.Len(t, Implied(), NumImplied)
	assert.Len(t, Suggested(), NumSuggested)
	assert.Len(t, Administrative(), NumAdmin)
	assert.Len(t, Benchmark(), NumBenchmark)
}
