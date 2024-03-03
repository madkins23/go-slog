package tests

import (
	"encoding/json"
	"fmt"
	"log/slog"

	testJSON "github.com/madkins23/go-slog/internal/json"
)

// -----------------------------------------------------------------------------
// Utility methods.

func (suite *SlogTestSuite) adjustExpected(expected, logMap map[string]any) {
	if t, found := logMap[slog.TimeKey]; found {
		expected["time"] = t
	}
	if l, found := logMap[slog.LevelKey]; found {
		expected["level"] = l
	}
	if _, found := logMap["message"]; found {
		expected["message"] = expected[slog.MessageKey]
		delete(expected, slog.MessageKey)
	}
}

// bufferReset clears the test suite's output capture buffer.
// This allows multiple log statements to be generated and evaluated in the same test.
// Otherwise, the log statements keep collecting in the buffer.
func (suite *SlogTestSuite) bufferReset() {
	suite.Buffer.Reset()
}

// fieldCounter returns a json.FieldCounter object for the output capture buffer.
func (suite *SlogTestSuite) fieldCounter() *testJSON.FieldCounter {
	return testJSON.NewFieldCounter(suite.Buffer.Bytes())
}

// logMap unmarshals JSON in the output capture buffer into a map[string]any.
// The buffer is sent to test logging output if the -debug=<level> flag is >= 1.
func (suite *SlogTestSuite) logMap() map[string]any {
	var results map[string]any
	err := json.Unmarshal(suite.Buffer.Bytes(), &results)
	if err != nil {
		err = fmt.Errorf("%w: '%s'", err, suite.Buffer.Bytes())
	}
	suite.Require().NoError(err)
	return results
}
