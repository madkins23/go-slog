package verify

import (
	"io"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/madkins23/go-slog/verify/test"
)

// Test_slog runs tests for the log/slog JSON handler.
func Test_slog(t *testing.T) {
	slogSuite := &test.SlogTestSuite{
		Creator: SlogHandlerCreator,
		Name:    "log/slog.JSONHandler",
	}
	if *test.UseWarnings {
		slogSuite.WarnOnly(test.WarnDuplicates)
	}
	suite.Run(t, slogSuite)
}

var _ test.CreateHandlerFn = SlogHandlerCreator

func SlogHandlerCreator(w io.Writer, options *slog.HandlerOptions) slog.Handler {
	return slog.NewJSONHandler(w, options)
}
