package verify

import (
	"io"
	"log/slog"
	"testing"

	"github.com/rs/zerolog"
	samber "github.com/samber/slog-zerolog/v2"
	"github.com/stretchr/testify/suite"

	"github.com/madkins23/go-slog/replace"
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
		sLogSuite.WarnOnly(test.WarnNanoDuration)
		sLogSuite.WarnOnly(test.WarnNanoTime)
		sLogSuite.WarnOnly(test.WarnNoReplAttrBasic)
		sLogSuite.WarnOnly(test.WarnResolver)
		sLogSuite.WarnOnly(test.WarnZeroPC)
		sLogSuite.WarnOnly(test.WarnZeroTime)
	}
	suite.Run(t, sLogSuite)
}

var _ test.HandlerCreator = &SlogSamberCreator{}

type SlogSamberCreator struct{}

func (creator *SlogSamberCreator) SimpleHandler(w io.Writer, level slog.Leveler, replAttr replace.AttrFn) slog.Handler {
	zeroLogger := zerolog.New(w)
	return samber.Option{
		Logger:      &zeroLogger,
		Level:       level,
		ReplaceAttr: replAttr,
	}.NewZerologHandler()
}

func (creator *SlogSamberCreator) SourceHandler(w io.Writer, level slog.Leveler, replAttr replace.AttrFn) slog.Handler {
	zeroLogger := zerolog.New(w)
	return samber.Option{
		Level:       level,
		Logger:      &zeroLogger,
		AddSource:   true,
		ReplaceAttr: replAttr,
	}.NewZerologHandler()
}

/*

// Test_slog_zerolog_samber runs tests for the samber zerolog handler.
// ReplaceAttr functions are used to fix known issues.
func Test_slog_zerolog_samber_replace(t *testing.T) {
	sLogSuite := &test.SlogTestSuite{
		Creator: &SlogSamberCreatorReplace{},
		Name:    "samber/slog-zerolog",
	}
	if *test.UseWarnings {
		//sLogSuite.WarnOnly(test.WarnMessageKey)
		sLogSuite.WarnOnly(test.WarnEmptyAttributes)
		sLogSuite.WarnOnly(test.WarnGroupInline)
		sLogSuite.WarnOnly(test.WarnLevelCase)
		sLogSuite.WarnOnly(test.WarnResolver)
		sLogSuite.WarnOnly(test.WarnZeroPC)
		sLogSuite.WarnOnly(test.WarnZeroTime)
	}
	suite.Run(t, sLogSuite)
}

var _ test.HandlerCreator = &SlogSamberCreatorReplace{}

type SlogSamberCreatorReplace struct{}

func (creator *SlogSamberCreatorReplace) SimpleHandler(w io.Writer, level slog.Leveler, _ replace.AttrFn) slog.Handler {
	zeroLogger := zerolog.New(w)
	return samber.Option{
		Logger:      &zeroLogger,
		Level:       level,
		ReplaceAttr: replace.MessageToMsg,
	}.NewZerologHandler()
}

func (creator *SlogSamberCreatorReplace) SourceHandler(w io.Writer, level slog.Leveler, _ replace.AttrFn) slog.Handler {
	zeroLogger := zerolog.New(w)
	return samber.Option{
		Level:       level,
		Logger:      &zeroLogger,
		AddSource:   true,
		ReplaceAttr: replace.MessageToMsg,
	}.NewZerologHandler()
}

*/
