package verify

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/madkins23/go-slog/creator/chanchalzap"
	warning2 "github.com/madkins23/go-slog/infra/warning"
	"github.com/madkins23/go-slog/verify/tests"
)

// TestVerifyChanchalZap runs tests for the chanchal/zaphandler handler.
func TestVerifyChanchalZap(t *testing.T) {
	slogSuite := tests.NewSlogTestSuite(chanchalzap.Creator())
	slogSuite.WarnOnly(warning2.CanceledContext)
	slogSuite.WarnOnly(warning2.Duplicates)
	slogSuite.WarnOnly(warning2.DurationSeconds)
	slogSuite.WarnOnly(warning2.GroupWithTop)
	slogSuite.WarnOnly(warning2.LevelCase)
	slogSuite.WarnOnly(warning2.LevelMath)
	slogSuite.WarnOnly(warning2.NoReplAttr)
	slogSuite.WarnOnly(warning2.NoReplAttrBasic)
	slogSuite.WarnOnly(warning2.SlogTest)
	slogSuite.WarnOnly(warning2.SourceCaller)
	slogSuite.WarnOnly(warning2.TimeSeconds)
	slogSuite.WarnOnly(warning2.WithGroup)
	slogSuite.WarnOnly(warning2.WithGroupEmpty)
	suite.Run(t, slogSuite)
}
