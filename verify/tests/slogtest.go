package tests

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"reflect"
	"strings"
	"time"

	"github.com/madkins23/go-slog/infra"
	"github.com/madkins23/go-slog/internal/warning"
)

// -----------------------------------------------------------------------------
// These tests are intended to mimic: src/testing/slogtest/slogtest.go (2024-01-07).

// TestAttributes tests whether attributes are logged properly.
//   - Implements slogtest "attrs" test.
func (suite *SlogTestSuite) TestAttributes() {
	logger := suite.Logger(infra.SimpleOptions())
	logger.Info(message, "first", "one", "second", 2, "pi", math.Pi)
	logMap := suite.logMap()
	suite.checkFieldCount(6, logMap)
	suite.Assert().Equal("one", logMap["first"])
	suite.Assert().Equal(float64(2), logMap["second"])
	suite.Assert().Equal(math.Pi, logMap["pi"])
}

// TestAttributesEmpty tests whether attributes with empty names and nil values are ignored.
//   - Based on the existing behavior of log/slog the field is hot created.
//   - Implements slogtest "empty-attr" test.
//   - From https://pkg.go.dev/log/slog@master#Handler
func (suite *SlogTestSuite) TestAttributesEmpty() {
	logger := suite.Logger(infra.SimpleOptions())
	logger.Info(message, "first", "one", "", nil, "pi", math.Pi)
	logMap := suite.logMap()
	suite.Assert().Equal("one", logMap["first"])
	suite.Assert().Equal(math.Pi, logMap["pi"])
	suite.checkNoEmptyAttribute(5, logMap)
}

// TestAttributesNotEmpty tests whether attributes with empty names and non-nil values are logged properly.
//   - Based on the existing behavior of log/slog the field IS created.
//   - https://github.com/golang/go/issues/59282
func (suite *SlogTestSuite) TestAttributesNotEmpty() {
	logger := suite.Logger(infra.SimpleOptions())
	logger.Info(message, "first", "one", "", "NOT NIL", "pi", math.Pi)
	logMap := suite.logMap()
	suite.Assert().Equal("one", logMap["first"])
	suite.Assert().Equal(math.Pi, logMap["pi"])
	suite.Assert().Equal("NOT NIL", logMap[""])
}

// TestAttributesWith tests whether attributes in With() are logged properly.
//   - Implements slogtest "WithAttrs" test.
func (suite *SlogTestSuite) TestAttributesWith() {
	logger := suite.Logger(infra.SimpleOptions())
	logger.With("first", "one", "second", 2).Info(message, "pi", math.Pi)
	logMap := suite.logMap()
	suite.checkFieldCount(6, logMap)
	suite.Assert().Equal("one", logMap["first"])
	suite.Assert().Equal(float64(2), logMap["second"])
	suite.Assert().Equal(math.Pi, logMap["pi"])
}

// TestAttributesWithEmpty tests whether empty attributes in With() are ignored.
//   - Based on the existing behavior of log/slog the field is hot created.
//   - From https://pkg.go.dev/log/slog@master#Handler
func (suite *SlogTestSuite) TestAttributesWithEmpty() {
	logger := suite.Logger(infra.SimpleOptions())
	logger.With("first", "one", "second", 2, "", nil).Info(message, "pi", math.Pi)
	logMap := suite.logMap()
	suite.Assert().Equal("one", logMap["first"])
	suite.Assert().Equal(float64(2), logMap["second"])
	suite.Assert().Equal(math.Pi, logMap["pi"])
	suite.checkNoEmptyAttribute(6, logMap)
}

// TestAttributesWithNotEmpty tests whether attribute with empty names but non-nil value in With() are properly logged.
//   - https://github.com/golang/go/issues/59282
func (suite *SlogTestSuite) TestAttributesWithNotEmpty() {
	logger := suite.Logger(infra.SimpleOptions())
	logger.With("first", "one", "second", 2, "", "NOT NIL").Info(message, "pi", math.Pi)
	logMap := suite.logMap()
	suite.checkFieldCount(7, logMap)
	suite.Assert().Equal("one", logMap["first"])
	suite.Assert().Equal(float64(2), logMap["second"])
	suite.Assert().Equal(math.Pi, logMap["pi"])
	suite.Assert().Equal("NOT NIL", logMap[""])
}

// TestGroup tests the use of a logging group.
//   - Implements slogtest "groups" test.
func (suite *SlogTestSuite) TestGroup() {
	logger := suite.Logger(infra.SimpleOptions())
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

// TestGroupEmpty tests logging an empty group added as a group attribute.
//   - Based on the existing behavior of log/slog the group field is not logged.
//   - Implements slogtest "empty-group" test.
//   - From https://pkg.go.dev/log/slog@master#Handler
func (suite *SlogTestSuite) TestGroupEmpty() {
	logger := suite.Logger(infra.SimpleOptions())
	logger.Info(message, slog.Group("group"))
	logMap := suite.logMap()
	if !suite.HasWarning(warning.GroupEmpty) {
		suite.checkFieldCount(3, logMap)
		_, found := logMap["group"]
		suite.Assert().False(found)
	} else {
		issues := make([]string, 0, 4)
		if len(logMap) > 3 {
			issues = append(issues, "too many fields")
			if grp, found := logMap["group"]; found {
				issues = append(issues, "found field")
				if group, ok := grp.(map[string]any); ok {
					issues = append(issues, "value is group")
					issues = append(issues, fmt.Sprintf("length: %d", len(group)))
				}
			}
		}
		if len(issues) > 0 {
			suite.AddWarning(warning.GroupEmpty, strings.Join(issues, ", "), suite.Buffer.String())
		} else {
			suite.AddUnused(warning.GroupEmpty, "")
		}
	}
}

// TestWithGroupEmpty tests logging an empty group added using WithGroup.
//   - Based on the existing behavior of log/slog the group field is not logged.
//   - From https://pkg.go.dev/log/slog@master#Handler
func (suite *SlogTestSuite) TestWithGroupEmpty() {
	logger := suite.Logger(infra.SimpleOptions())
	logger.WithGroup("group1").WithGroup("group2").Info(message, slog.Group("subGroup"))
	logMap := suite.logMap()
	if !suite.HasWarning(warning.WithGroupEmpty) {
		suite.checkFieldCount(3, logMap)
		_, found := logMap["group"]
		suite.Assert().False(found)
	} else {
		issues := make([]string, 0, 4)
		if len(logMap) > 3 {
			issues = append(issues, "too many fields")
			if grp, found := logMap["group1"]; found {
				issues = append(issues, "found field 'group1'")
				if group, ok := grp.(map[string]any); ok {
					issues = append(issues, "value is group")
					issues = append(issues, fmt.Sprintf("length: %d", len(group)))
				}
			}
		}
		if len(issues) > 0 {
			suite.AddWarning(warning.WithGroupEmpty, strings.Join(issues, ", "), suite.Buffer.String())
		} else {
			suite.AddUnused(warning.WithGroupEmpty, "")
		}
	}
}

// TestGroupInline tests the use of a group with an empty name.
//   - Based on the existing behavior of log/slog the group field is not logged and
//     the fields within the group are moved to the top level.
//   - Implements slogtest "inline-group" test.
//   - From https://pkg.go.dev/log/slog@master#Handler
func (suite *SlogTestSuite) TestGroupInline() {
	logger := suite.Logger(infra.SimpleOptions())
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
	if !suite.HasWarning(warning.GroupInline) {
		suite.checkFieldCount(8, logMap)
		checkFieldFn(logMap)
	} else {
		counter := suite.fieldCounter()
		suite.Require().NoError(counter.Parse())
		if counter.NumFields() == 6 {
			if group, ok := logMap[""].(map[string]any); ok {
				suite.Assert().Len(group, 3)
				checkFieldFn(group)
			} else {
				suite.Fail("Group not map[string]any")
			}
			suite.AddWarning(warning.GroupInline, "", suite.Buffer.String())
		} else {
			suite.AddUnused(warning.GroupInline, "")
		}
	}
}

// TestWithGroupInline tests the use of WithGroup() with an empty name.
//   - Based on the existing behavior of log/slog the group field is not logged and
//     the fields within the group are moved to the top level.
//   - From https://pkg.go.dev/log/slog@master#Handler
func (suite *SlogTestSuite) TestWithGroupInline() {
	logger := suite.Logger(infra.SimpleOptions())
	logger.WithGroup("").Info(message,
		"first", "one",
		"second", 2,
		slog.String("third", "3"),
		"fourth", "forth",
		"pi", math.Pi)
	logMap := suite.logMap()
	fmt.Printf(">>> %s\n", suite.String())
	suite.Assert().Equal("one", logMap["first"])
	suite.Assert().Equal(math.Pi, logMap["pi"])
	checkFieldFn := func(fieldMap map[string]any) {
		suite.Assert().Equal(float64(2), fieldMap["second"])
		suite.Assert().Equal("3", fieldMap["third"])
		suite.Assert().Equal("forth", fieldMap["fourth"])
	}
	if !suite.HasWarning(warning.GroupInline) {
		suite.checkFieldCount(8, logMap)
		checkFieldFn(logMap)
	} else {
		counter := suite.fieldCounter()
		suite.Require().NoError(counter.Parse())
		if counter.NumFields() == 6 {
			if group, ok := logMap[""].(map[string]any); ok {
				suite.Assert().Len(group, 3)
				checkFieldFn(group)
			} else {
				suite.Fail("Group not map[string]any")
			}
			suite.AddWarning(warning.GroupInline, "", suite.Buffer.String())
		} else {
			suite.AddUnused(warning.GroupInline, "")
		}
	}
}

// TestGroupWith tests the use of a logging group specified using WithGroup.
//   - Implements slogtest "WithGroup" test.
func (suite *SlogTestSuite) TestGroupWith() {
	logger := suite.Logger(infra.SimpleOptions())
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

// TestGroupWithMulti tests the use of multiple logging groups.
//   - Implements slogtest "multi-with" test.
func (suite *SlogTestSuite) TestGroupWithMulti() {
	logger := suite.Logger(infra.SimpleOptions())
	logger.With("first", "one").
		WithGroup("group").With("second", 2, "third", "3").
		WithGroup("subGroup").Info(message, "fourth", "forth", "pi", math.Pi)
	logMap := suite.logMap()
	expected := map[string]any{
		"level": "info",
		"msg":   "This is a message",
		"first": "one",
		"group": map[string]any{
			"second": float64(2),
			"third":  "3",
			"subGroup": map[string]any{
				"fourth": "forth",
				"pi":     math.Pi,
			},
		},
	}
	suite.adjustExpected(expected, logMap)
	if !suite.HasWarning(warning.WithGroup) {
		suite.checkFieldCount(5, logMap)
		suite.Assert().Equal(expected, logMap)
	} else if reflect.DeepEqual(expected, logMap) {
		suite.AddUnused(warning.WithGroup, "")
	} else {
		suite.AddWarning(warning.WithGroup, "", suite.Buffer.String())
	}
}

// TestGroupWithMultiSubEmpty tests the use of multiple logging groups when the subgroup is empty.
//   - Based on the existing behavior of log/slog the subgroup field is not logged.
//   - Implements slogtest "empty-group-record" test.
//   - From https://pkg.go.dev/log/slog@master#Handler
func (suite *SlogTestSuite) TestGroupWithMultiSubEmpty() {
	logger := suite.Logger(infra.SimpleOptions())
	logger.With("first", "one").
		WithGroup("group").With("second", 2, "third", "3").
		WithGroup("subGroup").Info(message)
	logMap := suite.logMap()
	expected := map[string]any{
		"level": "info",
		"msg":   "This is a message",
		"first": "one",
		"group": map[string]any{
			"second": float64(2),
			"third":  "3",
			// no subgroup here
		},
	}
	suite.adjustExpected(expected, logMap)
	if !suite.HasWarning(warning.WithGroup) && !suite.HasWarning(warning.GroupEmpty) {
		suite.checkFieldCount(5, logMap)
		suite.Assert().Equal(expected, logMap)
	} else if reflect.DeepEqual(expected, logMap) {
		suite.AddUnused(warning.WithGroup, "")
		suite.AddUnused(warning.GroupEmpty, "")
	} else if grpAny, found := logMap["group"]; !found {
		suite.AddWarning(warning.WithGroup, "no 'group' group", suite.Buffer.String())
	} else if group, ok := grpAny.(map[string]any); !ok {
		suite.AddWarning(warning.WithGroup, "'group' group not map", suite.Buffer.String())
	} else if _, found := group["subGroup"]; found {
		suite.AddWarning(warning.GroupEmpty, "found 'subGroup' group", suite.Buffer.String())
	} else {
		suite.AddWarning(warning.WithGroup, "", suite.Buffer.String())
	}
}

// TestKeys tests whether the three basic keys are present as their defined constants.
//   - Implements slogtest "built-ins" test.
func (suite *SlogTestSuite) TestKeys() {
	logger := suite.Logger(infra.SimpleOptions())
	logger.Info(message)
	logMap := suite.logMap()
	suite.checkFieldCount(3, logMap)
	suite.checkLevelKey("INFO", logMap)
	suite.checkMessageKey(message, logMap)
	suite.Assert().NotNil(logMap[slog.TimeKey])
}

// TestResolveGroup tests logging LogValuer objects within a group.
//   - Implements slogtest "resolve-groups" test.
//   - From https://pkg.go.dev/log/slog@master#Handler
func (suite *SlogTestSuite) TestResolveGroup() {
	logger := suite.Logger(infra.SimpleOptions())
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

// TestResolveGroupWith tests logging LogValuer objects within a group within a With().
//   - Implements slogtest "resolve-WithAttrs-groups" test.
//   - From https://pkg.go.dev/log/slog@master#Handler
func (suite *SlogTestSuite) TestResolveGroupWith() {
	logger := suite.Logger(infra.SimpleOptions())
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

// TestResolveWith tests logging LogValuer objects within a With().
//   - Implements slogtest "resolve-withAttrs" test.
//   - From https://pkg.go.dev/log/slog@master#Handler
func (suite *SlogTestSuite) TestResolveWith() {
	logger := suite.Logger(infra.SimpleOptions())
	logger.With("hidden", &hiddenValuer{v: "value"}).Info(message)
	logMap := suite.logMap()
	suite.checkFieldCount(4, logMap)
	suite.checkResolution("value", logMap["hidden"])
}

// TestResolveValuer tests logging LogValuer objects.
//   - Implements slogtest "resolve" test.
//   - From https://pkg.go.dev/log/slog@master#Handler
func (suite *SlogTestSuite) TestResolveValuer() {
	logger := suite.Logger(infra.SimpleOptions())
	hidden := &hiddenValuer{v: "something"}
	logger.Info(message, "hidden", hidden)
	logMap := suite.logMap()
	suite.checkFieldCount(4, logMap)
	suite.checkResolution("something", logMap["hidden"])
}

// TestZeroPC tests generation of a source key.
//   - Implements slogtest "empty-PC" test.
//   - From https://pkg.go.dev/log/slog@master#Handler
func (suite *SlogTestSuite) TestZeroPC() {
	logger := suite.Logger(infra.SourceOptions())
	record := slog.NewRecord(time.Now(), slog.LevelInfo, message, 0)
	suite.Require().NoError(logger.Handler().Handle(context.Background(), record))
	logMap := suite.logMap()
	suite.checkLevelKey("INFO", logMap)
	suite.checkMessageKey(message, logMap)
	suite.Assert().NotNil(logMap[slog.TimeKey])
	if suite.HasWarning(warning.ZeroPC) {
		if _, ok := logMap[slog.SourceKey].(map[string]any); ok {
			suite.AddWarning(warning.ZeroPC, "", suite.Buffer.String())
			return
		} else if _, ok := logMap["caller"]; ok {
			suite.AddWarning(warning.ZeroPC, "non-standard key 'caller'", suite.Buffer.String())
			return
		}
		suite.AddUnused(warning.ZeroPC, "")
	}

	suite.checkFieldCount(3, logMap)
}

// TestZeroTime tests whether a zero time in a slog.Record is output.
//   - Based on the existing behavior of log/slog the field is not logged.
//   - Implements slogtest "zero-time" test.
//   - From https://pkg.go.dev/log/slog@master#Handler
func (suite *SlogTestSuite) TestZeroTime() {
	logger := suite.Logger(infra.SimpleOptions())
	record := slog.NewRecord(time.Time{}, slog.LevelInfo, message, uintptr(0))
	suite.Require().NoError(logger.Handler().Handle(context.Background(), record))
	logMap := suite.logMap()
	suite.checkLevelKey("INFO", logMap)
	suite.checkMessageKey(message, logMap)
	if suite.HasWarning(warning.ZeroTime) {
		counter := suite.fieldCounter()
		suite.Require().NoError(counter.Parse())
		if counter.NumFields() == 3 {
			if timeAny, found := logMap[slog.TimeKey]; found {
				timeFound, ok := timeAny.(string)
				if !ok {
					timeFound = fmt.Sprintf("<bad type> %v", timeAny)
				}
				suite.AddWarning(warning.ZeroTime, timeFound, suite.Buffer.String())
				return
			}
		}
		suite.AddUnused(warning.ZeroTime, "")
	}
	suite.checkFieldCount(2, logMap)
	suite.Assert().Nil(logMap[slog.TimeKey])
}

// -----------------------------------------------------------------------------
// Tests extending slogTest tests.

// TestAttributeEmptyName tests whether attributes with empty names are logged properly.
//   - Based on the existing behavior of log/slog the field is created with a blank name.
//   - Extension of slogtest "empty-attr" test.
func (suite *SlogTestSuite) TestAttributeEmptyName() {
	logger := suite.Logger(infra.SimpleOptions())
	logger.Info(message, "first", "one", "", 2, "pi", math.Pi)
	logMap := suite.logMap()
	suite.checkFieldCount(6, logMap)
	suite.Assert().Equal("one", logMap["first"])
	value, found := logMap[""]
	suite.Assert().True(found)
	suite.Assert().Equal(float64(2), value)
	suite.Assert().Equal(math.Pi, logMap["pi"])
}

// TestAttributeNil tests whether attributes with nil values are logged properly.
//   - Based on the existing behavior of log/slog the field is created with a nil/null value.
//   - Extension of slogtest "empty-attr" test.
func (suite *SlogTestSuite) TestAttributeNil() {
	logger := suite.Logger(infra.SimpleOptions())
	logger.Info(message, "first", "one", "second", nil, "pi", math.Pi)
	logMap := suite.logMap()
	suite.checkFieldCount(6, logMap)
	suite.Assert().Equal("one", logMap["first"])
	suite.Assert().Nil(logMap["second"])
	suite.Assert().Equal(math.Pi, logMap["pi"])
}

// TestAttributeWithEmpty tests whether attributes with empty names and nil values
// specified in With() are logged properly.
//   - Based on the existing behavior of log/slog the field is not created.
//   - Extension of slogtest "WithAttrs" test.
func (suite *SlogTestSuite) TestAttributeWithEmpty() {
	logger := suite.Logger(infra.SimpleOptions())
	logger.With("", nil).Info(message, "first", "one", "pi", math.Pi)
	logMap := suite.logMap()
	suite.Assert().Equal("one", logMap["first"])
	suite.Assert().Equal(math.Pi, logMap["pi"])
	suite.checkNoEmptyAttribute(5, logMap)
}

// TestAttributeWithEmptyName tests whether With() attributes with empty names are logged properly.
//   - Based on the existing behavior of log/slog the field is created with a blank name.
//   - Extension of slogtest "WithAttrs" test.
func (suite *SlogTestSuite) TestAttributeWithEmptyName() {
	logger := suite.Logger(infra.SimpleOptions())
	logger.With("", 2).Info(message, "first", "one", "pi", math.Pi)
	logMap := suite.logMap()
	suite.checkFieldCount(6, logMap)
	suite.Assert().Equal("one", logMap["first"])
	value, found := logMap[""]
	suite.Assert().True(found)
	suite.Assert().Equal(float64(2), value)
	suite.Assert().Equal(math.Pi, logMap["pi"])
}

// TestAttributeWithNil tests whether With() attributes with nil values are logged properly.
//   - Based on the existing behavior of log/slog the field is created with a nil/null value.
//   - Extension of slogtest "WithAttrs" test.
func (suite *SlogTestSuite) TestAttributeWithNil() {
	logger := suite.Logger(infra.SimpleOptions())
	logger.With("second", nil).Info(message, "first", "one", "pi", math.Pi)
	logMap := suite.logMap()
	suite.checkFieldCount(6, logMap)
	suite.Assert().Equal("one", logMap["first"])
	suite.Assert().Nil(logMap["second"])
	suite.Assert().Equal(math.Pi, logMap["pi"])
}
