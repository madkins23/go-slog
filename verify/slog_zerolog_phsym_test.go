package verify

import (
	"io"
	"log/slog"
	"testing"

	"github.com/phsym/zeroslog"
	"github.com/stretchr/testify/suite"

	"github.com/madkins23/go-slog/verify/test"
)

// Test_slog_zerolog_phsym runs tests for the physym zerolog handler.
func Test_slog_zerolog_phsym(t *testing.T) {
	sLogSuite := &test.SlogTestSuite{
		Creator: &SlogPhsymCreator{},
		Name:    "phsym/zeroslog",
	}
	if *test.UseWarnings {
		sLogSuite.WarnOnly(test.WarnMessageKey)
		sLogSuite.WarnOnly(test.WarnEmptyAttributes)
		sLogSuite.WarnOnly(test.WarnGroupInline)
		sLogSuite.WarnOnly(test.WarnLevelCase)
		sLogSuite.WarnOnly(test.WarnSourceKey)
		sLogSuite.WarnOnly(test.WarnSubgroupEmpty)
		sLogSuite.WarnOnly(test.WarnZeroTime)
	}
	suite.Run(t, sLogSuite)
}

var _ test.LoggerCreator = &SlogPhsymCreator{}

type SlogPhsymCreator struct{}

func (creator *SlogPhsymCreator) SimpleLogger(w io.Writer) *slog.Logger {
	return slog.New(zeroslog.NewJsonHandler(w, nil))
}

func (creator *SlogPhsymCreator) SourceLogger(w io.Writer) *slog.Logger {
	return slog.New(zeroslog.NewJsonHandler(w, &zeroslog.HandlerOptions{AddSource: true}))
}
