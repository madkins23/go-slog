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
	slogSuite.WarnOnly(warning.GroupWithTop)
	slogSuite.WarnOnly(warning.NoReplAttr)
	slogSuite.WarnOnly(warning.SourceKey)
	slogSuite.WarnOnly(warning.WithGroup)

	// If group start strings are cached in prefix,
	// no way of knowing if group is empty or not
	// (so that it can be removed entirely)
	// until the final Handle() call.
	slogSuite.WarnOnly(warning.WithGroupEmpty)
	// If this situation can be detected
	// then the parent's prefix might be used
	// (if a parent link is kept).

	suite.Run(t, slogSuite)
}
