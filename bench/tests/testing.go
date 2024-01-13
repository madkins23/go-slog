package tests

import "log/slog"

// -----------------------------------------------------------------------------
// Benchmarks for testing the suite.
//
// Benchmark methods have names beginning with "benchmark" (all lowercase).
// They return a benchmark object containing the options for logger creation and
// the function to run during the benchmark.

func (suite *SlogBenchmarkSuite) Benchmark1() Benchmark {
	return NewBenchmark(&slog.HandlerOptions{}, func(logger *slog.Logger) {
		logger.Info("Benchmark1")
	})
}

func (suite *SlogBenchmarkSuite) Benchmark2() Benchmark {
	return NewBenchmark(&slog.HandlerOptions{}, func(logger *slog.Logger) {
		logger.Info("Benchmark2")
	})
}
