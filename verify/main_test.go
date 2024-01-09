package verify

import (
	"testing"

	"github.com/madkins23/go-slog/verify/test"
)

func TestMain(m *testing.M) {
	test.WithWarnings(m)
}
