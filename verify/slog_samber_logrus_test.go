package verify

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/madkins23/go-slog/creator/samberlogrus"
	"github.com/madkins23/go-slog/internal/warning"
	"github.com/madkins23/go-slog/verify/tests"
)

// TestVerifySamberLogrus runs tests for the samber/slog-logrus handler.
func TestVerifySamberLogrus(t *testing.T) {
	slogSuite := tests.NewSlogTestSuite(samberlogrus.Creator())
	slogSuite.WarnOnly(warning.EmptyAttributes)
	slogSuite.WarnOnly(warning.GroupInline)
	slogSuite.WarnOnly(warning.LevelCase)
	slogSuite.WarnOnly(warning.NoReplAttrBasic)
	slogSuite.WarnOnly(warning.Resolver)
	slogSuite.WarnOnly(warning.SlogTest)
	slogSuite.WarnOnly(warning.ZeroPC)
	slogSuite.WarnOnly(warning.ZeroTime)
	suite.Run(t, slogSuite)
}
