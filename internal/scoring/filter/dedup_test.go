package filter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/madkins23/go-slog/creator/madkinsflash"
	"github.com/madkins23/go-slog/creator/slogjson"
	"github.com/madkins23/go-slog/creator/veqryndedup"
)

func TestDedup(t *testing.T) {
	filter := Dedup()
	require.NotNil(t, filter)
	assert.False(t, filter.Keep(madkinsflash.Name))
	assert.True(t, filter.Keep(slogjson.Name))
	assert.False(t, filter.Keep(veqryndedup.Name(veqryndedup.Append)))
	assert.False(t, filter.Keep(veqryndedup.Name(veqryndedup.Ignore)))
	assert.False(t, filter.Keep(veqryndedup.Name(veqryndedup.Increment)))
	assert.True(t, filter.Keep(veqryndedup.Name(veqryndedup.Overwrite)))
}
