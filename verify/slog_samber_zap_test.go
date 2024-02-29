package verify

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/madkins23/go-slog/creator/samberzap"
	"github.com/madkins23/go-slog/internal/warning"
	"github.com/madkins23/go-slog/verify/tests"
)

// TestVerifySamberZap runs tests for the samber/slog-zap handler.
func TestVerifySamberZap(t *testing.T) {
	sLogSuite := tests.NewSlogTestSuite(samberzap.Creator())
	sLogSuite.WarnOnly(warning.DurationSeconds)
	sLogSuite.WarnOnly(warning.EmptyAttributes)
	sLogSuite.WarnOnly(warning.GroupDuration)
	sLogSuite.WarnOnly(warning.GroupInline)
	sLogSuite.WarnOnly(warning.LevelCase)
	sLogSuite.WarnOnly(warning.TimeMillis)
	sLogSuite.WarnOnly(warning.NoReplAttrBasic)
	sLogSuite.WarnOnly(warning.Resolver)
	sLogSuite.WarnOnly(warning.SourceKey)
	sLogSuite.WarnOnly(warning.ZeroPC)
	sLogSuite.WarnOnly(warning.ZeroTime)
	suite.Run(t, sLogSuite)
}
