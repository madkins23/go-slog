package verify

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/madkins23/go-slog/creator/madkinsflash"
	"github.com/madkins23/go-slog/infra/warning"
	"github.com/madkins23/go-slog/verify/tests"
)

// TestVerifyMadkinsFlash runs tests for the madkins/flash JSON handler.
func TestVerifyMadkinsFlash(t *testing.T) {
	slogSuite := tests.NewSlogTestSuite(madkinsflash.Creator())
	slogSuite.WarnOnly(warning.Duplicates)
	suite.Run(t, slogSuite)
}
