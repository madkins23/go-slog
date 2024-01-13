package tests

import (
	"bytes"
	"log/slog"

	"github.com/stretchr/testify/suite"

	"github.com/madkins23/go-slog/infra"
	"github.com/madkins23/go-slog/replace"
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

func NewSlogTestSuite(name string, fn infra.CreatorFn) *SlogTestSuite {
	return &SlogTestSuite{
		Creator:        infra.NewCreator(name, fn),
		WarningManager: infra.NewWarningManager(name),
	}
}

// -----------------------------------------------------------------------------
// Suite test configuration.

func (suite *SlogTestSuite) SetupTest() {
	suite.Buffer = &bytes.Buffer{}
}

// -----------------------------------------------------------------------------
// Handler/Logger creation.

// SimpleOptions returns a default, simple, slog.HandlerOptions.
func SimpleOptions() *slog.HandlerOptions {
	return &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
}

// LevelOptions returns a slog.HandlerOptions with the specified level.
func LevelOptions(level slog.Leveler) *slog.HandlerOptions {
	return &slog.HandlerOptions{
		Level: level,
	}
}

// SourceOptions returns a slog.HandlerOptions with the specified level
// and the AddSource field set to true.
func SourceOptions() *slog.HandlerOptions {
	return &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	}
}

// ReplaceAttrOptions returns a slog.HandlerOptions with the specified ReplaceAttr function.
func ReplaceAttrOptions(fn replace.AttrFn) *slog.HandlerOptions {
	return &slog.HandlerOptions{
		Level:       slog.LevelInfo,
		ReplaceAttr: fn,
	}
}

// Logger returns a new slog.Logger with the specified options.
func (suite *SlogTestSuite) Logger(options *slog.HandlerOptions) *slog.Logger {
	return slog.New(suite.Creator.NewHandle(suite.Buffer, options))
}
