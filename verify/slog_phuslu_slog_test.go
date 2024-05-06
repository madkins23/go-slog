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

	// For use when back testing with v1.93.0:
	//   go get github.com/phuslu/log@v1.0.93
	//slogSuite.WarnOnly(warning.DurationMillis)
	//slogSuite.WarnOnly(warning.LevelVar)
	//slogSuite.WarnOnly(warning.Mismatch)
	//slogSuite.WarnOnly(warning.StringAny)
	//slogSuite.WarnOnly(warning.TimeMillis)
	//slogSuite.WarnOnly(warning.ZeroPC)

	suite.Run(t, slogSuite)
}
