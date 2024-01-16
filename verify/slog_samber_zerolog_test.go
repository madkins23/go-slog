package verify

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/madkins23/go-slog/creator"
	"github.com/madkins23/go-slog/verify/tests"
)

// Test_slog_samber_zerolog runs tests for the samber zerolog handler.
func Test_slog_samber_zerolog(t *testing.T) {
	sLogSuite := tests.NewSlogTestSuite(creator.SlogSamberZerolog())
	sLogSuite.WarnOnly(tests.WarnDefaultLevel)
	sLogSuite.WarnOnly(tests.WarnDurationMillis)
	sLogSuite.WarnOnly(tests.WarnEmptyAttributes)
	sLogSuite.WarnOnly(tests.WarnGroupInline)
	sLogSuite.WarnOnly(tests.WarnLevelCase)
	sLogSuite.WarnOnly(tests.WarnMessageKey)
	sLogSuite.WarnOnly(tests.WarnTimeMillis)
	sLogSuite.WarnOnly(tests.WarnNoReplAttrBasic)
	sLogSuite.WarnOnly(tests.WarnResolver)
	sLogSuite.WarnOnly(tests.WarnZeroPC)
	sLogSuite.WarnOnly(tests.WarnZeroTime)
	suite.Run(t, sLogSuite)
}
