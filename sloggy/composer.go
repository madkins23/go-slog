package sloggy

import (
	"encoding"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net"
	"strconv"
	"time"

	"github.com/madkins23/go-slog/infra"
)

var (
	boolFalse   = []byte("false")
	boolTrue    = []byte("true")
	braceLeft   = []byte{'{'}
	braceRight  = []byte{'}'}
	colonSpace  = []byte{':', ' '}
	commaSpace  = []byte{',', ' '}
	doubleQuote = []byte{'"'}
	emptyString []byte
)

// -----------------------------------------------------------------------------

// composer handles writing attributes to an io.Writer.
type composer struct {
	io.Writer
	started bool
	replace infra.AttrFn
	groups  []string
}

func newComposer(writer io.Writer, started bool, replace infra.AttrFn, groups []string) *composer {
	return &composer{
		Writer:  writer,
		started: started,
		replace: replace,
		groups:  groups,
	}
}

func (c *composer) setStarted(started bool) {
	c.started = started
}

// -----------------------------------------------------------------------------

func (c *composer) begin() error {
	if _, err := c.Write(braceLeft); err != nil {
		return fmt.Errorf("begin brace: %w", err)
	}
	return nil
}

func (c *composer) addAttribute(attr slog.Attr, groups []string) error {
	if attr.Equal(infra.EmptyAttr()) {
		return nil
	}
	value := attr.Value.Resolve()
	if c.replace != nil {
		attr = c.replace(c.groups, attr)
		value = attr.Value
	}
	if attr.Equal(infra.EmptyAttr()) {
		return nil
	}
	if value.Kind() == slog.KindGroup {
		if emptyGroup(value.Group()) {
			return nil
		}
		if attr.Key == "" {
			if err := c.addAttributes(value.Group(), c.groups); err != nil {
				return fmt.Errorf("inline group attributes: %w", err)
			}
			return nil
		}
	}
	if !c.started {
		c.started = true
	} else if _, err := c.Write(commaSpace); err != nil {
		return fmt.Errorf("comma space: %w", err)
	}
	if err := c.addString(attr.Key); err != nil {
		return fmt.Errorf("key field name: %w", err)
	}
	if _, err := c.Write(colonSpace); err != nil {
		return fmt.Errorf("field separator: %w", err)
	}
	switch value.Kind() {
	case slog.KindGroup:
		return c.addGroup(attr.Key, value.Group())
	case slog.KindBool:
		return c.addBool(value.Bool())
	case slog.KindDuration:
		return c.addDuration(value.Duration())
	case slog.KindFloat64:
		return c.addFloat64(value.Float64())
	case slog.KindInt64:
		return c.addInt64(value.Int64())
	case slog.KindString:
		return c.addString(value.String())
	case slog.KindTime:
		return c.addTime(value.Time())
	case slog.KindUint64:
		return c.addUint64(value.Uint64())
	case slog.KindAny:
		fallthrough
	default:
		return c.addAny(value.Any())
	}
}

func (c *composer) addAttributes(attrs []slog.Attr, groups []string) error {
	for _, attr := range attrs {
		if err := c.addAttribute(attr, c.groups); err != nil {
			return fmt.Errorf("add attribute '%s': %w", attr.String(), err)
		}
	}
	return nil
}

func (c *composer) end() error {
	if _, err := c.Write(braceRight); err != nil {
		return fmt.Errorf("end brace: %w", err)
	}
	return nil
}

// -----------------------------------------------------------------------------

var boolImage = map[bool][]byte{
	false: boolFalse,
	true:  boolTrue,
}

func (c *composer) addAny(a any) error {
	switch v := a.(type) {
	case net.IP:
		return c.addIPAddress(v)
	case net.IPNet:
		return c.addIPNet(v)
	case net.HardwareAddr:
		return c.addMacAddress(v)
	case error:
		return c.addError(v)
	case fmt.Stringer:
		return c.addStringer(v)
	case json.Marshaler:
		return c.addJSONMarshaler(v)
	case encoding.TextMarshaler:
		return c.addTextMarshaler(v)
	default:
		// Everything else (e.g. an array of strings)
		// Note: Important stuff buried in some random structure may be ignored.
		//       For example, a LogValuer or a Stringer might show up as an empty map.
		if b, err := json.Marshal(a); err != nil {
			return fmt.Errorf("marshal %v: %w", a, err)
		} else if _, err := c.Write(b); err != nil {
			return fmt.Errorf("write bytes: %w", err)
		} else {
			return nil
		}
	}
}

func (c *composer) addBool(b bool) error {
	if _, err := c.Write(boolImage[b]); err != nil {
		return fmt.Errorf("boolean: %w", err)
	}
	return nil
}

func (c *composer) addDuration(d time.Duration) error {
	if _, err := c.Write(strconv.AppendInt([]byte{}, d.Nanoseconds(), 10)); err != nil {
		return fmt.Errorf("uint64: %w", err)
	}
	return nil
}

func (c *composer) addError(e error) error {
	if _, err := c.Write([]byte(e.Error())); err != nil {
		return fmt.Errorf("error '%v': %w", e, err)
	}
	return nil
}

func (c *composer) addFloat64(f float64) error {
	if _, err := c.Write(strconv.AppendFloat([]byte{}, f, 'f', -1, 64)); err != nil {
		return fmt.Errorf("float64: %w", err)
	}
	return nil
}

func (c *composer) addGroup(name string, attrs []slog.Attr) error {
	if err := c.begin(); err != nil {
		return fmt.Errorf("begin: %w", err)
	}
	// Local composer object resets started flag
	cg := newComposer(c.Writer, false, c.replace, c.groups)
	if err := cg.addAttributes(attrs, append(c.groups, name)); err != nil {
		return fmt.Errorf("add attributes: %w", err)
	}
	if err := c.end(); err != nil {
		return fmt.Errorf("end: %w", err)
	}
	return nil
}

func (c *composer) addInt64(i int64) error {
	if _, err := c.Write(strconv.AppendInt([]byte{}, i, 10)); err != nil {
		return fmt.Errorf("int64: %w", err)
	}
	return nil
}

func (c *composer) addIPAddress(ip net.IP) error {
	if err := c.addString(ip.String()); err != nil {
		return fmt.Errorf("IP address '%v': %w", ip, err)
	}
	return nil
}

func (c *composer) addIPNet(ip net.IPNet) error {
	if err := c.addString(ip.String()); err != nil {
		return fmt.Errorf("IP net '%v': %w", ip, err)
	}
	return nil
}

func (c *composer) addJSONMarshaler(m json.Marshaler) error {
	if txt, err := m.MarshalJSON(); err != nil {
		return c.addString("!ERROR:" + err.Error())
	} else {
		return c.addStringAsBytes(txt)
	}
}

func (c *composer) addMacAddress(mac net.HardwareAddr) error {
	if err := c.addString(mac.String()); err != nil {
		return fmt.Errorf("hardware (MAC) address '%v': %w", mac, err)
	}
	return nil
}

func (c *composer) addString(str string) error {
	if _, err := c.Write(doubleQuote); err != nil {
		return fmt.Errorf("left double quote: %w", err)
	}
	if _, err := c.Write([]byte(str)); err != nil {
		return fmt.Errorf("string '%s': %w", str, err)
	}
	if _, err := c.Write(doubleQuote); err != nil {
		return fmt.Errorf("right double quote: %w", err)
	}
	return nil
}

func (c *composer) addStringAsBytes(str []byte) error {
	if _, err := c.Write(doubleQuote); err != nil {
		return fmt.Errorf("left double quote: %w", err)
	}
	if _, err := c.Write(str); err != nil {
		return fmt.Errorf("string '%s': %w", str, err)
	}
	if _, err := c.Write(doubleQuote); err != nil {
		return fmt.Errorf("right double quote: %w", err)
	}
	return nil
}

func (c *composer) addStringer(s fmt.Stringer) error {
	if _, err := c.Write([]byte(s.String())); err != nil {
		return fmt.Errorf("stringer '%v': %w", s, err)
	}
	return nil
}

func (c *composer) addTextMarshaler(m encoding.TextMarshaler) error {
	if txt, err := m.MarshalText(); err != nil {
		return c.addString("!ERROR:" + err.Error())
	} else {
		return c.addStringAsBytes(txt)
	}
}

func (c *composer) addTime(t time.Time) error {
	if err := c.addString(t.Format(time.RFC3339Nano)); err != nil {
		return fmt.Errorf("time: %w", err)
	}
	return nil
}

func (c *composer) addUint64(i uint64) error {
	if _, err := c.Write(strconv.AppendUint([]byte{}, i, 10)); err != nil {
		return fmt.Errorf("uint64: %w", err)
	}
	return nil
}

// -----------------------------------------------------------------------------

func emptyGroup(attrs []slog.Attr) bool {
	for _, attr := range attrs {
		if attr.Equal(infra.EmptyAttr()) {
			continue
		}
		if attr.Value.Kind() == slog.KindGroup {
			if !emptyGroup(attr.Value.Group()) {
				return false
			}
		} else {
			// Attribute is not empty and not a group.
			return false
		}
	}
	return true
}
