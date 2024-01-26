package gin

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmpty(t *testing.T) {
	assert.IsType(t, map[Field]string{}, Empty())
	assert.Empty(t, Empty())
}

func TestNewParser(t *testing.T) {
	p := NewParser(nil)
	require.NotNil(t, p)
	assert.IsType(t, &parser{}, p)
}

func TestParser_Parse(t *testing.T) {

}
