package tests

import (
	"bytes"
	"fmt"
	"log/slog"
	"strings"

	"github.com/stretchr/testify/suite"

	"github.com/madkins23/go-slog/infra"
	"github.com/madkins23/go-slog/infra/warning"
)

// -----------------------------------------------------------------------------
// Top level definitions.

// SlogTestSuite provides various tests for a specified log/slog.Handler.
type SlogTestSuite struct {
	suite.Suite
	*warning.Manager

	// Creator creates a slog.Handler to be used in creating a slog.Logger for a test.
	// This field must be configured by test suites and shouldn't be changed later.
	Creator infra.Creator

	*bytes.Buffer
}

func NewSlogTestSuite(creator infra.Creator) *SlogTestSuite {
	return &SlogTestSuite{
		Creator: creator,
		Manager: NewWarningManager(creator.Name()),
	}
}

// -----------------------------------------------------------------------------
// Suite test configuration.

func (suite *SlogTestSuite) SetupSuite() {
	if suite.Creator.HasSummary() {
		fmt.Printf(":[ %s\n", suite.Creator.Name())
		for _, line := range strings.Split(suite.Creator.Summary(), "\n") {
			fmt.Printf(":: %s\n", line)
		}
		for name, link := range suite.Creator.Links() {
			fmt.Printf(":> %s --> %s\n", name, link)
		}
		fmt.Println(":]")
	}
}

func (suite *SlogTestSuite) SetupTest() {
	suite.Buffer = &bytes.Buffer{}
}

// -----------------------------------------------------------------------------
// Handler/Logger creation.

// Logger returns a new slog.Logger with the specified options.
func (suite *SlogTestSuite) Logger(options *slog.HandlerOptions) *slog.Logger {
	return suite.Creator.NewLogger(suite.Buffer, options)
}
