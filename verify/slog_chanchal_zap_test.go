package verify

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/madkins23/go-slog/creator/chanchalzap"
	"github.com/madkins23/go-slog/verify/tests"
	"github.com/madkins23/go-slog/warning"
)

// Test_slog_samber_zap runs tests for the samber zerolog handler.
func Test_slog_chanchal_zap(t *testing.T) {
	sLogSuite := tests.NewSlogTestSuite(chanchalzap.Creator())
	sLogSuite.WarnOnly(warning.CanceledContext)
	sLogSuite.WarnOnly(warning.Duplicates)
	sLogSuite.WarnOnly(warning.DurationSeconds)
	sLogSuite.WarnOnly(warning.LevelCase)
	sLogSuite.WarnOnly(warning.LevelMath)
	sLogSuite.WarnOnly(warning.TimeMillis)
	sLogSuite.WarnOnly(warning.NoReplAttr)
	sLogSuite.WarnOnly(warning.NoReplAttrBasic)
	sLogSuite.WarnOnly(warning.SourceKey)
	sLogSuite.WarnOnly(warning.WithGroup)
	suite.Run(t, sLogSuite)
}
