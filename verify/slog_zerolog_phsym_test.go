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

func (suite *SlogZerologPhsymTestSuite) SimpleLogger() *slog.Logger {
	return slog.New(zeroslog.NewJsonHandler(os.Stderr, nil))
}

func (suite *SlogZerologPhsymTestSuite) SourceLogger() *slog.Logger {
	return slog.New(zeroslog.NewJsonHandler(os.Stderr, &zeroslog.HandlerOptions{AddSource: true}))
}
