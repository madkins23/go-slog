package verify

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/madkins23/go-slog/creator/samberzerolog"
	warning2 "github.com/madkins23/go-slog/infra/warning"
	"github.com/madkins23/go-slog/verify/tests"
)

// TestVerifySamberZerolog runs tests for the samber/slog-zerolog handler.
func TestVerifySamberZerolog(t *testing.T) {
	slogSuite := tests.NewSlogTestSuite(samberzerolog.Creator())
	slogSuite.WarnOnly(warning2.DefaultLevel)
	slogSuite.WarnOnly(warning2.DurationMillis)
	slogSuite.WarnOnly(warning2.EmptyAttributes)
	slogSuite.WarnOnly(warning2.GroupDuration)
	slogSuite.WarnOnly(warning2.GroupInline)
	slogSuite.WarnOnly(warning2.LevelCase)
	slogSuite.WarnOnly(warning2.MessageKey)
	slogSuite.WarnOnly(warning2.NoReplAttrBasic)
	slogSuite.WarnOnly(warning2.Resolver)
	slogSuite.WarnOnly(warning2.SlogTest)
	slogSuite.WarnOnly(warning2.TimeSeconds)
	slogSuite.WarnOnly(warning2.ZeroPC)
	slogSuite.WarnOnly(warning2.ZeroTime)
	suite.Run(t, slogSuite)
}
