package verify

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/madkins23/go-slog/creator"
	"github.com/madkins23/go-slog/infra"
	"github.com/madkins23/go-slog/verify/tests"
)

// Test_slog_samber_zap runs tests for the samber zerolog handler.
func Test_slog_samber_zap(t *testing.T) {
	sLogSuite := tests.NewSlogTestSuite(creator.SlogSamberZap())
	sLogSuite.WarnOnly(infra.WarnEmptyAttributes)
	sLogSuite.WarnOnly(infra.WarnGroupInline)
	sLogSuite.WarnOnly(infra.WarnLevelCase)
	sLogSuite.WarnOnly(infra.WarnDurationSeconds)
	sLogSuite.WarnOnly(infra.WarnNanoDuration)
	sLogSuite.WarnOnly(infra.WarnNanoTime)
	sLogSuite.WarnOnly(infra.WarnNoReplAttrBasic)
	sLogSuite.WarnOnly(infra.WarnResolver)
	sLogSuite.WarnOnly(infra.WarnZeroPC)
	sLogSuite.WarnOnly(infra.WarnZeroTime)
	sLogSuite.WarnOnly(infra.WarnZeroTimeNow)
	suite.Run(t, sLogSuite)
}
