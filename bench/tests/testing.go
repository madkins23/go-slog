package tests

import (
	"log/slog"

	"github.com/madkins23/go-slog/infra"
)

// -----------------------------------------------------------------------------
// Benchmarks for testing the suite.
//
// Benchmark methods have names beginning with "benchmark" (all lowercase).
// They return a benchmark object containing the options for logger creation and
// the function to run during the benchmark.

func (suite *SlogBenchmarkSuite) Benchmark_Simple() Benchmark {
	return NewBenchmark(infra.SimpleOptions(), func(logger *slog.Logger) {
		logger.Info(message)
	})
}

func (suite *SlogBenchmarkSuite) Benchmark_Simple_Source() Benchmark {
	return NewBenchmark(infra.SourceOptions(), func(logger *slog.Logger) {
		logger.Info(message)
	})
}
