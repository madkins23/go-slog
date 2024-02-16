package bench

import (
	"testing"

	"github.com/madkins23/go-slog/internal/test"
)

func TestMain(m *testing.M) {
	test.WithWarnings(m)
}
