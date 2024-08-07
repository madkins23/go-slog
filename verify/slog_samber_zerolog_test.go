package verify

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/madkins23/go-slog/creator/samberzerolog"
	"github.com/madkins23/go-slog/infra/warning"
	"github.com/madkins23/go-slog/verify/tests"
)

// TestVerifySamberZerolog runs tests for the samber/slog-zerolog handler.
func TestVerifySamberZerolog(t *testing.T) {
	slogSuite := tests.NewSlogTestSuite(samberzerolog.Creator())
	slogSuite.WarnOnly(warning.DefaultLevel)
	slogSuite.WarnOnly(warning.DurationMillis)
	slogSuite.WarnOnly(warning.GroupDuration)
	slogSuite.WarnOnly(warning.GroupInline)
	slogSuite.WarnOnly(warning.LevelCase)
	slogSuite.WarnOnly(warning.MessageKey)
	slogSuite.WarnOnly(warning.NoEmptyName)
	slogSuite.WarnOnly(warning.NoNilValue)
	slogSuite.WarnOnly(warning.NoReplAttrBasic)
	slogSuite.WarnOnly(warning.Resolver)
	slogSuite.WarnOnly(warning.SlogTest)
	slogSuite.WarnOnly(warning.TimeSeconds)
	slogSuite.WarnOnly(warning.ZeroPC)
	slogSuite.WarnOnly(warning.ZeroTime)
	suite.Run(t, slogSuite)
}
