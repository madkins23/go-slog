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
	sLogSuite := tests.NewSlogTestSuite("phsym/zeroslog", SlogPhsymZerologHandlerCreatorFn)
	sLogSuite.WarnOnly(infra.WarnDuplicates)
	sLogSuite.WarnOnly(infra.WarnEmptyAttributes)
	sLogSuite.WarnOnly(infra.WarnGroupInline)
	sLogSuite.WarnOnly(infra.WarnLevelCase)
	sLogSuite.WarnOnly(infra.WarnMessageKey)
	sLogSuite.WarnOnly(infra.WarnNanoDuration)
	sLogSuite.WarnOnly(infra.WarnNanoTime)
	sLogSuite.WarnOnly(infra.WarnNoReplAttr)
	sLogSuite.WarnOnly(infra.WarnSourceKey)
	sLogSuite.WarnOnly(infra.WarnGroupEmpty)
	sLogSuite.WarnOnly(infra.WarnZeroTime)
	suite.Run(t, sLogSuite)
}

var _ infra.CreatorFn = SlogPhsymZerologHandlerCreatorFn

func SlogPhsymZerologHandlerCreatorFn(w io.Writer, options *slog.HandlerOptions) slog.Handler {
	return zeroslog.NewJsonHandler(w, &zeroslog.HandlerOptions{
		Level:     options.Level,
		AddSource: options.AddSource,
	})
}
