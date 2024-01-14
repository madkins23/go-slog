package tests

import (
	"context"
	"log/slog"

	"github.com/madkins23/go-slog/infra"
)

// -----------------------------------------------------------------------------
// Benchmarks for testing the suite.
//
// Benchmark methods have names beginning with "benchmark" (all lowercase).
// They return a benchmark object containing the options for logger creation and
// the function to run during the benchmark.

// -----------------------------------------------------------------------------
// Basic tests.

func (suite *SlogBenchmarkSuite) Benchmark_Disabled() Benchmark {
	mark := NewBenchmark(infra.SimpleOptions(), func(logger *slog.Logger) {
		logger.Info(message)
	}, nil)
	mark.SetDontCount(true)
	return mark
}

func (suite *SlogBenchmarkSuite) Benchmark_Simple() Benchmark {
	return NewBenchmark(infra.SimpleOptions(), func(logger *slog.Logger) {
		logger.Info(message)
	}, nil)
}

func (suite *SlogBenchmarkSuite) Benchmark_Simple_Source() Benchmark {
	return NewBenchmark(infra.SourceOptions(), func(logger *slog.Logger) {
		logger.Info(message)
	}, nil)
}

// -----------------------------------------------------------------------------

func (suite *SlogBenchmarkSuite) Benchmark_Attributes() Benchmark {
	return NewBenchmark(infra.SimpleOptions(), func(logger *slog.Logger) {
		logger.LogAttrs(context.Background(), slog.LevelInfo, message, allAttributes...)
	}, nil)
}

func (suite *SlogBenchmarkSuite) Benchmark_Key_Values() Benchmark {
	return NewBenchmark(infra.SimpleOptions(), func(logger *slog.Logger) {
		logger.Info(message, allKeyValues...)
	}, nil)
}

// -----------------------------------------------------------------------------

func (suite *SlogBenchmarkSuite) Benchmark_With_Attrs_Simple() Benchmark {
	return NewBenchmark(infra.SimpleOptions(), func(logger *slog.Logger) {
		logger.Info(message)
	}, withAllAttributes)
}

func (suite *SlogBenchmarkSuite) Benchmark_With_Attrs_Attributes() Benchmark {
	return NewBenchmark(infra.SimpleOptions(), func(logger *slog.Logger) {
		logger.LogAttrs(context.Background(), slog.LevelInfo, message, allAttributes...)
	}, withAllAttributes)
}

func (suite *SlogBenchmarkSuite) Benchmark_With_Attrs_Key_Values() Benchmark {
	return NewBenchmark(infra.SimpleOptions(), func(logger *slog.Logger) {
		logger.Info(message, allKeyValues...)
	}, withAllAttributes)
}

// -----------------------------------------------------------------------------

func (suite *SlogBenchmarkSuite) Benchmark_With_Group_Attributes() Benchmark {
	return NewBenchmark(infra.SimpleOptions(), func(logger *slog.Logger) {
		logger.LogAttrs(context.Background(), slog.LevelInfo, message, allAttributes...)
	}, withGroupAttributes)
}

func (suite *SlogBenchmarkSuite) Benchmark_With_Group_Key_Values() Benchmark {
	return NewBenchmark(infra.SimpleOptions(), func(logger *slog.Logger) {
		logger.Info(message, allKeyValues...)
	}, withGroupAttributes)
}

// -----------------------------------------------------------------------------

func (suite *SlogBenchmarkSuite) Benchmark_Big_Group() Benchmark {
	return NewBenchmark(infra.SimpleOptions(), func(logger *slog.Logger) {
		logger.Info(message, BigGroup())
	}, nil)
}

// -----------------------------------------------------------------------------

var _ HandlerFn = withAllAttributes

func withAllAttributes(handler slog.Handler) slog.Handler {
	return handler.WithAttrs(withAttributes)
}

var _ HandlerFn = withGroupAttributes

func withGroupAttributes(handler slog.Handler) slog.Handler {
	return handler.WithGroup("withGroup")
}
