package verify

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/madkins23/go-slog/creator/samberzerolog"
	"github.com/madkins23/go-slog/internal/warning"
	"github.com/madkins23/go-slog/verify/tests"
)

// TestVerifySamberZerolog runs tests for the samber/slog-zerolog handler.
func TestVerifySamberZerolog(t *testing.T) {
	sLogSuite := tests.NewSlogTestSuite(samberzerolog.Creator())
	sLogSuite.WarnOnly(warning.DefaultLevel)
	sLogSuite.WarnOnly(warning.DurationMillis)
	sLogSuite.WarnOnly(warning.EmptyAttributes)
	sLogSuite.WarnOnly(warning.GroupDuration)
	sLogSuite.WarnOnly(warning.GroupInline)
	sLogSuite.WarnOnly(warning.LevelCase)
	sLogSuite.WarnOnly(warning.MessageKey)
	sLogSuite.WarnOnly(warning.TimeMillis)
	sLogSuite.WarnOnly(warning.NoReplAttrBasic)
	sLogSuite.WarnOnly(warning.Resolver)
	sLogSuite.WarnOnly(warning.ZeroPC)
	sLogSuite.WarnOnly(warning.ZeroTime)
	suite.Run(t, sLogSuite)
}
