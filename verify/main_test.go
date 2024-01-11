package verify

import (
	"testing"

	"github.com/madkins23/go-slog/verify/tests"
)

func TestMain(m *testing.M) {
	tests.WithWarnings(m)
}
