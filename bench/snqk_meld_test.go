package bench

import (
	"testing"

	"github.com/madkins23/go-slog/bench/tests"
	"github.com/madkins23/go-slog/creator/snqkmeld"
)

// BenchmarkSnqkMeld runs benchmarks for the snqk/meld JSON handler.
func BenchmarkSnqkMeld(b *testing.B) {
	slogSuite := tests.NewSlogBenchmarkSuite(snqkmeld.Creator())
	tests.Run(b, slogSuite)
}
