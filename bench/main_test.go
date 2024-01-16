package bench

import (
	"testing"

	"github.com/madkins23/go-slog/infra"
)

func TestMain(m *testing.M) {
	infra.WithWarnings(m)
}
