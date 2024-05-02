package verify

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/madkins23/go-slog/creator/phsymzerolog"
	warning2 "github.com/madkins23/go-slog/infra/warning"
	"github.com/madkins23/go-slog/verify/tests"
)

// TestVerifyPhsymZerolog runs tests for the phsym/zeroslog handler.
func TestVerifyPhsymZerolog(t *testing.T) {
	slogSuite := tests.NewSlogTestSuite(phsymzerolog.Creator())
	slogSuite.WarnOnly(warning2.Duplicates)
	slogSuite.WarnOnly(warning2.DurationMillis)
	slogSuite.WarnOnly(warning2.EmptyAttributes)
	slogSuite.WarnOnly(warning2.GroupEmpty)
	slogSuite.WarnOnly(warning2.GroupInline)
	slogSuite.WarnOnly(warning2.LevelCase)
	slogSuite.WarnOnly(warning2.MessageKey)
	slogSuite.WarnOnly(warning2.NoReplAttr)
	slogSuite.WarnOnly(warning2.SlogTest)
	slogSuite.WarnOnly(warning2.SourceCaller)
	slogSuite.WarnOnly(warning2.TimeSeconds)
	slogSuite.WarnOnly(warning2.WithGroupEmpty)
	slogSuite.WarnOnly(warning2.ZeroTime)
	suite.Run(t, slogSuite)
}
