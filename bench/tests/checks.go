package tests

import (
	"bufio"
	"bytes"
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

// verifyLines applies the specified function(s) to each line in the captured bytes.
func verifyLines(fns ...VerifyFn) VerifyFn {
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
			if !result {
				break
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

func sorcerer(testName string) VerifyFn {
	return func(captured []byte, logMap map[string]any, manager *infra.WarningManager) bool {
		result := false
		text := testName
		logMap = getLogMap(captured, logMap, manager)
		if srcVal, found := logMap[slog.SourceKey]; found {
			if srcMap, ok := srcVal.(map[string]any); ok {
				missing := make([]string, 0)
				for _, field := range []string{"file", "function", "line"} {
					if _, found = srcMap[field]; !found {
						missing = append(missing, field)
					}
				}
				if len(missing) > 0 {
					text += ": " + strings.Join(missing, ",")
				} else {
					result = true
				}
			}
		}
		if !result {
			manager.AddWarning(warning.SourceKey, text, string(captured))
		}
		return result
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
