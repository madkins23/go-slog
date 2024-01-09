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
		sLogSuite.WarnOnly(test.WarnZeroTime)
	}
	suite.Run(t, sLogSuite)
}

var _ test.LoggerCreator = &SlogSamberCreator{}

type SlogSamberCreator struct{}

func (creator *SlogSamberCreator) SimpleLogger(w io.Writer) *slog.Logger {
	zeroLogger := zerolog.New(w)
	return slog.New(samber.Option{Logger: &zeroLogger, Level: slog.LevelInfo}.NewZerologHandler())
}

func (creator *SlogSamberCreator) SourceLogger(w io.Writer) *slog.Logger {
	zeroLogger := zerolog.New(w)
	return slog.New(samber.Option{
		Logger:    &zeroLogger,
		AddSource: true,
	}.NewZerologHandler())
}
