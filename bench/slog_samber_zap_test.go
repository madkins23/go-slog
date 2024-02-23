package bench

import (
	"testing"

	"github.com/madkins23/go-slog/bench/tests"
	"github.com/madkins23/go-slog/creator/zap/samber"
)

// Benchmark_slog_samber_zap runs benchmarks for the darvaza/slog-zap handler.
func Benchmark_slog_samber_zap(b *testing.B) {
	slogSuite := tests.NewSlogBenchmarkSuite(samber.SlogSamberZap())
	tests.Run(b, slogSuite)
}
