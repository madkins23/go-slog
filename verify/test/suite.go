package test

import (
	"bytes"
	"io"
	"log/slog"

	"github.com/stretchr/testify/suite"

	"github.com/madkins23/go-slog/replace"
)

// -----------------------------------------------------------------------------
// Top level definitions.

// CreateHandlerFn is responsible for generating log/slog handler instances.
// Define one for a given test file and use to instantiate SlogTestSuite.
type CreateHandlerFn func(w io.Writer, options *slog.HandlerOptions) slog.Handler

// SlogTestSuite provides various tests for a specified log/slog.Handler.
type SlogTestSuite struct {
	suite.Suite
	*bytes.Buffer
	warn     map[string]bool
	warnings map[string]*Warning

	// Creator creates a slog.Handler to be used in creating a slog.Logger for a test.
	// This field must be configured by test suites.
	Creator CreateHandlerFn

	// Name of Handler for warnings display.
	Name string
}

// -----------------------------------------------------------------------------
// Suite test configuration.

const duplicateFieldsNotError = true

// suites captures all suites tested together into an array.
// This array is used when showing warnings.
var suites = make([]*SlogTestSuite, 0)

func (suite *SlogTestSuite) SetupSuite() {
	suites = append(suites, suite)
	if duplicateFieldsNotError {
		// There doesn't seem to be a rule about this in https://pkg.go.dev/log/slog@master#Handler.
		suite.WarnOnly(WarnDuplicates)
	}
}

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

// Logger returns a slog.Logger with the specified options.
func (suite *SlogTestSuite) Logger(options *slog.HandlerOptions) *slog.Logger {
	return slog.New(suite.Creator(suite.Buffer, options))
}
