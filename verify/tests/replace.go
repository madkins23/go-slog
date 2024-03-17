package tests

import (
	"fmt"
	"log/slog"
	"math"
	"reflect"
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
	logger := suite.Logger(infra.ReplaceAttrOptions(func(_ []string, a slog.Attr) slog.Attr {
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
		suite.AddUnused(warning.NoReplAttr, suite.String())
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
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
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
			suite.AddWarning(warnings[0], strings.Join(issues, "\n"), "")
			return
		}
		suite.AddUnused(warnings[0], suite.String())
	}
	suite.checkFieldCount(3, logMap)
	suite.Assert().Nil(logMap[slog.TimeKey])
	suite.Assert().Equal("Tilted", logMap[slog.LevelKey])
	suite.Assert().Equal(message, logMap["Message"])
	suite.Assert().Equal("all knowledge", logMap[slog.SourceKey])
}

// TestReplaceAttrGroup tests the groups argument passed to a HandlerOptions.ReplaceAttr function.
// This checks to see if group names are properly tracked and passed.
//   - https://pkg.go.dev/log/slog@master#HandlerOptions
func (suite *SlogTestSuite) TestReplaceAttrGroup() {
	if suite.HasWarning(warning.NoReplAttr) {
		// Nothing to see here, move along.
		return
	}
	logger := suite.Logger(&slog.HandlerOptions{
		Level: slog.LevelInfo,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			var current string
			if len(groups) > 0 {
				current = groups[len(groups)-1]
			}
			switch current {
			case "":
				switch a.Key {
				case "alpha":
					return slog.Float64("alpha", math.Pi)
				case "bravo":
					return infra.EmptyAttr()
				case slog.TimeKey, slog.LevelKey, slog.MessageKey:
					// Do nothing here.
				default:
					suite.Assert().Fail("unexpected key", "'%s' group: '%s'", a.Key, current)
				}
			case "group":
				switch a.Key {
				case "charlie":
					return slog.Float64("charlie", math.E)
				case "delta":
					return infra.EmptyAttr()
				default:
					suite.Assert().Fail("unexpected group key", "'%s' group: '%s'", a.Key, current)
				}
			case "subGroup":
				switch a.Key {
				case "echo":
					return slog.Float64("echo", math.SqrtPi)
				case "foxtrot":
					return slog.Float64("foxtrot", math.SqrtE)
				case "golf", "hotel":
					return infra.EmptyAttr()
				default:
					suite.Assert().Failf("unexpected subGroup key", "'%s' group: '%s'", a.Key, current)
				}
			}
			return a
		},
	})
	logger.With("alpha", 1, "bravo", 2).
		WithGroup("group").With("charlie", 3, "delta", 4).
		WithGroup("subGroup").With("echo", 5, "foxtrot", 6).
		Info(message, "golf", 7, "hotel", 8)
	logMap := suite.logMap()
	delete(logMap, slog.TimeKey)
	expected := map[string]any{
		slog.LevelKey:   slog.LevelInfo.String(),
		slog.MessageKey: message,
		"alpha":         math.Pi,
		"group": map[string]any{
			"charlie": math.E,
			"subGroup": map[string]any{
				"echo":    math.SqrtPi,
				"foxtrot": math.SqrtE,
			},
		},
	}
	if suite.HasWarning(warning.EmptyAttributes) {
		stripEmptyAttr(logMap)
	}
	if suite.HasWarning(warning.LevelCase) {
		if lvl, ok := logMap[slog.LevelKey].(string); ok {
			logMap[slog.LevelKey] = strings.ToUpper(lvl)
		}
	}
	if suite.HasWarning(warning.MessageKey) {
		if msg, ok := logMap["message"].(string); ok {
			logMap[slog.MessageKey] = msg
			delete(logMap, "message")
		}
	}
	if !suite.HasWarning(warning.ReplAttrGroup) {
		suite.Assert().Equal(expected, logMap)
	} else if reflect.DeepEqual(expected, logMap) {
		suite.AddUnused(warning.ReplAttrGroup, suite.String())
	} else {
		suite.AddWarning(warning.ReplAttrGroup, "", suite.String())
	}
}

func stripEmptyAttr(logMap map[string]any) {
	for key, value := range logMap {
		if key == "" && value == nil {
			delete(logMap, key)
			continue
		}
		if group, ok := value.(map[string]any); ok {
			stripEmptyAttr(group)
		}
	}
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
			suite.AddWarning(warnings[0], strings.Join(issues, "\n"), "")
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
			suite.AddWarning(warning.NoReplAttr, strings.Join(issues, "\n"), "")
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

// TestReplaceAttrFnChangeKey tests the RemoveEmptyKey function.
func (suite *SlogTestSuite) TestReplaceAttrFnChangeKey() {
	logger := suite.Logger(infra.SimpleOptions())
	logger.Info(message)
	logMap := suite.logMap()
	if _, found := logMap["message"]; !found {
		// Can't run test on this handler.
		// Note: If all handlers start using slog.MessageKey
		// the rest of this test will fail to test anything.
		return
	}
	suite.bufferReset()
	logger = suite.Logger(infra.ReplaceAttrOptions(
		replace.ChangeKey("message", slog.MessageKey, false, replace.TopCheck)))
	logger.Info(message)
	logMap = suite.logMap()

	warnings := suite.HasWarnings(warning.NoReplAttrBasic, warning.NoReplAttr)
	if len(warnings) > 0 {
		issues := make([]string, 0, 3)
		value, found := logMap[slog.MessageKey]
		if len(logMap) < 4 {
			issues = append(issues, fmt.Sprintf("too few attributes: %d", len(logMap)))
		}
		if !found {
			issues = append(issues, "no message key")
		}
		if str, ok := value.(string); !ok {
			issues = append(issues, "message not string")
		} else if message != str {
			issues = append(issues, "wrong message: '"+str+"'")
		}
		if len(issues) > 0 {
			suite.AddWarning(warnings[0], strings.Join(issues, "\n"), "")
			return
		}
		suite.AddUnused(warnings[0], "")
	} else {
		_, found := logMap["message"]
		suite.Assert().False(found)
		msg, found := logMap[slog.MessageKey]
		suite.Assert().True(found)
		str, ok := msg.(string)
		suite.Assert().True(ok)
		suite.Assert().Equal(message, str)
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
	logger = suite.Logger(infra.ReplaceAttrOptions(
		replace.RemoveKey(slog.TimeKey, false, replace.TopCheck)))
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
			suite.AddWarning(warnings[0], strings.Join(issues, "\n"), "")
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
