package verify

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/madkins23/go-slog/creator/slogjson"
	"github.com/madkins23/go-slog/verify/tests"
	"github.com/madkins23/go-slog/warning"
)

// Test_slog runs tests for the log/slog JSON handler.
func Test_slog(t *testing.T) {
	slogSuite := tests.NewSlogTestSuite(slogjson.Creator())
	slogSuite.WarnOnly(warning.Duplicates)
	suite.Run(t, slogSuite)
}
