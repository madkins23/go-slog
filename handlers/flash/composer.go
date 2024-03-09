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

var composerPool = newGenPool[composer]()

// -----------------------------------------------------------------------------

// composer handles writing attributes to an io.Writer.
type composer struct {
	buffer  []byte
	started bool
	replace infra.AttrFn
	groups  []string
}

func newComposer(buffer []byte, started bool, replace infra.AttrFn, groups []string) *composer {
	comp := composerPool.get()
	comp.buffer = buffer
	comp.started = started
	comp.replace = replace
	comp.groups = groups
	return comp
}

func reuseComposer(comp *composer) {
	composerPool.put(comp)
}

func (c *composer) reset() {
	c.started = false
}

func (c *composer) getBytes() []byte {
	return c.buffer
}

// -----------------------------------------------------------------------------

func (c *composer) addAttribute(attr slog.Attr) error {
	if attr.Value.Kind() == slog.KindLogValuer {
		attr.Value = attr.Value.Resolve()
	}
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
		if value.Bool() {
			c.buffer = append(c.buffer, "true"...)
		} else {
			c.buffer = append(c.buffer, "false"...)
		}
	case slog.KindDuration:
		c.buffer = strconv.AppendInt(c.buffer, value.Duration().Nanoseconds(), 10)
	case slog.KindFloat64:
		c.buffer = strconv.AppendFloat(c.buffer, value.Float64(), 'f', -1, 64)
	case slog.KindInt64:
		c.buffer = strconv.AppendInt(c.buffer, value.Int64(), 10)
	case slog.KindString:
		c.addString(value.String())
	case slog.KindTime:
		c.buffer = append(c.buffer, '"')
		c.buffer = value.Time().AppendFormat(c.buffer, time.RFC3339Nano)
		c.buffer = append(c.buffer, '"')
	case slog.KindUint64:
		c.buffer = strconv.AppendUint(c.buffer, value.Uint64(), 10)
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

func (c *composer) addBytes(b ...byte) {
	c.buffer = append(c.buffer, b...)
}

func (c *composer) addByteArray(b []byte) {
	c.buffer = append(c.buffer, b...)
}

func (c *composer) addGroup(attrs []slog.Attr) error {
	var err error
	c.buffer = append(c.buffer, '{')
	c.reset() // Reset composer (started = false) to avoid comma.
	if err = c.addAttributes(attrs); err != nil {
		return fmt.Errorf("add attributes: %w", err)
	}
	c.addBytes('}')
	return nil
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
	c.buffer = append(c.buffer, '"')
	c.buffer = append(c.buffer, str...)
	c.buffer = append(c.buffer, '"')
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
