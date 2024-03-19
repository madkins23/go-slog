package json

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	x := Parse([]byte(`
		{
			"alpha": 23,
			"bravo": "Cheers!"
		}`))
	assert.NotNil(t, x)
	assert.IsType(t, map[string]any{}, x)
	alpha, found := x["alpha"]
	assert.True(t, found)
	assert.Equal(t, float64(23), alpha)
	bravo, found := x["bravo"]
	assert.True(t, found)
	assert.Equal(t, "Cheers!", bravo)
}
