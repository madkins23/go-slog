package bench

import (
	"testing"

	"github.com/madkins23/go-slog/bench/tests"
	"github.com/madkins23/go-slog/creator"
)

// Benchmark_slog_chanchal_zap runs benchmarks for the chanchal/zaphandler handler.
func Benchmark_slog_chanchal_zap(b *testing.B) {
	slogSuite := tests.NewSlogBenchmarkSuite(creator.SlogChanchalZapHandler())
	tests.Run(b, slogSuite)
}
