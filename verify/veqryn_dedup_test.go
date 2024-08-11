package verify

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/madkins23/go-slog/creator/veqryndedup"
	"github.com/madkins23/go-slog/infra/warning"
	"github.com/madkins23/go-slog/verify/tests"
)

// TestVerifyVeqrynDedupGroup runs tests for the veqryn/dedup JSON handler in Ignore mode.
func TestVerifyVeqrynDedupGroup(t *testing.T) {
	slogSuite := tests.NewSlogTestSuite(veqryndedup.Creator(veqryndedup.Append))
	slogSuite.WarnOnly(warning.SkipDedup)
	suite.Run(t, slogSuite)
}

// TestVerifyVeqrynDedupIgnore runs tests for the veqryn/dedup JSON handler in Ignore mode.
func TestVerifyVeqrynDedupIgnore(t *testing.T) {
	slogSuite := tests.NewSlogTestSuite(veqryndedup.Creator(veqryndedup.Ignore))
	slogSuite.WarnOnly(warning.SkipDedup)
	suite.Run(t, slogSuite)
}

// TestVerifyVeqrynDedupIgnore runs tests for the veqryn/dedup JSON handler in Ignore mode.
func TestVerifyVeqrynDedupIncr(t *testing.T) {
	slogSuite := tests.NewSlogTestSuite(veqryndedup.Creator(veqryndedup.Increment))
	slogSuite.WarnOnly(warning.SkipDedup)
	suite.Run(t, slogSuite)
}

// TestVerifyVeqrynDedupOverwrite runs tests for the veqryn/dedup JSON handler in Overwrite mode.
func TestVerifyVeqrynDedupOverwrite(t *testing.T) {
	slogSuite := tests.NewSlogTestSuite(veqryndedup.Creator(veqryndedup.Overwrite))
	suite.Run(t, slogSuite)
}
