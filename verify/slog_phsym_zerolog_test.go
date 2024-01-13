package verify

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/madkins23/go-slog/creator"
	"github.com/madkins23/go-slog/infra"
	"github.com/madkins23/go-slog/verify/tests"
)

// Test_slog_zerolog_phsym runs tests for the physym zerolog handler.
func Test_slog_zerolog_phsym(t *testing.T) {
	slogSuite := tests.NewSlogTestSuite(creator.SlogPhsymZerolog())
	slogSuite.WarnOnly(infra.WarnDuplicates)
	slogSuite.WarnOnly(infra.WarnEmptyAttributes)
	slogSuite.WarnOnly(infra.WarnGroupInline)
	slogSuite.WarnOnly(infra.WarnLevelCase)
	slogSuite.WarnOnly(infra.WarnMessageKey)
	slogSuite.WarnOnly(infra.WarnNanoDuration)
	slogSuite.WarnOnly(infra.WarnNanoTime)
	slogSuite.WarnOnly(infra.WarnNoReplAttr)
	slogSuite.WarnOnly(infra.WarnSourceKey)
	slogSuite.WarnOnly(infra.WarnGroupEmpty)
	slogSuite.WarnOnly(infra.WarnZeroTime)
	suite.Run(t, slogSuite)
}
