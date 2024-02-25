package verify

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/madkins23/go-slog/creator/chanchalzap"
	"github.com/madkins23/go-slog/internal/warning"
	"github.com/madkins23/go-slog/verify/tests"
)

// TestVerifyChanchalZap runs tests for the chanchal/zaphandler handler.
func TestVerifyChanchalZap(t *testing.T) {
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
