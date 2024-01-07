package verify

import (
	"log/slog"
	"os"
	"testing"

	"github.com/phsym/zeroslog"
	"github.com/stretchr/testify/suite"
)

type SlogZerologPhsymTestSuite struct {
	SlogTestSuite
}

// Test_slog_zerolog_samber runs tests for the slog-zerolog handler.
func Test_slog_zerolog_phsym(t *testing.T) {
	suite.Run(t, &SlogZerologSamberTestSuite{})
}

func (suite *SlogZerologPhsymTestSuite) SimpleHandler() slog.Handler {
	return zeroslog.NewJsonHandler(os.Stderr, nil)
}
