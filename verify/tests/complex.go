package tests

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"math"
	"reflect"
	"strings"
	"time"

	"github.com/madkins23/go-slog/handlers/trace"
	"github.com/madkins23/go-slog/infra"
	warning2 "github.com/madkins23/go-slog/infra/warning"
	intTest "github.com/madkins23/go-slog/internal/test"
)

// Test algorithms represented as strings.
var tests = []string{
	"               M",
	"               M+A",
	"               M+B",
	"W+A            M",
	"W+A            M+B",
	"G1             M",
	"G1             M+A", // phuslu/slog works properly here
	"G1             M+B", // phuslu/slog works properly here
	"G1+A           M",
	"G1+A           M+B", // phuslu/slog seems to lose track of group
	"G1+C           M+B", // phuslu/slog seems to lose track of group
	"G1   G2        M",
	"G1   G2        M+C",
	"G1+A G2        M",
	"G1+A G2        M+C",
	"G1   G2+B      M",
	"G1   G2+B      M+C", // phuslu/slog seems to lose track of group
	"G1+A G2+B      M",
	"G1+A G2+B      M+C", // phuslu/slog seems to lose track of group
	"G1   G2   G3   M",
	"G1   G2   G3   M+C",
	"G1+A G2   G3   M",
	"G1+A G2   G3   M+C",
	"G1   G2+B G3   M",
	"G1   G2+B G3   M+C",
	"G1   G2   G3+C M",
	"G1   G2   G3+C M+B", // phuslu/slog seems to lose track of group
	"G1+A G2+B G3   M",
	"G1+A G2+B G3   M+C",
	"G1+A G2   G3+C M",
	"G1+A G2   G3+C M+B", // phuslu/slog seems to lose track of group
	"G1   G2+B G3+C M",
	"G1   G2+B G3+C M+A", // phuslu/slog seems to lose track of group
	"G1+A G2+B G3+C M",
	"G1+A G2+B G3+C M+D", // phuslu/slog seems to lose track of group
}

// TestComplexCases executes a series of algorithmically generated test cases.
// Each test case is represented by a sequence of characters that is a 'program'.
// For each test case two things are done:
//   - a log statement is generated and
//   - an expected data structure of nested map[string]any is generated.
//
// The expectation is that the logMap from the log statement matches the expected data.
// Various special code bits are used to handle discrepancies marked by warnings.
//
// This may seem a bit like cheating as the log statement and expected result are
// generated from the same algorithm for each test case.
// However, once the algorithmic strings/programs were working it was really simple
// to generate a bunch of tests by pattern of characters.
// This in turn surfaced a bunch of new discrepancies between logger implementations,
// resulting in several new warnings.
func (suite *SlogTestSuite) TestComplexCases() {
	logger := suite.Logger(infra.SimpleOptions())
	mismatches := make([]string, 0)
	for _, test := range tests {
		intTest.Debugf(1, "Complex: %s\n", test)
		if intTest.DebugLevel() > 4 {
			fmt.Println("  Trace:")
			logger := slog.New(trace.NewHandler("    "))
			parser := newParser(suite.Manager, logger, test)
			suite.Assert().NoError(parser.parse())
		}
		suite.Buffer.Reset()
		parser := newParser(suite.Manager, logger, test)
		suite.Assert().NoError(parser.parse())
		expected := parser.expected()
		actual := suite.logMap()
		parser.fixActual(actual)
		if intTest.DebugLevel() > 1 {
			suite.Assert().NoError(show(expected, actual))
		}
		if !reflect.DeepEqual(expected, actual) {
			mismatches = append(mismatches, test)
		}
		if !suite.HasWarning(warning2.Mismatch) {
			suite.Assert().Equal(expected, actual, test)
		}
	}
	if len(mismatches) > 0 {
		failed := strings.Join(mismatches, " | ")
		intTest.Debugf(1, ">>> %d Mismatches: %s", len(mismatches), failed)
		suite.AddWarning(warning2.Mismatch, failed, suite.Buffer.String())
	}
}

// -----------------------------------------------------------------------------

type parser struct {
	*warning2.Manager
	inGroup, inWith  bool
	name, definition string
	logger           *slog.Logger
	logMap, ptrMap   map[string]any
}

func newParser(manager *warning2.Manager, logger *slog.Logger, definition string) *parser {
	p := &parser{
		Manager:    manager,
		name:       definition,
		definition: definition,
		logger:     logger,
		logMap:     make(map[string]any),
	}
	p.ptrMap = p.logMap
	return p
}

func (p *parser) currLog() *slog.Logger {
	return p.logger
}

func (p *parser) currMap() map[string]any {
	return p.ptrMap
}

// expected returns the original logMap as the expected result.
func (p *parser) expected() map[string]any {
	p.removeEmptyGroups(p.logMap, 1)
	return p.logMap
}

func (p *parser) pushLog(logger *slog.Logger) {
	p.logger = logger
}

func (p *parser) pushMap(logMap map[string]any) {
	p.ptrMap = logMap
}

func (p *parser) parse() error {
	for len(p.definition) > 0 {
		definition := p.definition
		if err := p.execute(); err != nil {
			return fmt.Errorf("execute \"%s\": %w", definition, err)
		}
	}
	return nil
}

func (p *parser) execute() error {
	instruction := p.definition[0]
	p.definition = p.definition[1:]
	var attrs []slog.Attr
	var err error
	switch instruction {
	case ' ':
		return nil
	case 'G':
		grpName := "group"
		if len(p.definition) > 0 {
			grpName += string(p.definition[0])
			p.definition = p.definition[1:]
		}
		newMap := make(map[string]any)
		p.currMap()[grpName] = newMap
		p.pushLog(p.currLog().WithGroup(grpName))
		p.pushMap(newMap)
		p.inGroup = true
		p.inWith = false
		attrs, err = p.getAttrs()
		if err != nil {
			return fmt.Errorf("get attributes: %w", err)
		}
		if len(attrs) > 0 {
			p.pushLog(p.currLog().With(anyList(attrs)...))
			p.inWith = true
			if p.HasWarning(warning2.GroupWithTop) {
				err = p.addAttrsToMap(p.logMap, attrs...)
			} else {
				err = p.addAttrs(attrs...)
			}
			if err != nil {
				return fmt.Errorf("add attributes: %w", err)
			}
		}
	case 'M':
		attrs, err = p.getAttrs()
		if err != nil {
			return fmt.Errorf("get attributes: %w", err)
		}
		p.currLog().Info(message, anyList(attrs)...)
		p.logMap[slog.LevelKey] = "INFO"
		p.logMap[slog.MessageKey] = message
		if p.HasWarning(warning2.GroupAttrMsgTop) && (!p.inGroup || p.inWith) {
			//	"G1             M+A", fails here
			//	"G1+A           M+B", succeeds here
			if err = p.addAttrsToMap(p.logMap, attrs...); err != nil {
				return fmt.Errorf("add attributes: %w", err)
			}
		} else {
			//	"G1             M+A", succeeds here
			//	"G1+A           M+B", fails here
			if err = p.addAttrs(attrs...); err != nil {
				return fmt.Errorf("add attributes: %w", err)
			}
		}
	case 'W':
		attrs, err = p.getAttrs()
		if err != nil {
			return fmt.Errorf("get attributes: %w", err)
		}
		p.pushLog(p.currLog().With(anyList(attrs)...))
		if err := p.addAttrs(attrs...); err != nil {
			return fmt.Errorf("add attributes: %w", err)
		}
	default:
		return fmt.Errorf("bad test case instruction '%c' : \"%s\"", instruction, p.definition)
	}

	return nil
}

func (p *parser) fixActual(actual map[string]any) {
	delete(actual, slog.TimeKey)
	if _, found := actual["message"]; found {
		// Handler uses incorrect message key.
		actual[slog.MessageKey] = actual["message"]
		delete(actual, "message")
	}
	if p.HasWarning(warning2.LevelCase) {
		if lvl, found := actual[slog.LevelKey]; found {
			if level, ok := lvl.(string); ok {
				actual[slog.LevelKey] = strings.ToUpper(level)
			}
		}
	}
	p.removeEmptyGroups(actual, 0)
}

func (p *parser) removeEmptyGroups(logMap map[string]any, depth int) {
	for key, val := range logMap {
		if group, ok := val.(map[string]any); ok {
			if len(group) > 0 {
				p.removeEmptyGroups(group, depth+1)
			}
			if depth > 0 || !p.HasWarning(warning2.GroupEmpty) {
				if len(group) < 1 {
					delete(logMap, key)
				}
			} else if p.HasWarning(warning2.GroupEmpty) {
				if len(group) < 1 {
					delete(logMap, key)
				}
			}
		} else if _, ok := val.(slog.LogValuer); ok && p.HasWarning(warning2.Resolver) {
			logMap[key] = make(map[string]any)
		}
	}
}

// -----------------------------------------------------------------------------

var attributes = map[byte][]slog.Attr{
	'A': {
		slog.String("string", "value"),
		slog.Int("int", -13),
		slog.Uint64("uint", 23),
	},
	'B': {
		slog.String("aTime", time.Now().Format(time.RFC3339)),
		slog.Bool("bool", true),
		slog.Duration("duration", time.Hour+3*time.Minute+22*time.Second),
	},
	'C': {
		slog.Group("groupC",
			"name", "Goober Snoofus",
			"skidoo", 23,
			"pi", math.Pi),
		slog.Bool("bool", false),
		slog.Any("valuer", &hiddenValuer{"Big Tree"}),
		slog.Duration("duration", 23*time.Minute+49*time.Second),
	},
	'D': {
		slog.Float64("E", math.E),
		slog.Uint64("uint64", 79),
	},
}

func (p *parser) addAttrs(attrs ...slog.Attr) error {
	return p.addAttrsToMap(p.currMap(), attrs...)
}

func (p *parser) addAttrToMap(logMap map[string]any, attr slog.Attr) error {
	var value any
	switch attr.Value.Kind() {
	case slog.KindAny:
		value = attr.Value.Any()
	case slog.KindBool:
		value = attr.Value.Bool()
	case slog.KindDuration:
		if p.HasWarning(warning2.GroupDuration) && p.inGroup {
			value = float64(attr.Value.Duration().Nanoseconds())
		} else if p.Manager.HasWarning(warning2.DurationSeconds) {
			value = attr.Value.Duration().Seconds()
		} else if p.Manager.HasWarning(warning2.DurationMillis) {
			value = float64(attr.Value.Duration().Milliseconds())
		} else {
			value = float64(attr.Value.Duration().Nanoseconds())
		}
	case slog.KindFloat64:
		value = attr.Value.Float64()
	case slog.KindInt64:
		// JSON converts all numbers to float64.
		value = float64(attr.Value.Int64())
	case slog.KindString:
		value = attr.Value.String()
	case slog.KindTime:
		value = attr.Value.Time()
	case slog.KindUint64:
		// JSON converts all numbers to float64.
		value = float64(attr.Value.Uint64())
	case slog.KindGroup:
		subMap := make(map[string]any)
		if err := p.addAttrsToMap(subMap, attr.Value.Group()...); err != nil {
			return fmt.Errorf("add attributes: %w", err)
		}
		value = subMap
	case slog.KindLogValuer:
		if p.HasWarning(warning2.Resolver) {
			return nil
		}
		value = attr.Value.LogValuer().LogValue().String()
	default:
		slog.Warn("Unknown slog.Attr.Value.Kind", "kind", attr.Value.Kind().String())
	}
	logMap[attr.Key] = value
	return nil
}

func (p *parser) addAttrsToMap(logMap map[string]any, attrs ...slog.Attr) error {
	for _, attr := range attrs {
		if err := p.addAttrToMap(logMap, attr); err != nil {
			return fmt.Errorf("add attribute: %w", err)
		}
	}
	return nil
}

func (p *parser) getAttrs() ([]slog.Attr, error) {
	result := make([]slog.Attr, 0)
	if len(p.definition) > 1 && p.definition[0] == '+' {
		x := p.definition[1]
		p.definition = p.definition[2:]
		if attrs, found := attributes[x]; !found {
			return nil, fmt.Errorf("non-existent attribute list '%c'", x)
		} else {
			result = append(result, attrs...)
		}
	}
	return result, nil
}

// -----------------------------------------------------------------------------

func anyList(attributes []slog.Attr) []any {
	result := make([]any, len(attributes))
	for i, attr := range attributes {
		result[i] = attr
	}
	return result
}

func show(expected map[string]any, actual map[string]any) error {
	if err := showX("Expected", expected); err != nil {
		return fmt.Errorf("show Expected: %w", err)
	}
	if err := showX("Actual", actual); err != nil {
		return fmt.Errorf("show Actual: %w", err)
	}
	return nil
}

func showX(name string, logMap map[string]any) error {
	var b []byte
	var err error
	if intTest.DebugLevel() > 3 {
		b, err = json.MarshalIndent(logMap, "    ", "  ")
	} else {
		b, err = json.Marshal(logMap)
	}
	if err != nil {
		return fmt.Errorf("marshal %s: %w", strings.ToLower(name), err)
	} else {
		fmt.Printf("  %s:\n    %s\n", name, string(b))
		return nil
	}
}
