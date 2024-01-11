package verify

import (
	"io"
	"log/slog"
	"testing"

	"github.com/rs/zerolog"
	samber "github.com/samber/slog-zerolog/v2"
	"github.com/stretchr/testify/suite"

	"github.com/madkins23/go-slog/verify/tests"
)

// Test_slog_samber_zerolog runs tests for the samber zerolog handler.
func Test_slog_samber_zerolog(t *testing.T) {
	sLogSuite := &tests.SlogTestSuite{
		Creator: SlogSamberZerologHandlerCreator,
		Name:    "samber/slog-zerolog",
	}
	if *tests.UseWarnings {
		sLogSuite.WarnOnly(tests.WarnDefaultLevel)
		sLogSuite.WarnOnly(tests.WarnEmptyAttributes)
		sLogSuite.WarnOnly(tests.WarnGroupInline)
		sLogSuite.WarnOnly(tests.WarnLevelCase)
		sLogSuite.WarnOnly(tests.WarnMessageKey)
		sLogSuite.WarnOnly(tests.WarnNanoDuration)
		sLogSuite.WarnOnly(tests.WarnNanoTime)
		sLogSuite.WarnOnly(tests.WarnNoReplAttrBasic)
		sLogSuite.WarnOnly(tests.WarnResolver)
		sLogSuite.WarnOnly(tests.WarnZeroPC)
		sLogSuite.WarnOnly(tests.WarnZeroTime)
	}
	suite.Run(t, sLogSuite)
}

var _ tests.CreateHandlerFn = SlogSamberZerologHandlerCreator

func SlogSamberZerologHandlerCreator(w io.Writer, options *slog.HandlerOptions) slog.Handler {
	zeroLogger := zerolog.New(w)
	return samber.Option{
		Logger:      &zeroLogger,
		Level:       options.Level,
		AddSource:   options.AddSource,
		ReplaceAttr: options.ReplaceAttr,
	}.NewZerologHandler()
}