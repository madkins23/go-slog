package bench

import (
	"testing"

	"github.com/madkins23/go-slog/bench/tests"
	"github.com/madkins23/go-slog/creator"
)

// Benchmark_slog_samber_zerolog runs benchmarks for the samber/slog-zerolog handler.
func Benchmark_slog_samber_zerolog(b *testing.B) {
	slogSuite := tests.NewSlogBenchmarkSuite(creator.SlogSamberZerolog())
	tests.Run(b, slogSuite)
}