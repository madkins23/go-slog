package tests

import (
	"context"
	"log/slog"

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
		func(captured []byte, logMap map[string]any, manager *infra.WarningManager) bool {
			if len(captured) > 0 {
				manager.AddWarning(warning.NotDisabled, "Disabled", string(captured))
				return false
			}
			return true
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
		matcher("Simple", expectedBasic()))
}

//goland:noinspection GoSnakeCaseUsage
func (suite *SlogBenchmarkSuite) Benchmark_Simple_Source() Benchmark {
	return NewBenchmark(infra.SourceOptions(),
		func(logger *slog.Logger) {
			logger.Info(message)
		},
		nil,
		verify(
			finder("Simple_Source", expectedBasic()),
			sorcerer("Simple_Source")))
}

// -----------------------------------------------------------------------------

//goland:noinspection GoSnakeCaseUsage
func (suite *SlogBenchmarkSuite) Benchmark_Attributes() Benchmark {
	return NewBenchmark(infra.SimpleOptions(),
		func(logger *slog.Logger) {
			logger.LogAttrs(context.Background(), slog.LevelInfo, message, allAttributes...)
		},
		nil,
		verify(
			finder("Attributes", expectedBasic()),
			finder("Attributes", allValuesMap()),
			noDuplicates("Attributes"),
		))
}

//goland:noinspection GoSnakeCaseUsage
func (suite *SlogBenchmarkSuite) Benchmark_Key_Values() Benchmark {
	return NewBenchmark(infra.SimpleOptions(),
		func(logger *slog.Logger) {
			logger.Info(message, allKeyValues...)
		},
		nil,
		verify(
			finder("Key_Values", expectedBasic()),
			finder("Key_Values", allValuesMap()),
			noDuplicates("Key_Values"),
		))
}

// -----------------------------------------------------------------------------

//goland:noinspection GoSnakeCaseUsage
func (suite *SlogBenchmarkSuite) Benchmark_With_Attrs_Simple() Benchmark {
	return NewBenchmark(infra.SimpleOptions(),
		func(logger *slog.Logger) {
			logger.Info(message)
		},
		withAllAttributes,
		verify(
			finder("With_Attrs_Simple:Basic", expectedBasic()),
			finder("With_Attrs_Simple:With", withValuesMap()),
			noDuplicates("With_Attrs_Simple"),
		))
}

//goland:noinspection GoSnakeCaseUsage
func (suite *SlogBenchmarkSuite) Benchmark_With_Attrs_Attributes() Benchmark {
	return NewBenchmark(infra.SimpleOptions(),
		func(logger *slog.Logger) {
			logger.LogAttrs(context.Background(), slog.LevelInfo, message, allAttributes...)
		},
		withAllAttributes,
		verify(
			finder("With_Attrs_Attributes:Basic", expectedBasic()),
			finder("With_Attrs_Attributes:All", allValuesMap()),
			finder("With_Attrs_Attributes:With", withValuesMap()),
			noDuplicates("With_Attrs_Attributes"),
		))
}

//goland:noinspection GoSnakeCaseUsage
func (suite *SlogBenchmarkSuite) Benchmark_With_Attrs_Key_Values() Benchmark {
	return NewBenchmark(infra.SimpleOptions(),
		func(logger *slog.Logger) {
			logger.Info(message, allKeyValues...)
		},
		withAllAttributes,
		verify(
			finder("With_Attrs_Key_Values:Basic", expectedBasic()),
			finder("With_Attrs_Key_Values:All", allValuesMap()),
			finder("With_Attrs_Key_Values:With", withValuesMap()),
			noDuplicates("With_Attrs_Key_Values"),
		))
}

// -----------------------------------------------------------------------------

//goland:noinspection GoSnakeCaseUsage
func (suite *SlogBenchmarkSuite) Benchmark_With_Group_Attributes() Benchmark {
	return NewBenchmark(infra.SimpleOptions(),
		func(logger *slog.Logger) {
			logger.LogAttrs(context.Background(), slog.LevelInfo, message, allAttributes...)
		},
		withGroup,
		verify(
			finder("With_Group_Attributes:Basic", expectedBasic()),
			finder("With_Group_Attributes:All", map[string]any{
				"withGroup": allValuesMap(),
			}),
			noDuplicates("With_Group_Attributes"),
		))
}

//goland:noinspection GoSnakeCaseUsage
func (suite *SlogBenchmarkSuite) Benchmark_With_Group_Key_Values() Benchmark {
	return NewBenchmark(infra.SimpleOptions(),
		func(logger *slog.Logger) {
			logger.Info(message, allKeyValues...)
		},
		withGroup,
		verify(
			finder("With_Group_Key_Values:Basic", expectedBasic()),
			finder("With_Group_Key_Values:All", map[string]any{
				"withGroup": allValuesMap(),
			}),
			noDuplicates("With_Group_Key_Values"),
		))
}

// -----------------------------------------------------------------------------
// Large/Long tests.

//goland:noinspection GoSnakeCaseUsage
func (suite *SlogBenchmarkSuite) Benchmark_Big_Group() Benchmark {
	return NewBenchmark(infra.SimpleOptions(),
		func(logger *slog.Logger) {
			bg := BigGroup()
			logger.Info(message, bg)
		},
		nil,
		bigGroupChecker("Big_Group"))
}

//goland:noinspection GoSnakeCaseUsage
func (suite *SlogBenchmarkSuite) Benchmark_Logging() Benchmark {
	bm := NewBenchmark(infra.SimpleOptions(),
		func(logger *slog.Logger) {
			for _, logData := range logData() {
				logger.Info("Handle", logData...)
			}
		},
		nil,
		verifyLines(fields("Logging", "level", "msg", "code", "duration", "method", "sys", "url")))
	bm.SetDontCount(true)
	return bm
}

// -----------------------------------------------------------------------------

// getLogMap returns the specified logMap, if not empty, or a new one created from the captured bytes.
// If a new logMap is created it is run through fixLogMap before returning it.
func getLogMap(captured []byte, logMap map[string]any, manager *infra.WarningManager) map[string]any {
	var err error
	if logMap == nil {
		if logMap, err = parseLogMap(captured); err != nil {
			manager.AddWarning(warning.TestError, err.Error(), string(captured))
		}
		fixLogMap(logMap)
	}
	return logMap
}
