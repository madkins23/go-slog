package verify

import (
	"testing"

	"github.com/madkins23/go-slog/infra/warning"
)

func TestMain(m *testing.M) {
	warning.WithWarnings(m)
}
