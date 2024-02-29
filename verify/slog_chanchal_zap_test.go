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
	slogSuite := tests.NewSlogTestSuite(chanchalzap.Creator())
	slogSuite.WarnOnly(warning.CanceledContext)
	slogSuite.WarnOnly(warning.Duplicates)
	slogSuite.WarnOnly(warning.DurationSeconds)
	slogSuite.WarnOnly(warning.GroupEmpty)
	slogSuite.WarnOnly(warning.GroupInline)
	slogSuite.WarnOnly(warning.GroupWithTop)
	slogSuite.WarnOnly(warning.LevelCase)
	slogSuite.WarnOnly(warning.LevelMath)
	slogSuite.WarnOnly(warning.TimeMillis)
	slogSuite.WarnOnly(warning.NoReplAttr)
	slogSuite.WarnOnly(warning.NoReplAttrBasic)
	slogSuite.WarnOnly(warning.SourceKey)
	slogSuite.WarnOnly(warning.WithGroup)
	suite.Run(t, slogSuite)
}
