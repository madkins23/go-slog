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

// -----------------------------------------------------------------------------

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

// -----------------------------------------------------------------------------
// Calls made during context-specific setup of warning manager.

// managers captures all managers tested together into an array.
// This array is used when showing warnings.
var managers = make([]*WarningManager, 0)

func NewWarningManager(name string) *WarningManager {
	mgr := &WarningManager{Name: name}
	managers = append(managers, mgr)
	mgr.Predefine(WarnSkippingTest, WarnUndefined, WarnUnused)
	return mgr
}

// Predefine warnings that can be referenced during testing.
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
		slog.Warn("Undefined warning", "warning", warning.Name)
	}
	if wrnMgr.warnOnly == nil {
		wrnMgr.warnOnly = make(map[string]bool)
	}
	wrnMgr.warnOnly[warning.Name] = true
}

// -----------------------------------------------------------------------------
// Calls made during testing.

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

// SkipTest adds warnings for a test that is being skipped.
// The first warning is for skipping a test with the text set to the 'because' warning argument.
// The second warning is for the 'because' warning with the text set to skipping the test.
func (wrnMgr *WarningManager) SkipTest(because *Warning) {
	wrnMgr.AddWarning(WarnSkippingTest, because.Name, "")
	wrnMgr.AddWarning(because, WarnSkippingTest.Name, "")
}

// SkipTestIf checks the warnings provided to see if any have been set in the suite,
// adding skipTest warnings for the first one and returning true.
// False is returned if none of the warnings are found.
func (wrnMgr *WarningManager) SkipTestIf(warnings ...*Warning) bool {
	for _, warning := range warnings {
		if wrnMgr.warnOnly[warning.Name] {
			wrnMgr.AddWarning(WarnSkippingTest, warning.Name, "")
			wrnMgr.AddWarning(warning, WarnSkippingTest.Name, "")
			return true
		}
	}
	return false
}

// -----------------------------------------------------------------------------
// Result display functionality.

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

// Track handlers that invoke warnings for use in TestMain.
var byWarning = make(map[*Warning]map[string]bool)

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

	forHandler := ""
	if wrnMgr.Name != "" {
		forHandler = " for " + wrnMgr.Name
	}

	warnings := wrnMgr.GetWarnings()
	if warnings != nil && len(warnings) > 0 {
		// Warnings grouped by level.
		var warningLevel = warnLevelUnused
		_, _ = fmt.Fprintf(output, "\nWarnings%s:\n", forHandler)
		for _, warning := range warnings {
			// Track handlers that invoke warnings for use in TestMain.
			warn := wrnMgr.predefined[warning.Name]
			if byWarning[warn] == nil {
				byWarning[warn] = make(map[string]bool)
			}
			byWarning[warn][wrnMgr.Name] = true

			// Show warnings.
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
// If multiple SlogTestSuite instances are defined in separate files in the same package:
//   - an addition list of warnings and which handlers throw them will be shown and
//   - the TestMain function must be moved to a separate file, as it can only be defined once.
func WithWarnings(m *testing.M) {
	flag.Parse()
	exitVal := m.Run()

	for _, manager := range managers {
		manager.ShowWarnings(nil)
	}

	if len(managers) > 1 {
		ShowHandlersByWarning()
	}

	os.Exit(exitVal)
}

// ShowHandlersByWarning uses the global byWarning map to
// show the handlers that issue each warning.
func ShowHandlersByWarning() {
	fmt.Printf("\nHandlers by warning:\n")
	byLevel := make(map[WarningLevel][]*Warning)
	byName := make(map[string]*Warning)
	for warning := range byWarning {
		if byLevel[warning.Level] == nil {
			byLevel[warning.Level] = make([]*Warning, 0)
		}
		byLevel[warning.Level] = append(byLevel[warning.Level], warning)
		byName[warning.Name] = warning
	}

	for level := WarnLevelRequired; level >= WarnLevelAdmin; level-- {
		if warnings, found := byLevel[level]; found {
			fmt.Printf("  %s\n", warningLevelNames[level])
			names := make([]string, 0, len(warnings))
			for _, warning := range warnings {
				names = append(names, warning.Name)
			}
			sort.Strings(names)
			for _, name := range names {
				fmt.Printf("    %s\n", name)
				handlers := byWarning[byName[name]]
				hdlrNames := make([]string, 0, len(handlers))
				for handler := range handlers {
					hdlrNames = append(hdlrNames, handler)
				}
				for _, hdlrName := range hdlrNames {
					fmt.Printf("      %s\n", hdlrName)
				}
			}
		}
	}
}
