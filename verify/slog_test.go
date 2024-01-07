package verify

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// Test_slog runs tests for the log/slog JSON handler.
// No overrides are required.
func Test_slog(t *testing.T) {
	suite.Run(t, &SlogTestSuite{})
}
