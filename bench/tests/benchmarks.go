package tests

import (
	"context"
	"log/slog"

	"github.com/madkins23/go-slog/infra"
	warning2 "github.com/madkins23/go-slog/infra/warning"
)

// -----------------------------------------------------------------------------
// Benchmarks for testing the suite.
//
// Benchmark methods have names beginning with "benchmark" (all lowercase).
// They return a benchmark object containing the options for logger creation and
// the function to run during the benchmark.

// -----------------------------------------------------------------------------

// BenchmarkFn defines a function that executes a benchmark test.
type BenchmarkFn func(logger *slog.Logger)

// HandlerFn defines a function that adjusts a slog.Handler prior to using it
// to generate a slog.Logger.
// The general use for this is to apply WithAttrs and/or WithGroup methods to the handler.
type HandlerFn func(handler slog.Handler) slog.Handler

// VerifyFn defines a function that is used to verify the functionality of a benchmark test.
// Verification makes certain that there is actually something happening in the benchmark
// (and that it is generating the expected output) to avoid having any zombie benchmarks
// that don't actually do anything.
type VerifyFn func(captured []byte, logMap map[string]any, manager *warning2.Manager) error

// Benchmark objects are used to define benchmark tests.
type Benchmark struct {
	// Options is a required pointer to a preloaded slog.HandlerOptions object (e.g. infra.SimpleOptions).
	Options *slog.HandlerOptions

	// BenchmarkFn is a required BenchmarkFn which executes the actual benchmark test.
	BenchmarkFn

	// HandlerFn is an optional HandlerFn used to adjust the slog.Handler object
	// (if available from the infra.Creator for the benchmark) prior to using it to generate a slog.Logger object.
	HandlerFn

	// VerifyFn is an optional VerifyFn used to verify test log results.
	VerifyFn

	// DontCount is used to avoid double-checking the number of log lines written during the benchmark.
	// Normally a single line is generated for each benchmark execution,
	// so the log lines generated matches the number of executions of the BenchmarkFn,
	// but in some cases multiple lines are generated per execution.
	// This field is optional as its default value is false (number of lines/executions must match).
	DontCount bool
}

// -----------------------------------------------------------------------------
// Basic tests.

// BenchmarkDisabled logs at the wrong level and nothing comes out.
// The default logger is set to INFO, a DEBUG log message should be ignored.
func (suite *SlogBenchmarkSuite) BenchmarkDisabled() *Benchmark {
	return &Benchmark{
		Options: infra.SimpleOptions(),
		BenchmarkFn: func(logger *slog.Logger) {
			logger.Debug(message)
		},
		VerifyFn: func(captured []byte, logMap map[string]any, manager *warning2.Manager) error {
			if len(captured) > 0 {
				manager.AddWarning(warning2.NotDisabled, "Disabled", string(captured))
				return warning2.NotDisabled
			}
			return nil
		},
		DontCount: true,
	}
}

// BenchmarkSimple logs a simple line with just a message.
func (suite *SlogBenchmarkSuite) BenchmarkSimple() *Benchmark {
	return &Benchmark{
		Options: infra.SimpleOptions(),
		BenchmarkFn: func(logger *slog.Logger) {
			logger.Info(message)
		},
		VerifyFn: matcher("Simple", expectedBasic()),
	}
}

// BenchmarkSimpleSource logs a simple line with just a message
// and the AddSource option set so the slog.SourceKey group is created.
func (suite *SlogBenchmarkSuite) BenchmarkSimpleSource() *Benchmark {
	return &Benchmark{
		Options: infra.SourceOptions(),
		BenchmarkFn: func(logger *slog.Logger) {
			logger.Info(message)
		},
		VerifyFn: verify(
			finder("SimpleSource", expectedBasic()),
			sourcerer("SimpleSource")),
	}
}

// -----------------------------------------------------------------------------

// BenchmarkAttributes logs a message with a lot of attributes.
func (suite *SlogBenchmarkSuite) BenchmarkAttributes() *Benchmark {
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

// BenchmarkKeyValues logs a message with a lot of key value pairs.
func (suite *SlogBenchmarkSuite) BenchmarkKeyValues() *Benchmark {
	return &Benchmark{
		Options: infra.SimpleOptions(),
		BenchmarkFn: func(logger *slog.Logger) {
			logger.Info(message, allKeyValues...)
		},
		VerifyFn: verify(
			finder("KeyValues", expectedBasic()),
			finder("KeyValues", allValuesMap()),
			noDuplicates("KeyValues"),
		),
	}
}

// -----------------------------------------------------------------------------

// BenchmarkWithAttrsSimple logs a simple message to a logger created from
// a handler configure with a lot of attributes.
func (suite *SlogBenchmarkSuite) BenchmarkWithAttrsSimple() *Benchmark {
	return &Benchmark{
		Options: infra.SimpleOptions(),
		BenchmarkFn: func(logger *slog.Logger) {
			logger.Info(message)
		},
		HandlerFn: withAllAttributes,
		VerifyFn: verify(
			finder("WithAttrsSimple:Basic", expectedBasic()),
			finder("WithAttrsSimple:With", withValuesMap()),
			noDuplicates("WithAttrsSimple"),
		),
	}
}

// BenchmarkWithAttrsAttributes logs a message and a lot of attributes to a logger created from
// a handler configure with a different lot of attributes.
func (suite *SlogBenchmarkSuite) BenchmarkWithAttrsAttributes() *Benchmark {
	return &Benchmark{
		Options: infra.SimpleOptions(),
		BenchmarkFn: func(logger *slog.Logger) {
			logger.LogAttrs(context.Background(), slog.LevelInfo, message, allAttributes...)
		},
		HandlerFn: withAllAttributes,
		VerifyFn: verify(
			finder("WithAttrsAttributes:Basic", expectedBasic()),
			finder("WithAttrsAttributes:All", allValuesMap()),
			finder("WithAttrsAttributes:With", withValuesMap()),
			noDuplicates("WithAttrsAttributes"),
		),
	}
}

// BenchmarkWithAttrsKeyValues logs a message and a lot of attributes to a logger created from
// a handler configure with a different lot of key value pairs.
func (suite *SlogBenchmarkSuite) BenchmarkWithAttrsKeyValues() *Benchmark {
	return &Benchmark{
		Options: infra.SimpleOptions(),
		BenchmarkFn: func(logger *slog.Logger) {
			logger.Info(message, allKeyValues...)
		},
		HandlerFn: withAllAttributes,
		VerifyFn: verify(
			finder("WithAttrsKeyValues:Basic", expectedBasic()),
			finder("WithAttrsKeyValues:All", allValuesMap()),
			finder("WithAttrsKeyValues:With", withValuesMap()),
			noDuplicates("WithAttrsKeyValues"),
		),
	}
}

// -----------------------------------------------------------------------------

// BenchmarkWithGroupAttributes logs a message and a lot of attributes to
// a logger created from a handler configured with an open group.
func (suite *SlogBenchmarkSuite) BenchmarkWithGroupAttributes() *Benchmark {
	return &Benchmark{
		Options: infra.SimpleOptions(),
		BenchmarkFn: func(logger *slog.Logger) {
			logger.LogAttrs(context.Background(), slog.LevelInfo, message, allAttributes...)
		},
		HandlerFn: withGroup,
		VerifyFn: verify(
			finder("WithGroupAttributes:Basic", expectedBasic()),
			finder("WithGroupAttributes:All", map[string]any{
				"withGroup": allValuesMap(),
			}),
			noDuplicates("With_Group_Attributes"),
		),
	}
}

// BenchmarkWithGroupKeyValues logs a message and a lot of key value pairs to
// a logger created from a handler configured with an open group.
func (suite *SlogBenchmarkSuite) BenchmarkWithGroupKeyValues() *Benchmark {
	return &Benchmark{
		Options: infra.SimpleOptions(),
		BenchmarkFn: func(logger *slog.Logger) {
			logger.Info(message, allKeyValues...)
		},
		HandlerFn: withGroup,
		VerifyFn: verify(
			finder("WithGroupKeyValues:Basic", expectedBasic()),
			finder("WithGroupKeyValues:All", map[string]any{
				"withGroup": allValuesMap(),
			}),
			noDuplicates("WithGroupKeyValues"),
		),
	}
}

// -----------------------------------------------------------------------------
// Large/Long tests.

// BenchmarkBigGroup logs several levels of nested groups.
func (suite *SlogBenchmarkSuite) BenchmarkBigGroup() *Benchmark {
	return &Benchmark{
		Options: infra.SimpleOptions(),
		BenchmarkFn: func(logger *slog.Logger) {
			logger.Info(message, BigGroup())
		},
		VerifyFn: bigGroupChecker("BigGroup"),
	}
}

// BenchmarkLogging logs a series of lines taken from Gin server output.
func (suite *SlogBenchmarkSuite) BenchmarkLogging() *Benchmark {
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
func getLogMap(captured []byte, logMap map[string]any, manager *warning2.Manager) map[string]any {
	var err error
	if logMap == nil {
		if logMap, err = parseLogMap(captured); err != nil {
			manager.AddWarning(warning2.TestError, err.Error(), string(captured))
		}
		fixLogMap(logMap)
	}
	return logMap
}
