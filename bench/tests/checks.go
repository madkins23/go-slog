package tests

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"reflect"
	"sort"
	"strings"

	intJSON "github.com/madkins23/go-slog/internal/json"
	"github.com/madkins23/go-slog/internal/test"
	"github.com/madkins23/go-slog/internal/warning"
)

// -----------------------------------------------------------------------------

// bigGroupChecker checks the captured/logMap data to see if it is a big group.
func bigGroupChecker(testName string) VerifyFn {
	return func(captured []byte, logMap map[string]any, manager *warning.Manager) error {
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
	return func(captured []byte, logMap map[string]any, manager *warning.Manager) error {
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
	return func(captured []byte, logMap map[string]any, manager *warning.Manager) error {
		logMap = getLogMap(captured, logMap, manager)
		badFields, strFields, err := finderDeep(expected, logMap, "")
		if err != nil {
			return fmt.Errorf("finderDeep: %w", err)
		}
		if len(badFields) > 0 {
			text := testName + ": " + strings.Join(badFields, ",")
			manager.AddWarningFn(warning.Mismatch, text, string(captured))
		} else if manager.HasWarning(warning.Mismatch) {
			manager.AddUnused(warning.Mismatch, string(captured))
		}
		if len(strFields) > 0 {
			text := testName + ": " + strings.Join(strFields, ",")
			manager.AddWarningFn(warning.StringAny, text, string(captured))
		} else if manager.HasWarning(warning.StringAny) {
			manager.AddUnused(warning.StringAny, string(captured))
		}
		return nil
	}
}

// finderDeep matches the parts of the actual map against what is expected.
// The actual map can have other unspecified fields.
// This function is the recursive workhorse for finder.
func finderDeep(expected map[string]any, actual map[string]any, prefix string) ([]string, []string, error) {
	badFields := make([]string, 0)
	strFields := make([]string, 0)
	for field, value := range expected {
		expMap, expOk := value.(map[string]any)
		actMap, actOk := actual[field].(map[string]any)
		if expOk && actOk {
			if bf, sf, err := finderDeep(expMap, actMap, prefix+field+"."); err != nil {
				return nil, nil, err
			} else {
				badFields = append(badFields, bf...)
				strFields = append(strFields, sf...)
			}
		} else if !reflect.DeepEqual(value, actual[field]) {
			if actVal, ok := actual[field].(string); ok {
				if valJSON, err := json.Marshal(value); err != nil {
					return nil, nil, fmt.Errorf("marshal JSON: %w", err)
				} else if actVal == string(valJSON) {
					strFields = append(strFields, field)
					continue
				} else {
					test.Debugf(3, ">>> %s: %s =?= %s (%s)\n", field, value, actual[field], string(valJSON))
				}
			}
			badFields = append(badFields, field)
		}
	}
	sort.Strings(badFields)
	sort.Strings(strFields)
	return badFields, strFields, nil
}

// matcher matches the parts entirety of the actual map against
// the entirety of the expected map using reflect.DeepEqual.
func matcher(testName string, expected map[string]any) VerifyFn {
	return func(captured []byte, logMap map[string]any, manager *warning.Manager) error {
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
	return func(captured []byte, logMap map[string]any, manager *warning.Manager) error {
		logMap = getLogMap(captured, logMap, manager)
		counter := intJSON.NewFieldCounter(captured)
		if len(counter.Duplicates()) > 0 {
			manager.AddWarningFn(warning.Duplicates, testName, string(captured))
			return warning.Duplicates.ErrorExtra(testName)
		}
		return nil
	}
}

// sourcerer checks to see if the slog.SourceKey is present and properly configured.
func sourcerer(testName string) VerifyFn {
	return func(captured []byte, logMap map[string]any, manager *warning.Manager) error {
		logMap = getLogMap(captured, logMap, manager)
		if srcVal, found := logMap[slog.SourceKey]; found {
			if srcMap, ok := srcVal.(map[string]any); !ok {
				text := "source not map"
				manager.AddWarningFnText(warning.SourceKey, testName, text, string(captured))
				return warning.SourceKey.ErrorExtra(text)
			} else {
				missing := make([]string, 0)
				for _, field := range []string{"file", "function", "line"} {
					if _, found = srcMap[field]; !found {
						missing = append(missing, field)
					}
				}
				if len(missing) > 0 {
					text := fmt.Sprintf("missing fields: %s", strings.Join(missing, ","))
					manager.AddWarningFnText(warning.SourceKey, testName, text, string(captured))
					return warning.SourceKey.ErrorExtra(text)
				}
			}
		} else if callVal, found := logMap["caller"]; found {
			var text string
			if caller, ok := callVal.(string); !ok {
				text = "caller not string"
			} else if parts := strings.Split(caller, ":"); len(parts) != 2 {
				text = "caller not string"
			} else {
				text = caller
			}
			manager.AddWarningFnText(warning.SourceCaller, testName, text, string(captured))
			return warning.SourceCaller.ErrorExtra(text)
		} else {
			text := fmt.Sprintf("no '%s' key", slog.SourceKey)
			manager.AddWarningFnText(warning.SourceKey, testName, text, string(captured))
			return warning.SourceKey.ErrorExtra(text)
		}
		return nil
	}
}

// verify runs the specified functions against the captured/logMap data.
func verify(fns ...VerifyFn) VerifyFn {
	return func(captured []byte, logMap map[string]any, manager *warning.Manager) error {
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
	return func(captured []byte, logMap map[string]any, manager *warning.Manager) error {
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
