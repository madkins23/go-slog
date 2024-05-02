package verify

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/madkins23/go-slog/creator/slogjson"
	"github.com/madkins23/go-slog/infra/warning"
	"github.com/madkins23/go-slog/verify/tests"
)

// TestVerifySlogJSON runs tests for the slog/JSONHandler JSON handler.
func TestVerifySlogJSON(t *testing.T) {
	slogSuite := tests.NewSlogTestSuite(slogjson.Creator())
	slogSuite.WarnOnly(warning.Duplicates)
	suite.Run(t, slogSuite)
}
