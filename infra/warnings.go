package infra

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
	WarnNoReplAttr      = "HandlerOptions.ReplaceAttr not available"
	WarnNoReplAttrBasic = "HandlerOptions.ReplaceAttr not available for basic fields"
	WarnResolver        = "LogValuer objects are not resolved"
	WarnSkippingTest    = "Skipping test"
	WarnSourceKey       = "Source data not logged when AddSource flag set"
	WarnGroupEmpty      = "Empty (sub)group(s) logged"
	WarnUnused          = "Unused Warning(s)"
	WarnZeroPC          = "SourceKey logged for zero PC"
	WarnZeroTime        = "Zero time is logged"
)

// useWarnings is the flag value for enabling warnings instead of known errors.
// Command line setting:
//
//	go test ./... -args -useWarnings
//
// This flag will automatically set WarnLevelCase.
// Other behavior must be activated in specific handler test managers, for example:
//
//	sLogSuite := &test.SlogTestSuite{Creator: &SlogCreator{}}
//	sLogSuite.WarnOnly(test.WarnMessageKey)
//	suite.Run(t, slogSuite)
var useWarnings = flag.Bool("useWarnings", false, "Show warnings instead of known errors")

// WarningManager manages the warning set for a test run.
type WarningManager struct {
	// Name of Handler for warnings display.
	Name string

	warn     map[string]bool
	warnings map[string]*Warning
}

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

// managers captures all managers tested together into an array.
// This array is used when showing warnings.
var managers = make([]*WarningManager, 0)

func NewWarningManager(name string) *WarningManager {
	mgr := &WarningManager{Name: name}
	managers = append(managers, mgr)
	return mgr
}

// WarnOnly sets a flag to collect warnings instead of failing tests.
// The warn argument is one of the global constants beginning with 'Warn'.
func (wrnMgr *WarningManager) WarnOnly(warning string) {
	if wrnMgr.warn == nil {
		wrnMgr.warn = make(map[string]bool)
	}
	wrnMgr.warn[warning] = true
}

// GetWarnings returns an array of Warning records sorted by warn text.
// If there are no warnings the result array will be nil.
// Use this method if manual processing of warnings is required,
// otherwise use the WithWarnings method.
func (wrnMgr *WarningManager) GetWarnings() []*Warning {
	if wrnMgr.warnings == nil || len(wrnMgr.warnings) < 1 {
		return nil
	}
	if unused, found := wrnMgr.warnings[WarnUnused]; found {
		// Clean up WarnUnused warning instances.
		really := make([]WarningInstance, 0)
		for _, instance := range unused.Data {
			if _, found := wrnMgr.warnings[instance.Text]; !found {
				// OK, there are no such wrnMgr.
				really = append(really, instance)
			}
		}
		if len(really) > 0 {
			unused.Data = really
			unused.Count = uint(len(really))
		} else {
			delete(wrnMgr.warnings, WarnUnused)
		}
	}
	// Sort wrnMgr by warning string.
	warningStrings := make([]string, 0, len(wrnMgr.warnings))
	for warning := range wrnMgr.warnings {
		warningStrings = append(warningStrings, warning)
	}
	sort.Strings(warningStrings)
	w := make([]*Warning, len(warningStrings))
	for i, warning := range warningStrings {
		w[i] = wrnMgr.warnings[warning]
	}
	return w
}

// ShowWarnings prints any warnings to Stdout in a preformatted manner.
// Use the WarningManager method if more control over output is required.
//
// Note: Both Stdout and Stderr are captured by the the 'go test' command and
// shunted into Stdout (see https://pkg.go.dev/cmd/go#hdr-Test_packages).
// This output stream is only visible when the 'go test -v flag' is used.
func (wrnMgr *WarningManager) ShowWarnings(output io.Writer) {
	if output == nil {
		output = os.Stdout
	}
	warnings := wrnMgr.GetWarnings()
	if warnings != nil && len(warnings) > 0 {
		forHandler := ""
		if wrnMgr.Name != "" {
			forHandler = " for " + wrnMgr.Name
		}
		_, _ = fmt.Fprintf(output, "GetWarnings%s:\n", forHandler)
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
// This will cause the ShowWarnings method to be called on all test managers
// after all other output has been done, instead of buried in the middle.
// To use, add the following to a '_test' file:
//
//	func TestMain(m *testing.M) {
//	    test.WithWarnings(m)
//	}
//
// This step can be omitted if warnings are being sent to an output file.
//
// Note: The TestMain function can only be defined once in a package.
// If multiple SlogTestSuite instances are created in separate files in
// the same package, TestMain can be moved into a single main_test.go file
// as is done in the go-slog/verify package.
func WithWarnings(m *testing.M) {
	flag.Parse()
	exitVal := m.Run()
	for _, testSuite := range managers {
		testSuite.ShowWarnings(nil)
	}
	os.Exit(exitVal)
}

// AddWarning to results list, specifying warning string and optional extra text.
// If the addLogRecord flag is true the current log record JSON is also stored.
// The current function name is acquired from the CurrentFunctionName() and stored.
func (wrnMgr *WarningManager) AddWarning(warning string, text string, logRecordJSON string) {
	if wrnMgr.warnings == nil {
		wrnMgr.warnings = make(map[string]*Warning)
	}
	record, found := wrnMgr.warnings[warning]
	if !found {
		record = &Warning{Name: warning}
		wrnMgr.warnings[warning] = record
	}
	record.Count++
	if record.Data == nil {
		record.Data = make([]WarningInstance, 0)
	}
	instance := WarningInstance{
		Function: CurrentFunctionName(),
		Text:     text,
	}
	if logRecordJSON != "" {
		instance.Record = strings.TrimRight(logRecordJSON, "\n")
	}
	record.Data = append(record.Data, instance)
}

// HasWarning returns true if the specified warning has been set in the test suite.
func (wrnMgr *WarningManager) HasWarning(warning string) bool {
	return *useWarnings && wrnMgr.warn[warning]
}

// HasWarnings checks all specified warnings and returns an array of
// any that have been set in the test suite in the same order.
// If none are found an empty array is returned.
func (wrnMgr *WarningManager) HasWarnings(warnings ...string) []string {
	found := make([]string, 0, len(warnings))
	if *useWarnings {
		for _, warning := range warnings {
			if wrnMgr.HasWarning(warning) {
				found = append(found, warning)
			}
		}
	}
	return found
}

// skipTest adds warnings for a test that is being skipped.
// The first warning is for skipping a test with the text set to the 'because' warning argument.
// The second warning is for the 'because' warning with the text set to skipping the test.
func (wrnMgr *WarningManager) skipTest(because string) {
	wrnMgr.AddWarning(WarnSkippingTest, because, "")
	wrnMgr.AddWarning(because, WarnSkippingTest, "")
}

// skipTestIf checks the warnings provided to see if any have been set in the suite,
// adding skipTest warnings for the first one and returning true.
// False is returned if none of the warnings are found.
func (wrnMgr *WarningManager) skipTestIf(warnings ...string) bool {
	for _, warning := range warnings {
		if wrnMgr.warn[warning] {
			wrnMgr.AddWarning(WarnSkippingTest, warning, "")
			wrnMgr.AddWarning(warning, WarnSkippingTest, "")
			return true
		}
	}
	return false
}
