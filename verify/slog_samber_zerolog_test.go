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

// Test_slog_samber_zerolog runs tests for the samber zerolog handler.
func Test_slog_samber_zerolog(t *testing.T) {
	sLogSuite := &test.SlogTestSuite{
		Creator: SlogSamberZerologHandlerCreator,
		Name:    "samber/slog-zerolog",
	}
	if *test.UseWarnings {
		sLogSuite.WarnOnly(test.WarnDefaultLevel)
		sLogSuite.WarnOnly(test.WarnEmptyAttributes)
		sLogSuite.WarnOnly(test.WarnGroupInline)
		sLogSuite.WarnOnly(test.WarnLevelCase)
		sLogSuite.WarnOnly(test.WarnMessageKey)
		sLogSuite.WarnOnly(test.WarnNanoDuration)
		sLogSuite.WarnOnly(test.WarnNanoTime)
		sLogSuite.WarnOnly(test.WarnNoReplAttrBasic)
		sLogSuite.WarnOnly(test.WarnResolver)
		sLogSuite.WarnOnly(test.WarnZeroPC)
		sLogSuite.WarnOnly(test.WarnZeroTime)
	}
	suite.Run(t, sLogSuite)
}

var _ test.CreateHandlerFn = SlogSamberZerologHandlerCreator

func SlogSamberZerologHandlerCreator(w io.Writer, options *slog.HandlerOptions) slog.Handler {
	zeroLogger := zerolog.New(w)
	return samber.Option{
		Logger:      &zeroLogger,
		Level:       options.Level,
		AddSource:   options.AddSource,
		ReplaceAttr: options.ReplaceAttr,
	}.NewZerologHandler()
}
