package tests

import (
	"bytes"
	"log/slog"

	"github.com/stretchr/testify/suite"

	"github.com/madkins23/go-slog/infra"
)

// -----------------------------------------------------------------------------
// Top level definitions.

// SlogTestSuite provides various tests for a specified log/slog.Handler.
type SlogTestSuite struct {
	suite.Suite
	*bytes.Buffer
	*infra.WarningManager

	// Creator creates a slog.Handler to be used in creating a slog.Logger for a test.
	// This field must be configured by test suites and shouldn't be changed later.
	Creator infra.Creator
}

func NewSlogTestSuite(creator infra.Creator) *SlogTestSuite {
	return &SlogTestSuite{
		Creator:        creator,
		WarningManager: infra.NewWarningManager(creator.Name()),
	}
}

// -----------------------------------------------------------------------------
// Suite test configuration.

func (suite *SlogTestSuite) SetupTest() {
	suite.Buffer = &bytes.Buffer{}
}

// -----------------------------------------------------------------------------
// Handler/Logger creation.

// Logger returns a new slog.Logger with the specified options.
func (suite *SlogTestSuite) Logger(options *slog.HandlerOptions) *slog.Logger {
	return slog.New(suite.Creator.NewHandle(suite.Buffer, options))
}
