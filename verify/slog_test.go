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

var _ test.HandlerCreator = &SlogCreator{}

type SlogCreator struct{}

func (creator *SlogCreator) SimpleHandler(w io.Writer, level slog.Leveler) slog.Handler {
	return slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level: level,
	})
}

func (creator *SlogCreator) SourceHandler(w io.Writer, level slog.Leveler) slog.Handler {
	return slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
	})
}
