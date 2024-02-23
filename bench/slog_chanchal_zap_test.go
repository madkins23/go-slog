package bench

import (
	"testing"

	"github.com/madkins23/go-slog/bench/tests"
	"github.com/madkins23/go-slog/creator/zap/chanchal"
)

// Benchmark_slog_chanchal_zap runs benchmarks for the chanchal/zap handler.
func Benchmark_slog_chanchal_zap(b *testing.B) {
	slogSuite := tests.NewSlogBenchmarkSuite(chanchal.SlogChanchalZapHandler())
	tests.Run(b, slogSuite)
}
