package verify

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/madkins23/go-slog/creator"
	"github.com/madkins23/go-slog/verify/tests"
)

// Test_slog_zerolog_phsym runs tests for the physym zerolog handler.
func Test_slog_zerolog_phsym(t *testing.T) {
	slogSuite := tests.NewSlogTestSuite(creator.SlogPhsymZerolog())
	slogSuite.WarnOnly(tests.WarnDuplicates)
	slogSuite.WarnOnly(tests.WarnDurationMillis)
	slogSuite.WarnOnly(tests.WarnEmptyAttributes)
	slogSuite.WarnOnly(tests.WarnGroupEmpty)
	slogSuite.WarnOnly(tests.WarnGroupInline)
	slogSuite.WarnOnly(tests.WarnLevelCase)
	slogSuite.WarnOnly(tests.WarnMessageKey)
	slogSuite.WarnOnly(tests.WarnTimeMillis)
	slogSuite.WarnOnly(tests.WarnNoReplAttr)
	slogSuite.WarnOnly(tests.WarnSourceKey)
	slogSuite.WarnOnly(tests.WarnZeroTime)
	suite.Run(t, slogSuite)
}
