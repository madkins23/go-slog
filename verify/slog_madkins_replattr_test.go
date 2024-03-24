package verify

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/madkins23/go-slog/creator/madkinsreplattr"
	"github.com/madkins23/go-slog/internal/warning"
	"github.com/madkins23/go-slog/verify/tests"
)

// TestVerifyMadkinsReplAttr runs tests for the madkins/replattr JSON handler.
func TestVerifyMadkinsReplAttr(t *testing.T) {
	slogSuite := tests.NewSlogTestSuite(madkinsreplattr.Creator())
	slogSuite.WarnOnly(warning.Duplicates)
	suite.Run(t, slogSuite)
}
