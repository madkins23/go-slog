package verify

import (
	"io"
	"log/slog"
	"testing"

	"github.com/rs/zerolog"
	samber "github.com/samber/slog-zerolog/v2"
	"github.com/stretchr/testify/suite"

	"github.com/madkins23/go-slog/verify/test"
)

// Test_slog_zerolog_samber runs tests for the samber zerolog handler.
func Test_slog_zerolog_samber(t *testing.T) {
	sLogSuite := &test.SlogTestSuite{
		Creator: &SlogSamberCreator{},
		Name:    "samber/slog-zerolog",
	}
	if *test.UseWarnings {
		sLogSuite.WarnOnly(test.WarnMessageKey)
		sLogSuite.WarnOnly(test.WarnEmptyAttributes)
		sLogSuite.WarnOnly(test.WarnGroupInline)
		sLogSuite.WarnOnly(test.WarnLevelCase)
		sLogSuite.WarnOnly(test.WarnResolver)
		sLogSuite.WarnOnly(test.WarnZeroPC)
		sLogSuite.WarnOnly(test.WarnZeroTime)
	}
	suite.Run(t, sLogSuite)
}

var _ test.HandlerCreator = &SlogSamberCreator{}

type SlogSamberCreator struct{}

func (creator *SlogSamberCreator) SimpleHandler(w io.Writer, level slog.Leveler) slog.Handler {
	zeroLogger := zerolog.New(w)
	return samber.Option{
		Logger: &zeroLogger,
		Level:  level,
	}.NewZerologHandler()
}

func (creator *SlogSamberCreator) SourceHandler(w io.Writer, level slog.Leveler) slog.Handler {
	zeroLogger := zerolog.New(w)
	return samber.Option{
		Level:     level,
		Logger:    &zeroLogger,
		AddSource: true,
	}.NewZerologHandler()
}
