package tests

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"reflect"
	"sort"
	"strings"

	"github.com/madkins23/go-slog/infra"
	"github.com/madkins23/go-slog/internal/test"
	"github.com/madkins23/go-slog/json"
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
		matcher("Simple_Source", expectedSource()))
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
		liner("Logging",
			fields("Logging", "level", "msg", "code", "duration", "method", "sys", "url")))
	bm.SetDontCount(true)
	return bm
}

// -----------------------------------------------------------------------------

func bigGroupChecker(testName string) VerifyFn {
	return func(captured []byte, logMap map[string]any, manager *infra.WarningManager) bool {
		logMap = getLogMap(captured, logMap, manager)
		if groupMap, found := logMap[valGroupName]; found {
			if group, ok := groupMap.(map[string]any); !ok {
				manager.AddWarning(warning.Mismatch, testName+": not a group", "")
			} else if maxDepth, err := bigGroupCheck(group, 0, bigGroupLimit, valGroupName); err != nil {
				manager.AddWarning(warning.Mismatch, testName+": "+err.Error(), "")
			} else if maxDepth != bigGroupLimit {
				manager.AddWarning(warning.Mismatch,
					fmt.Sprintf("%s: maxDepth %d != %d limit", testName, maxDepth, bigGroupLimit),
					"")
			} else {
				return true
			}
		}
		return false
	}
}

// fields checks the captured/logMap data to see if the specified fields exist.
func fields(testName string, fields ...string) VerifyFn {
	return func(captured []byte, logMap map[string]any, manager *infra.WarningManager) bool {
		logMap = getLogMap(captured, logMap, manager)
		missing := make([]string, 0, len(fields))
		for _, field := range fields {
			if _, found := logMap[field]; !found {
				missing = append(missing, field)
			}
		}
		if len(missing) > 0 {
			manager.AddWarning(warning.Mismatch, testName+": "+strings.Join(missing, ","), string(captured))
			return false
		}
		return true
	}
}

// finder matches the parts of the actual map against what is expected.
// The actual map can have other unspecified fields.
func finder(testName string, expected map[string]any) VerifyFn {
	return func(captured []byte, logMap map[string]any, manager *infra.WarningManager) bool {
		logMap = getLogMap(captured, logMap, manager)
		badFields := finderDeep(expected, logMap, "")
		if len(badFields) > 0 {
			for _, field := range badFields {
				test.Debugf(2, ">?>   %s: %v != %v\n", field, expected[field], logMap[field])
			}
			manager.AddWarningFn(warning.Mismatch,
				testName+": "+strings.Join(badFields, ","),
				string(captured))
			return false
		}
		return true
	}
}

func finderDeep(expected map[string]any, actual map[string]any, prefix string) []string {
	badFields := make([]string, 0)
	for field, value := range expected {
		expMap, expOk := value.(map[string]any)
		actMap, actOk := actual[field].(map[string]any)
		if expOk && actOk {
			badFields = append(badFields, finderDeep(expMap, actMap, prefix+field+".")...)
		} else if !reflect.DeepEqual(value, actual[field]) {
			badFields = append(badFields, field)
		}
	}
	sort.Strings(badFields)
	return badFields
}

// liner applies the specified function(s) to each line in the captured bytes.
func liner(testName string, fns ...VerifyFn) VerifyFn {
	return func(captured []byte, logMap map[string]any, manager *infra.WarningManager) bool {
		var buffer bytes.Buffer
		result := true
		lineReader := bufio.NewReader(bytes.NewReader(captured))
		for {
			chunk, isPrefix, err := lineReader.ReadLine()
			if err != nil {
				if !errors.Is(err, io.EOF) {
					slog.Error("Read line", "err", err)
					result = false
				}
				break
			}
			buffer.Write(chunk)
			if isPrefix {
				continue
			}
			logMap = getLogMap(buffer.Bytes(), nil, manager)
			for _, fn := range fns {
				if !fn(buffer.Bytes(), logMap, manager) {
					result = false
				}
			}
			buffer.Reset()
		}
		return result
	}
}

// matcher matches the parts entirety of the actual map against
// the entirety of the expected map using reflect.DeepEqual.
func matcher(testName string, expected map[string]any) VerifyFn {
	return func(captured []byte, logMap map[string]any, manager *infra.WarningManager) bool {
		logMap = getLogMap(captured, logMap, manager)
		if !reflect.DeepEqual(expected, fixLogMap(logMap)) {
			test.Debugf(2, ">?> %v\n", expected)
			test.Debugf(2, ">=> %v\n", fixLogMap(logMap))
			manager.AddWarningFn(warning.Mismatch, testName, string(captured))
			return false
		}
		return true
	}
}

// noDuplicates checks to see if there are duplicate fields in the logMap.
func noDuplicates(testName string) VerifyFn {
	return func(captured []byte, logMap map[string]any, manager *infra.WarningManager) bool {
		logMap = getLogMap(captured, logMap, manager)
		counter := json.NewFieldCounter(captured)
		if len(counter.Duplicates()) > 0 {
			manager.AddWarning(warning.Duplicates, testName, string(captured))
			return false
		}
		return true
	}
}

// verify runs the specified functions against the captured/logMap data.
func verify(fns ...VerifyFn) VerifyFn {
	return func(captured []byte, logMap map[string]any, manager *infra.WarningManager) bool {
		logMap = getLogMap(captured, logMap, manager)
		result := true
		for _, fn := range fns {
			if !fn(captured, logMap, manager) {
				result = false
			}
		}
		return result
	}
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
