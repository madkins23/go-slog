package verify

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/madkins23/go-slog/creator"
	"github.com/madkins23/go-slog/verify/tests"
	"github.com/madkins23/go-slog/warning"
)

// Test_slog_samber_zap runs tests for the samber zerolog handler.
func Test_slog_samber_zap(t *testing.T) {
	sLogSuite := tests.NewSlogTestSuite(creator.SlogSamberZap())
	sLogSuite.WarnOnly(warning.DurationMillis)
	sLogSuite.WarnOnly(warning.DurationSeconds)
	sLogSuite.WarnOnly(warning.EmptyAttributes)
	sLogSuite.WarnOnly(warning.GroupInline)
	sLogSuite.WarnOnly(warning.LevelCase)
	sLogSuite.WarnOnly(warning.TimeMillis)
	sLogSuite.WarnOnly(warning.NoReplAttrBasic)
	sLogSuite.WarnOnly(warning.Resolver)
	sLogSuite.WarnOnly(warning.ZeroPC)
	sLogSuite.WarnOnly(warning.ZeroTime)
	suite.Run(t, sLogSuite)
}
