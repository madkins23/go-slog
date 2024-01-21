package tests

import (
	"context"
	"log/slog"
	"reflect"

	"github.com/madkins23/go-slog/infra"
	"github.com/madkins23/go-slog/warning"
)

// -----------------------------------------------------------------------------
// Benchmarks for testing the suite.
//
// Benchmark methods have names beginning with "benchmark" (all lowercase).
// They return a benchmark object containing the options for logger creation and
// the function to run during the benchmark.

// -----------------------------------------------------------------------------
// Basic tests.

//goland:noinspection GoSnakeCaseUsage
func (suite *SlogBenchmarkSuite) Benchmark_Disabled() Benchmark {
	bm := NewBenchmark(infra.SimpleOptions(),
		func(logger *slog.Logger) {
			logger.Debug(message)
		},
		nil,
		func(captured []byte, manager *infra.WarningManager) {
			if len(captured) > 0 {
				manager.AddWarning(NotDisabled, "Disabled", string(captured))
			}
		})
	bm.SetDontCount(true)
	return bm
}

//goland:noinspection GoSnakeCaseUsage
func (suite *SlogBenchmarkSuite) Benchmark_Simple() Benchmark {
	return NewBenchmark(infra.SimpleOptions(),
		func(logger *slog.Logger) {
			logger.Info(message)
		},
		nil,
		matcher("Simple", map[string]any{
			slog.LevelKey:   slog.LevelInfo.String(),
			slog.MessageKey: message,
		}))
}

//goland:noinspection GoSnakeCaseUsage
func (suite *SlogBenchmarkSuite) Benchmark_Simple_Source() Benchmark {
	return NewBenchmark(infra.SourceOptions(),
		func(logger *slog.Logger) {
			logger.Info(message)
		},
		nil,
		matcher("Simple_Source", map[string]any{
			slog.LevelKey:   slog.LevelInfo.String(),
			slog.MessageKey: message,
			slog.SourceKey:  []string{"file", "function", "line"},
		}))
}

// -----------------------------------------------------------------------------

//goland:noinspection GoSnakeCaseUsage
func (suite *SlogBenchmarkSuite) Benchmark_Attributes() Benchmark {
	return NewBenchmark(infra.SimpleOptions(), func(logger *slog.Logger) {
		logger.LogAttrs(context.Background(), slog.LevelInfo, message, allAttributes...)
	}, nil, nil)
}

//goland:noinspection GoSnakeCaseUsage
func (suite *SlogBenchmarkSuite) Benchmark_Key_Values() Benchmark {
	return NewBenchmark(infra.SimpleOptions(), func(logger *slog.Logger) {
		logger.Info(message, allKeyValues...)
	}, nil, nil)
}

// -----------------------------------------------------------------------------

//goland:noinspection GoSnakeCaseUsage
func (suite *SlogBenchmarkSuite) Benchmark_With_Attrs_Simple() Benchmark {
	return NewBenchmark(infra.SimpleOptions(), func(logger *slog.Logger) {
		logger.Info(message)
	}, withAllAttributes, nil)
}

//goland:noinspection GoSnakeCaseUsage
func (suite *SlogBenchmarkSuite) Benchmark_With_Attrs_Attributes() Benchmark {
	return NewBenchmark(infra.SimpleOptions(), func(logger *slog.Logger) {
		logger.LogAttrs(context.Background(), slog.LevelInfo, message, allAttributes...)
	}, withAllAttributes, nil)
}

//goland:noinspection GoSnakeCaseUsage
func (suite *SlogBenchmarkSuite) Benchmark_With_Attrs_Key_Values() Benchmark {
	return NewBenchmark(infra.SimpleOptions(), func(logger *slog.Logger) {
		logger.Info(message, allKeyValues...)
	}, withAllAttributes, nil)
}

// -----------------------------------------------------------------------------

//goland:noinspection GoSnakeCaseUsage
func (suite *SlogBenchmarkSuite) Benchmark_With_Group_Attributes() Benchmark {
	return NewBenchmark(infra.SimpleOptions(), func(logger *slog.Logger) {
		logger.LogAttrs(context.Background(), slog.LevelInfo, message, allAttributes...)
	}, withGroup, nil)
}

//goland:noinspection GoSnakeCaseUsage
func (suite *SlogBenchmarkSuite) Benchmark_With_Group_Key_Values() Benchmark {
	return NewBenchmark(infra.SimpleOptions(), func(logger *slog.Logger) {
		logger.Info(message, allKeyValues...)
	}, withGroup, nil)
}

// -----------------------------------------------------------------------------
// Large/Long tests.

//goland:noinspection GoSnakeCaseUsage
func (suite *SlogBenchmarkSuite) Benchmark_Big_Group() Benchmark {
	return NewBenchmark(infra.SimpleOptions(), func(logger *slog.Logger) {
		logger.Info(message, BigGroup())
	}, nil, nil)
}

//goland:noinspection GoSnakeCaseUsage
func (suite *SlogBenchmarkSuite) Benchmark_Logging() Benchmark {
	test := NewBenchmark(infra.SimpleOptions(), func(logger *slog.Logger) {
		for _, logData := range logData() {
			logger.Info("Handle", logData...)
		}
	}, nil, nil)
	test.SetDontCount(true)
	return test
}

// -----------------------------------------------------------------------------

func matcher(testName string, expected map[string]any) VerifyFn {
	return func(captured []byte, manager *infra.WarningManager) {
		if logMap, err := logMap(captured); err != nil {
			manager.AddWarning(warning.TestError, err.Error(), string(captured))
		} else if !reflect.DeepEqual(expected, fixLogMap(logMap)) {
			manager.AddWarningFn(Mismatch, testName, string(captured))
		}
	}
}
