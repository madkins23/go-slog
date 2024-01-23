package verify

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/madkins23/go-slog/creator"
	"github.com/madkins23/go-slog/verify/tests"
	"github.com/madkins23/go-slog/warning"
)

// Test_slog_samber_zerolog runs tests for the samber zerolog handler.
func Test_slog_darvaza_zerolog(t *testing.T) {
	slogSuite := tests.NewSlogTestSuite(creator.SlogDarvazaZerolog())
	slogSuite.WarnOnly(warning.DefaultLevel)
	slogSuite.WarnOnly(warning.DurationMillis)
	slogSuite.WarnOnly(warning.EmptyAttributes)
	slogSuite.WarnOnly(warning.GroupInline)
	slogSuite.WarnOnly(warning.LevelCase)
	slogSuite.WarnOnly(warning.MessageKey)
	slogSuite.WarnOnly(warning.TimeMillis)
	slogSuite.WarnOnly(warning.NoReplAttrBasic)
	slogSuite.WarnOnly(warning.Resolver)
	slogSuite.WarnOnly(warning.ZeroPC)
	slogSuite.WarnOnly(warning.ZeroTime)
	suite.Run(t, slogSuite)
}
