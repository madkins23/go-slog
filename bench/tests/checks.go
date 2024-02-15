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
	return func(captured []byte, logMap map[string]any, manager *infra.WarningManager) error {
		logMap = getLogMap(captured, logMap, manager)
		if groupMap, found := logMap[valGroupName]; found {
			if group, ok := groupMap.(map[string]any); !ok {
				text := testName + ": not a group"
				manager.AddWarning(warning.Mismatch, text, "")
				return warning.Mismatch.ErrorExtra(text)
			} else if maxDepth, err := bigGroupCheck(group, 0, bigGroupLimit, valGroupName); err != nil {
				text := testName + ": " + err.Error()
				manager.AddWarning(warning.Mismatch, text, "")
				return warning.Mismatch.ErrorExtra(text)
			} else if maxDepth != bigGroupLimit {
				text := fmt.Sprintf("%s: maxDepth %d != %d limit", testName, maxDepth, bigGroupLimit)
				manager.AddWarning(warning.Mismatch, text, "")
				return warning.Mismatch.ErrorExtra(text)
			}
		}
		return nil
	}
}

// fields checks the captured/logMap data to see if the specified fields exist.
func fields(testName string, fields ...string) VerifyFn {
	return func(captured []byte, logMap map[string]any, manager *infra.WarningManager) error {
		logMap = getLogMap(captured, logMap, manager)
		missing := make([]string, 0, len(fields))
		for _, field := range fields {
			if _, found := logMap[field]; !found {
				missing = append(missing, field)
			}
		}
		if len(missing) > 0 {
			text := testName + ": " + strings.Join(missing, ",")
			manager.AddWarning(warning.Mismatch, text, string(captured))
			return warning.Mismatch.ErrorExtra(text)
		}
		return nil
	}
}

// finder matches the parts of the actual map against what is expected.
// The actual map can have other unspecified fields.
func finder(testName string, expected map[string]any) VerifyFn {
	return func(captured []byte, logMap map[string]any, manager *infra.WarningManager) error {
		logMap = getLogMap(captured, logMap, manager)
		badFields := finderDeep(expected, logMap, "")
		if len(badFields) > 0 {
			for _, field := range badFields {
				test.Debugf(2, ">?>   %s: %v != %v\n", field, expected[field], logMap[field])
			}
			text := testName + ": " + strings.Join(badFields, ",")
			manager.AddWarningFn(warning.Mismatch, text, string(captured))
			return warning.Mismatch.ErrorExtra(text)
		}
		return nil
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

// matcher matches the parts entirety of the actual map against
// the entirety of the expected map using reflect.DeepEqual.
func matcher(testName string, expected map[string]any) VerifyFn {
	return func(captured []byte, logMap map[string]any, manager *infra.WarningManager) error {
		logMap = getLogMap(captured, logMap, manager)
		if !reflect.DeepEqual(expected, fixLogMap(logMap)) {
			test.Debugf(2, ">?> %v\n", expected)
			test.Debugf(2, ">=> %v\n", fixLogMap(logMap))
			manager.AddWarningFn(warning.Mismatch, testName, string(captured))
			return warning.Mismatch.ErrorExtra(testName)
		}
		return nil
	}
}

// noDuplicates checks to see if there are duplicate fields in the logMap.
func noDuplicates(testName string) VerifyFn {
	return func(captured []byte, logMap map[string]any, manager *infra.WarningManager) error {
		logMap = getLogMap(captured, logMap, manager)
		counter := json.NewFieldCounter(captured)
		if len(counter.Duplicates()) > 0 {
			manager.AddWarning(warning.Duplicates, testName, string(captured))
			return warning.Duplicates.ErrorExtra(testName)
		}
		return nil
	}
}

func sorcerer(testName string) VerifyFn {
	return func(captured []byte, logMap map[string]any, manager *infra.WarningManager) error {
		logMap = getLogMap(captured, logMap, manager)
		if srcVal, found := logMap[slog.SourceKey]; !found {
			text := fmt.Sprintf("%s: no %s key", testName, slog.SourceKey)
			manager.AddWarning(warning.SourceKey, text, string(captured))
			return warning.SourceKey.ErrorExtra(text)
		} else if srcMap, ok := srcVal.(map[string]any); !ok {
			text := fmt.Sprintf("%s: source not map", testName)
			manager.AddWarning(warning.SourceKey, text, string(captured))
			return warning.SourceKey.ErrorExtra(text)
		} else {
			missing := make([]string, 0)
			for _, field := range []string{"file", "function", "line"} {
				if _, found = srcMap[field]; !found {
					missing = append(missing, field)
				}
			}
			if len(missing) > 0 {
				text := fmt.Sprintf("%s: missing fields: %s", testName, strings.Join(missing, ","))
				manager.AddWarning(warning.SourceKey, text, string(captured))
				return warning.SourceKey.ErrorExtra(text)
			}
		}
		return nil
	}
}

// verify runs the specified functions against the captured/logMap data.
func verify(fns ...VerifyFn) VerifyFn {
	return func(captured []byte, logMap map[string]any, manager *infra.WarningManager) error {
		logMap = getLogMap(captured, logMap, manager)
		errs := make([]error, 0, len(fns))
		for _, fn := range fns {
			if err := fn(captured, logMap, manager); err != nil {
				errs = append(errs, err)
			}
		}
		if len(errs) > 0 {
			return errors.Join(errs...)
		}
		return nil
	}
}

// verifyLines applies the specified function(s) to each line in the captured bytes.
func verifyLines(fns ...VerifyFn) VerifyFn {
	return func(captured []byte, logMap map[string]any, manager *infra.WarningManager) error {
		var buffer bytes.Buffer
		lineReader := bufio.NewReader(bytes.NewReader(captured))
		for {
			chunk, isPrefix, err := lineReader.ReadLine()
			if err != nil {
				if !errors.Is(err, io.EOF) {
					return fmt.Errorf("read line: %w", err)
				}
				break
			}
			buffer.Write(chunk)
			if isPrefix {
				continue
			}
			logMap = getLogMap(buffer.Bytes(), nil, manager)
			errs := make([]error, 0, len(fns))
			for _, fn := range fns {
				if err := fn(buffer.Bytes(), logMap, manager); err != nil {
					errs = append(errs, err)
				}
			}
			if len(errs) > 0 {
				return errors.Join(errs...)
			}
			buffer.Reset()
		}
		return nil
	}
}