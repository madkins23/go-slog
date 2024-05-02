package verify

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/madkins23/go-slog/creator/madkinssloggy"
	"github.com/madkins23/go-slog/infra/warning"
	"github.com/madkins23/go-slog/verify/tests"
)

// TestVerifyMadkinsSloggy runs tests for the madkins/sloggy JSON handler.
func TestVerifyMadkinsSloggy(t *testing.T) {
	slogSuite := tests.NewSlogTestSuite(madkinssloggy.Creator())
	slogSuite.WarnOnly(warning.Duplicates)
	suite.Run(t, slogSuite)
}
