package verify

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/madkins23/go-slog/creator/madkinssloggy"
	"github.com/madkins23/go-slog/internal/warning"
	"github.com/madkins23/go-slog/verify/tests"
)

// TestVerifyMadkinsSloggy runs tests for the slog/JSONHandler JSON handler.
func TestVerifyMadkinsSloggy(t *testing.T) {
	slogSuite := tests.NewSlogTestSuite(madkinssloggy.Creator())
	slogSuite.WarnOnly(warning.Duplicates)
	slogSuite.WarnOnly(warning.GroupInline)
	slogSuite.WarnOnly(warning.GroupWithTop)
	slogSuite.WarnOnly(warning.NoReplAttr)
	slogSuite.WarnOnly(warning.SourceKey)
	slogSuite.WarnOnly(warning.WithGroup)
	slogSuite.WarnOnly(warning.WithGroupEmpty)
	suite.Run(t, slogSuite)
}
