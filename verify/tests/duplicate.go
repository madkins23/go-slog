package tests

import "github.com/madkins23/go-slog/infra"

// -----------------------------------------------------------------------------
// Duplicate testing, which isn't currently regarded as an error.
// This issue is under discussion in https://github.com/golang/go/issues/59365.

// TestAttributeDuplicate tests whether duplicate attributes are logged properly.
//   - Based on the existing behavior of log/slog the second occurrence overrides the first.
//   - See https://github.com/golang/go/issues/59365
func (suite *SlogTestSuite) TestAttributeDuplicate() {
	logger := suite.Logger(infra.SimpleOptions())
	logger.Info(message,
		"alpha", "one", "alpha", 2, "bravo", "hurrah",
		"charlie", "brown", "charlie", 3, "charlie", 23.79)
	logMap := suite.logMap()
	suite.checkFieldCount(6, logMap)
}

// TestAttributeWithDuplicate tests whether duplicate attributes are logged properly
// when the duplicate is introduced by With() and then the main call.
//   - Based on the existing behavior of log/slog the second occurrence overrides the first.
//   - See https://github.com/golang/go/issues/59365
func (suite *SlogTestSuite) TestAttributeWithDuplicate() {
	logger := suite.Logger(infra.SimpleOptions())
	logger.
		With("alpha", "one", "bravo", "hurrah", "charlie", "brown", "charlie", "jones").
		Info(message, "alpha", 2, "charlie", 23.70)
	logMap := suite.logMap()
	suite.checkFieldCount(6, logMap)
}
