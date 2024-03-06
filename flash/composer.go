package flash

import (
	"encoding"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/madkins23/go-slog/infra"
)

// -----------------------------------------------------------------------------

// composer handles writing attributes to an io.Writer.
type composer struct {
	buffer  []byte
	started bool
	replace infra.AttrFn
	groups  []string
}

func newComposer(buffer []byte, started bool, replace infra.AttrFn, groups []string) *composer {
	return &composer{
		buffer:  buffer,
		started: started,
		replace: replace,
		groups:  groups,
	}
}

func (c *composer) setStarted(started bool) {
	c.started = started
}

func (c *composer) getBytes() []byte {
	return c.buffer
}

// -----------------------------------------------------------------------------

func (c *composer) addAttribute(attr slog.Attr) error {
	attr.Value = attr.Value.Resolve()
	if c.replace != nil {
		attr = c.replace(c.groups, attr)
	}
	if attr.Equal(infra.EmptyAttr()) {
		return nil
	}
	value := attr.Value
	if value.Kind() == slog.KindGroup {
		if emptyGroup(value.Group()) {
			return nil
		}
		if attr.Key == "" {
			if err := c.addAttributes(value.Group()); err != nil {
				return fmt.Errorf("inline group attributes: %w", err)
			}
			return nil
		}
	}
	if !c.started {
		c.started = true
	} else {
		c.buffer = append(c.buffer, ',', ' ')
	}
	c.addString(attr.Key)
	c.buffer = append(c.buffer, ':', ' ')
	switch value.Kind() {
	case slog.KindGroup:
		return c.addGroup(value.Group())
	case slog.KindBool:
		c.addBool(value.Bool())
	case slog.KindDuration:
		c.addDuration(value.Duration())
	case slog.KindFloat64:
		c.addFloat64(value.Float64())
	case slog.KindInt64:
		c.addInt64(value.Int64())
	case slog.KindString:
		c.addString(value.String())
	case slog.KindTime:
		c.addTime(value.Time())
	case slog.KindUint64:
		c.addUint64(value.Uint64())
	case slog.KindAny:
		fallthrough
	default:
		return c.addAny(value.Any())
	}
	return nil
}

func (c *composer) addAttributes(attrs []slog.Attr) error {
	for _, attr := range attrs {
		if err := c.addAttribute(attr); err != nil {
			return fmt.Errorf("add attribute '%s': %w", attr.String(), err)
		}
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
	case fmt.Stringer:
		c.addString(v.String())
	case error:
		c.addString(v.Error())
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
		} else {
			c.buffer = append(c.buffer, b...)
			return nil
		}
	}
	return nil
}

var (
	boolFalse = []byte("false")
	boolTrue  = []byte("true")
)

func (c *composer) addBool(b bool) {
	c.buffer = append(c.buffer, boolImage[b]...)
}

func (c *composer) addBytes(b ...byte) {
	c.buffer = append(c.buffer, b...)
}

func (c *composer) addByteString(b []byte) {
	c.buffer = append(c.buffer, b...)
}

func (c *composer) addDuration(d time.Duration) {
	c.buffer = strconv.AppendInt(c.buffer, d.Nanoseconds(), 10)
}

func (c *composer) addFloat64(f float64) {
	c.buffer = strconv.AppendFloat(c.buffer, f, 'f', -1, 64)
}

func (c *composer) addGroup(attrs []slog.Attr) error {
	var err error
	c.buffer = append(c.buffer, '{')
	c.setStarted(false)
	if err = c.addAttributes(attrs); err != nil {
		return fmt.Errorf("add attributes: %w", err)
	}
	c.addBytes('}')
	return nil
}

func (c *composer) addInt64(i int64) {
	c.buffer = strconv.AppendInt(c.buffer, i, 10)
}

func (c *composer) addJSONMarshaler(m json.Marshaler) error {
	if txt, err := m.MarshalJSON(); err != nil {
		slog.Error("MarshalJSON error", "err", err)
		c.addString("!ERROR:" + err.Error())
		return fmt.Errorf("marshal JSON: %w", err)
	} else {
		c.addStringAsBytes(txt)
		return nil
	}
}

func (c *composer) addString(str string) {
	c.addStringAsBytes([]byte(str))
}

func (c *composer) addStringAsBytes(str []byte) {
	c.buffer = append(c.buffer, '"')
	c.buffer = append(c.buffer, str...)
	c.buffer = append(c.buffer, '"')
}

func (c *composer) addTextMarshaler(m encoding.TextMarshaler) error {
	if txt, err := m.MarshalText(); err != nil {
		slog.Error("MarshalText error", "err", err)
		c.addString("!ERROR:" + err.Error())
		return fmt.Errorf("marshal text: %w", err)
	} else {
		c.addStringAsBytes(txt)
		return nil
	}
}

func (c *composer) addTime(t time.Time) {
	c.addString(t.Format(time.RFC3339Nano))
}

func (c *composer) addUint64(i uint64) {
	c.buffer = strconv.AppendUint(c.buffer, i, 10)
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
