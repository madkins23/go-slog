package tests

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"runtime"
	"time"
)

// -----------------------------------------------------------------------------
// Tests created from reviewing log/slog documentation:
//  https://pkg.go.dev/log/slog@master#Handler
//  https://github.com/golang/example/blob/master/slog-handler-guide/README.md

// TestCancelledContext verifies that a cancelled context will not affect logging.
//   - https://github.com/golang/example/blob/master/slog-handler-guide/README.md
func (suite *SlogTestSuite) TestCancelledContext() {
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

// TestDefaultLevel tests whether the handler under test
// is created by default with slog.LevelInfo.
//   - Implied by https://pkg.go.dev/log/slog@master#Handler:
//     "First, we wanted the default level to be Info,
//     Since Levels are ints, Info is the default value for int, zero."
func (suite *SlogTestSuite) TestDefaultLevel() {
	for _, options := range []*slog.HandlerOptions{
		{},
		{AddSource: true},
	} {
		ctx := context.Background()
		logger := suite.Logger(options)
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
				where := ""
				if options.AddSource {
					where = " with AddSource"
				}
				suite.addWarning(WarnDefaultLevel,
					fmt.Sprintf("defaultlevel%s is '%s'", where, name), false)
				continue
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
}

// TestDerivedInvariantWith tests to see if
// deriving another handler via With() changes the original handler.
//   - https://github.com/golang/example/blob/master/slog-handler-guide/README.md
func (suite *SlogTestSuite) TestDerivedInvariantWith() {
	simpleLogger := suite.Logger(SimpleOptions())
	simpleLogger.Info(message)
	origLogMap := suite.logMap()
	delete(origLogMap, slog.TimeKey)
	suite.bufferReset()
	withLogger := simpleLogger.With("alpha", "omega")
	withLogger.Info(message)
	suite.bufferReset()
	simpleLogger.Info(message)
	currLogMap := suite.logMap()
	delete(currLogMap, slog.TimeKey)
	suite.Assert().Equal(origLogMap, currLogMap)
}

// TestDerivedInvariantWithGroup tests to see if
// deriving another handler via WithGroup() changes the original handler.
//   - https://github.com/golang/example/blob/master/slog-handler-guide/README.md
func (suite *SlogTestSuite) TestDerivedInvariantWithGroup() {
	simpleLogger := suite.Logger(SimpleOptions())
	simpleLogger.Info(message)
	origLogMap := suite.logMap()
	delete(origLogMap, slog.TimeKey)
	suite.bufferReset()
	withGroupLogger := simpleLogger.With("alpha", "omega")
	withGroupLogger.Info(message)
	suite.bufferReset()
	simpleLogger.Info(message)
	currLogMap := suite.logMap()
	delete(currLogMap, slog.TimeKey)
	suite.Assert().Equal(origLogMap, currLogMap)
}

// TestDisabled tests whether logging is disabled by level.
//   - https://pkg.go.dev/log/slog@master#hdr-Levels
func (suite *SlogTestSuite) TestDisabled() {
	logger := suite.Logger(SimpleOptions())
	logger.Debug(message)
	suite.Assert().Empty(suite.Buffer)
}

// TestKey tests generation of a source key.
//   - https://pkg.go.dev/log/slog@master#HandlerOptions
//   - https://pkg.go.dev/log/slog@master#Source
func (suite *SlogTestSuite) TestKey() {
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

// TestKeyCase tests whether level keys are properly cased.
//   - Based on the existing behavior of log/slog they should be uppercase.
func (suite *SlogTestSuite) TestKeyCase() {
	ctx := context.Background()
	logger := suite.Logger(LevelOptions(slog.LevelDebug))
	for name, level := range logLevels {
		logger.Log(ctx, level, message)
		logMap := suite.logMap()
		suite.checkLevelKey(name, logMap)
		suite.bufferReset()
	}
}

// TestLevelVar tests the use of a slog.LevelVar.
//   - https://pkg.go.dev/log/slog@master#hdr-Levels
//   - https://pkg.go.dev/log/slog@master#LevelVar
func (suite *SlogTestSuite) TestLevelVar() {
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

// TestLogAttributes tests the LogAttrs call with all attribute objects.
//   - https://pkg.go.dev/log/slog@master#Logger.LogAttrs
//   - https://pkg.go.dev/log/slog@master#Attr
func (suite *SlogTestSuite) TestLogAttributes() {
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