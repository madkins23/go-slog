package verify

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/madkins23/go-slog/creator/snqkmeld"
	"github.com/madkins23/go-slog/verify/tests"
)

// TestVerifySnqkMeld runs tests for the snqk/meld JSON handler.
func TestVerifySnqkMeld(t *testing.T) {
	slogSuite := tests.NewSlogTestSuite(snqkmeld.Creator())
	suite.Run(t, slogSuite)
}
