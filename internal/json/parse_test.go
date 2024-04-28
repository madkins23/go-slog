package json

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testJSON = `
		{
			"alpha": 23,
			"bravo": "Cheers!",
			"pi":    3.14159
		}`
	badJSON = `
		{
			"alpha": 23:
			"bravo": "Cheers!";
			"pi":    3.14159
		}`
	badError = "unmarshal json: invalid character ':' after object key:value pair"
)

var (
	testMap = map[string]any{
		"alpha": float64(23),
		"bravo": "Cheers!",
		"pi":    math.Round(math.Pi*100_000) / 100_000,
	}
	badMap = map[string]any{
		"error": badError,
	}
)

func TestExpect(t *testing.T) {
	assert.Equal(t, testMap, Expect(testJSON))
}

func TestExpect_error(t *testing.T) {
	assert.Equal(t, badMap, Expect(badJSON))
}

func TestParse(t *testing.T) {
	logMap, err := Parse([]byte(testJSON))
	require.NoError(t, err)
	assert.Equal(t, testMap, logMap)
}

func TestParse_error(t *testing.T) {
	_, err := Parse([]byte(badJSON))
	require.ErrorContains(t, err, badError)
}
