package verify

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/madkins23/go-slog/creator/phusluslog"
	warning2 "github.com/madkins23/go-slog/infra/warning"
	"github.com/madkins23/go-slog/verify/tests"
)

// TestVerifyPhusluSlog runs tests for the phuslu/slog handler.
func TestVerifyPhusluSlog(t *testing.T) {
	slogSuite := tests.NewSlogTestSuite(phusluslog.Creator())
	slogSuite.WarnOnly(warning2.Duplicates)
	suite.Run(t, slogSuite)
}
