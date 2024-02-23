package bench

import (
	"testing"

	"github.com/madkins23/go-slog/warning"
)

// TestMain captures the Go test harness to show warning results after testing.
// This function is defined separately from the other test files
// because it can only be defined once in the package.
func TestMain(m *testing.M) {
	warning.WithWarnings(m)
}
