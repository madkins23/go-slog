package verify

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/madkins23/go-slog/creator"
	"github.com/madkins23/go-slog/infra"
	"github.com/madkins23/go-slog/verify/tests"
)

// Test_slog runs tests for the log/slog JSON handler.
func Test_slog(t *testing.T) {
	slogSuite := tests.NewSlogTestSuite(creator.Slog())
	slogSuite.WarnOnly(infra.WarnDuplicates)
	suite.Run(t, slogSuite)
}
