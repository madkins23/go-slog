package verify

import (
	"io"
	"log/slog"
	"testing"

	"github.com/rs/zerolog"
	samber "github.com/samber/slog-zerolog/v2"
	"github.com/stretchr/testify/suite"

	"github.com/madkins23/go-slog/infra"
	"github.com/madkins23/go-slog/verify/tests"
)

// Test_slog_samber_zerolog runs tests for the samber zerolog handler.
func Test_slog_samber_zerolog(t *testing.T) {
	sLogSuite := tests.NewSlogTestSuite("samber/slog-zerolog", SlogSamberZerologHandlerCreatorFn)
	sLogSuite.WarnOnly(infra.WarnDefaultLevel)
	sLogSuite.WarnOnly(infra.WarnEmptyAttributes)
	sLogSuite.WarnOnly(infra.WarnGroupInline)
	sLogSuite.WarnOnly(infra.WarnLevelCase)
	sLogSuite.WarnOnly(infra.WarnMessageKey)
	sLogSuite.WarnOnly(infra.WarnNanoDuration)
	sLogSuite.WarnOnly(infra.WarnNanoTime)
	sLogSuite.WarnOnly(infra.WarnNoReplAttrBasic)
	sLogSuite.WarnOnly(infra.WarnResolver)
	sLogSuite.WarnOnly(infra.WarnZeroPC)
	sLogSuite.WarnOnly(infra.WarnZeroTime)
	suite.Run(t, sLogSuite)
}

var _ infra.CreatorFn = SlogSamberZerologHandlerCreatorFn

func SlogSamberZerologHandlerCreatorFn(w io.Writer, options *slog.HandlerOptions) slog.Handler {
	zeroLogger := zerolog.New(w)
	return samber.Option{
		Logger:      &zeroLogger,
		Level:       options.Level,
		AddSource:   options.AddSource,
		ReplaceAttr: options.ReplaceAttr,
	}.NewZerologHandler()
}
