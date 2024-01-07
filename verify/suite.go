package verify

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"math"
	"time"

	"github.com/stretchr/testify/suite"

	testJSON "github.com/madkins23/go-slog/json"
	"github.com/madkins23/go-slog/test"
)

const (
	message = "This is a message"
)

// SlogTestSuite provides various tests for a specified log/slog.Hander.
type SlogTestSuite struct {
	suite.Suite
	*bytes.Buffer
}

func (suite *SlogTestSuite) SetupTest() {
	suite.Buffer = &bytes.Buffer{}
}

// SimpleHandler returns a simple handler to be used in an slog.Logger.
// Override this method to test other types of slog JSON handlers.
// TODO: Can this be moved into simpleLogger()?
func (suite *SlogTestSuite) SimpleHandler() slog.Handler {
	return slog.NewJSONHandler(suite.Buffer, nil)
}

func (suite *SlogTestSuite) simpleLogger() *slog.Logger {
	return slog.New(suite.SimpleHandler())
}

// TestSimpleAttributes tests whether attributes are logged properly.
func (suite *SlogTestSuite) TestSimpleAttributes() {
	logger := suite.simpleLogger()
	logger.Info(message, "first", "one", "second", 2, "pi", math.Pi)
	logMap := suite.logMap()
	suite.assertFieldCount(6, logMap)
	suite.Assert().Equal("one", logMap["first"])
	suite.Assert().Equal(float64(2), logMap["second"])
	suite.Assert().Equal(math.Pi, logMap["pi"])
}

// TestSimpleAttributeDuplicate tests whether duplicate attributes are logged properly.
// Based on the existing behavior of log/slog the second occurrence overrides the first.
func (suite *SlogTestSuite) TestSimpleAttributeDuplicate() {
	logger := suite.simpleLogger()
	logger.Info(message, "alpha", "one", "alpha", 2)
	logMap := suite.logMap()
	suite.Assert().Len(logMap, 4)
	counter := suite.fieldCounter()
	suite.Require().NoError(counter.Parse())
	suite.Assert().Equal(uint(4), counter.NumFields())
	duplicates := counter.Duplicates()
	suite.Assert().Len(duplicates, 1)
	suite.Assert().Equal(uint(2), duplicates["alpha"])
	suite.Assert().Equal(float64(2), logMap["alpha"])
}

// TestSimpleAttributeEmpty tests whether attributes with empty names and nil values are logged properly.
// Based on the existing behavior of log/slog the field is hot created.
func (suite *SlogTestSuite) TestSimpleAttributeEmpty() {
	logger := suite.simpleLogger()
	logger.Info(message, "first", "one", "", nil, "pi", math.Pi)
	logMap := suite.logMap()
	suite.assertFieldCount(5, logMap)
	suite.Assert().Equal("one", logMap["first"])
	_, found := logMap[""]
	suite.Assert().False(found)
	suite.Assert().Equal(math.Pi, logMap["pi"])
}

// TestSimpleAttributeEmptyName tests whether attributes with empty names are logged properly.
// Based on the existing behavior of log/slog the field is created with a blank name.
func (suite *SlogTestSuite) TestSimpleAttributeEmptyName() {
	logger := suite.simpleLogger()
	logger.Info(message, "first", "one", "", 2, "pi", math.Pi)
	logMap := suite.logMap()
	suite.assertFieldCount(6, logMap)
	suite.Assert().Equal("one", logMap["first"])
	value, found := logMap[""]
	suite.Assert().True(found)
	suite.Assert().Equal(float64(2), value)
	suite.Assert().Equal(math.Pi, logMap["pi"])
}

// TestSimpleAttributeNil tests whether attributes with nil values are logged properly.
// Based on the existing behavior of log/slog the field is created with a nil/null value.
func (suite *SlogTestSuite) TestSimpleAttributeNil() {
	logger := suite.simpleLogger()
	logger.Info(message, "first", "one", "second", nil, "pi", math.Pi)
	logMap := suite.logMap()
	suite.assertFieldCount(6, logMap)
	suite.Assert().Equal("one", logMap["first"])
	suite.Assert().Nil(logMap["second"])
	suite.Assert().Equal(math.Pi, logMap["pi"])
}

// TestSimpleAttributesWith tests whether attributes in With() are logged properly.
func (suite *SlogTestSuite) TestSimpleAttributesWith() {
	logger := suite.simpleLogger()
	logger.With("first", "one", "second", 2).Info(message, "pi", math.Pi)
	logMap := suite.logMap()
	suite.assertFieldCount(6, logMap)
	suite.Assert().Equal("one", logMap["first"])
	suite.Assert().Equal(float64(2), logMap["second"])
	suite.Assert().Equal(math.Pi, logMap["pi"])
}

// TestSimpleAttributeWithDuplicate tests whether duplicate attributes are logged properly
// when the duplicate is introduced by With() and then the main call.
// Based on the existing behavior of log/slog the second occurrence overrides the first.
func (suite *SlogTestSuite) TestSimpleAttributeWithDuplicate() {
	logger := suite.simpleLogger()
	logger.With("alpha", "one").Info(message, "alpha", 2)
	logMap := suite.logMap()
	suite.Assert().Len(logMap, 4)
	counter := suite.fieldCounter()
	suite.Require().NoError(counter.Parse())
	suite.Assert().Equal(uint(4), counter.NumFields())
	duplicates := counter.Duplicates()
	suite.Assert().Len(duplicates, 1)
	suite.Assert().Equal(uint(2), duplicates["alpha"])
	suite.Assert().Equal(float64(2), logMap["alpha"])
}

// TestSimpleAttributeWithEmpty tests whether attributes with empty names and nil values
// specified in With() are logged properly.
// Based on the existing behavior of log/slog the field is hot created.
func (suite *SlogTestSuite) TestSimpleAttributeWithEmpty() {
	logger := suite.simpleLogger()
	logger.With("", nil).Info(message, "first", "one", "pi", math.Pi)
	logMap := suite.logMap()
	suite.assertFieldCount(5, logMap)
	suite.Assert().Equal("one", logMap["first"])
	_, found := logMap[""]
	suite.Assert().False(found)
	suite.Assert().Equal(math.Pi, logMap["pi"])
}

// TestSimpleAttributeWithEmptyName tests whether With() attributes with empty names are logged properly.
// Based on the existing behavior of log/slog the field is created with a blank name.
func (suite *SlogTestSuite) TestSimpleAttributeWithEmptyName() {
	logger := suite.simpleLogger()
	logger.With("", 2).Info(message, "first", "one", "pi", math.Pi)
	logMap := suite.logMap()
	suite.assertFieldCount(6, logMap)
	suite.Assert().Equal("one", logMap["first"])
	value, found := logMap[""]
	suite.Assert().True(found)
	suite.Assert().Equal(float64(2), value)
	suite.Assert().Equal(math.Pi, logMap["pi"])
}

// TestSimpleAttributeWithNil tests whether With() attributes with nil values are logged properly.
// Based on the existing behavior of log/slog the field is created with a nil/null value.
func (suite *SlogTestSuite) TestSimpleAttributeWithNil() {
	logger := suite.simpleLogger()
	logger.With("second", nil).Info(message, "first", "one", "pi", math.Pi)
	logMap := suite.logMap()
	suite.assertFieldCount(6, logMap)
	suite.Assert().Equal("one", logMap["first"])
	suite.Assert().Nil(logMap["second"])
	suite.Assert().Equal(math.Pi, logMap["pi"])
}

// TestSimpleDisabled tests whether logging is disabled by level.
func (suite *SlogTestSuite) TestSimpleDisabled() {
	logger := suite.simpleLogger()
	logger.Debug(message)
	suite.Assert().Empty(suite.Buffer)
}

// TestSimpleGroup tests the use of a logging group.
func (suite *SlogTestSuite) TestSimpleGroup() {
	logger := suite.simpleLogger()
	logger.Info(message, "first", "one",
		slog.Group("group", "second", 2, slog.String("third", "3"), "fourth", "forth"),
		"pi", math.Pi)
	logMap := suite.logMap()
	suite.assertFieldCount(6, logMap)
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
func (suite *SlogTestSuite) TestSimpleGroupEmpty() {
	logger := suite.simpleLogger()
	logger.Info(message, slog.Group("group"))
	logMap := suite.logMap()
	suite.assertFieldCount(3, logMap)
	_, found := logMap["group"]
	suite.Assert().False(found)
}

// TestSimpleGroupInline tests the use of a group with an empty name.
// Based on the existing behavior of log/slog the group field is not logged and
// the fields within the group are moved to the top level.
func (suite *SlogTestSuite) TestSimpleGroupInline() {
	logger := suite.simpleLogger()
	logger.Info(message, "first", "one",
		slog.Group("", "second", 2, slog.String("third", "3"), "fourth", "forth"),
		"pi", math.Pi)
	logMap := suite.logMap()
	suite.assertFieldCount(8, logMap)
	suite.Assert().Equal("one", logMap["first"])
	suite.Assert().Equal(float64(2), logMap["second"])
	suite.Assert().Equal("3", logMap["third"])
	suite.Assert().Equal("forth", logMap["fourth"])
	suite.Assert().Equal(math.Pi, logMap["pi"])
}

// TestSimpleKeys tests whether the three basic keys are present as their defined constants.
func (suite *SlogTestSuite) TestSimpleKeys() {
	logger := suite.simpleLogger()
	logger.Info(message)
	logMap := suite.logMap()
	suite.assertFieldCount(3, logMap)
	suite.Assert().Equal("INFO", logMap[slog.LevelKey])
	suite.Assert().Equal(message, logMap[slog.MessageKey])
	suite.Assert().NotNil(logMap[slog.TimeKey])
}

// TestSimpleLevel tests whether the simple logger is created with slog.LevelInfo.
// Other tests (e.g. TestSimpleDisabled) depend on this.
func (suite *SlogTestSuite) TestSimpleLevel() {
	logger := suite.simpleLogger()
	suite.Assert().False(logger.Enabled(context.TODO(), -1))
	suite.Assert().True(logger.Enabled(context.TODO(), slog.LevelInfo))
	suite.Assert().True(logger.Enabled(context.TODO(), 1))
	suite.Assert().True(logger.Enabled(context.TODO(), slog.LevelWarn))
	suite.Assert().True(logger.Enabled(context.TODO(), slog.LevelError))
}

// TestSimpleZeroTime tests whether a zero time in a slog.Record is output.
// Based on the existing behavior of log/slog the field is not logged.
func (suite *SlogTestSuite) TestSimpleZeroTime() {
	logger := suite.simpleLogger()
	record := slog.NewRecord(time.Time{}, slog.LevelInfo, message, uintptr(0))
	suite.Require().NoError(logger.Handler().Handle(context.TODO(), record))
	logMap := suite.logMap()
	suite.assertFieldCount(2, logMap)
	suite.Assert().Equal("INFO", logMap[slog.LevelKey])
	suite.Assert().Equal(message, logMap[slog.MessageKey])
}

// -----------------------------------------------------------------------------

// assertFieldCount checks whether the prescribed number of fields exist at the top level.
// In addition to using the logMap generated by unmarshaling the JSON log data,
// the custom test.FieldCounter is used to make sure there are no duplicates.
func (suite *SlogTestSuite) assertFieldCount(count int, logMap map[string]any) {
	suite.Assert().Len(logMap, count)
	// Double check to make sure there are no duplicate fields at the top level.
	counter := suite.fieldCounter()
	suite.Require().NoError(counter.Parse())
	suite.Assert().Equal(uint(count), counter.NumFields())
	suite.Assert().Empty(counter.Duplicates())
}

// -----------------------------------------------------------------------------

func (suite *SlogTestSuite) fieldCounter() *testJSON.FieldCounter {
	return testJSON.NewFieldCounter(suite.Buffer.Bytes())
}

func (suite *SlogTestSuite) logMap() map[string]any {
	test.Debugf(1, "JSON: %s", suite.Buffer.Bytes())
	var results map[string]any
	suite.Require().NoError(json.Unmarshal(suite.Buffer.Bytes(), &results))
	return results
}

// -----------------------------------------------------------------------------

func (suite *SlogTestSuite) TestXXX() {
	test.Debugf(0, "Dummy test")
}
