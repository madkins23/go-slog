package verify

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/madkins23/go-slog/creator/svcrunnerjsonlog"
	"github.com/madkins23/go-slog/infra/warning"
	"github.com/madkins23/go-slog/verify/tests"
)

// TestVerifySvcrunnerJsonlog runs tests for the svcrunner/jsonlog handler.
func TestVerifySvcrunnerJsonlog(t *testing.T) {
	slogSuite := tests.NewSlogTestSuite(svcrunnerjsonlog.Creator())
	slogSuite.WarnOnly(warning.Duplicates)
	slogSuite.WarnOnly(warning.DurationString)
	slogSuite.WarnOnly(warning.LevelVar)
	slogSuite.WarnOnly(warning.MessageKey)
	slogSuite.WarnOnly(warning.NoEmptyName)
	slogSuite.WarnOnly(warning.NoReplAttr)
	slogSuite.WarnOnly(warning.SlogTest)
	slogSuite.WarnOnly(warning.SourceKey)
	suite.Run(t, slogSuite)
}
