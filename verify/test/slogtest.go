package test

import (
	"context"
	"log/slog"
	"math"
	"time"
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
	if suite.hasWarning(WarnGroupInline) {
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
		if suite.hasWarning(WarnSubgroupEmpty) {
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
	if suite.hasWarning(WarnZeroTime) {
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
	if suite.hasWarning(WarnZeroPC) {
		if _, ok := logMap[slog.SourceKey].(map[string]any); ok {
			suite.addWarning(WarnZeroPC, "", true)
			return
		}
		suite.addWarning(WarnUnused, WarnZeroPC, false)
	}

	suite.checkFieldCount(3, logMap)
}
