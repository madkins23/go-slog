package tests

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"strings"
	"time"
)

// -----------------------------------------------------------------------------
// These tests are intended to mimic: src/testing/slogtest/slogtest.go (2024-01-07).

// TestAttributes tests whether attributes are logged properly.
//   - Implements slogtest "attrs" test.
func (suite *SlogTestSuite) TestAttributes() {
	logger := suite.Logger(SimpleOptions())
	logger.Info(message, "first", "one", "second", 2, "pi", math.Pi)
	logMap := suite.logMap()
	suite.checkFieldCount(6, logMap)
	suite.Assert().Equal("one", logMap["first"])
	suite.Assert().Equal(float64(2), logMap["second"])
	suite.Assert().Equal(math.Pi, logMap["pi"])
}

// TestAttributesEmpty tests whether attributes with empty names and nil values are logged properly.
//   - Based on the existing behavior of log/slog the field is hot created.
//   - Implements slogtest "empty-attr" test.
//   - From https://pkg.go.dev/log/slog@master#Handler
func (suite *SlogTestSuite) TestAttributesEmpty() {
	logger := suite.Logger(SimpleOptions())
	logger.Info(message, "first", "one", "", nil, "pi", math.Pi)
	logMap := suite.logMap()
	suite.Assert().Equal("one", logMap["first"])
	suite.Assert().Equal(math.Pi, logMap["pi"])
	suite.checkNoEmptyAttribute(5, logMap)
}

// TestAttributesWith tests whether attributes in With() are logged properly.
//   - Implements slogtest "WithAttrs" test.
func (suite *SlogTestSuite) TestAttributesWith() {
	logger := suite.Logger(SimpleOptions())
	logger.With("first", "one", "second", 2).Info(message, "pi", math.Pi)
	logMap := suite.logMap()
	suite.checkFieldCount(6, logMap)
	suite.Assert().Equal("one", logMap["first"])
	suite.Assert().Equal(float64(2), logMap["second"])
	suite.Assert().Equal(math.Pi, logMap["pi"])
}

// TestGroup tests the use of a logging group.
//   - Implements slogtest "groups" test.
func (suite *SlogTestSuite) TestGroup() {
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

// TestGroupEmpty tests logging an empty group.
//   - Based on the existing behavior of log/slog the group field is not logged.
//   - Implements slogtest "empty-group" test.
//   - From https://pkg.go.dev/log/slog@master#Handler
func (suite *SlogTestSuite) TestGroupEmpty() {
	logger := suite.Logger(SimpleOptions())
	logger.Info(message, slog.Group("group"))
	logMap := suite.logMap()
	if suite.HasWarning(WarnGroupEmpty) {
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
			suite.addWarning(WarnGroupEmpty, strings.Join(issues, ", "), suite.Buffer.String())
			return
		}
		suite.addWarning(WarnUnused, WarnGroupEmpty, "")
	}
	suite.checkFieldCount(3, logMap)
	_, found := logMap["group"]
	suite.Assert().False(found)
}

// TestGroupInline tests the use of a group with an empty name.
//   - Based on the existing behavior of log/slog the group field is not logged and
//
// the fields within the group are moved to the top level.
//   - Implements slogtest "inline-group" test.
//   - From https://pkg.go.dev/log/slog@master#Handler
func (suite *SlogTestSuite) TestGroupInline() {
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
	if suite.HasWarning(WarnGroupInline) {
		counter := suite.fieldCounter()
		suite.Require().NoError(counter.Parse())
		if counter.NumFields() == 6 {
			if group, ok := logMap[""].(map[string]any); ok {
				suite.Assert().Len(group, 3)
				checkFieldFn(group)
			} else {
				suite.Fail("Group not map[string]any")
			}
			suite.addWarning(WarnGroupInline, "", suite.Buffer.String())
			return
		}
		suite.addWarning(WarnUnused, WarnGroupInline, "")
	}
	suite.checkFieldCount(8, logMap)
	checkFieldFn(logMap)
}

// TestGroupWith tests the use of a logging group specified using WithGroup.
//   - Implements slogtest "WithGroup" test.
func (suite *SlogTestSuite) TestGroupWith() {
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

// TestGroupWithMulti tests the use of multiple logging groups.
//   - Implements slogtest "multi-with" test.
func (suite *SlogTestSuite) TestGroupWithMulti() {
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

// TestGroupWithMultiSubEmpty tests the use of multiple logging groups when the subgroup is empty.
//   - Based on the existing behavior of log/slog the subgroup field is not logged.
//   - Implements slogtest "empty-group-record" test.
//   - From https://pkg.go.dev/log/slog@master#Handler
func (suite *SlogTestSuite) TestGroupWithMultiSubEmpty() {
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
		if suite.HasWarning(WarnGroupEmpty) {
			if len(group) > 2 {
				if subGroup, found := group["subGroup"]; found {
					if sg, ok := subGroup.(map[string]any); ok && len(sg) < 1 {
						suite.addWarning(WarnGroupEmpty, "", suite.Buffer.String())
						return
					}
				}
			}
			suite.addWarning(WarnUnused, WarnGroupEmpty, "")
		}
		suite.Assert().Len(group, 2)
		_, found := group["subGroup"]
		suite.Assert().False(found)
	} else {
		suite.Fail("Group not map[string]any")
	}
}

// TestKeys tests whether the three basic keys are present as their defined constants.
//   - Implements slogtest "built-ins" test.
func (suite *SlogTestSuite) TestKeys() {
	logger := suite.Logger(SimpleOptions())
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

// TestResolveGroupWith tests logging LogValuer objects within a group within a With().
//   - Implements slogtest "resolve-WithAttrs-groups" test.
//   - From https://pkg.go.dev/log/slog@master#Handler
func (suite *SlogTestSuite) TestResolveGroupWith() {
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

// TestResolveWith tests logging LogValuer objects within a With().
//   - Implements slogtest "resolve-withAttrs" test.
//   - From https://pkg.go.dev/log/slog@master#Handler
func (suite *SlogTestSuite) TestResolveWith() {
	logger := suite.Logger(SimpleOptions())
	logger.With("hidden", &hiddenValuer{v: "value"}).Info(message)
	logMap := suite.logMap()
	suite.checkFieldCount(4, logMap)
	suite.checkResolution("value", logMap["hidden"])
}

// TestResolveValuer tests logging LogValuer objects.
//   - Implements slogtest "resolve" test.
//   - From https://pkg.go.dev/log/slog@master#Handler
func (suite *SlogTestSuite) TestResolveValuer() {
	logger := suite.Logger(SimpleOptions())
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
	logger := suite.Logger(SourceOptions())
	record := slog.NewRecord(time.Now(), slog.LevelInfo, message, 0)
	suite.Require().NoError(logger.Handler().Handle(context.Background(), record))
	logMap := suite.logMap()
	suite.checkLevelKey("INFO", logMap)
	suite.checkMessageKey(message, logMap)
	suite.Assert().NotNil(logMap[slog.TimeKey])
	if suite.HasWarning(WarnZeroPC) {
		if _, ok := logMap[slog.SourceKey].(map[string]any); ok {
			suite.addWarning(WarnZeroPC, "", suite.Buffer.String())
			return
		}
		suite.addWarning(WarnUnused, WarnZeroPC, "")
	}

	suite.checkFieldCount(3, logMap)
}

// TestZeroTime tests whether a zero time in a slog.Record is output.
//   - Based on the existing behavior of log/slog the field is not logged.
//   - Implements slogtest "zero-time" test.
//   - From https://pkg.go.dev/log/slog@master#Handler
func (suite *SlogTestSuite) TestZeroTime() {
	logger := suite.Logger(SimpleOptions())
	record := slog.NewRecord(time.Time{}, slog.LevelInfo, message, uintptr(0))
	suite.Require().NoError(logger.Handler().Handle(context.Background(), record))
	logMap := suite.logMap()
	suite.checkLevelKey("INFO", logMap)
	suite.checkMessageKey(message, logMap)
	if suite.HasWarning(WarnZeroTime) {
		counter := suite.fieldCounter()
		suite.Require().NoError(counter.Parse())
		if counter.NumFields() == 3 {
			if timeAny, found := logMap[slog.TimeKey]; found {
				timeZero := suite.parseTime(timeAny)
				suite.Assert().Equal(time.Time{}, timeZero, "time should be zero")
				suite.addWarning(WarnZeroTime, "", suite.Buffer.String())
				return
			}
		}
		suite.addWarning(WarnUnused, WarnZeroTime, "")
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

// TestAttributeNil tests whether attributes with nil values are logged properly.
//   - Based on the existing behavior of log/slog the field is created with a nil/null value.
//   - Extension of slogtest "empty-attr" test.
func (suite *SlogTestSuite) TestAttributeNil() {
	logger := suite.Logger(SimpleOptions())
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
	logger := suite.Logger(SimpleOptions())
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

// TestAttributeWithNil tests whether With() attributes with nil values are logged properly.
//   - Based on the existing behavior of log/slog the field is created with a nil/null value.
//   - Extension of slogtest "WithAttrs" test.
func (suite *SlogTestSuite) TestAttributeWithNil() {
	logger := suite.Logger(SimpleOptions())
	logger.With("second", nil).Info(message, "first", "one", "pi", math.Pi)
	logMap := suite.logMap()
	suite.checkFieldCount(6, logMap)
	suite.Assert().Equal("one", logMap["first"])
	suite.Assert().Nil(logMap["second"])
	suite.Assert().Equal(math.Pi, logMap["pi"])
}
