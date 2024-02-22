package tests

import (
	"context"
	"log/slog"

	"github.com/madkins23/go-slog/infra"
	"github.com/madkins23/go-slog/internal/test"
	"github.com/madkins23/go-slog/warning"
)

// -----------------------------------------------------------------------------
// Benchmarks for testing the suite.
//
// Benchmark methods have names beginning with "benchmark" (all lowercase).
// They return a benchmark object containing the options for logger creation and
// the function to run during the benchmark.

// -----------------------------------------------------------------------------

type BenchmarkFn func(logger *slog.Logger)
type HandlerFn func(handler slog.Handler) slog.Handler
type VerifyFn func(captured []byte, logMap map[string]any, manager *test.WarningManager) error

type Benchmark struct {
	Options     *slog.HandlerOptions
	BenchmarkFn BenchmarkFn
	HandlerFn   HandlerFn
	VerifyFn    VerifyFn
	DontCount   bool
}

// -----------------------------------------------------------------------------
// Basic tests.

// Benchmark_Disabled logs at the wrong level and nothing comes out.
// The default logger is set to INFO, a DEBUG log message should be ignored.
//
//goland:noinspection GoSnakeCaseUsage
func (suite *SlogBenchmarkSuite) Benchmark_Disabled() *Benchmark {
	return &Benchmark{
		Options: infra.SimpleOptions(),
		BenchmarkFn: func(logger *slog.Logger) {
			logger.Debug(message)
		},
		VerifyFn: func(captured []byte, logMap map[string]any, manager *test.WarningManager) error {
			if len(captured) > 0 {
				manager.AddWarning(warning.NotDisabled, "Disabled", string(captured))
				return warning.NotDisabled
			}
			return nil
		},
		DontCount: true,
	}
}

// Benchmark_Simple logs a simple line with just a message.
//
//goland:noinspection GoSnakeCaseUsage
func (suite *SlogBenchmarkSuite) Benchmark_Simple() *Benchmark {
	return &Benchmark{
		Options: infra.SimpleOptions(),
		BenchmarkFn: func(logger *slog.Logger) {
			logger.Info(message)
		},
		VerifyFn: matcher("Simple", expectedBasic()),
	}
}

// Benchmark_Simple_Source logs a simple line with just a message
// and the AddSource option set so the slog.SourceKey group is created.
//
//goland:noinspection GoSnakeCaseUsage
func (suite *SlogBenchmarkSuite) Benchmark_Simple_Source() *Benchmark {
	return &Benchmark{
		Options: infra.SourceOptions(),
		BenchmarkFn: func(logger *slog.Logger) {
			logger.Info(message)
		},
		VerifyFn: verify(
			finder("Simple_Source", expectedBasic()),
			sourcerer("Simple_Source")),
	}
}

// -----------------------------------------------------------------------------

// Benchmark_Attributes logs a message with a lot of attributes.
//
//goland:noinspection GoSnakeCaseUsage
func (suite *SlogBenchmarkSuite) Benchmark_Attributes() *Benchmark {
	return &Benchmark{
		Options: infra.SimpleOptions(),
		BenchmarkFn: func(logger *slog.Logger) {
			logger.LogAttrs(context.Background(), slog.LevelInfo, message, allAttributes...)
		},
		VerifyFn: verify(
			finder("Attributes", expectedBasic()),
			finder("Attributes", allValuesMap()),
			noDuplicates("Attributes"),
		),
	}
}

// Benchmark_Key_Values logs a message with a lot of key value pairs.
//
//goland:noinspection GoSnakeCaseUsage
func (suite *SlogBenchmarkSuite) Benchmark_Key_Values() *Benchmark {
	return &Benchmark{
		Options: infra.SimpleOptions(),
		BenchmarkFn: func(logger *slog.Logger) {
			logger.Info(message, allKeyValues...)
		},
		VerifyFn: verify(
			finder("Key_Values", expectedBasic()),
			finder("Key_Values", allValuesMap()),
			noDuplicates("Key_Values"),
		),
	}
}

// -----------------------------------------------------------------------------

// Benchmark_With_Attrs_Simple logs a simple message to a logger created from
// a handler configure with a lot of attributes.
//
//goland:noinspection GoSnakeCaseUsage
func (suite *SlogBenchmarkSuite) Benchmark_With_Attrs_Simple() *Benchmark {
	return &Benchmark{
		Options: infra.SimpleOptions(),
		BenchmarkFn: func(logger *slog.Logger) {
			logger.Info(message)
		},
		HandlerFn: withAllAttributes,
		VerifyFn: verify(
			finder("With_Attrs_Simple:Basic", expectedBasic()),
			finder("With_Attrs_Simple:With", withValuesMap()),
			noDuplicates("With_Attrs_Simple"),
		),
	}
}

// Benchmark_With_Attrs_Attributes logs a message and a lot of attributes to a logger created from
// a handler configure with a different lot of attributes.
//
//goland:noinspection GoSnakeCaseUsage
func (suite *SlogBenchmarkSuite) Benchmark_With_Attrs_Attributes() *Benchmark {
	return &Benchmark{
		Options: infra.SimpleOptions(),
		BenchmarkFn: func(logger *slog.Logger) {
			logger.LogAttrs(context.Background(), slog.LevelInfo, message, allAttributes...)
		},
		HandlerFn: withAllAttributes,
		VerifyFn: verify(
			finder("With_Attrs_Attributes:Basic", expectedBasic()),
			finder("With_Attrs_Attributes:All", allValuesMap()),
			finder("With_Attrs_Attributes:With", withValuesMap()),
			noDuplicates("With_Attrs_Attributes"),
		),
	}
}

// Benchmark_With_Attrs_Key_Values logs a message and a lot of attributes to a logger created from
// a handler configure with a different lot of key value pairs.
//
//goland:noinspection GoSnakeCaseUsage
func (suite *SlogBenchmarkSuite) Benchmark_With_Attrs_Key_Values() *Benchmark {
	return &Benchmark{
		Options: infra.SimpleOptions(),
		BenchmarkFn: func(logger *slog.Logger) {
			logger.Info(message, allKeyValues...)
		},
		HandlerFn: withAllAttributes,
		VerifyFn: verify(
			finder("With_Attrs_Key_Values:Basic", expectedBasic()),
			finder("With_Attrs_Key_Values:All", allValuesMap()),
			finder("With_Attrs_Key_Values:With", withValuesMap()),
			noDuplicates("With_Attrs_Key_Values"),
		),
	}
}

// -----------------------------------------------------------------------------

// Benchmark_With_Group_Attributes logs a message and a lot of attributes to
// a logger created from a handler configured with an open group.
//
//goland:noinspection GoSnakeCaseUsage
func (suite *SlogBenchmarkSuite) Benchmark_With_Group_Attributes() *Benchmark {
	return &Benchmark{
		Options: infra.SimpleOptions(),
		BenchmarkFn: func(logger *slog.Logger) {
			logger.LogAttrs(context.Background(), slog.LevelInfo, message, allAttributes...)
		},
		HandlerFn: withGroup,
		VerifyFn: verify(
			finder("With_Group_Attributes:Basic", expectedBasic()),
			finder("With_Group_Attributes:All", map[string]any{
				"withGroup": allValuesMap(),
			}),
			noDuplicates("With_Group_Attributes"),
		),
	}
}

// Benchmark_With_Group_Key_Values logs a message and a lot of key value pairs to
// a logger created from a handler configured with an open group.
//
//goland:noinspection GoSnakeCaseUsage
func (suite *SlogBenchmarkSuite) Benchmark_With_Group_Key_Values() *Benchmark {
	return &Benchmark{
		Options: infra.SimpleOptions(),
		BenchmarkFn: func(logger *slog.Logger) {
			logger.Info(message, allKeyValues...)
		},
		HandlerFn: withGroup,
		VerifyFn: verify(
			finder("With_Group_Key_Values:Basic", expectedBasic()),
			finder("With_Group_Key_Values:All", map[string]any{
				"withGroup": allValuesMap(),
			}),
			noDuplicates("With_Group_Key_Values"),
		),
	}
}

// -----------------------------------------------------------------------------
// Large/Long tests.

// Benchmark_Big_Group logs several levels of nested groups.
//
//goland:noinspection GoSnakeCaseUsage
func (suite *SlogBenchmarkSuite) Benchmark_Big_Group() *Benchmark {
	return &Benchmark{
		Options: infra.SimpleOptions(),
		BenchmarkFn: func(logger *slog.Logger) {
			bg := BigGroup()
			logger.Info(message, bg)
		},
		VerifyFn: bigGroupChecker("Big_Group"),
	}
}

// Benchmark_Logging logs a series of lines taken from Gin server output.
//
//goland:noinspection GoSnakeCaseUsage
func (suite *SlogBenchmarkSuite) Benchmark_Logging() *Benchmark {
	return &Benchmark{
		Options: infra.SimpleOptions(),
		BenchmarkFn: func(logger *slog.Logger) {
			for _, logData := range logData() {
				logger.Info("Handle", logData...)
			}
		},
		VerifyFn:  verifyLines(fields("Logging", "level", "msg", "code", "elapsed", "method", "url")),
		DontCount: true,
	}
}

// -----------------------------------------------------------------------------

// getLogMap returns the specified logMap, if not empty, or a new one created from the captured bytes.
// If a new logMap is created it is run through fixLogMap before returning it.
func getLogMap(captured []byte, logMap map[string]any, manager *test.WarningManager) map[string]any {
	var err error
	if logMap == nil {
		if logMap, err = parseLogMap(captured); err != nil {
			manager.AddWarning(warning.TestError, err.Error(), string(captured))
		}
		fixLogMap(logMap)
	}
	return logMap
}
