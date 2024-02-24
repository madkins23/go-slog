package bench

import (
	"testing"

	"github.com/madkins23/go-slog/bench/tests"
	"github.com/madkins23/go-slog/creator/slogjson"
)

// Benchmark_slog runs benchmarks for the log/slog JSON handler.
func Benchmark_slog(b *testing.B) {
	slogSuite := tests.NewSlogBenchmarkSuite(slogjson.Creator())
	tests.Run(b, slogSuite)
}
