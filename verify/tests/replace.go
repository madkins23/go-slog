package tests

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/madkins23/go-slog/infra"
	"github.com/madkins23/go-slog/internal/warning"
	"github.com/madkins23/go-slog/replace"
)

// -----------------------------------------------------------------------------
// Tests of slog.HandlerOptions.ReplaceAttr functionality.

// TestReplaceAttr tests the use of HandlerOptions.ReplaceAttr.
//   - https://pkg.go.dev/log/slog@master#HandlerOptions
func (suite *SlogTestSuite) TestReplaceAttr() {
	logger := suite.Logger(infra.ReplaceAttrOptions(func(groups []string, a slog.Attr) slog.Attr {
		switch a.Key {
		case "alpha":
			return slog.String(a.Key, "omega")
		case "change":
			return slog.String("bravo", a.Value.String())
		case "remove":
			return infra.EmptyAttr()
		}
		return a
	}))
	logger.Info(message, "alpha", "beta", "change", "my key", "remove", "me")
	logMap := suite.logMap()
	if !suite.HasWarning(warning.NoReplAttr) {
		suite.Assert().Equal("omega", logMap["alpha"])
		suite.Assert().Equal("my key", logMap["bravo"])
		suite.Assert().Nil(logMap["remove"])
	} else {
		issues := make([]string, 0, 4)
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
			suite.AddWarning(warning.NoReplAttr, strings.Join(issues, ", "), "")
			return
		}
		suite.AddUnused(warning.NoReplAttr, suite.Buffer.String())
	}
	if suite.HasWarning(warning.EmptyAttributes) {
		suite.checkFieldCount(6, logMap)
	} else {
		suite.checkFieldCount(5, logMap)
	}
}

// TestReplaceAttrBasic tests the use of HandlerOptions.ReplaceAttr
// on basic attributes (time, level, message, source).
//   - https://pkg.go.dev/log/slog@master#HandlerOptions
func (suite *SlogTestSuite) TestReplaceAttrBasic() {
	logger := suite.Logger(&slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			switch a.Key {
			case slog.TimeKey:
				return infra.EmptyAttr()
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
	warnings := suite.HasWarnings(warning.NoReplAttr, warning.NoReplAttrBasic)
	if len(warnings) > 0 {
		issues := make([]string, 0, 5)
		if len(logMap) > 3 {
			issues = append(issues, fmt.Sprintf("too many attributes: %d", len(logMap)))
		}
		if logMap[slog.TimeKey] != nil {
			issues = append(issues, slog.TimeKey+" field still exists")
		}
		if logMap[slog.MessageKey] != nil {
			issues = append(issues, slog.MessageKey+" field still exists")
		} else if suite.HasWarning(warning.MessageKey) && logMap["message"] != nil {
			issues = append(issues, "message field still exists")
		}
		// TODO: This one may still work, in samber it's apparently a separate field from basic.
		if value, ok := logMap[slog.SourceKey].(string); !ok || value != "all knowledge" {
			issues = append(issues, fmt.Sprintf("%s == %v", slog.SourceKey, logMap[slog.SourceKey]))
		}
		if len(issues) > 0 {
			suite.AddWarning(warnings[0], strings.Join(issues, ", "), "")
			return
		}
		suite.AddUnused(warnings[0], suite.Buffer.String())
	}
	suite.checkFieldCount(3, logMap)
	suite.Assert().Nil(logMap[slog.TimeKey])
	suite.Assert().Equal("Tilted", logMap[slog.LevelKey])
	suite.Assert().Equal(message, logMap["Message"])
	suite.Assert().Equal("all knowledge", logMap[slog.SourceKey])
}

// -----------------------------------------------------------------------------
// Tests of go-slog/replace ReplaceAttr functions.

// TestReplaceAttrFnLevelCase tests the Level[Lower,Upper]Case functions.
func (suite *SlogTestSuite) TestReplaceAttrFnLevelCase() {
	start := "INFO"
	fixed := "info"
	attrFn := replace.LevelLowerCase
	if suite.HasWarning(warning.LevelCase) {
		start = "info"
		fixed = "INFO"
		attrFn = replace.LevelUpperCase
	}

	logger := suite.Logger(infra.SimpleOptions())
	logger.Info(message)
	logMap := suite.logMap()
	level, ok := logMap[slog.LevelKey].(string)
	suite.Require().True(ok)
	suite.Assert().Equal(start, level)
	suite.bufferReset()
	logger = suite.Logger(infra.ReplaceAttrOptions(attrFn))
	logger.Info(message)
	logMap = suite.logMap()
	level, ok = logMap[slog.LevelKey].(string)
	suite.Require().True(ok)
	warnings := suite.HasWarnings(warning.NoReplAttrBasic, warning.NoReplAttr)
	if len(warnings) > 0 {
		issues := make([]string, 0, 3)
		if len(logMap) < 3 {
			issues = append(issues, fmt.Sprintf("too few attributes: %d", len(logMap)))
		}
		if !ok {
			issues = append(issues, "no level key")
		}
		if level != "" {
			issues = append(issues, "level value not null")
		}
		if len(issues) > 0 {
			suite.AddWarning(warnings[0], strings.Join(issues, ", "), "")
			return
		}
		suite.AddUnused(warnings[0], "")
	}
	suite.Assert().Equal(fixed, level)
}

// TestReplaceAttrFnRemoveEmptyKey tests the RemoveEmptyKey function.
func (suite *SlogTestSuite) TestReplaceAttrFnRemoveEmptyKey() {
	logger := suite.Logger(infra.SimpleOptions())
	logger.Info(message, "", "garbage")
	logMap := suite.logMap()
	value, ok := logMap[""]
	suite.Require().True(ok)
	suite.Require().Equal("garbage", value)
	suite.bufferReset()
	logger = suite.Logger(infra.ReplaceAttrOptions(replace.RemoveEmptyKey))
	logger.Info(message, "", nil)
	logMap = suite.logMap()
	value, ok = logMap[""]
	if suite.HasWarning(warning.NoReplAttr) {
		issues := make([]string, 0, 3)
		if len(logMap) < 4 {
			issues = append(issues, fmt.Sprintf("too few attributes: %d", len(logMap)))
		}
		if !ok {
			issues = append(issues, "no empty key")
		}
		if value != nil {
			issues = append(issues, "empty key value not null")
		}
		if len(issues) > 0 {
			suite.AddWarning(warning.NoReplAttr, strings.Join(issues, ", "), "")
			return
		}
		suite.AddUnused(warning.NoReplAttr, "")
	}
	if suite.HasWarning(warning.EmptyAttributes) {
		suite.Assert().Len(logMap, 4)
		suite.Assert().True(ok)
		suite.Assert().Nil(value)
	} else {
		suite.Assert().Len(logMap, 3)
		suite.Assert().False(ok)
	}
}

// TestReplaceAttrFnRemoveTime tests the RemoveEmptyKey function.
func (suite *SlogTestSuite) TestReplaceAttrFnRemoveTime() {
	logger := suite.Logger(infra.SimpleOptions())
	logger.Info(message)
	logMap := suite.logMap()
	suite.Require().Len(logMap, 3)
	value, ok := logMap[slog.TimeKey]
	suite.Require().True(ok)
	suite.NotNil(value)
	suite.bufferReset()
	logger = suite.Logger(infra.ReplaceAttrOptions(replace.RemoveTime))
	logger.Info(message)
	logMap = suite.logMap()
	value, ok = logMap[slog.TimeKey].(string)
	warnings := suite.HasWarnings(warning.NoReplAttrBasic, warning.NoReplAttr)
	if len(warnings) > 0 {
		issues := make([]string, 0, 3)
		if len(logMap) < 3 {
			issues = append(issues, fmt.Sprintf("too few attributes: %d", len(logMap)))
		}
		if !ok {
			issues = append(issues, "no time key")
		}
		if value != "" {
			issues = append(issues, "time value not empty string")
		}
		if len(issues) > 0 {
			suite.AddWarning(warnings[0], strings.Join(issues, ", "), "")
			return
		}
		suite.AddUnused(warnings[0], "")
	}
	if suite.HasWarning(warning.EmptyAttributes) {
		suite.Require().Len(logMap, 3)
		suite.Assert().True(ok)
		suite.Assert().Nil(value)
	} else {
		suite.Require().Len(logMap, 2)
		suite.Assert().False(ok)
	}
}
