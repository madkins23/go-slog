package sloggy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/madkins23/go-slog/handlers/sloggy/test"
	"github.com/madkins23/go-slog/infra"
)

const (
	message = "This is a message. No, really!"
)

// -----------------------------------------------------------------------------
// Top level definitions.

// HandlerTestSuite provides various tests for a specified log/slog.Handler.
type HandlerTestSuite struct {
	suite.Suite
	*bytes.Buffer
}

func NewHandlerTestSuite() *HandlerTestSuite {
	return &HandlerTestSuite{}
}

func TestHandlerSuite(t *testing.T) {
	suite.Run(t, NewHandlerTestSuite())
}

// -----------------------------------------------------------------------------
// Suite test configuration.

func (suite *HandlerTestSuite) SetupTest() {
	suite.Buffer = &bytes.Buffer{}
}

// -----------------------------------------------------------------------------

// logMap unmarshals JSON in the output capture buffer into a map[string]any.
// The buffer is sent to test logging output if the -debug=<level> flag is >= 1.
func (suite *HandlerTestSuite) logMap() map[string]any {
	var results map[string]any
	err := json.Unmarshal(suite.Bytes(), &results)
	if err != nil {
		err = fmt.Errorf("unmarshal '%s': %w", suite.Bytes(), err)
	}
	suite.Require().NoError(err)
	return results
}

func (suite *HandlerTestSuite) newHandler(options *slog.HandlerOptions) *Handler {
	hdlr := NewHandler(suite.Buffer, options)
	suite.Require().NotNil(hdlr)
	return hdlr
}

// -----------------------------------------------------------------------------

func (suite *HandlerTestSuite) TestEnabled() {
	ctx := context.Background()
	hdlr := suite.newHandler(nil)
	suite.Assert().False(hdlr.Enabled(ctx, slog.LevelDebug-1))
	suite.Assert().False(hdlr.Enabled(ctx, slog.LevelDebug))
	suite.Assert().False(hdlr.Enabled(ctx, slog.LevelDebug+1))
	suite.Assert().False(hdlr.Enabled(ctx, slog.LevelInfo-1))
	suite.Assert().True(hdlr.Enabled(ctx, slog.LevelInfo))
	suite.Assert().True(hdlr.Enabled(ctx, slog.LevelInfo+1))
	suite.Assert().True(hdlr.Enabled(ctx, slog.LevelWarn-1))
	suite.Assert().True(hdlr.Enabled(ctx, slog.LevelWarn))
	suite.Assert().True(hdlr.Enabled(ctx, slog.LevelWarn+1))
	suite.Assert().True(hdlr.Enabled(ctx, slog.LevelError-1))
	suite.Assert().True(hdlr.Enabled(ctx, slog.LevelError))
	suite.Assert().True(hdlr.Enabled(ctx, slog.LevelError+1))
}

func (suite *HandlerTestSuite) TestBasicAttributes() {
	hdlr := suite.newHandler(nil)
	suite.Assert().NoError(hdlr.Handle(context.Background(),
		slog.NewRecord(time.Now(), slog.LevelInfo, message, 0)))
	logMap := suite.logMap()
	suite.Assert().IsType("string", logMap[slog.TimeKey])
	suite.Require().Equal(slog.LevelInfo.String(), logMap[slog.LevelKey])
	suite.Require().Equal(message, logMap[slog.MessageKey])
}

func (suite *HandlerTestSuite) TestAttributes() {
	hdlr := suite.newHandler(nil)
	record := slog.NewRecord(time.Now(), slog.LevelInfo, message, 0)
	record.AddAttrs(test.Attributes...)
	suite.Assert().NoError(hdlr.Handle(context.Background(), record))
	logMap := suite.logMap()
	// Basic fields tested in Test_Enabled.
	delete(logMap, slog.TimeKey)
	delete(logMap, slog.LevelKey)
	delete(logMap, slog.MessageKey)
	suite.Assert().Equal(test.AttributeMap, logMap)
}

func (suite *HandlerTestSuite) TestWithAttrs() {
	hdlr := suite.newHandler(nil).
		WithAttrs([]slog.Attr{
			slog.String("make", "Ford"),
			infra.EmptyAttr(),
			slog.Int("year", 1957)})
	record := slog.NewRecord(time.Now(), slog.LevelInfo, message, 0)
	suite.Assert().NoError(hdlr.Handle(context.Background(), record))
	logMap := suite.logMap()
	suite.Assert().Len(logMap, 5)
	// Basic fields tested in Test_Enabled.
	suite.Assert().Equal("Ford", logMap["make"])
	suite.Assert().Equal(float64(1957), logMap["year"])
	// Add another layer.
	hdlr = hdlr.WithAttrs([]slog.Attr{
		infra.EmptyAttr(),
		slog.Float64("price", 3456.98),
		slog.String("owner", "Elvis Presley"),
		infra.EmptyAttr()})
	suite.Reset()
	suite.Assert().NoError(hdlr.Handle(context.Background(), record))
	logMap = suite.logMap()
	suite.Assert().Len(logMap, 7)
	suite.Assert().Equal("Ford", logMap["make"])
	suite.Assert().Equal(float64(1957), logMap["year"])
	suite.Assert().Equal(3456.98, logMap["price"])
	suite.Assert().Equal("Elvis Presley", logMap["owner"])
}

func (suite *HandlerTestSuite) TestWithGroup() {
	hdlr := suite.newHandler(nil).WithGroup("group")
	record := slog.NewRecord(time.Now(), slog.LevelInfo, message, 0)
	record.AddAttrs(
		infra.EmptyAttr(),
		slog.String("Goober", "Snoofus"),
		infra.EmptyAttr(),
		slog.Float64("pi", math.Pi),
		infra.EmptyAttr())
	suite.Assert().NoError(hdlr.Handle(context.Background(), record))
	logMap := suite.logMap()
	suite.Assert().Len(logMap, 4)
	// Basic fields tested in Test_Enabled.
	grp, found := logMap["group"]
	suite.Assert().True(found)
	group, ok := grp.(map[string]any)
	suite.Assert().True(ok)
	suite.Assert().Len(group, 2)
	suite.Assert().Equal("Snoofus", group["Goober"])
	suite.Assert().Equal(math.Pi, group["pi"])
	// Add another layer.
	hdlr = hdlr.WithGroup("subGroup")
	suite.Reset()
	suite.Assert().NoError(hdlr.Handle(context.Background(), record))
	logMap = suite.logMap()
	suite.Assert().Len(logMap, 4)
	grp, found = logMap["group"]
	suite.Assert().True(found)
	group, ok = grp.(map[string]any)
	suite.Assert().True(ok)
	suite.Assert().Len(group, 1)
	sub, found := group["subGroup"]
	suite.Assert().True(found)
	subGroup, ok := sub.(map[string]any)
	suite.Assert().True(ok)
	suite.Assert().Len(subGroup, 2)
	suite.Assert().Equal("Snoofus", subGroup["Goober"])
	suite.Assert().Equal(math.Pi, subGroup["pi"])
}

func (suite *HandlerTestSuite) TestWithGroupAttr() {
	hdlr := suite.newHandler(nil).
		WithAttrs([]slog.Attr{slog.String("first", "one")}).
		WithGroup("group").
		WithAttrs([]slog.Attr{slog.Int("second", 2), slog.String("third", "3")}).
		WithGroup("subGroup")
	record := slog.NewRecord(time.Now(), slog.LevelInfo, message, 0)
	record.AddAttrs(
		slog.String("fourth", "forth"),
		slog.Float64("pi", math.Pi))
	suite.Assert().NoError(hdlr.Handle(context.Background(), record))
	logMap := suite.logMap()
	suite.Assert().Len(logMap, 5)
	// Basic fields tested in Test_Enabled.
	suite.Assert().Equal("one", logMap["first"])
	grp, found := logMap["group"]
	suite.Assert().True(found)
	group, ok := grp.(map[string]any)
	suite.Assert().True(ok)
	suite.Assert().Len(group, 3)
	suite.Assert().Equal(float64(2), group["second"])
	suite.Assert().Equal("3", group["third"])
	sub, found := group["subGroup"]
	suite.Assert().True(found)
	subGroup, ok := sub.(map[string]any)
	suite.Assert().True(ok)
	suite.Assert().Len(subGroup, 2)
	suite.Assert().Equal("forth", subGroup["fourth"])
	suite.Assert().Equal(math.Pi, subGroup["pi"])
}

func (suite *HandlerTestSuite) TestWithGroupAttrSubEmpty() {
	hdlr := suite.newHandler(nil).
		WithAttrs([]slog.Attr{slog.String("first", "one")}).
		WithGroup("group").
		WithAttrs([]slog.Attr{slog.Int("second", 2), slog.String("third", "3")}).
		WithGroup("subGroup")
	record := slog.NewRecord(time.Now(), slog.LevelInfo, message, 0)
	suite.Assert().NoError(hdlr.Handle(context.Background(), record))
	logMap := suite.logMap()
	suite.Assert().Len(logMap, 5)
	// Basic fields tested in Test_Enabled.
	suite.Assert().Equal("one", logMap["first"])
	grp, found := logMap["group"]
	suite.Assert().True(found)
	group, ok := grp.(map[string]any)
	suite.Assert().True(ok)
	suite.Assert().Len(group, 2)
	suite.Assert().Equal(float64(2), group["second"])
	suite.Assert().Equal("3", group["third"])
}
