package verify

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/madkins23/go-slog/creator/samberzap"
	"github.com/madkins23/go-slog/infra/warning"
	"github.com/madkins23/go-slog/verify/tests"
)

// TestVerifySamberZap runs tests for the samber/slog-zap handler.
func TestVerifySamberZap(t *testing.T) {
	slogSuite := tests.NewSlogTestSuite(samberzap.Creator())
	slogSuite.WarnOnly(warning.DurationSeconds)
	slogSuite.WarnOnly(warning.GroupDuration)
	slogSuite.WarnOnly(warning.GroupInline)
	slogSuite.WarnOnly(warning.LevelCase)
	slogSuite.WarnOnly(warning.NoEmptyName)
	slogSuite.WarnOnly(warning.NoNilValue)
	slogSuite.WarnOnly(warning.NoReplAttrBasic)
	slogSuite.WarnOnly(warning.Resolver)
	slogSuite.WarnOnly(warning.SlogTest)
	slogSuite.WarnOnly(warning.SourceCaller)
	slogSuite.WarnOnly(warning.TimeSeconds)
	slogSuite.WarnOnly(warning.ZeroPC)
	slogSuite.WarnOnly(warning.ZeroTime)
	suite.Run(t, slogSuite)
}
