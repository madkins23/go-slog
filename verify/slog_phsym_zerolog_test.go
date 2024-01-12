package verify

import (
	"io"
	"log/slog"
	"testing"

	"github.com/phsym/zeroslog"
	"github.com/stretchr/testify/suite"

	"github.com/madkins23/go-slog/infra"
	"github.com/madkins23/go-slog/verify/tests"
)

// Test_slog_zerolog_phsym runs tests for the physym zerolog handler.
func Test_slog_zerolog_phsym(t *testing.T) {
	sLogSuite := &tests.SlogTestSuite{
		Name:    "phsym/zeroslog",
		Creator: infra.NewCreator(SlogPhsymZerologHandlerCreatorFn),
	}
	sLogSuite.WarnOnly(tests.WarnDuplicates)
	sLogSuite.WarnOnly(tests.WarnEmptyAttributes)
	sLogSuite.WarnOnly(tests.WarnGroupInline)
	sLogSuite.WarnOnly(tests.WarnLevelCase)
	sLogSuite.WarnOnly(tests.WarnMessageKey)
	sLogSuite.WarnOnly(tests.WarnNanoDuration)
	sLogSuite.WarnOnly(tests.WarnNanoTime)
	sLogSuite.WarnOnly(tests.WarnNoReplAttr)
	sLogSuite.WarnOnly(tests.WarnSourceKey)
	sLogSuite.WarnOnly(tests.WarnGroupEmpty)
	sLogSuite.WarnOnly(tests.WarnZeroTime)
	suite.Run(t, sLogSuite)
}

var _ infra.CreatorFn = SlogPhsymZerologHandlerCreatorFn

func SlogPhsymZerologHandlerCreatorFn(w io.Writer, options *slog.HandlerOptions) slog.Handler {
	return zeroslog.NewJsonHandler(w, &zeroslog.HandlerOptions{
		Level:     options.Level,
		AddSource: options.AddSource,
	})
}
