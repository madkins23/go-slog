package verify

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/madkins23/go-slog/creator/phusluslog"
	"github.com/madkins23/go-slog/infra/warning"
	"github.com/madkins23/go-slog/verify/tests"
)

// TestVerifyPhusluSlog runs tests for the phuslu/slog handler.
func TestVerifyPhusluSlog(t *testing.T) {
	slogSuite := tests.NewSlogTestSuite(phusluslog.Creator())
	slogSuite.WarnOnly(warning.Duplicates)
	suite.Run(t, slogSuite)
}
