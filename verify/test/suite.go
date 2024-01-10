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
	"github.com/madkins23/go-slog/replace"
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
// Use a handler from HandlerCreator to create a logger.

func SimpleOptions() *slog.HandlerOptions {
	return &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
}

func LevelOptions(level slog.Leveler) *slog.HandlerOptions {
	return &slog.HandlerOptions{
		Level: level,
	}
}

func SourceOptions() *slog.HandlerOptions {
	return &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	}
}

// Logger returns a slog.Logger with the specified options.
func (suite *SlogTestSuite) Logger(options *slog.HandlerOptions) *slog.Logger {
	return slog.New(suite.Creator(suite.Buffer, options))
}

// -----------------------------------------------------------------------------
// Warning mechanism to trade test failure for warning list at end of tests.

const (
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
// Constant data used for tests.

const (
	message = "This is a message"
)

// -----------------------------------------------------------------------------
// These tests are intended to mimic: src/testing/slogtest/slogtest.go (2024-01-07).

// TestSimpleAttributes tests whether attributes are logged properly.
// Implements slogtest "attrs" test.
func (suite *SlogTestSuite) TestSimpleAttributes() {
	logger := suite.Logger(SimpleOptions())
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
	logger := suite.Logger(SimpleOptions())
	logger.Info(message, "first", "one", "", nil, "pi", math.Pi)
	logMap := suite.logMap()
	suite.Assert().Equal("one", logMap["first"])
	suite.Assert().Equal(math.Pi, logMap["pi"])
	suite.checkNoEmptyAttribute(5, logMap)
}

// TestSimpleAttributesWith tests whether attributes in With() are logged properly.
// Implements slogtest "WithAttrs" test.
func (suite *SlogTestSuite) TestSimpleAttributesWith() {
	logger := suite.Logger(SimpleOptions())
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
	logger := suite.Logger(SimpleOptions())
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
	logger := suite.Logger(SimpleOptions())
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
	logger := suite.Logger(SimpleOptions())
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
	logger := suite.Logger(SimpleOptions())
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
	logger := suite.Logger(SimpleOptions())
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
	logger := suite.Logger(SimpleOptions())
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

// TestSimpleKeys tests whether the three basic keys are present as their defined constants.
// Implements slogtest "built-ins" test.
func (suite *SlogTestSuite) TestSimpleKeys() {
	logger := suite.Logger(SimpleOptions())
	logger.Info(message)
	logMap := suite.logMap()
	suite.checkFieldCount(3, logMap)
	suite.checkLevelKey("INFO", logMap)
	suite.checkMessageKey(message, logMap)
	suite.Assert().NotNil(logMap[slog.TimeKey])
}

// TestSimpleResolveValuer tests logging LogValuer objects.
// Implements slogtest "resolve" test.
func (suite *SlogTestSuite) TestSimpleResolveValuer() {
	logger := suite.Logger(SimpleOptions())
	hidden := &hiddenValuer{v: "something"}
	logger.Info(message, "hidden", hidden)
	logMap := suite.logMap()
	suite.checkFieldCount(4, logMap)
	suite.checkResolution("something", logMap["hidden"])
}

// TestSimpleResolveGroup tests logging LogValuer objects within a group.
// Implements slogtest "resolve-groups" test.
func (suite *SlogTestSuite) TestSimpleResolveGroup() {
	logger := suite.Logger(SimpleOptions())
	logger.Info(message, slog.Group("group",
		slog.Float64("pi", math.Pi), slog.Any("hidden", &hiddenValuer{v: "value"})))
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
	logger := suite.Logger(SimpleOptions())
	logger.With("hidden", &hiddenValuer{v: "value"}).Info(message)
	logMap := suite.logMap()
	suite.checkFieldCount(4, logMap)
	suite.checkResolution("value", logMap["hidden"])
}

// TestSimpleResolveGroupWith tests logging LogValuer objects within a group within a With().
// Implements slogtest "resolve-WithAttrs-groups" test.
func (suite *SlogTestSuite) TestSimpleResolveGroupWith() {
	logger := suite.Logger(SimpleOptions())
	logger.With(slog.Group("group",
		slog.Float64("pi", math.Pi), slog.Any("hidden", &hiddenValuer{v: "value"}))).
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
	logger := suite.Logger(SimpleOptions())
	record := slog.NewRecord(time.Time{}, slog.LevelInfo, message, uintptr(0))
	suite.Require().NoError(logger.Handler().Handle(context.Background(), record))
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
	logger := suite.Logger(SimpleOptions())
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
	logger := suite.Logger(SimpleOptions())
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
	logger := suite.Logger(SimpleOptions())
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
	logger := suite.Logger(SimpleOptions())
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
	logger := suite.Logger(SimpleOptions())
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
	logger := suite.Logger(SourceOptions())
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
	logger := suite.Logger(SimpleOptions())
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
	logger := suite.Logger(SimpleOptions())
	logger.
		With("alpha", "one", "bravo", "hurrah", "charlie", "brown", "charlie", "jones").
		Info(message, "alpha", 2, "charlie", 23.70)
	logMap := suite.logMap()
	suite.checkFieldCount(6, logMap)
}

// -----------------------------------------------------------------------------
// Additional tests.

// TestSimpleContextCancelled verifies that a cancelled context will not affect logging.
func (suite *SlogTestSuite) TestSimpleContextCancelled() {
	logger := suite.Logger(SimpleOptions())
	ctx, cancelFn := context.WithCancel(context.Background())
	logger.InfoContext(ctx, message)
	logMap := suite.logMap()
	suite.checkFieldCount(3, logMap)
	suite.checkLevelKey("INFO", logMap)
	suite.checkMessageKey(message, logMap)
	suite.Assert().NotNil(logMap[slog.TimeKey])
	cancelFn()
	suite.bufferReset()
	logger.InfoContext(ctx, message)
	logMap = suite.logMap()
	suite.checkFieldCount(3, logMap)
	suite.checkLevelKey("INFO", logMap)
	suite.checkMessageKey(message, logMap)
	suite.Assert().NotNil(logMap[slog.TimeKey])
}

// TestSimpleLogAttributes tests the LogAttrs call with all attribute objects.
func (suite *SlogTestSuite) TestSimpleLogAttributes() {
	logger := suite.Logger(SimpleOptions())
	t := time.Now()
	logger.LogAttrs(context.Background(), slog.LevelInfo, message,
		slog.Time("when", t),
		slog.Duration("howLong", time.Minute),
		slog.String("goober", "snoofus"),
		slog.Bool("boolean", true),
		slog.Float64("pi", math.Pi),
		slog.Int("skidoo", 23),
		slog.Int64("minus", -64),
		slog.Uint64("unsigned", 79),
		slog.Any("any", []string{"alpha", "omega"}))
	logMap := suite.logMap()
	suite.checkFieldCount(12, logMap)
	when, ok := logMap["when"].(string)
	suite.True(ok)
	if suite.warn[WarnNanoTime] {
		// Some handlers log times as RFC3339 instead of RFC3339Nano
		suite.Equal(t.Format(time.RFC3339), when)
	} else {
		// Based on the existing behavior of log/slog it should be RFC3339Nano.
		suite.Equal(t.Format(time.RFC3339Nano), when)
	}
	howLong, ok := logMap["howLong"].(float64)
	suite.True(ok)
	if suite.warn[WarnNanoDuration] {
		// Some handlers push out milliseconds instead of nanoseconds.
		suite.Equal(float64(60000), howLong)
	} else {
		// Based on the existing behavior of log/slog it should be nanoseconds.
		//goland:noinspection GoRedundantConversion
		suite.Equal(float64(6e+10), howLong)
	}
	suite.Equal("snoofus", logMap["goober"])
	suite.Equal(true, logMap["boolean"])
	// All numeric attributes come back as float64 due to JSON formatting and parsing.
	suite.Equal(math.Pi, logMap["pi"])
	suite.Equal(float64(23), logMap["skidoo"])
	suite.Equal(float64(-64), logMap["minus"])
	suite.Equal(float64(79), logMap["unsigned"])
	fixed, ok := logMap["any"].([]any)
	suite.True(ok)
	array := make([]string, 0)
	for _, x := range fixed {
		str, ok := x.(string)
		suite.True(ok)
		array = append(array, str)
	}
	suite.Equal([]string{"alpha", "omega"}, array)
}

// TestSimpleDisabled tests whether logging is disabled by level.
func (suite *SlogTestSuite) TestSimpleDisabled() {
	logger := suite.Logger(SimpleOptions())
	logger.Debug(message)
	suite.Assert().Empty(suite.Buffer)
}

// TestSimpleKeyCase tests whether level keys are properly cased.
// Based on the existing behavior of log/slog they should be uppercase.
func (suite *SlogTestSuite) TestSimpleKeyCase() {
	logger := suite.Logger(LevelOptions(slog.LevelDebug))
	for name, level := range map[string]slog.Level{
		"DEBUG": slog.LevelDebug,
		"INFO":  slog.LevelInfo,
		"WARN":  slog.LevelWarn,
		"ERROR": slog.LevelError,
	} {
		logger.Log(context.Background(), level, message)
		logMap := suite.logMap()
		suite.checkLevelKey(name, logMap)
		suite.bufferReset()
	}
}

// TestSimpleLevel tests whether the simple logger is created by default with slog.LevelInfo.
// Other tests (e.g. TestSimpleDisabled) depend on this.
func (suite *SlogTestSuite) TestSimpleLevel() {
	logger := suite.Logger(SimpleOptions())
	suite.Assert().False(logger.Enabled(context.Background(), -1))
	suite.Assert().True(logger.Enabled(context.Background(), slog.LevelInfo))
	suite.Assert().True(logger.Enabled(context.Background(), 1))
	suite.Assert().True(logger.Enabled(context.Background(), slog.LevelWarn))
	suite.Assert().True(logger.Enabled(context.Background(), slog.LevelError))
}

// TestSimpleLevelVar tests the use of a slog.LevelVar.
func (suite *SlogTestSuite) TestSimpleLevelVar() {
	var programLevel = new(slog.LevelVar)
	logger := suite.Logger(LevelOptions(programLevel))
	// Should be INFO by default.
	suite.Assert().Equal(slog.LevelInfo, programLevel.Level())
	suite.Assert().False(logger.Enabled(context.Background(), -1))
	suite.Assert().True(logger.Enabled(context.Background(), slog.LevelInfo))
	suite.Assert().True(logger.Enabled(context.Background(), 1))
	// Change the level.
	programLevel.Set(slog.LevelWarn)
	suite.Assert().Equal(slog.LevelWarn, programLevel.Level())
	suite.Assert().False(logger.Enabled(context.Background(), 3))
	suite.Assert().True(logger.Enabled(context.Background(), slog.LevelWarn))
	suite.Assert().True(logger.Enabled(context.Background(), 5))
}

// TestSimpleLevelDifferent tests whether the simple logger is created with slog.LevelWarn.
// This verifies the test suite can change the level when creating a logger.
// It also verifies changing the level via the handler.
func (suite *SlogTestSuite) TestSimpleLevelDifferent() {
	logger := suite.Logger(LevelOptions(slog.LevelWarn))
	suite.Assert().False(logger.Enabled(context.Background(), -1))
	suite.Assert().False(logger.Enabled(context.Background(), slog.LevelInfo))
	suite.Assert().False(logger.Enabled(context.Background(), 3))
	suite.Assert().True(logger.Enabled(context.Background(), slog.LevelWarn))
	suite.Assert().True(logger.Enabled(context.Background(), 5))
	suite.Assert().True(logger.Enabled(context.Background(), slog.LevelError))
}

// TestSimpleTimestampFormat tests whether a timestamp can be parsed.
// Based on the existing behavior of log/slog the timestamp format is RFC3339.
func (suite *SlogTestSuite) TestSimpleTimestampFormat() {
	logger := suite.Logger(SimpleOptions())
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
	logger := suite.Logger(SourceOptions())
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

// TestSourceLevel tests whether the source logger is created by default with slog.LevelInfo.
func (suite *SlogTestSuite) TestSourceLevel() {
	logger := suite.Logger(SourceOptions())
	suite.Assert().False(logger.Enabled(context.Background(), -1))
	suite.Assert().True(logger.Enabled(context.Background(), slog.LevelInfo))
	suite.Assert().True(logger.Enabled(context.Background(), 1))
	suite.Assert().True(logger.Enabled(context.Background(), slog.LevelWarn))
	suite.Assert().True(logger.Enabled(context.Background(), slog.LevelError))
}

// TestSourceLevelDifferent tests whether the source logger is created with slog.LevelWarn.
// This verifies the test suite can change the level when creating a logger.
// It also verifies changing the level via the handler.
func (suite *SlogTestSuite) TestSourceLevelDifferent() {
	logger := suite.Logger(&slog.HandlerOptions{
		Level:     slog.LevelWarn,
		AddSource: true,
	})
	suite.Assert().False(logger.Enabled(context.Background(), -1))
	suite.Assert().False(logger.Enabled(context.Background(), slog.LevelInfo))
	suite.Assert().False(logger.Enabled(context.Background(), 1))
	suite.Assert().True(logger.Enabled(context.Background(), slog.LevelWarn))
	suite.Assert().True(logger.Enabled(context.Background(), 5))
	suite.Assert().True(logger.Enabled(context.Background(), slog.LevelError))
}

// -----------------------------------------------------------------------------
// Tests of slog.HandlerOptions.ReplaceAttr functionality.

// TestSimpleReplaceAttr tests the use of HandlerOptions.ReplaceAttr.
func (suite *SlogTestSuite) TestSimpleReplaceAttr() {
	logger := suite.Logger(&slog.HandlerOptions{
		Level: slog.LevelInfo,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			switch a.Key {
			case "alpha":
				return slog.String(a.Key, "omega")
			case "change":
				return slog.String("bravo", a.Value.String())
			case "remove":
				return replace.EmptyAttr
			}
			return a
		},
	})
	logger.Info(message, "alpha", "beta", "change", "my key", "remove", "me")
	logMap := suite.logMap()
	if suite.warn[WarnNoReplAttr] {
		issues := make([]string, 4)
		if len(logMap) > 5 {
			issues = append(issues, fmt.Sprintf("too many attributes: %d", len(logMap)))
		}
		value, ok := logMap["alpha"].(string)
		suite.Require().True(ok)
		if value != "omega" {
			issues = append(issues, fmt.Sprintf("alpha == %s", value))
		}
		if logMap["change"] != nil {
			issues = append(issues, "change still exists")
		}
		if logMap["remove"] != nil {
			issues = append(issues, "remove still exists")
		}
		if len(issues) > 0 {
			suite.addWarning(WarnNoReplAttr, strings.Join(issues, ", "), false)
			return
		}
		suite.addWarning(WarnUnused, WarnNoReplAttr, true)
	}
	if suite.warn[WarnEmptyAttributes] {
		suite.checkFieldCount(6, logMap)
	} else {
		suite.checkFieldCount(5, logMap)
	}
	suite.Assert().Equal("omega", logMap["alpha"])
	suite.Assert().Equal("my key", logMap["bravo"])
	suite.Assert().Nil(logMap["remove"])
}

// TestSourceReplaceAttrBasic tests the use of HandlerOptions.ReplaceAttr
// on basic attributes (time, level, message, source).
func (suite *SlogTestSuite) TestSourceReplaceAttrBasic() {
	logger := suite.Logger(&slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			switch a.Key {
			case slog.TimeKey:
				return replace.EmptyAttr
			case slog.LevelKey:
				return slog.String(slog.LevelKey, "Tilted")
			case slog.MessageKey:
				return slog.String("Message", a.Value.String())
			case slog.SourceKey:
				return slog.String(slog.SourceKey, "all knowledge")
			}
			return a
		},
	})
	logger.Info(message)
	logMap := suite.logMap()
	if suite.warn[WarnNoReplAttr] || suite.warn[WarnNoReplAttrBasic] {
		issues := make([]string, 0, 5)
		if len(logMap) > 3 {
			issues = append(issues, fmt.Sprintf("too many attributes: %d", len(logMap)))
		}
		if logMap[slog.TimeKey] != nil {
			issues = append(issues, slog.TimeKey+" field still exists")
		}
		if logMap[slog.MessageKey] != nil {
			issues = append(issues, slog.MessageKey+" field still exists")
		} else if suite.warn[WarnMessageKey] && logMap["message"] != nil {
			issues = append(issues, "message field still exists")
		}
		// TODO: This one may still work, in samber it's apparently a separate field from basic.
		if value, ok := logMap[slog.SourceKey].(string); !ok || value != "all knowledge" {
			issues = append(issues, fmt.Sprintf("%s == %v", slog.SourceKey, logMap[slog.SourceKey]))
		}
		if len(issues) > 0 {
			suite.addWarning(WarnNoReplAttrBasic, strings.Join(issues, ", "), false)
			return
		}
		suite.addWarning(WarnUnused, WarnNoReplAttrBasic, true)
	}
	suite.checkFieldCount(3, logMap)
	suite.Assert().Nil(logMap[slog.TimeKey])
	suite.Assert().Equal("Tilted", logMap[slog.LevelKey])
	suite.Assert().Equal(message, logMap["Message"])
	suite.Assert().Equal("all knowledge", logMap[slog.SourceKey])
}

// -----------------------------------------------------------------------------
// Methods for checks common to various tests or too long to put in a test.

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
// Utility methods.

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

func (suite *SlogTestSuite) bufferReset() {
	suite.Buffer.Reset()
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

type hiddenValuer struct {
	v any
}

func (r *hiddenValuer) LogValue() slog.Value {
	return slog.AnyValue(r.v)
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
