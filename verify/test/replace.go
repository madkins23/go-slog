package test

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/madkins23/go-slog/replace"
)

// -----------------------------------------------------------------------------
// Tests of slog.HandlerOptions.ReplaceAttr functionality.

// TestSimpleReplaceAttr tests the use of HandlerOptions.ReplaceAttr.
func (suite *SlogTestSuite) TestSimpleReplaceAttr() {
	logger := suite.Logger(&slog.HandlerOptions{
		Level: slog.LevelInfo,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			switch a.Key {
			case "alpha":
				return slog.String(a.Key, "omega")
			case "change":
				return slog.String("bravo", a.Value.String())
			case "remove":
				return replace.EmptyAttr
			}
			return a
		},
	})
	logger.Info(message, "alpha", "beta", "change", "my key", "remove", "me")
	logMap := suite.logMap()
	if suite.hasWarning(WarnNoReplAttr) {
		issues := make([]string, 4)
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
			suite.addWarning(WarnNoReplAttr, strings.Join(issues, ", "), false)
			return
		}
		suite.addWarning(WarnUnused, WarnNoReplAttr, true)
	}
	if suite.hasWarning(WarnEmptyAttributes) {
		suite.checkFieldCount(6, logMap)
	} else {
		suite.checkFieldCount(5, logMap)
	}
	suite.Assert().Equal("omega", logMap["alpha"])
	suite.Assert().Equal("my key", logMap["bravo"])
	suite.Assert().Nil(logMap["remove"])
}

// TestSourceReplaceAttrBasic tests the use of HandlerOptions.ReplaceAttr
// on basic attributes (time, level, message, source).
func (suite *SlogTestSuite) TestSourceReplaceAttrBasic() {
	logger := suite.Logger(&slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			switch a.Key {
			case slog.TimeKey:
				return replace.EmptyAttr
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
	if suite.hasWarning(WarnNoReplAttr) || suite.hasWarning(WarnNoReplAttrBasic) {
		issues := make([]string, 0, 5)
		if len(logMap) > 3 {
			issues = append(issues, fmt.Sprintf("too many attributes: %d", len(logMap)))
		}
		if logMap[slog.TimeKey] != nil {
			issues = append(issues, slog.TimeKey+" field still exists")
		}
		if logMap[slog.MessageKey] != nil {
			issues = append(issues, slog.MessageKey+" field still exists")
		} else if suite.hasWarning(WarnMessageKey) && logMap["message"] != nil {
			issues = append(issues, "message field still exists")
		}
		// TODO: This one may still work, in samber it's apparently a separate field from basic.
		if value, ok := logMap[slog.SourceKey].(string); !ok || value != "all knowledge" {
			issues = append(issues, fmt.Sprintf("%s == %v", slog.SourceKey, logMap[slog.SourceKey]))
		}
		if len(issues) > 0 {
			suite.addWarning(WarnNoReplAttrBasic, strings.Join(issues, ", "), false)
			return
		}
		suite.addWarning(WarnUnused, WarnNoReplAttrBasic, true)
	}
	suite.checkFieldCount(3, logMap)
	suite.Assert().Nil(logMap[slog.TimeKey])
	suite.Assert().Equal("Tilted", logMap[slog.LevelKey])
	suite.Assert().Equal(message, logMap["Message"])
	suite.Assert().Equal("all knowledge", logMap[slog.SourceKey])
}
