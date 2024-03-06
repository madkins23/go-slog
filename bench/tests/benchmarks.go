package tests

import (
	"context"
	"log/slog"

	"github.com/madkins23/go-slog/infra"
	"github.com/madkins23/go-slog/internal/warning"
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
type VerifyFn func(captured []byte, logMap map[string]any, manager *warning.Manager) error

type Benchmark struct {
	Options     *slog.HandlerOptions
	BenchmarkFn BenchmarkFn
	HandlerFn   HandlerFn
	VerifyFn    VerifyFn
	DontCount   bool
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
		VerifyFn: func(captured []byte, logMap map[string]any, manager *warning.Manager) error {
			if len(captured) > 0 {
				manager.AddWarning(warning.NotDisabled, "Disabled", string(captured))
				return warning.NotDisabled
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
func getLogMap(captured []byte, logMap map[string]any, manager *warning.Manager) map[string]any {
	var err error
	if logMap == nil {
		if logMap, err = parseLogMap(captured); err != nil {
			manager.AddWarning(warning.TestError, err.Error(), string(captured))
		}
		fixLogMap(logMap)
	}
	return logMap
}
