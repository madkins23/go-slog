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
		Creator: SlogPhsymHandlerCreator,
		Name:    "phsym/zeroslog",
	}
	if *test.UseWarnings {
		sLogSuite.WarnOnly(test.WarnMessageKey)
		sLogSuite.WarnOnly(test.WarnEmptyAttributes)
		sLogSuite.WarnOnly(test.WarnGroupInline)
		sLogSuite.WarnOnly(test.WarnLevelCase)
		sLogSuite.WarnOnly(test.WarnNanoDuration)
		sLogSuite.WarnOnly(test.WarnNanoTime)
		sLogSuite.WarnOnly(test.WarnNoReplAttr)
		sLogSuite.WarnOnly(test.WarnSourceKey)
		sLogSuite.WarnOnly(test.WarnSubgroupEmpty)
		sLogSuite.WarnOnly(test.WarnZeroTime)
	}
	suite.Run(t, sLogSuite)
}

var _ test.CreateHandlerFn = SlogPhsymHandlerCreator

func SlogPhsymHandlerCreator(w io.Writer, options *slog.HandlerOptions) slog.Handler {
	return zeroslog.NewJsonHandler(w, &zeroslog.HandlerOptions{
		Level:     options.Level,
		AddSource: options.AddSource,
	})
}
