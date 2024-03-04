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
	suite.Assert().IsType(1.23, logMap[slog.TimeKey])
	suite.Require().Equal(slog.LevelInfo.String(), logMap[slog.LevelKey])
	suite.Require().Equal(message, logMap[slog.MessageKey])
}

func (suite *HandlerTestSuite) TestAttributes() {
	anything := []any{"alpha", "omega"}
	hdlr := suite.newHandler(nil)
	now := time.Now()
	record := slog.NewRecord(time.Now(), slog.LevelInfo, message, 0)
	record.AddAttrs(
		slog.Time("when", now),
		slog.Duration("howLong", time.Minute),
		slog.String("goober", "snoofus"),
		slog.Bool("boolean", true),
		slog.Float64("pi", math.Pi),
		slog.Int("skidoo", 23),
		slog.Int64("minus", -64),
		slog.Uint64("unsigned", 79),
		slog.Any("any", anything),
		slog.Group("group",
			slog.String("name", "Beatles"),
			infra.EmptyAttr(),
			slog.Float64("pi", math.Pi)))
	suite.Assert().NoError(hdlr.Handle(context.Background(), record))
	logMap := suite.logMap()
	// Basic fields tested in Test_Enabled.
	suite.Assert().Len(logMap, 13)
	suite.Assert().Equal(float64(now.Nanosecond()), logMap["when"])
	suite.Assert().Equal(float64(time.Minute.Nanoseconds()), logMap["howLong"])
	suite.Assert().Equal("snoofus", logMap["goober"])
	suite.Assert().Equal(true, logMap["boolean"])
	suite.Assert().Equal(math.Pi, logMap["pi"])
	suite.Assert().Equal(float64(23), logMap["skidoo"])
	suite.Assert().Equal(float64(-64), logMap["minus"])
	suite.Assert().Equal(float64(79), logMap["unsigned"])
	suite.Assert().Equal(anything, logMap["any"])
	grp, found := logMap["group"]
	suite.Assert().True(found)
	group, ok := grp.(map[string]any)
	suite.Assert().True(ok)
	suite.Assert().Len(group, 2)
	suite.Assert().Equal("Beatles", group["name"])
	suite.Assert().Equal(math.Pi, group["pi"])
}
