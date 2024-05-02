package verify

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/madkins23/go-slog/creator/madkinsreplattr"
	"github.com/madkins23/go-slog/infra/warning"
	"github.com/madkins23/go-slog/verify/tests"
)

// TestVerifyMadkinsReplAttr runs tests for the madkins/replattr JSON handler.
//
// This test was constructed to show the performance of madkins/flash
// when configured with verification errors that are fixed using ReplaceAttr functions.
// The result shows up on the cmd/server scores chart.
func TestVerifyMadkinsReplAttr(t *testing.T) {
	slogSuite := tests.NewSlogTestSuite(madkinsreplattr.Creator())
	slogSuite.WarnOnly(warning.Duplicates)
	suite.Run(t, slogSuite)
}
