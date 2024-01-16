package infra

import (
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"sort"
	"strings"
	"testing"
)

// -----------------------------------------------------------------------------
// Warnings mechanism to trade test failure for warning list at end of tests.

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

// WarningLevel for warnings used mainly to organize warnings on output.
type WarningLevel uint

const (
	warnLevelUnused WarningLevel = iota
	WarnLevelAdmin
	WarnLevelSuggested
	WarnLevelImplied
	WarnLevelRequired
)

var warningLevelNames = map[WarningLevel]string{
	WarnLevelAdmin:     "Administrative",
	WarnLevelSuggested: "Suggested",
	WarnLevelImplied:   "Implied",
	WarnLevelRequired:  "Required",
}

// Warning definition.
type Warning struct {
	// Level is the warning level.
	Level WarningLevel

	// Name of the warning.
	Name string
}

var (
	WarnSkippingTest = &Warning{
		Level: WarnLevelAdmin,
		Name:  "Skipping test",
	}
	WarnUndefined = &Warning{
		Level: WarnLevelAdmin,
		Name:  "Undefined Warnings(s)",
	}
	WarnUnused = &Warning{
		Level: WarnLevelAdmin,
		Name:  "Unused Warnings(s)",
	}
)

// WarningManager manages the warning set for a test run.
type WarningManager struct {
	// Name of Handler for warnings display.
	Name string

	predefined map[string]*Warning
	warnOnly   map[string]bool
	warnings   map[string]*Warnings
}

// Warnings gathers instances for a specific Warning.
type Warnings struct {
	// Level is the warning level.
	Level WarningLevel

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
	mgr.Predefine(WarnSkippingTest, WarnUndefined, WarnUnused)
	return mgr
}

func (wrnMgr *WarningManager) Predefine(warnings ...*Warning) {
	if wrnMgr.predefined == nil {
		wrnMgr.predefined = make(map[string]*Warning, len(warnings))
	}
	for _, warning := range warnings {
		wrnMgr.predefined[warning.Name] = warning
	}
}

// WarnOnly sets a flag to collect warnings instead of failing tests.
// The warning argument is one of the global constants beginning with 'Warn'
// and it must be predefined to the manager.
func (wrnMgr *WarningManager) WarnOnly(warning *Warning) {
	if _, found := wrnMgr.predefined[warning.Name]; !found {
		wrnMgr.AddWarning(WarnUndefined, warning.Name, "")
		slog.Warn("Undefined warning '%s'", warning.Name)
	}
	if wrnMgr.warnOnly == nil {
		wrnMgr.warnOnly = make(map[string]bool)
	}
	wrnMgr.warnOnly[warning.Name] = true
}

// GetWarnings returns an array of Warnings records sorted by warning level and text.
// If there are no warnings the result array will be nil.
// Use this method if manual processing of warnings is required,
// otherwise use the WithWarnings method.
func (wrnMgr *WarningManager) GetWarnings() []*Warnings {
	if wrnMgr.warnings == nil || len(wrnMgr.warnings) < 1 {
		return nil
	}

	if unused, found := wrnMgr.warnings[WarnUnused.Name]; found {
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
			delete(wrnMgr.warnings, WarnUnused.Name)
		}
	}

	// Sort warnings by warning level and string.
	warningStrings := make([]string, 0, len(wrnMgr.warnings))
	for warning := range wrnMgr.warnings {
		warningStrings = append(warningStrings, warning)
	}
	sort.Slice(warningStrings, func(i, j int) bool {
		iWarning := wrnMgr.warnings[warningStrings[i]]
		jWarning := wrnMgr.warnings[warningStrings[j]]
		if iWarning.Level > jWarning.Level {
			return true
		} else if iWarning.Level < jWarning.Level {
			return false
		}
		return iWarning.Name < jWarning.Name
	})

	w := make([]*Warnings, len(warningStrings))
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
		var warningLevel = warnLevelUnused
		_, _ = fmt.Fprintf(output, "Warnings%s:\n", forHandler)
		for _, warning := range warnings {
			if warning.Level != warningLevel {
				warningLevel = warning.Level
				levelName, found := warningLevelNames[warningLevel]
				if !found {
					levelName = fmt.Sprintf("Unknown level %d", warningLevel)
				}
				_, _ = fmt.Fprintf(output, "  %s\n", levelName)
			}
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
		_, _ = fmt.Fprintln(output)
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

// AddUnused adds a WarnUnused warning to the results list.
// The warning added is WarnUnused and the extra text is the name of the specified warning.
func (wrnMgr *WarningManager) AddUnused(warning *Warning, logRecordJSON string) {
	wrnMgr.AddWarning(WarnUnused, warning.Name, logRecordJSON)
}

// AddWarning to results list, specifying warning string and optional extra text.
// If the addLogRecord flag is true the current log record JSON is also stored.
// The current function name is acquired from the CurrentFunctionName() and stored.
func (wrnMgr *WarningManager) AddWarning(warning *Warning, text string, logRecordJSON string) {
	if wrnMgr.warnings == nil {
		wrnMgr.warnings = make(map[string]*Warnings)
	}
	record, found := wrnMgr.warnings[warning.Name]
	if !found {
		record = &Warnings{Name: warning.Name, Level: warning.Level}
		wrnMgr.warnings[warning.Name] = record
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
func (wrnMgr *WarningManager) HasWarning(warning *Warning) bool {
	return *useWarnings && wrnMgr.warnOnly[warning.Name]
}

// HasWarnings checks all specified warnings and returns an array of
// any that have been set in the test suite in the same order.
// If none are found an empty array is returned.
func (wrnMgr *WarningManager) HasWarnings(warnings ...*Warning) []*Warning {
	found := make([]*Warning, 0, len(warnings))
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
func (wrnMgr *WarningManager) skipTest(because *Warning) {
	wrnMgr.AddWarning(WarnSkippingTest, because.Name, "")
	wrnMgr.AddWarning(because, WarnSkippingTest.Name, "")
}

// skipTestIf checks the warnings provided to see if any have been set in the suite,
// adding skipTest warnings for the first one and returning true.
// False is returned if none of the warnings are found.
func (wrnMgr *WarningManager) skipTestIf(warnings ...*Warning) bool {
	for _, warning := range warnings {
		if wrnMgr.warnOnly[warning.Name] {
			wrnMgr.AddWarning(WarnSkippingTest, warning.Name, "")
			wrnMgr.AddWarning(warning, WarnSkippingTest.Name, "")
			return true
		}
	}
	return false
}
