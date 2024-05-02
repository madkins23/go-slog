package verify

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/madkins23/go-slog/creator/samberzap"
	warning2 "github.com/madkins23/go-slog/infra/warning"
	"github.com/madkins23/go-slog/verify/tests"
)

// TestVerifySamberZap runs tests for the samber/slog-zap handler.
func TestVerifySamberZap(t *testing.T) {
	slogSuite := tests.NewSlogTestSuite(samberzap.Creator())
	slogSuite.WarnOnly(warning2.DurationSeconds)
	slogSuite.WarnOnly(warning2.EmptyAttributes)
	slogSuite.WarnOnly(warning2.GroupDuration)
	slogSuite.WarnOnly(warning2.GroupInline)
	slogSuite.WarnOnly(warning2.LevelCase)
	slogSuite.WarnOnly(warning2.NoReplAttrBasic)
	slogSuite.WarnOnly(warning2.Resolver)
	slogSuite.WarnOnly(warning2.SlogTest)
	slogSuite.WarnOnly(warning2.SourceCaller)
	slogSuite.WarnOnly(warning2.TimeSeconds)
	slogSuite.WarnOnly(warning2.ZeroPC)
	slogSuite.WarnOnly(warning2.ZeroTime)
	suite.Run(t, slogSuite)
}
