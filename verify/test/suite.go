package test

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"math"
	"os"
	"runtime"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	testJSON "github.com/madkins23/go-slog/json"
	"github.com/madkins23/go-slog/test"
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

// -----------------------------------------------------------------------------
// Define top-level structs.

// SlogTestSuite provides various tests for a specified log/slog.Handler.
type SlogTestSuite struct {
	suite.Suite
	*bytes.Buffer
	warn     map[string]bool
	warnings map[string]*Warning

	// Creator is an object that generates slog.Logger objects for tests.
	// This field must be configured by test suites.
	Creator HandlerCreator

	// Name of Handler for warnings display.
	Name string
}

// HandlerCreator is responsible for generating log/slog handler instances.
// Define one for a given test file and use to instantiate SlogTestSuite.
type HandlerCreator interface {
	// SimpleHandler creates the simplest possible handler.
	SimpleHandler(w io.Writer, level slog.Leveler) slog.Handler

	// SourceHandler creates a simple handler with source logging turned on.
	SourceHandler(w io.Writer, level slog.Leveler) slog.Handler
}

// -----------------------------------------------------------------------------
// Suite configuration methods.

const duplicateFieldsNotError = true

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
// Warning mechanism to trade test failure for warning list at end of tests.

const (
	WarnDuplicates      = "Duplicate field(s) found"
	WarnEmptyAttributes = "Empty attribute(s) logged (\"\":null)"
	WarnGroupInline     = "Group with empty key does not inline subfields"
	WarnLevelCase       = "Log level in lowercase"
	WarnMessageKey      = "Wrong message key (should be 'msg')"
	WarnResolver        = "LogValuer objects are not resolved"
	WarnSourceKey       = "Source data not logged when AddSource flag set"
	WarnSubgroupEmpty   = "Empty subgroup(s) logged"
	WarnUnused          = "Unused Warning(s)"
	WarnZeroPC          = "SourceKey logged for zero PC"
	WarnZeroTime        = "Zero time is logged"
)

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

var suites = make([]*SlogTestSuite, 0)

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

// -----------------------------------------------------------------------------
// Use a handler from HandlerCreator to create a logger.

// SimpleLogger returns a simple handler within a slog.Logger.
// Override this method to test other types of slog JSON handlers.
func (suite *SlogTestSuite) SimpleLogger(level slog.Leveler) *slog.Logger {
	if level == nil {
		level = slog.LevelInfo
	}
	return slog.New(suite.Creator.SimpleHandler(suite.Buffer, level))
}

// SourceLogger returns a simple handler with the source key activated
// wrapped within a slog.logger.
// Override this method to test other types of slog JSON handlers.
func (suite *SlogTestSuite) SourceLogger(level slog.Leveler) *slog.Logger {
	if level == nil {
		level = slog.LevelInfo
	}
	return slog.New(suite.Creator.SourceHandler(suite.Buffer, level))
}

// -----------------------------------------------------------------------------
// Constant data used for tests.

const (
	message = "This is a message"
)

// -----------------------------------------------------------------------------
// These tests are intended to mimic: src/testing/slogtest/slogtest.go (2024-01-07).

// TestSimpleAttributes tests whether attributes are logged properly.
// Implements slogtest "attrs" test.
func (suite *SlogTestSuite) TestSimpleAttributes() {
	logger := suite.SimpleLogger(nil)
	logger.Info(message, "first", "one", "second", 2, "pi", math.Pi)
	logMap := suite.logMap()
	suite.checkFieldCount(6, logMap)
	suite.Assert().Equal("one", logMap["first"])
	suite.Assert().Equal(float64(2), logMap["second"])
	suite.Assert().Equal(math.Pi, logMap["pi"])
}

// TestSimpleAttributeEmpty tests whether attributes with empty names and nil values are logged properly.
// Based on the existing behavior of log/slog the field is hot created.
// Implements slogtest "empty-attr" test.
func (suite *SlogTestSuite) TestSimpleAttributeEmpty() {
	logger := suite.SimpleLogger(nil)
	logger.Info(message, "first", "one", "", nil, "pi", math.Pi)
	logMap := suite.logMap()
	suite.Assert().Equal("one", logMap["first"])
	suite.Assert().Equal(math.Pi, logMap["pi"])
	suite.checkNoEmptyAttribute(5, logMap)
}

// TestSimpleAttributesWith tests whether attributes in With() are logged properly.
// Implements slogtest "WithAttrs" test.
func (suite *SlogTestSuite) TestSimpleAttributesWith() {
	logger := suite.SimpleLogger(nil)
	logger.With("first", "one", "second", 2).Info(message, "pi", math.Pi)
	logMap := suite.logMap()
	suite.checkFieldCount(6, logMap)
	suite.Assert().Equal("one", logMap["first"])
	suite.Assert().Equal(float64(2), logMap["second"])
	suite.Assert().Equal(math.Pi, logMap["pi"])
}

// TestSimpleGroup tests the use of a logging group.
// Implements slogtest "groups" test.
func (suite *SlogTestSuite) TestSimpleGroup() {
	logger := suite.SimpleLogger(nil)
	logger.Info(message, "first", "one",
		slog.Group("group", "second", 2, slog.String("third", "3"), "fourth", "forth"),
		"pi", math.Pi)
	logMap := suite.logMap()
	suite.checkFieldCount(6, logMap)
	suite.Assert().Equal("one", logMap["first"])
	if group, ok := logMap["group"].(map[string]any); ok {
		suite.Assert().Len(group, 3)
		suite.Assert().Equal(float64(2), group["second"])
		suite.Assert().Equal("3", group["third"])
		suite.Assert().Equal("forth", group["fourth"])
	} else {
		suite.Fail("Group not map[string]any")
	}
	suite.Assert().Equal(math.Pi, logMap["pi"])
}

// TestSimpleGroupEmpty tests logging an empty group.
// Based on the existing behavior of log/slog the group field is not logged.
// Implements slogtest "empty-group" test.
func (suite *SlogTestSuite) TestSimpleGroupEmpty() {
	logger := suite.SimpleLogger(nil)
	logger.Info(message, slog.Group("group"))
	logMap := suite.logMap()
	suite.checkFieldCount(3, logMap)
	_, found := logMap["group"]
	suite.Assert().False(found)
}

// TestSimpleGroupInline tests the use of a group with an empty name.
// Based on the existing behavior of log/slog the group field is not logged and
// the fields within the group are moved to the top level.
// Implements slogtest "inline-group" test.
func (suite *SlogTestSuite) TestSimpleGroupInline() {
	logger := suite.SimpleLogger(nil)
	logger.Info(message, "first", "one",
		slog.Group("", "second", 2, slog.String("third", "3"), "fourth", "forth"),
		"pi", math.Pi)
	logMap := suite.logMap()
	suite.Assert().Equal("one", logMap["first"])
	suite.Assert().Equal(math.Pi, logMap["pi"])
	checkFieldFn := func(fieldMap map[string]any) {
		suite.Assert().Equal(float64(2), fieldMap["second"])
		suite.Assert().Equal("3", fieldMap["third"])
		suite.Assert().Equal("forth", fieldMap["fourth"])
	}
	if suite.warn[WarnGroupInline] {
		counter := suite.fieldCounter()
		suite.Require().NoError(counter.Parse())
		if counter.NumFields() == 6 {
			if group, ok := logMap[""].(map[string]any); ok {
				suite.Assert().Len(group, 3)
				checkFieldFn(group)
			} else {
				suite.Fail("Group not map[string]any")
			}
			suite.addWarning(WarnGroupInline, "", true)
			return
		}
		suite.addWarning(WarnUnused, WarnGroupInline, false)
	}
	suite.checkFieldCount(8, logMap)
	checkFieldFn(logMap)
}

// TestSimpleGroupWith tests the use of a logging group specified using WithGroup.
// Implements slogtest "WithGroup" test.
func (suite *SlogTestSuite) TestSimpleGroupWith() {
	logger := suite.SimpleLogger(nil)
	logger.WithGroup("group").Info(message, "first", "one", "pi", math.Pi)
	logMap := suite.logMap()
	suite.checkFieldCount(4, logMap)
	if group, ok := logMap["group"].(map[string]any); ok {
		suite.Assert().Len(group, 2)
		suite.Assert().Equal("one", group["first"])
		suite.Assert().Equal(math.Pi, group["pi"])
	} else {
		suite.Fail("Group not map[string]any")
	}
}

// TestSimpleGroupWithMulti tests the use of multiple logging groups.
// Implements slogtest "multi-with" test.
func (suite *SlogTestSuite) TestSimpleGroupWithMulti() {
	logger := suite.SimpleLogger(nil)
	logger.With("first", "one").
		WithGroup("group").With("second", 2, "third", "3").
		WithGroup("subGroup").Info(message, "fourth", "forth", "pi", math.Pi)
	logMap := suite.logMap()
	suite.checkFieldCount(5, logMap)
	if group, ok := logMap["group"].(map[string]any); ok {
		suite.Assert().Len(group, 3)
		suite.Assert().Equal(float64(2), group["second"])
		suite.Assert().Equal("3", group["third"])
		if group, ok := group["subGroup"].(map[string]any); ok {
			suite.Assert().Len(group, 2)
			suite.Assert().Equal("forth", group["fourth"])
			suite.Assert().Equal(math.Pi, group["pi"])
		} else {
			suite.Fail("Sub-group not map[string]any")
		}
	} else {
		suite.Fail("Group not map[string]any")
	}
}

// TestSimpleGroupWithMultiSubEmpty tests the use of multiple logging groups when the subgroup is empty.
// Based on the existing behavior of log/slog the subgroup field is not logged.
// Implements slogtest "empty-group-record" test.
func (suite *SlogTestSuite) TestSimpleGroupWithMultiSubEmpty() {
	logger := suite.SimpleLogger(nil)
	logger.With("first", "one").
		WithGroup("group").With("second", 2, "third", "3").
		WithGroup("subGroup").Info(message)
	logMap := suite.logMap()
	suite.checkFieldCount(5, logMap)
	_, found := logMap["subGroup"]
	suite.Assert().False(found, "subGroup found at top level")
	if group, ok := logMap["group"].(map[string]any); ok {
		suite.Assert().Equal(float64(2), group["second"])
		suite.Assert().Equal("3", group["third"])
		if suite.warn[WarnSubgroupEmpty] {
			if len(group) > 2 {
				if subGroup, found := group["subGroup"]; found {
					if sg, ok := subGroup.(map[string]any); ok && len(sg) < 1 {
						suite.addWarning(WarnSubgroupEmpty, "", true)
						return
					}
				}
			}
			suite.addWarning(WarnUnused, WarnSubgroupEmpty, false)
		}
		suite.Assert().Len(group, 2)
		_, found := group["subGroup"]
		suite.Assert().False(found)
	} else {
		suite.Fail("Group not map[string]any")
	}
}

// TestSimpleKeyCase tests whether level keys are properly cased.
// TODO:
//func (suite *SlogTestSuite) TestSimpleKeyCase() {
//	logger := suite.SimpleLogger()
//	logger.
//	for name, level := range map[string]int{} {
//
//	}
//	logger.Info(message)
//	logMap := suite.logMap()
//	suite.checkLevelKey("INFO", logMap)
//	suite.checkLevelKey("INFO", logMap)
//	suite.checkLevelKey("INFO", logMap)
//	suite.checkLevelKey("INFO", logMap)
//}

// TestSimpleKeys tests whether the three basic keys are present as their defined constants.
// Implements slogtest "built-ins" test.
func (suite *SlogTestSuite) TestSimpleKeys() {
	logger := suite.SimpleLogger(nil)
	logger.Info(message)
	logMap := suite.logMap()
	suite.checkFieldCount(3, logMap)
	suite.checkLevelKey("INFO", logMap)
	suite.checkMessageKey(message, logMap)
	suite.Assert().NotNil(logMap[slog.TimeKey])
}

// TestSimpleResolve tests logging LogValuer objects.
// Implements slogtest "resolve" test.
func (suite *SlogTestSuite) TestSimpleResolve() {
	logger := suite.SimpleLogger(nil)
	logger.Info(message, "hidden", &hiddenValue{v: "value"})
	logMap := suite.logMap()
	suite.checkFieldCount(4, logMap)
	suite.checkResolution("value", logMap["hidden"])
}

// TestSimpleResolveGroup tests logging LogValuer objects within a group.
// Implements slogtest "resolve-groups" test.
func (suite *SlogTestSuite) TestSimpleResolveGroup() {
	logger := suite.SimpleLogger(nil)
	logger.Info(message, slog.Group("group",
		slog.Float64("pi", math.Pi), slog.Any("hidden", &hiddenValue{v: "value"})))
	logMap := suite.logMap()
	suite.checkFieldCount(4, logMap)
	if group, ok := logMap["group"].(map[string]any); ok {
		suite.Assert().Len(group, 2)
		suite.Assert().Equal(math.Pi, group["pi"])
		suite.checkResolution("value", group["hidden"])
	} else {
		suite.Fail("Group not map[string]any")
	}
}

// TestSimpleResolveWith tests logging LogValuer objects within a With().
// Implements slogtest "resolve-withAttrs" test.
func (suite *SlogTestSuite) TestSimpleResolveWith() {
	logger := suite.SimpleLogger(nil)
	logger.With("hidden", &hiddenValue{v: "value"}).Info(message)
	logMap := suite.logMap()
	suite.checkFieldCount(4, logMap)
	suite.checkResolution("value", logMap["hidden"])
}

// TestSimpleResolveGroupWith tests logging LogValuer objects within a group within a With().
// Implements slogtest "resolve-WithAttrs-groups" test.
func (suite *SlogTestSuite) TestSimpleResolveGroupWith() {
	logger := suite.SimpleLogger(nil)
	logger.With(slog.Group("group",
		slog.Float64("pi", math.Pi), slog.Any("hidden", &hiddenValue{v: "value"}))).
		Info(message)
	logMap := suite.logMap()
	suite.checkFieldCount(4, logMap)
	if group, ok := logMap["group"].(map[string]any); ok {
		suite.Assert().Len(group, 2)
		suite.Assert().Equal(math.Pi, group["pi"])
		suite.checkResolution("value", group["hidden"])
	} else {
		suite.Fail("Group not map[string]any")
	}
}

// TestSimpleZeroTime tests whether a zero time in a slog.Record is output.
// Based on the existing behavior of log/slog the field is not logged.
// Implements slogtest "zero-time" test.
func (suite *SlogTestSuite) TestSimpleZeroTime() {
	logger := suite.SimpleLogger(nil)
	record := slog.NewRecord(time.Time{}, slog.LevelInfo, message, uintptr(0))
	suite.Require().NoError(logger.Handler().Handle(context.TODO(), record))
	logMap := suite.logMap()
	suite.checkLevelKey("INFO", logMap)
	suite.checkMessageKey(message, logMap)
	if suite.warn[WarnZeroTime] {
		counter := suite.fieldCounter()
		suite.Require().NoError(counter.Parse())
		if counter.NumFields() == 3 {
			if timeAny, found := logMap[slog.TimeKey]; found {
				timeZero := suite.parseTime(timeAny)
				suite.Assert().Equal(time.Time{}, timeZero, "time should be zero")
				suite.addWarning(WarnZeroTime, "", true)
				return
			}
		}
		suite.addWarning(WarnUnused, WarnZeroTime, false)
	}
	suite.checkFieldCount(2, logMap)
	suite.Assert().Nil(logMap[slog.TimeKey])
}

// -----------------------------------------------------------------------------
// Tests extending slogTest tests.

// TestSimpleAttributeEmptyName tests whether attributes with empty names are logged properly.
// Based on the existing behavior of log/slog the field is created with a blank name.
// Extension of slogtest "empty-attr" test.
func (suite *SlogTestSuite) TestSimpleAttributeEmptyName() {
	logger := suite.SimpleLogger(nil)
	logger.Info(message, "first", "one", "", 2, "pi", math.Pi)
	logMap := suite.logMap()
	suite.checkFieldCount(6, logMap)
	suite.Assert().Equal("one", logMap["first"])
	value, found := logMap[""]
	suite.Assert().True(found)
	suite.Assert().Equal(float64(2), value)
	suite.Assert().Equal(math.Pi, logMap["pi"])
}

// TestSimpleAttributeWithEmpty tests whether attributes with empty names and nil values
// specified in With() are logged properly.
// Based on the existing behavior of log/slog the field is hot created.
// Extension of slogtest "WithAttrs" test.
func (suite *SlogTestSuite) TestSimpleAttributeWithEmpty() {
	logger := suite.SimpleLogger(nil)
	logger.With("", nil).Info(message, "first", "one", "pi", math.Pi)
	logMap := suite.logMap()
	suite.Assert().Equal("one", logMap["first"])
	suite.Assert().Equal(math.Pi, logMap["pi"])
	suite.checkNoEmptyAttribute(5, logMap)
}

// TestSimpleAttributeWithEmptyName tests whether With() attributes with empty names are logged properly.
// Based on the existing behavior of log/slog the field is created with a blank name.
// Extension of slogtest "WithAttrs" test.
func (suite *SlogTestSuite) TestSimpleAttributeWithEmptyName() {
	logger := suite.SimpleLogger(nil)
	logger.With("", 2).Info(message, "first", "one", "pi", math.Pi)
	logMap := suite.logMap()
	suite.checkFieldCount(6, logMap)
	suite.Assert().Equal("one", logMap["first"])
	value, found := logMap[""]
	suite.Assert().True(found)
	suite.Assert().Equal(float64(2), value)
	suite.Assert().Equal(math.Pi, logMap["pi"])
}

// TestSimpleAttributeNil tests whether attributes with nil values are logged properly.
// Based on the existing behavior of log/slog the field is created with a nil/null value.
// Extension of slogtest "empty-attr" test.
func (suite *SlogTestSuite) TestSimpleAttributeNil() {
	logger := suite.SimpleLogger(nil)
	logger.Info(message, "first", "one", "second", nil, "pi", math.Pi)
	logMap := suite.logMap()
	suite.checkFieldCount(6, logMap)
	suite.Assert().Equal("one", logMap["first"])
	suite.Assert().Nil(logMap["second"])
	suite.Assert().Equal(math.Pi, logMap["pi"])
}

// TestSimpleAttributeWithNil tests whether With() attributes with nil values are logged properly.
// Based on the existing behavior of log/slog the field is created with a nil/null value.
// Extension of slogtest "WithAttrs" test.
func (suite *SlogTestSuite) TestSimpleAttributeWithNil() {
	logger := suite.SimpleLogger(nil)
	logger.With("second", nil).Info(message, "first", "one", "pi", math.Pi)
	logMap := suite.logMap()
	suite.checkFieldCount(6, logMap)
	suite.Assert().Equal("one", logMap["first"])
	suite.Assert().Nil(logMap["second"])
	suite.Assert().Equal(math.Pi, logMap["pi"])
}

// TestSourceZeroPC tests generation of a source key.
// Implements slogtest "empty-PC" test.
func (suite *SlogTestSuite) TestSourceZeroPC() {
	logger := suite.SourceLogger(nil)
	record := slog.NewRecord(time.Now(), slog.LevelInfo, message, 0)
	suite.Require().NoError(logger.Handler().Handle(context.Background(), record))
	logMap := suite.logMap()
	suite.checkLevelKey("INFO", logMap)
	suite.checkMessageKey(message, logMap)
	suite.Assert().NotNil(logMap[slog.TimeKey])
	if suite.warn[WarnZeroPC] {
		if _, ok := logMap[slog.SourceKey].(map[string]any); ok {
			suite.addWarning(WarnZeroPC, "", true)
			return
		}
		suite.addWarning(WarnUnused, WarnZeroPC, false)
	}

	suite.checkFieldCount(3, logMap)
}

// -----------------------------------------------------------------------------
// Duplicate testing, which isn't currently regarded as an error.
// This issue is under discussion in https://github.com/golang/go/issues/59365.

// TestSimpleAttributeDuplicate tests whether duplicate attributes are logged properly.
// Based on the existing behavior of log/slog the second occurrence overrides the first.
func (suite *SlogTestSuite) TestSimpleAttributeDuplicate() {
	logger := suite.SimpleLogger(nil)
	logger.Info(message,
		"alpha", "one", "alpha", 2, "bravo", "hurrah",
		"charlie", "brown", "charlie", 3, "charlie", 23.79)
	logMap := suite.logMap()
	suite.checkFieldCount(6, logMap)
}

// TestSimpleAttributeWithDuplicate tests whether duplicate attributes are logged properly
// when the duplicate is introduced by With() and then the main call.
// Based on the existing behavior of log/slog the second occurrence overrides the first.
func (suite *SlogTestSuite) TestSimpleAttributeWithDuplicate() {
	logger := suite.SimpleLogger(nil)
	logger.
		With("alpha", "one", "bravo", "hurrah", "charlie", "brown", "charlie", "jones").
		Info(message, "alpha", 2, "charlie", 23.70)
	logMap := suite.logMap()
	suite.checkFieldCount(6, logMap)
}

// -----------------------------------------------------------------------------
// Additional tests.

// TestSimpleDisabled tests whether logging is disabled by level.
func (suite *SlogTestSuite) TestSimpleDisabled() {
	logger := suite.SimpleLogger(nil)
	logger.Debug(message)
	suite.Assert().Empty(suite.Buffer)
}

// TestSimpleLevel tests whether the simple logger is created with slog.LevelInfo.
// Other tests (e.g. TestSimpleDisabled) depend on this.
func (suite *SlogTestSuite) TestSimpleLevel() {
	logger := suite.SimpleLogger(nil)
	suite.Assert().False(logger.Enabled(context.TODO(), -1))
	suite.Assert().True(logger.Enabled(context.TODO(), slog.LevelInfo))
	suite.Assert().True(logger.Enabled(context.TODO(), 1))
	suite.Assert().True(logger.Enabled(context.TODO(), slog.LevelWarn))
	suite.Assert().True(logger.Enabled(context.TODO(), slog.LevelError))
}

// TestSimpleTimestampFormat tests whether a timestamp can be parsed.
// Based on the existing behavior of log/slog the timestamp format is RFC3339.
func (suite *SlogTestSuite) TestSimpleTimestampFormat() {
	logger := suite.SimpleLogger(nil)
	logger.Info(message)
	logMap := suite.logMap()
	suite.checkFieldCount(3, logMap)
	timeObj := suite.parseTime(logMap[slog.TimeKey])
	suite.Assert().Equal(time.Now().Year(), timeObj.Year())
	suite.Assert().Equal(time.Now().Month(), timeObj.Month())
	suite.Assert().Equal(time.Now().Day(), timeObj.Day())
}

// TestSourceKey tests generation of a source key.
func (suite *SlogTestSuite) TestSourceKey() {
	logger := suite.SourceLogger(nil)
	var pcs [1]uintptr
	runtime.Callers(2, pcs[:])
	record := slog.NewRecord(time.Now(), slog.LevelInfo, message, pcs[0])
	suite.Require().NoError(logger.Handler().Handle(context.Background(), record))
	logMap := suite.logMap()
	suite.checkLevelKey("INFO", logMap)
	suite.checkMessageKey(message, logMap)
	suite.Assert().NotNil(logMap[slog.TimeKey])
	suite.checkSourceKey(4, logMap)
}

// -----------------------------------------------------------------------------

// checkFieldCount checks whether the prescribed number of fields exist at the top level.
// In addition to using the logMap generated by unmarshaling the JSON log data,
// the custom test.FieldCounter is used to make sure there are no duplicates.
func (suite *SlogTestSuite) checkFieldCount(fieldCount uint, logMap map[string]any) {
	if suite.warn[WarnDuplicates] {
		counter := suite.fieldCounter()
		suite.Require().NoError(counter.Parse())
		if len(counter.Duplicates()) > 0 {
			suite.addWarning(WarnDuplicates, fmt.Sprintf("%v", counter.Duplicates()), false)
			return
		}
		//goland:noinspection GoBoolExpressions
		if !duplicateFieldsNotError {
			// Don't WarnUnused since WarnDuplicates is currently always set.
			suite.addWarning(WarnUnused, WarnDuplicates, false)
		}
	}
	suite.Assert().Len(logMap, int(fieldCount))
	// Double check to make sure there are no duplicate fields at the top level.
	counter := suite.fieldCounter()
	suite.Require().NoError(counter.Parse())
	suite.Assert().Equal(fieldCount, counter.NumFields())
	suite.Assert().Empty(counter.Duplicates())
}

func (suite *SlogTestSuite) checkLevelKey(level string, logMap map[string]any) {
	// The log/slog.JSONHandler generates uppercase.
	level = strings.ToUpper(level)
	if suite.warn[WarnLevelCase] {
		if logLevel, ok := logMap[slog.LevelKey].(string); ok {
			if suite.Assert().Equal(level, strings.ToUpper(logLevel)) && level != logLevel {
				suite.addWarning(WarnLevelCase, "'"+logLevel+"'", false)
				return
			}
		}
		suite.addWarning(WarnUnused, WarnLevelCase, false)
	}
	suite.Assert().Equal(level, logMap[slog.LevelKey])
}

func (suite *SlogTestSuite) checkMessageKey(message string, logMap map[string]any) {
	if suite.warn[WarnMessageKey] {
		if _, found := logMap[slog.MessageKey]; found {
			// Something exists for the proper key so fall through to test assertion.
		} else if msg, found := logMap["message"]; found {
			// Found something on the known alternate key.
			if message == msg {
				suite.addWarning(WarnMessageKey, "`message`", false)
				return
			}
		}
		suite.addWarning(WarnUnused, WarnMessageKey, false)
	}
	suite.Assert().Equal(message, logMap[slog.MessageKey])
}

func (suite *SlogTestSuite) checkNoEmptyAttribute(fieldCount uint, logMap map[string]any) {
	if suite.warn[WarnEmptyAttributes] {
		// Warn for logging of empty attribute.
		counter := suite.fieldCounter()
		suite.Require().NoError(counter.Parse())
		if counter.NumFields() == fieldCount+1 {
			if _, found := logMap[""]; found {
				suite.addWarning(WarnEmptyAttributes, "", true)
				return
			}
		}
		suite.addWarning(WarnUnused, WarnEmptyAttributes, false)
	}
	suite.checkFieldCount(fieldCount, logMap)
	_, found := logMap[""]
	suite.Assert().False(found)
}

func (suite *SlogTestSuite) checkResolution(value any, actual any) {
	if suite.warn[WarnResolver] {
		if value != actual {
			suite.addWarning(WarnResolver, "", true)
			return
		}
		suite.addWarning(WarnUnused, WarnResolver, false)
	}
	suite.Assert().Equal(value, actual)
}

var sourceKeys = map[string]any{
	"file":     "",
	"function": "",
	"line":     123.456,
}

func (suite *SlogTestSuite) checkSourceKey(fieldCount uint, logMap map[string]any) {
	if suite.warn[WarnSourceKey] {
		sourceData := logMap[slog.SourceKey]
		if sourceData == nil {
			suite.addWarning(WarnSourceKey, "no 'source' key", true)
			return
		}
		source, ok := sourceData.(map[string]any)
		if !ok {
			suite.addWarning(WarnSourceKey, "'source' key not a group", true)
			return
		}
		var text strings.Builder
		sep := ""
		for field := range sourceKeys {
			var state string
			value := source[field]
			if value == nil {
				state = "missing"
			} else if _, ok := value.(string); !ok {
				state = "not a string"
			}
			if state != "" {
				text.WriteString(fmt.Sprintf("%s%s: %s", sep, field, state))
				sep = ", "
			}
		}
		if text.Len() > 0 {
			suite.addWarning(WarnSourceKey, text.String(), true)
		}
		suite.addWarning(WarnUnused, WarnSourceKey, false)
	}

	suite.checkFieldCount(fieldCount, logMap)
	if group, ok := logMap[slog.SourceKey].(map[string]any); ok {
		suite.Assert().Len(group, 3)
		for field, exemplar := range sourceKeys {
			suite.Assert().NotNil(group[field])
			suite.Assert().IsType(exemplar, group[field], "key: "+field)
		}
	} else {
		suite.Fail("Group not map[string]any")
	}
}

func (suite *SlogTestSuite) parseTime(timeAny any) time.Time {
	suite.Assert().NotNil(timeAny)
	timeStr, ok := timeAny.(string)
	suite.Assert().True(ok)
	timeObj, err := time.Parse(time.RFC3339, timeStr)
	suite.Assert().NoError(err)
	return timeObj
}

// -----------------------------------------------------------------------------

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

func (suite *SlogTestSuite) fieldCounter() *testJSON.FieldCounter {
	return testJSON.NewFieldCounter(suite.Buffer.Bytes())
}

func (suite *SlogTestSuite) logMap() map[string]any {
	test.Debugf(1, ">>> JSON: %s", suite.Buffer.Bytes())
	var results map[string]any
	suite.Require().NoError(json.Unmarshal(suite.Buffer.Bytes(), &results))
	return results
}

// -----------------------------------------------------------------------------

type hiddenValue struct {
	v any
}

func (r *hiddenValue) LogValue() slog.Value {
	return slog.AnyValue(r.v)
}

func (r *hiddenValue) String() string {
	return fmt.Sprintf("<hiddenValue(%v)>", r.v)
}

// -----------------------------------------------------------------------------

func currentFunctionName() string {
	pc := make([]uintptr, 10)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	more := true
	for more {
		var frame runtime.Frame
		frame, more = frames.Next()
		parts := strings.Split(frame.Function, ".")
		functionName := parts[len(parts)-1]
		if strings.HasPrefix(functionName, "Test") {
			return functionName
		}
	}
	return "Unknown"
}
