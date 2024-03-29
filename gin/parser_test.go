package gin

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	args, err := Parse(`200 |    2.512908ms |             ::1 | GET      "/handler?tag=samber_zap"`)
	assert.NoError(t, err)
	assert.Len(t, args, 5)
	for _, thing := range args {
		if attr, ok := thing.(slog.Attr); assert.True(t, ok) {
			switch Field(attr.Key) {
			case Code:
				assert.Equal(t, "code=200", attr.String())
			case Elapsed:
				assert.Equal(t, "elapsed=2.512908ms", attr.String())
			case Client:
				assert.Equal(t, "client=::1", attr.String())
			case Method:
				assert.Equal(t, "method=GET", attr.String())
			case Url:
				assert.Equal(t, "url=/handler?tag=samber_zap", attr.String())
			default:
				assert.Fail(t, "unknown attribute key '%s'", attr.Key)
			}
		}
	}
}

func TestParse_Error_Split(t *testing.T) {
	args, err := Parse(`2XX |         ::1 | GET      "/handler?tag=samber_zap"`)
	assert.ErrorContains(t, err, "wrong number of parts")
	assert.ErrorContains(t, err, "3")
	assert.Nil(t, args)
}

func TestParse_Error_CodeNotNum(t *testing.T) {
	args, err := Parse(`2XX |    512908ms |             ::1 | GET      "/handler?tag=samber_zap"`)
	assert.ErrorContains(t, err, "parse Code")
	assert.ErrorContains(t, err, "2XX")
	assert.Nil(t, args)
}
