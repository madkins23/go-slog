package test

import (
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"testing"
)

// -----------------------------------------------------------------------------
// Warning mechanism to trade test failure for warning list at end of tests.

const (
	WarnDefaultLevel    = "Handler doesn't default to slog.LevelInfo"
	WarnDuplicates      = "Duplicate field(s) found"
	WarnEmptyAttributes = "Empty attribute(s) logged (\"\":null)"
	WarnGroupInline     = "Group with empty key does not inline subfields"
	WarnLevelCase       = "Log level in lowercase"
	WarnMessageKey      = "Wrong message key (should be 'msg')"
	WarnNanoDuration    = "slog.Duration() doesn't log nanoseconds"
	WarnNanoTime        = "slog.Time() doesn't log nanoseconds"
	WarnNoReplAttr      = "HandlerOptions.ReplAttr not available"
	WarnNoReplAttrBasic = "HandlerOptions.ReplAttr not available for basic fields"
	WarnResolver        = "LogValuer objects are not resolved"
	WarnSourceKey       = "Source data not logged when AddSource flag set"
	WarnSubgroupEmpty   = "Empty subgroup(s) logged"
	WarnUnused          = "Unused Warning(s)"
	WarnZeroPC          = "SourceKey logged for zero PC"
	WarnZeroTime        = "Zero time is logged"
)

// UseWarnings is the flag value for enabling warnings instead of known errors.
// Command line setting:
//
//	go test ./... -args -useWarnings
//
// This flag will automatically set WarnLevelCase.
// Other behavior must be activated in specific handler test suites, for example:
//
//	sLogSuite := &test.SlogTestSuite{Creator: &SlogCreator{}}
//	if *test.UseWarnings {
//		sLogSuite.WarnOnly(test.WarnMessageKey)
//	}
//	suite.Run(t, slogSuite)
var UseWarnings = flag.Bool("useWarnings", false, "Show warnings instead of known errors")

// Warning encapsulates data from non-error warnings.
type Warning struct {
	// Name of warning.
	Name string

	// Count of times warning is issued.
	Count uint

	// Data associated with the specific instances of the warning, if any.
	Data []WarningInstance
}

// WarningInstance encapsulates data for a specific warning instance.
type WarningInstance struct {
	Function string
	Record   string
	Text     string
}

// WarnOnly sets a flag to collect warnings instead of failing tests.
// The warn argument is one of the global constants beginning with 'Warn'.
func (suite *SlogTestSuite) WarnOnly(warning string) {
	if suite.warn == nil {
		suite.warn = make(map[string]bool)
	}
	suite.warn[warning] = true
}

// Warnings returns an array of Warning records sorted by warn text.
// If there are no warnings the result array will be nil.
// Use this method if manual processing of warnings is required,
// otherwise use the WithWarnings method.
func (suite *SlogTestSuite) Warnings() []*Warning {
	if suite.warnings == nil || len(suite.warnings) < 1 {
		return nil
	}
	if unused, found := suite.warnings[WarnUnused]; found {
		// Clean up WarnUnused warning instances.
		really := make([]WarningInstance, 0)
		for _, instance := range unused.Data {
			if _, found := suite.warnings[instance.Text]; !found {
				// OK, there are no such warnings.
				really = append(really, instance)
			}
		}
		if len(really) > 0 {
			unused.Data = really
			unused.Count = uint(len(really))
		} else {
			delete(suite.warnings, WarnUnused)
		}
	}
	// Sort warnings by warning string.
	warningStrings := make([]string, 0, len(suite.warnings))
	for warning := range suite.warnings {
		warningStrings = append(warningStrings, warning)
	}
	sort.Strings(warningStrings)
	warnings := make([]*Warning, len(warningStrings))
	for i, warning := range warningStrings {
		warnings[i] = suite.warnings[warning]
	}
	return warnings
}

// ShowWarnings prints any warnings to Stdout in a preformatted manner.
// Use the Warnings method if more control over output is required.
// Note: Both Stdout and Stderr are captured by the the 'go test' command and
// shunted into Stdout (see https://pkg.go.dev/cmd/go#hdr-Test_packages).
// This output stream is only visible when the 'go test -v flag' is used.
func (suite *SlogTestSuite) ShowWarnings(output io.Writer) {
	if output == nil {
		output = os.Stdout
	}
	warnings := suite.Warnings()
	if warnings != nil && len(warnings) > 0 {
		forHandler := ""
		if suite.Name != "" {
			forHandler = " for " + suite.Name
		}
		_, _ = fmt.Fprintf(output, "Warnings%s:\n", forHandler)
		for _, warning := range warnings {
			_, _ = fmt.Fprintf(output, "  %4d %s\n", warning.Count, warning.Name)
			for _, data := range warning.Data {
				text := data.Function
				if data.Text != "" {
					text += ": " + data.Text
				}
				_, _ = fmt.Fprintf(output, "       %s\n", text)
				if data.Record != "" {
					_, _ = fmt.Fprintf(output, "         %s\n", data.Record)
				}
			}
		}
	}
}

// WithWarnings implements the guts of TestMain (see https://pkg.go.dev/testing#hdr-Main).
// This will cause the ShowWarnings method to be called on all test suites
// after all other output has been done, instead of buried in the middle.
// To use, add the following to a '_test' file:
//
//	func TestMain(m *testing.M) {
//	    test.WithWarnings(m)
//	}
//
// This step can be omitted if warnings are being sent to an output file.
// Note: The TestMain function can only be defined once in a package.
// If multiple SlogTestSuite instances are created in separate files in
// the same package, TestMain can be moved into a single main_test.go file
// as is done in the go-slog/verify package.
func WithWarnings(m *testing.M) {
	flag.Parse()
	exitVal := m.Run()
	for _, testSuite := range suites {
		testSuite.ShowWarnings(nil)
	}
	os.Exit(exitVal)
}

func (suite *SlogTestSuite) addWarning(warning string, text string, addLogRecord bool) {
	if suite.warnings == nil {
		suite.warnings = make(map[string]*Warning)
	}
	record, found := suite.warnings[warning]
	if !found {
		record = &Warning{Name: warning}
		suite.warnings[warning] = record
	}
	record.Count++
	if record.Data == nil {
		record.Data = make([]WarningInstance, 0)
	}
	instance := WarningInstance{
		Function: currentFunctionName(),
		Text:     text,
	}
	if addLogRecord {
		instance.Record = strings.TrimRight(suite.Buffer.String(), "\n")
	}
	record.Data = append(record.Data, instance)
}

func (suite *SlogTestSuite) hasWarning(warning string) bool {
	return suite.warn[warning]
}