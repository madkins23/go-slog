package filter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/madkins23/go-slog/creator/madkinsflash"
	"github.com/madkins23/go-slog/creator/madkinsreplattr"
	"github.com/madkins23/go-slog/creator/slogjson"
	"github.com/madkins23/go-slog/creator/veqryndedup"
)

func TestReplAttr(t *testing.T) {
	filter := ReplAttr()
	require.NotNil(t, filter)
	assert.True(t, filter.Keep(madkinsflash.Name))
	assert.True(t, filter.Keep(madkinsreplattr.Name))
	assert.False(t, filter.Keep(slogjson.Name))
	assert.False(t, filter.Keep(veqryndedup.Name(veqryndedup.Overwrite)))
}
