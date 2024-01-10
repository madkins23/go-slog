package test

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

func (suite *SlogTestSuite) bufferReset() {
	suite.Buffer.Reset()
}

func (suite *SlogTestSuite) fieldCounter() *testJSON.FieldCounter {
	return testJSON.NewFieldCounter(suite.Buffer.Bytes())
}

func (suite *SlogTestSuite) logMap() map[string]any {
	test.Debugf(1, ">>> JSON: %s", suite.Buffer.Bytes())
	var results map[string]any
	suite.Require().NoError(json.Unmarshal(suite.Buffer.Bytes(), &results))
	return results
}

func (suite *SlogTestSuite) showLog() {
	fmt.Printf(">>> %s", suite.Buffer)
}

// -----------------------------------------------------------------------------
// Utility functions.

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
