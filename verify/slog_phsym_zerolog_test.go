package verify

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/madkins23/go-slog/creator/phsymzerolog"
	"github.com/madkins23/go-slog/internal/warning"
	"github.com/madkins23/go-slog/verify/tests"
)

// TestVerifyPhsymZerolog runs tests for the phsym/zeroslog handler.
func TestVerifyPhsymZerolog(t *testing.T) {
	slogSuite := tests.NewSlogTestSuite(phsymzerolog.Creator())
	slogSuite.WarnOnly(warning.Duplicates)
	slogSuite.WarnOnly(warning.DurationMillis)
	slogSuite.WarnOnly(warning.EmptyAttributes)
	slogSuite.WarnOnly(warning.GroupEmpty)
	slogSuite.WarnOnly(warning.GroupInline)
	slogSuite.WarnOnly(warning.LevelCase)
	slogSuite.WarnOnly(warning.MessageKey)
	slogSuite.WarnOnly(warning.TimeMillis)
	slogSuite.WarnOnly(warning.NoReplAttr)
	slogSuite.WarnOnly(warning.SourceKey)
	slogSuite.WarnOnly(warning.WithGroupEmpty)
	slogSuite.WarnOnly(warning.ZeroTime)
	suite.Run(t, slogSuite)
}
