package verify

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/madkins23/go-slog/creator/chanchalzap"
	"github.com/madkins23/go-slog/infra/warning"
	"github.com/madkins23/go-slog/verify/tests"
)

// TestVerifyChanchalZap runs tests for the chanchal/zaphandler handler.
func TestVerifyChanchalZap(t *testing.T) {
	slogSuite := tests.NewSlogTestSuite(chanchalzap.Creator())
	slogSuite.WarnOnly(warning.CanceledContext)
	slogSuite.WarnOnly(warning.Duplicates)
	slogSuite.WarnOnly(warning.DurationSeconds)
	slogSuite.WarnOnly(warning.GroupWithTop)
	slogSuite.WarnOnly(warning.LevelCase)
	slogSuite.WarnOnly(warning.LevelMath)
	slogSuite.WarnOnly(warning.NoReplAttr)
	slogSuite.WarnOnly(warning.NoReplAttrBasic)
	slogSuite.WarnOnly(warning.SlogTest)
	slogSuite.WarnOnly(warning.SourceCaller)
	slogSuite.WarnOnly(warning.TimeSeconds)
	slogSuite.WarnOnly(warning.WithGroup)
	slogSuite.WarnOnly(warning.WithGroupEmpty)
	suite.Run(t, slogSuite)
}
