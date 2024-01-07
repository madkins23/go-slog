package verify

import (
	"log/slog"
	"testing"

	"github.com/rs/zerolog"
	samber "github.com/samber/slog-zerolog/v2"
	"github.com/stretchr/testify/suite"
)

type SlogZerologSamberTestSuite struct {
	SlogTestSuite
}

// Test_slog_zerolog_samber runs tests for the slog-zerolog handler.
func Test_slog_zerolog_samber(t *testing.T) {
	suite.Run(t, &SlogZerologSamberTestSuite{})
}

func (suite *SlogZerologSamberTestSuite) SimpleHandler() slog.Handler {
	zeroLogger := zerolog.New(suite.Buffer)
	return samber.Option{Logger: &zeroLogger}.NewZerologHandler()
}
