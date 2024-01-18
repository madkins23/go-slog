package tests

import (
	"encoding/json"
	"fmt"

	"github.com/madkins23/go-slog/internal/test"
	testJSON "github.com/madkins23/go-slog/json"
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
