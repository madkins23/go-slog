package tests

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"

	testJSON "github.com/madkins23/go-slog/json"
	"github.com/madkins23/go-slog/test"
)

// -----------------------------------------------------------------------------
// Utility methods.

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
	test.Debugf(1, ">>> JSON: %s", suite.Buffer.Bytes())
	var results map[string]any
	suite.Require().NoError(json.Unmarshal(suite.Buffer.Bytes(), &results))
	return results
}

// showLog sends the output capture buffer to test logging output.
// This is not dependent on the use of the -debug flag.
func (suite *SlogTestSuite) showLog() {
	fmt.Printf(">>> %s", suite.Buffer)
}

// -----------------------------------------------------------------------------
// Utility functions.

// currentFunctionName checks up the call stack for the name of the current test function.
// Only the last part of the function name (after the last period) is returned.
// The function name is found by checking for a prefix of "Test".
// If no test function is found "Unknown" is returned.
func currentFunctionName() string {
	pc := make([]uintptr, 10)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	more := true
	for more {
		var frame runtime.Frame
		frame, more = frames.Next()
		parts := strings.Split(frame.Function, ".")
		functionName := parts[len(parts)-1]
		if strings.HasPrefix(functionName, "Test") {
			return functionName
		}
	}
	return "Unknown"
}
