package tests

import (
	"fmt"
	"log/slog"
	"strings"
	"time"
)

// -----------------------------------------------------------------------------
// Methods for checks common to various tests or too long to put in a test.

// checkFieldCount checks whether the prescribed number of fields exist at the top level.
// In addition to using the logMap generated by unmarshaling the JSON log data,
// the custom test.FieldCounter is used to make sure there are no duplicates.
func (suite *SlogTestSuite) checkFieldCount(fieldCount uint, logMap map[string]any) {
	if suite.HasWarning(WarnDuplicates) {
		counter := suite.fieldCounter()
		suite.Require().NoError(counter.Parse())
		if len(counter.Duplicates()) > 0 {
			suite.addWarning(WarnDuplicates, fmt.Sprintf("%v", counter.Duplicates()), "")
			return
		}
	}
	suite.Assert().Len(logMap, int(fieldCount))
	// Double check to make sure there are no duplicate fields at the top level.
	counter := suite.fieldCounter()
	suite.Require().NoError(counter.Parse())
	suite.Assert().Equal(fieldCount, counter.NumFields())
	suite.Assert().Empty(counter.Duplicates())
}

func (suite *SlogTestSuite) checkLevelKey(level string, logMap map[string]any) {
	// The log/slog.JSONHandler generates uppercase.
	level = strings.ToUpper(level)
	if suite.HasWarning(WarnLevelCase) {
		if logLevel, ok := logMap[slog.LevelKey].(string); ok {
			if suite.Assert().Equal(level, strings.ToUpper(logLevel)) && level != logLevel {
				suite.addWarning(WarnLevelCase, "'"+logLevel+"'", "")
				return
			}
		}
		suite.addWarning(WarnUnused, WarnLevelCase, "")
	}
	suite.Assert().Equal(level, logMap[slog.LevelKey])
}

func (suite *SlogTestSuite) checkMessageKey(message string, logMap map[string]any) {
	if suite.HasWarning(WarnMessageKey) {
		if _, found := logMap[slog.MessageKey]; found {
			// Something exists for the proper key so fall through to test assertion.
		} else if msg, found := logMap["message"]; found {
			// Found something on the known alternate key.
			if message == msg {
				suite.addWarning(WarnMessageKey, "`message`", "")
				return
			}
		}
		suite.addWarning(WarnUnused, WarnMessageKey, "")
	}
	suite.Assert().Equal(message, logMap[slog.MessageKey])
}

func (suite *SlogTestSuite) checkNoEmptyAttribute(fieldCount uint, logMap map[string]any) {
	if suite.HasWarning(WarnEmptyAttributes) {
		// Warn for logging of empty attribute.
		counter := suite.fieldCounter()
		suite.Require().NoError(counter.Parse())
		if counter.NumFields() == fieldCount+1 {
			if _, found := logMap[""]; found {
				suite.addWarning(WarnEmptyAttributes, "", suite.Buffer.String())
				return
			}
		}
		suite.addWarning(WarnUnused, WarnEmptyAttributes, "")
	}
	suite.checkFieldCount(fieldCount, logMap)
	_, found := logMap[""]
	suite.Assert().False(found)
}

func (suite *SlogTestSuite) checkResolution(value any, actual any) {
	if suite.HasWarning(WarnResolver) {
		if value != actual {
			suite.addWarning(WarnResolver, "", suite.Buffer.String())
			return
		}
		suite.addWarning(WarnUnused, WarnResolver, "")
	}
	suite.Assert().Equal(value, actual)
}

var sourceKeys = map[string]any{
	"file":     "",
	"function": "",
	"line":     123.456,
}

func (suite *SlogTestSuite) checkSourceKey(fieldCount uint, logMap map[string]any) {
	if suite.HasWarning(WarnSourceKey) {
		sourceData := logMap[slog.SourceKey]
		if sourceData == nil {
			suite.addWarning(WarnSourceKey, "no 'source' key", suite.Buffer.String())
			return
		}
		source, ok := sourceData.(map[string]any)
		if !ok {
			suite.addWarning(WarnSourceKey, "'source' key not a group", suite.Buffer.String())
			return
		}
		var text strings.Builder
		sep := ""
		for field := range sourceKeys {
			var state string
			value := source[field]
			if value == nil {
				state = "missing"
			} else if _, ok := value.(string); !ok {
				state = "not a string"
			}
			if state != "" {
				text.WriteString(fmt.Sprintf("%s%s: %s", sep, field, state))
				sep = ", "
			}
		}
		if text.Len() > 0 {
			suite.addWarning(WarnSourceKey, text.String(), suite.Buffer.String())
		}
		suite.addWarning(WarnUnused, WarnSourceKey, "")
	}

	suite.checkFieldCount(fieldCount, logMap)
	if group, ok := logMap[slog.SourceKey].(map[string]any); ok {
		suite.Assert().Len(group, 3)
		for field, exemplar := range sourceKeys {
			suite.Assert().NotNil(group[field])
			suite.Assert().IsType(exemplar, group[field], "key: "+field)
		}
	} else {
		suite.Fail("Group not map[string]any")
	}
}

func (suite *SlogTestSuite) parseTime(timeAny any) time.Time {
	suite.Assert().NotNil(timeAny)
	timeStr, ok := timeAny.(string)
	suite.Assert().True(ok)
	timeObj, err := time.Parse(time.RFC3339, timeStr)
	suite.Assert().NoError(err)
	return timeObj
}
