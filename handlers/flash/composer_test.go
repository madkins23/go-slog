package flash

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/madkins23/go-slog/internal/test"
)

func TestComposer_addEscape(t *testing.T) {
	for escStr, expStr := range test.EscapeCases {
		c := newComposer([]byte{}, true, nil, nil, fixExtras(nil))
		c.addEscaped([]byte(escStr))
		x := string(c.buffer)
		assert.Equal(t, expStr, x)
	}
}
