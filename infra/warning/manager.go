package warning

import (
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/madkins23/go-slog/internal/misc"
)

// -----------------------------------------------------------------------------
// Warnings mechanism to trade test failure for warning list at end of tests.

// useWarnings is the flag value for enabling warning instead of known errors.
// The default value for this flag is true.
// In the rare event that warnings should be disabled
// (resulting in actual errors in the test harness and no warning results) use:
//
//	go test ./... -args -useWarnings=false
//
// This flag will automatically set WarnLevelCase.
// Other behavior must be activated in specific handler test managers, for example:
//
//	sLogSuite := &test.SlogTestSuite{Creator: &SlogCreator{}}
//	sLogSuite.WarnOnly(test.WarnMessageKey)
//	suite.Run(t, slogSuite)
var useWarnings = flag.Bool("useWarnings", true, "Show warning instead of known errors")

// -----------------------------------------------------------------------------

// Manager manages the warning set for a test run.
type Manager struct {
	// Name of Handler for warning display.
	Name string

	fnPrefix   string
	showPrefix string
	predefined map[string]*Warning
	warnOnly   map[string]bool
	warnings   map[string]*Instances
}

// Instances gathers instances for a specific Warning.
type Instances struct {
	// Level is the warning level.
	Level Level

	// Name of warning.
	Name string

	// Summary of warning.
	Summary string

	// Count of times warning is issued.
	Count uint

	// Data associated with the specific instances of the warning, if any.
	Data []Instance
}

// Instance encapsulates data for a specific warning instance.
type Instance struct {
	Function string
	Record   string
	Text     string
}

// -----------------------------------------------------------------------------
// Calls made during context-specific setup of warning manager.

// managers captures all managers tested together into a map by manager name.
// Due to the way benchmarking works there will be more than one with the same name,
// but this doesn't matter as they should all be the same.
// This map is used when showing data at the end.
var managers = make(map[string]*Manager)

func NewWarningManager(name string, fnPrefix string, showPrefix string) *Manager {
	mgr := &Manager{
		Name:       name,
		fnPrefix:   fnPrefix,
		showPrefix: showPrefix,
	}
	mgr.Predefine(Administrative()...)
	managers[name] = mgr // Overwrite each time per above comment.
	return mgr
}

// Predefine warning that can be referenced during testing.
func (mgr *Manager) Predefine(warnings ...*Warning) {
	if mgr.predefined == nil {
		mgr.predefined = make(map[string]*Warning, len(warnings))
	}
	for _, w := range warnings {
		mgr.predefined[w.Name] = w
	}
}

// WarnOnly sets a flag to collect warning instead of failing tests.
// The warning argument is one of the global constants beginning with 'Warn'
// and it must be predefined to the manager.
func (mgr *Manager) WarnOnly(w *Warning) {
	if _, found := mgr.predefined[w.Name]; !found {
		mgr.AddWarning(Undefined, w.Summary, "")
		slog.Warn("Undefined warning", "warning", w.Summary)
	}
	if mgr.warnOnly == nil {
		mgr.warnOnly = make(map[string]bool)
	}
	mgr.warnOnly[w.Name] = true
}

// -----------------------------------------------------------------------------
// Calls to be made during testing.

// AddUnused adds an Unused warning to the results list.
// The warning added is Unused and the extra text is the name of the specified warning.
func (mgr *Manager) AddUnused(w *Warning, logRecordJSON string) {
	mgr.AddWarning(Unused, w.Name, logRecordJSON)
}

// AddWarning to results list, specifying warning string and optional extra text.
// If the addLogRecord flag is true the current log record JSON is also stored.
// The current function name is acquired from the currentFunctionName() and stored.
func (mgr *Manager) AddWarning(w *Warning, text string, logRecordJSON string) {
	mgr.AddWarningFnText(w, misc.CurrentFunctionName(mgr.fnPrefix), text, logRecordJSON)
}

// AddWarningFn adds a warning to the results list, specifying warning string and function name.
// If the addLogRecord flag is true the current log record JSON is also stored.
// The current function name is acquired from the currentFunctionName() and stored.
func (mgr *Manager) AddWarningFn(w *Warning, fnName string, logRecordJSON string) {
	mgr.AddWarningFnText(w, fnName, "", logRecordJSON)
}

// AddWarningFnText to results list, specifying warning string, function name, and optional extra text.
// If the addLogRecord flag is true the current log record JSON is also stored.
// The current function name is acquired from the currentFunctionName() and stored.
func (mgr *Manager) AddWarningFnText(w *Warning, fnName string, text string, logRecordJSON string) {
	if mgr.warnings == nil {
		mgr.warnings = make(map[string]*Instances)
	}
	record, found := mgr.warnings[w.Name]
	if !found {
		record = &Instances{
			Level:   w.Level,
			Name:    w.Name,
			Summary: w.Summary,
		}
		mgr.warnings[w.Name] = record
	}
	record.Count++
	if record.Data == nil {
		record.Data = make([]Instance, 0)
	}
	instance := Instance{
		Function: fnName,
		Text:     text,
	}
	if logRecordJSON != "" {
		instance.Record = strings.TrimRight(logRecordJSON, "\n")
	}
	record.Data = append(record.Data, instance)
}

// HasWarning returns true if the specified warning has been set in the test suite.
func (mgr *Manager) HasWarning(w *Warning) bool {
	return *useWarnings && mgr.warnOnly[w.Name]
}

// HasWarnings checks all specified warning and returns an array of
// any that have been set in the test suite in the same order.
// If none are found an empty array is returned.
func (mgr *Manager) HasWarnings(warnings ...*Warning) []*Warning {
	found := make([]*Warning, 0, len(warnings))
	if *useWarnings {
		for _, w := range warnings {
			if mgr.HasWarning(w) {
				found = append(found, w)
			}
		}
	}
	return found
}

// SkipTest adds warning for a test that is being skipped.
// The first warning is for skipping a test with the text set to the 'because' warning argument.
// The second warning is for the 'because' warning with the text set to skipping the test.
func (mgr *Manager) SkipTest(because *Warning) {
	mgr.AddWarning(SkippingTest, because.Summary, "")
	mgr.AddWarning(because, SkippingTest.Summary, "")
}

// SkipTestIf checks the warning provided to see if any have been set in the suite,
// adding SkipTest warning for the first one and returning true.
// False is returned if none of the warning are found.
func (mgr *Manager) SkipTestIf(warns ...*Warning) bool {
	for _, w := range warns {
		if mgr.warnOnly[w.Name] {
			mgr.AddWarning(SkippingTest, w.Summary, "")
			mgr.AddWarning(w, SkippingTest.Summary, "")
			return true
		}
	}
	return false
}

// -----------------------------------------------------------------------------
// Display warning at end of testing.

// GetWarnings returns an array of Instances records sorted by warning level and text.
// If there are no warning the result array will be nil.
// Use this method if manual processing of warning is required,
// otherwise use the WithWarnings method.
func (mgr *Manager) GetWarnings() []*Instances {
	if mgr.warnings == nil || len(mgr.warnings) < 1 {
		return nil
	}

	if unused, found := mgr.warnings[Unused.Name]; found {
		// Clean up Unused warning instances.
		really := make([]Instance, 0)
		for _, instance := range unused.Data {
			if _, found := mgr.warnings[instance.Text]; !found {
				// OK, there are no such mgr.
				really = append(really, instance)
			}
		}
		if len(really) > 0 {
			unused.Data = really
			unused.Count = uint(len(really))
		} else {
			delete(mgr.warnings, Unused.Name)
		}
	}

	// Sort warning by warning level and string.
	warningStrings := make([]string, 0, len(mgr.warnings))
	for w := range mgr.warnings {
		warningStrings = append(warningStrings, w)
	}
	sort.Slice(warningStrings, func(i, j int) bool {
		iWarning := mgr.warnings[warningStrings[i]]
		jWarning := mgr.warnings[warningStrings[j]]
		if iWarning.Level > jWarning.Level {
			return true
		} else if iWarning.Level < jWarning.Level {
			return false
		}
		return iWarning.Name < jWarning.Name
	})

	warnings := make([]*Instances, len(warningStrings))
	for i, w := range warningStrings {
		warnings[i] = mgr.warnings[w]
	}
	return warnings
}

// Track handlers that invoke warning for use in TestMain.
var byWarning = make(map[*Warning]map[string]bool)

// ShowWarnings prints any warning to Stdout in a preformatted manner.
// Use the Manager method if more control over output is required.
//
// Note: Both Stdout and Stderr are captured by the the 'go test' command and
// shunted into Stdout (see https://pkg.go.dev/cmd/go#hdr-Test_packages).
// This output stream is only visible when the 'go test -v flag' is used.
func (mgr *Manager) ShowWarnings(output io.Writer) {
	if output == nil {
		output = os.Stdout
	}

	forHandler := ""
	if mgr.Name != "" {
		forHandler = " for " + mgr.Name
	}

	// Must always have this line, even if there are no warnings.
	// The internal/data/ParseWarningData function depends on this line to get the handler name.
	_, _ = fmt.Fprintf(output, "%s\n%sWarnings%s:\n", mgr.showPrefix, mgr.showPrefix, forHandler)

	if warnings := mgr.GetWarnings(); warnings == nil || len(warnings) == 0 {
		_, _ = fmt.Fprintf(output, "%s  None\n", mgr.showPrefix)
	} else {
		// Warnings grouped by level.
		warningTree := make(map[Level][]*Instances)
		for _, w := range warnings {
			list, found := warningTree[w.Level]
			if !found {
				list = make([]*Instances, 0)
			}
			warningTree[w.Level] = append(list, w)
		}
		for _, list := range warningTree {
			sort.Slice(list, func(i, j int) bool {
				return list[i].Name < list[j].Name
			})
		}
		for _, level := range LevelOrder {
			if list, ok := warningTree[level]; ok {
				_, _ = fmt.Fprintf(output, "%s  %s\n", mgr.showPrefix, level.String())
				for _, w := range list {
					// Track handlers that invoke warning for use in TestMain.
					warn := mgr.predefined[w.Name]
					if byWarning[warn] == nil {
						byWarning[warn] = make(map[string]bool)
					}
					byWarning[warn][mgr.Name] = true

					_, _ = fmt.Fprintf(output, "%s  %4d [%s] %s\n", mgr.showPrefix, w.Count, w.Name, w.Summary)
					lineSep := "\n" + mgr.showPrefix + "           +"
					for _, data := range w.Data {
						text := strings.Join(strings.Split(data.Text, "\n"), lineSep)
						if data.Text != "" {
							_, _ = fmt.Fprintf(output, "%s         %s: %s\n", mgr.showPrefix, data.Function, text)
						} else {
							_, _ = fmt.Fprintf(output, "%s         %s\n", mgr.showPrefix, data.Function)
						}
						if data.Record != "" {
							_, _ = fmt.Fprintf(output, "%s           %s\n", mgr.showPrefix, data.Record)
						}
					}
				}
			}
		}
	}
}

// ShowHandlersByWarning uses the global byWarning map to
// show the handlers that issue each warning.
func ShowHandlersByWarning(showPrefix string) {
	byLevel := make(map[Level][]*Warning)
	byName := make(map[string]*Warning)
	for w := range byWarning {
		if byLevel[w.Level] == nil {
			byLevel[w.Level] = make([]*Warning, 0)
		}
		byLevel[w.Level] = append(byLevel[w.Level], w)
		byName[w.Name] = w
	}
	if len(byLevel) < 1 || len(byName) < 1 {
		return
	}
	fmt.Printf("%s\n%s Handlers by warning:\n", showPrefix, showPrefix)
	for _, level := range LevelOrder {
		if warnings, found := byLevel[level]; found {
			fmt.Printf("%s  %s\n", showPrefix, level.String())
			names := make([]string, 0, len(warnings))
			for _, w := range warnings {
				names = append(names, w.Name)
			}
			sort.Strings(names)
			for _, name := range names {
				fmt.Printf("%s    [%s] %s\n", showPrefix, name, byName[name].Summary)
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
func WithWarnings(m *testing.M) {
	flag.Parse()
	exitVal := m.Run()

	managerNames := make([]string, 0, len(managers))
	for name := range managers {
		managerNames = append(managerNames, name)
	}
	sort.Strings(managerNames)
	var showPrefix string
	for _, name := range managerNames {
		manager := managers[name]
		manager.ShowWarnings(nil)
		if showPrefix == "" {
			showPrefix = manager.showPrefix
		}
	}

	ShowHandlersByWarning(showPrefix)

	os.Exit(exitVal)
}
