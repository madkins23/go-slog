package test

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"runtime"
	"time"
)

// -----------------------------------------------------------------------------
// Tests created from reviewing log/slog documentation.

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

// TestSimpleDefaultLevel tests whether the simple logger is created by default with slog.LevelInfo.
// Other tests (e.g. TestSimpleDisabled) depend on this.
func (suite *SlogTestSuite) TestSimpleDefaultLevel() {
	ctx := context.Background()
	logger := suite.Logger(&slog.HandlerOptions{})
	if suite.hasWarning(WarnDefaultLevel) {
		level := slog.Level(100)
		name := ""

		for _, logLevel := range logLevels {
			lvl := logLevel.Level()
			if logger.Enabled(ctx, lvl) {
				if lvl < level {
					level = lvl
					name = logLevel.String()
				}
			}
		}
		if name != "" {
			suite.addWarning(WarnDefaultLevel, fmt.Sprintf("defaultlevel is '%s'", name), false)
			return
		}
		suite.addWarning(WarnUnused, WarnDefaultLevel, false)
	}
	suite.Assert().False(logger.Enabled(ctx, slog.LevelDebug-1))
	suite.Assert().False(logger.Enabled(ctx, slog.LevelDebug))
	suite.Assert().False(logger.Enabled(ctx, slog.LevelInfo-1))
	suite.Assert().True(logger.Enabled(ctx, slog.LevelInfo))
	suite.Assert().True(logger.Enabled(ctx, slog.LevelInfo+1))
	suite.Assert().True(logger.Enabled(ctx, slog.LevelWarn))
	suite.Assert().True(logger.Enabled(ctx, slog.LevelError))
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
	if suite.hasWarning(WarnNanoTime) {
		// Some handlers log times as RFC3339 instead of RFC3339Nano
		suite.Equal(t.Format(time.RFC3339), when)
	} else {
		// Based on the existing behavior of log/slog it should be RFC3339Nano.
		suite.Equal(t.Format(time.RFC3339Nano), when)
	}
	howLong, ok := logMap["howLong"].(float64)
	suite.True(ok)
	if suite.hasWarning(WarnNanoDuration) {
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
	ctx := context.Background()
	logger := suite.Logger(LevelOptions(slog.LevelDebug))
	for name, level := range logLevels {
		logger.Log(ctx, level, message)
		logMap := suite.logMap()
		suite.checkLevelKey(name, logMap)
		suite.bufferReset()
	}
}

// TestSimpleLevelVar tests the use of a slog.LevelVar.
func (suite *SlogTestSuite) TestSimpleLevelVar() {
	ctx := context.Background()
	var programLevel = new(slog.LevelVar)
	logger := suite.Logger(LevelOptions(programLevel))
	// Should be INFO by default.
	suite.Assert().Equal(slog.LevelInfo, programLevel.Level())
	suite.Assert().False(logger.Enabled(ctx, -1))
	suite.Assert().True(logger.Enabled(ctx, slog.LevelInfo))
	suite.Assert().True(logger.Enabled(ctx, 1))
	// Change the level.
	programLevel.Set(slog.LevelWarn)
	suite.Assert().Equal(slog.LevelWarn, programLevel.Level())
	suite.Assert().False(logger.Enabled(ctx, 3))
	suite.Assert().True(logger.Enabled(ctx, slog.LevelWarn))
	suite.Assert().True(logger.Enabled(ctx, 5))
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
