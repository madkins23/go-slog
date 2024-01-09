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
		Creator: &SlogCreator{},
		Name:    "log/slog.JSONHandler",
	}
	if *test.UseWarnings {
		slogSuite.WarnOnly(test.WarnDuplicates)
	}
	suite.Run(t, slogSuite)
}

var _ test.LoggerCreator = &SlogCreator{}

type SlogCreator struct{}

func (creator *SlogCreator) SimpleLogger(w io.Writer) *slog.Logger {
	return slog.New(slog.NewJSONHandler(w, nil))
}

func (creator *SlogCreator) SourceLogger(w io.Writer) *slog.Logger {
	return slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{AddSource: true}))
}
