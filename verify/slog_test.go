package verify

import (
	"io"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/madkins23/go-slog/infra"
	"github.com/madkins23/go-slog/verify/tests"
)

// Test_slog runs tests for the log/slog JSON handler.
func Test_slog(t *testing.T) {
	slogSuite := tests.NewSlogTestSuite("log/slog.JSONHandler", SlogHandlerCreator)
	slogSuite.WarnOnly(infra.WarnDuplicates)
	suite.Run(t, slogSuite)
}

var _ infra.CreatorFn = SlogHandlerCreator

func SlogHandlerCreator(w io.Writer, options *slog.HandlerOptions) slog.Handler {
	return slog.NewJSONHandler(w, options)
}
