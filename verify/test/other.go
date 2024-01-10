package test

import (
	"context"
	"log/slog"
	"time"
)

// -----------------------------------------------------------------------------
// Other tests.

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
