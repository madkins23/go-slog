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

	"github.com/madkins23/go-slog/warning"
)

// -----------------------------------------------------------------------------
// Warnings mechanism to trade test failure for warning list at end of tests.

// useWarnings is the flag value for enabling warning instead of known errors.
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
var useWarnings = flag.Bool("useWarnings", false, "Show warning instead of known errors")

// -----------------------------------------------------------------------------

// WarningManager manages the warning set for a test run.
type WarningManager struct {
	// Name of Handler for warning display.
	Name string

	fnPrefix   string
	showPrefix string
	predefined map[string]*warning.Warning
	warnOnly   map[string]bool
	warnings   map[string]*Warnings
}

// Warnings gathers instances for a specific Warning.
type Warnings struct {
	// Level is the warning level.
	Level warning.Level

	// Name of warning.
	Name string

	// Description of warning.
	Description string

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
// This array is used when showing warning.
var managers = make([]*WarningManager, 0)

func NewWarningManager(name string, fnPrefix string, showPrefix string) *WarningManager {
	mgr := &WarningManager{
		Name:       name,
		fnPrefix:   fnPrefix,
		showPrefix: showPrefix,
	}
	managers = append(managers, mgr)
	mgr.Predefine(warning.Administrative()...)
	return mgr
}

// Predefine warning that can be referenced during testing.
func (wrnMgr *WarningManager) Predefine(warnings ...*warning.Warning) {
	if wrnMgr.predefined == nil {
		wrnMgr.predefined = make(map[string]*warning.Warning, len(warnings))
	}
	for _, w := range warnings {
		wrnMgr.predefined[w.Name] = w
	}
}

// WarnOnly sets a flag to collect warning instead of failing tests.
// The warning argument is one of the global constants beginning with 'Warn'
// and it must be predefined to the manager.
func (wrnMgr *WarningManager) WarnOnly(w *warning.Warning) {
	if _, found := wrnMgr.predefined[w.Name]; !found {
		wrnMgr.AddWarning(warning.Undefined, w.Description, "")
		slog.Warn("Undefined warning", "warning", w.Description)
	}
	if wrnMgr.warnOnly == nil {
		wrnMgr.warnOnly = make(map[string]bool)
	}
	wrnMgr.warnOnly[w.Name] = true
}

// -----------------------------------------------------------------------------
// Calls to be made during testing.

// AddUnused adds a Unused warning to the results list.
// The warning added is Unused and the extra text is the name of the specified warning.
func (wrnMgr *WarningManager) AddUnused(w *warning.Warning, logRecordJSON string) {
	wrnMgr.AddWarning(warning.Unused, w.Name, logRecordJSON)
}

// AddWarning to results list, specifying warning string and optional extra text.
// If the addLogRecord flag is true the current log record JSON is also stored.
// The current function name is acquired from the CurrentFunctionName() and stored.
func (wrnMgr *WarningManager) AddWarning(w *warning.Warning, text string, logRecordJSON string) {
	wrnMgr.addWarning(w, CurrentFunctionName(wrnMgr.fnPrefix), text, logRecordJSON)
}

// AddWarningFn to results list, specifying warning string and function name.
// If the addLogRecord flag is true the current log record JSON is also stored.
// The current function name is acquired from the CurrentFunctionName() and stored.
func (wrnMgr *WarningManager) AddWarningFn(w *warning.Warning, fnName string, logRecordJSON string) {
	wrnMgr.addWarning(w, fnName, "", logRecordJSON)
}

// addWarning to results list, specifying warning string, function name, and optional extra text.
// If the addLogRecord flag is true the current log record JSON is also stored.
// The current function name is acquired from the CurrentFunctionName() and stored.
func (wrnMgr *WarningManager) addWarning(w *warning.Warning, fnName string, text string, logRecordJSON string) {
	if wrnMgr.warnings == nil {
		wrnMgr.warnings = make(map[string]*Warnings)
	}
	record, found := wrnMgr.warnings[w.Name]
	if !found {
		record = &Warnings{
			Level:       w.Level,
			Name:        w.Name,
			Description: w.Description,
		}
		wrnMgr.warnings[w.Name] = record
	}
	record.Count++
	if record.Data == nil {
		record.Data = make([]WarningInstance, 0)
	}
	instance := WarningInstance{
		Function: fnName,
		Text:     text,
	}
	if logRecordJSON != "" {
		instance.Record = strings.TrimRight(logRecordJSON, "\n")
	}
	record.Data = append(record.Data, instance)
}

// HasWarning returns true if the specified warning has been set in the test suite.
func (wrnMgr *WarningManager) HasWarning(w *warning.Warning) bool {
	return *useWarnings && wrnMgr.warnOnly[w.Name]
}

// HasWarnings checks all specified warning and returns an array of
// any that have been set in the test suite in the same order.
// If none are found an empty array is returned.
func (wrnMgr *WarningManager) HasWarnings(warnings ...*warning.Warning) []*warning.Warning {
	found := make([]*warning.Warning, 0, len(warnings))
	if *useWarnings {
		for _, w := range warnings {
			if wrnMgr.HasWarning(w) {
				found = append(found, w)
			}
		}
	}
	return found
}

// SkipTest adds warning for a test that is being skipped.
// The first warning is for skipping a test with the text set to the 'because' warning argument.
// The second warning is for the 'because' warning with the text set to skipping the test.
func (wrnMgr *WarningManager) SkipTest(because *warning.Warning) {
	wrnMgr.AddWarning(warning.SkippingTest, because.Description, "")
	wrnMgr.AddWarning(because, warning.SkippingTest.Description, "")
}

// SkipTestIf checks the warning provided to see if any have been set in the suite,
// adding SkipTest warning for the first one and returning true.
// False is returned if none of the warning are found.
func (wrnMgr *WarningManager) SkipTestIf(warns ...*warning.Warning) bool {
	for _, w := range warns {
		if wrnMgr.warnOnly[w.Name] {
			wrnMgr.AddWarning(warning.SkippingTest, w.Description, "")
			wrnMgr.AddWarning(w, warning.SkippingTest.Description, "")
			return true
		}
	}
	return false
}

// -----------------------------------------------------------------------------
// Display warning at end of testing.

// GetWarnings returns an array of Warnings records sorted by warning level and text.
// If there are no warning the result array will be nil.
// Use this method if manual processing of warning is required,
// otherwise use the WithWarnings method.
func (wrnMgr *WarningManager) GetWarnings() []*Warnings {
	if wrnMgr.warnings == nil || len(wrnMgr.warnings) < 1 {
		return nil
	}

	if unused, found := wrnMgr.warnings[warning.Unused.Name]; found {
		// Clean up Unused warning instances.
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
			delete(wrnMgr.warnings, warning.Unused.Name)
		}
	}

	// Sort warning by warning level and string.
	warningStrings := make([]string, 0, len(wrnMgr.warnings))
	for w := range wrnMgr.warnings {
		warningStrings = append(warningStrings, w)
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

	warnings := make([]*Warnings, len(warningStrings))
	for i, w := range warningStrings {
		warnings[i] = wrnMgr.warnings[w]
	}
	return warnings
}

// Track handlers that invoke warning for use in TestMain.
var byWarning = make(map[*warning.Warning]map[string]bool)

// ShowWarnings prints any warning to Stdout in a preformatted manner.
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
		warningTree := make(map[warning.Level][]*Warnings)
		for _, w := range warnings {
			list, found := warningTree[w.Level]
			if !found {
				list = make([]*Warnings, 0)
			}
			warningTree[w.Level] = append(list, w)
		}
		for _, list := range warningTree {
			sort.Slice(list, func(i, j int) bool {
				return list[i].Name < list[j].Name
			})
		}
		_, _ = fmt.Fprintf(output, "\n%sWarnings%s:\n", wrnMgr.showPrefix, forHandler)
		for _, level := range warning.LevelOrder {
			if list, ok := warningTree[level]; ok {
				_, _ = fmt.Fprintf(output, "%s  %s\n", wrnMgr.showPrefix, level.String())
				for _, w := range list {
					// Track handlers that invoke warning for use in TestMain.
					warn := wrnMgr.predefined[w.Name]
					if byWarning[warn] == nil {
						byWarning[warn] = make(map[string]bool)
					}
					byWarning[warn][wrnMgr.Name] = true

					_, _ = fmt.Fprintf(output, "%s  %4d [%s] %s\n", wrnMgr.showPrefix, w.Count, w.Name, w.Description)
					for _, data := range w.Data {
						text := data.Function
						if data.Text != "" {
							text += ": " + data.Text
						}
						_, _ = fmt.Fprintf(output, "%s         %s\n", wrnMgr.showPrefix, text)
						if data.Record != "" {
							_, _ = fmt.Fprintf(output, "%s           %s\n", wrnMgr.showPrefix, data.Record)
						}
					}
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
// This step can be omitted if warning are being sent to an output file.
//
// If multiple SlogTestSuite instances are defined in separate files in the same package:
//   - an addition list of warning and which handlers throw them will be shown and
//   - the TestMain function must be moved to a separate file, as it can only be defined once.
func WithWarnings(m *testing.M) {
	flag.Parse()
	exitVal := m.Run()

	var showPrefix string
	for _, manager := range managers {
		manager.ShowWarnings(nil)
		if showPrefix == "" {
			showPrefix = manager.showPrefix
		}
	}

	if len(managers) > 1 {
		ShowHandlersByWarning(showPrefix)
	}

	os.Exit(exitVal)
}

// ShowHandlersByWarning uses the global byWarning map to
// show the handlers that issue each warning.
func ShowHandlersByWarning(showPrefix string) {
	fmt.Printf("%s\n%s Handlers by warning:\n", showPrefix, showPrefix)
	byLevel := make(map[warning.Level][]*warning.Warning)
	byName := make(map[string]*warning.Warning)
	for w := range byWarning {
		if byLevel[w.Level] == nil {
			byLevel[w.Level] = make([]*warning.Warning, 0)
		}
		byLevel[w.Level] = append(byLevel[w.Level], w)
		byName[w.Name] = w
	}
	for _, level := range warning.LevelOrder {
		if warnings, found := byLevel[level]; found {
			fmt.Printf("%s  %s\n", showPrefix, level.String())
			names := make([]string, 0, len(warnings))
			for _, w := range warnings {
				names = append(names, w.Name)
			}
			sort.Strings(names)
			for _, name := range names {
				fmt.Printf("%s    [%s] %s\n", showPrefix, name, byName[name].Description)
				handlers := byWarning[byName[name]]
				hdlrNames := make([]string, 0, len(handlers))
				for handler := range handlers {
					hdlrNames = append(hdlrNames, handler)
				}
				for _, hdlrName := range hdlrNames {
					fmt.Printf("%s      %s\n", showPrefix, hdlrName)
				}
			}
		}
	}
}
