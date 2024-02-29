package tests

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"math"
	"reflect"
	"strings"
	"time"

	"github.com/madkins23/go-slog/infra"
	intTest "github.com/madkins23/go-slog/internal/test"
	"github.com/madkins23/go-slog/internal/warning"
	"github.com/madkins23/go-slog/trace"
)

var tests = []string{
	"     M",
	"     M+A",
	"     M+B",
	"W+A  M",
	"W+A  M+B",
	"G1   M",
	"G1   M+A",
	"G1   M+B",
	"G1+A M",
	"G1+A M+B",
	"G1+C M+B",
	// TODO: chanchal/zerolog fails these
	//       Presumably the +C at the end of the last one keeps G2 open in actual.
	//"G1 G2 M", // slog/json &C
	"G1 G2 M+C",
	//"G1+A G2   M", // chanchal/zerolog
	"G1+A G2   M+C", // chanchal/zerolog
	//"G1 G2+B M",     // chanchal/zerolog
	"G1 G2+B M+C",
	//"G1+A G2+B M", // chanchal/zerolog
	"G1+A G2+B M+C",
	//"G1 G2 G3 M",   // slog/json &c
	"G1 G2 G3 M+C",
	//"G1+A G2 G3 M", // slog/json &c
	"G1+A G2 G3 M+C",
	//"G1 G2+B G3 M",   // chanchal/zerolog
	"G1 G2+B G3 M+C",
	//"G1 G2 G3+C M",   // chanchal/zerolog
	"G1 G2 G3+C M+B",
	//"G1+A G2+B G3 M", // chanchal/zerolog
	"G1+A G2+B G3 M+C",
	//"G1+A G2 G3+C M", // chanchal/zerolog
	"G1+A G2 G3+C M+B",
	//"G1 G2+B G3+C M", // chanchal/zerolog
	"G1 G2+B G3+C M+A",
	//"G1+A G2+B G3+C M", // chanchal/zerolog
	//"G1+A G2+B G3+C M+D", // ???
}

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
		if !suite.HasWarning(warning.Mismatch) {
			suite.Assert().Equal(expected, actual, test)
		}
	}
	if len(mismatches) > 0 {
		failed := strings.Join(mismatches, " | ")
		intTest.Debugf(1, ">>> Mismatches: %s", failed)
		suite.AddWarning(warning.Mismatch, failed, suite.Buffer.String())
	}
}

// -----------------------------------------------------------------------------

type parser struct {
	*warning.Manager
	inGroup          bool
	name, definition string
	logger           *slog.Logger
	logMap, ptrMap   map[string]any
}

func newParser(manager *warning.Manager, logger *slog.Logger, definition string) *parser {
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
	if !p.HasWarning(warning.GroupEmpty) {
		p.removeEmptyGroups(p.logMap)
	}
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
		attrs, err = p.getAttrs()
		if err != nil {
			return fmt.Errorf("get attributes: %w", err)
		}
		if len(attrs) > 0 {
			p.pushLog(p.currLog().With(anyList(attrs)...))
			if p.HasWarning(warning.GroupWithTop) {
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
		if err = p.addAttrs(attrs...); err != nil {
			return fmt.Errorf("add attributes: %w", err)
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
	if p.HasWarning(warning.LevelCase) {
		if lvl, found := actual[slog.LevelKey]; found {
			if level, ok := lvl.(string); ok {
				actual[slog.LevelKey] = strings.ToUpper(level)
			}
		}
	}
}

func (p *parser) removeEmptyGroups(logMap map[string]any) {
	for key, val := range logMap {
		if group, ok := val.(map[string]any); ok {
			if len(group) < 1 {
				delete(logMap, key)
			} else {
				p.removeEmptyGroups(group)
			}
		} else if _, ok := val.(slog.LogValuer); ok && p.HasWarning(warning.Resolver) {
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
		if p.HasWarning(warning.GroupDuration) && p.inGroup {
			value = float64(attr.Value.Duration().Nanoseconds())
		} else if p.Manager.HasWarning(warning.DurationSeconds) {
			value = attr.Value.Duration().Seconds()
		} else if p.Manager.HasWarning(warning.DurationMillis) {
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
		if p.HasWarning(warning.Resolver) {
			value = attr.Value.LogValuer()
		} else {
			value = attr.Value.LogValuer().LogValue().String()
		}
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
