package verify

import (
	"testing"

	"github.com/madkins23/go-slog/internal/warning"
)

func TestMain(m *testing.M) {
	warning.WithWarnings(m)
}
