package tests

import (
	"bufio"
	"bytes"
	"context"
	_ "embed"
	"errors"
	"io"
	"log/slog"
	"regexp"
	"strconv"
	"strings"

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

//goland:noinspection GoSnakeCaseUsage
func (suite *SlogBenchmarkSuite) Benchmark_Disabled() Benchmark {
	test := NewBenchmark(infra.SimpleOptions(), func(logger *slog.Logger) {
		logger.Info(message)
	}, nil)
	test.SetDontCount(true)
	return test
}

//goland:noinspection GoSnakeCaseUsage
func (suite *SlogBenchmarkSuite) Benchmark_Simple() Benchmark {
	return NewBenchmark(infra.SimpleOptions(), func(logger *slog.Logger) {
		logger.Info(message)
	}, nil)
}

//goland:noinspection GoSnakeCaseUsage
func (suite *SlogBenchmarkSuite) Benchmark_Simple_Source() Benchmark {
	return NewBenchmark(infra.SourceOptions(), func(logger *slog.Logger) {
		logger.Info(message)
	}, nil)
}

// -----------------------------------------------------------------------------

//goland:noinspection GoSnakeCaseUsage
func (suite *SlogBenchmarkSuite) Benchmark_Attributes() Benchmark {
	return NewBenchmark(infra.SimpleOptions(), func(logger *slog.Logger) {
		logger.LogAttrs(context.Background(), slog.LevelInfo, message, allAttributes...)
	}, nil)
}

//goland:noinspection GoSnakeCaseUsage
func (suite *SlogBenchmarkSuite) Benchmark_Key_Values() Benchmark {
	return NewBenchmark(infra.SimpleOptions(), func(logger *slog.Logger) {
		logger.Info(message, allKeyValues...)
	}, nil)
}

// -----------------------------------------------------------------------------

//goland:noinspection GoSnakeCaseUsage
func (suite *SlogBenchmarkSuite) Benchmark_With_Attrs_Simple() Benchmark {
	return NewBenchmark(infra.SimpleOptions(), func(logger *slog.Logger) {
		logger.Info(message)
	}, withAllAttributes)
}

//goland:noinspection GoSnakeCaseUsage
func (suite *SlogBenchmarkSuite) Benchmark_With_Attrs_Attributes() Benchmark {
	return NewBenchmark(infra.SimpleOptions(), func(logger *slog.Logger) {
		logger.LogAttrs(context.Background(), slog.LevelInfo, message, allAttributes...)
	}, withAllAttributes)
}

//goland:noinspection GoSnakeCaseUsage
func (suite *SlogBenchmarkSuite) Benchmark_With_Attrs_Key_Values() Benchmark {
	return NewBenchmark(infra.SimpleOptions(), func(logger *slog.Logger) {
		logger.Info(message, allKeyValues...)
	}, withAllAttributes)
}

// -----------------------------------------------------------------------------

//goland:noinspection GoSnakeCaseUsage
func (suite *SlogBenchmarkSuite) Benchmark_With_Group_Attributes() Benchmark {
	return NewBenchmark(infra.SimpleOptions(), func(logger *slog.Logger) {
		logger.LogAttrs(context.Background(), slog.LevelInfo, message, allAttributes...)
	}, withGroupAttributes)
}

//goland:noinspection GoSnakeCaseUsage
func (suite *SlogBenchmarkSuite) Benchmark_With_Group_Key_Values() Benchmark {
	return NewBenchmark(infra.SimpleOptions(), func(logger *slog.Logger) {
		logger.Info(message, allKeyValues...)
	}, withGroupAttributes)
}

// -----------------------------------------------------------------------------
// Large/Long tests.

//goland:noinspection GoSnakeCaseUsage
func (suite *SlogBenchmarkSuite) Benchmark_Big_Group() Benchmark {
	return NewBenchmark(infra.SimpleOptions(), func(logger *slog.Logger) {
		logger.Info(message, BigGroup())
	}, nil)
}

//goland:noinspection GoSnakeCaseUsage
func (suite *SlogBenchmarkSuite) Benchmark_Logging() Benchmark {
	test := NewBenchmark(infra.SimpleOptions(), func(logger *slog.Logger) {
		for _, logData := range getLogData() {
			logger.Info("Handle", logData...)
		}
	}, nil)
	test.SetDontCount(true)
	return test
}

// -----------------------------------------------------------------------------
// Acquire a bunch of log statements to use in Benchmark_Logging.

//go:embed logging.txt
var logging []byte

var logData [][]any

var (
	ptnCode  = regexp.MustCompile(`\s(\d+)\s*$`)
	ptnSplit = regexp.MustCompile(`\s+`)
)

// Read log data from embedded data, construct array of arg arrays for logging.
func getLogData() [][]any {
	if logData == nil {
		reader := bufio.NewReader(bytes.NewReader(logging))
		var line bytes.Buffer
		for {
			if chunk, isPrefix, err := reader.ReadLine(); errors.Is(err, io.EOF) {
				break
			} else if err != nil {
				slog.Warn("Error reading logging data line", "err", err)
			} else {
				line.Write(chunk)
				if isPrefix {
					continue
				}
			}
			// Example line:
			//  07:55:52 INF 200 |    9.522199ms |             ::1 | GET      "/chart.svg?tag=samber_zap&item=MemAllocs" sys=gin
			if parts := strings.Split(string(line.Bytes()), "|"); len(parts) != 4 {
				slog.Warn("Wrong number of parts", "num", len(parts), "line", line, "func", "getLogData")
			} else {
				var args []any
				// Parse parts[0]:
				//  07:55:52 INF 200
				if matches := ptnCode.FindStringSubmatch(parts[0]); len(matches) != 2 {
					slog.Warn("Unable to parse code", "from", parts[0])
				} else if num, err := strconv.ParseInt(matches[1], 10, 64); err != nil {
					slog.Warn("Unable to parse int", "from", parts[0], "func", "getLogData")
				} else {
					args = append(args, "code", num)
				}
				args = append(args, "duration", strings.Trim(parts[1], " "))
				// Ignore parts[2] (::1) since I don't know what it is.
				// Parse parts[3]:
				//  GET      "/chart.svg?tag=samber_zap&item=MemAllocs" sys=gin
				parts = ptnSplit.Split(strings.Trim(parts[3], " "), -1)
				if len(parts) == 3 {
					args = append(args, "method", parts[0])
					args = append(args, "url", strings.Trim(parts[1], "\""))
				}
				args = append(args, "sys", "gin")
				logData = append(logData, args)
			}
			line.Reset()
		}
	}
	return logData
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
