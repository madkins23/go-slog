package json

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"regexp/syntax"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/gertd/go-pluralize"
)

// FieldCounter manually parses JSON to count fields.
// The main purpose of this is to find duplicate fields.
// The encoding/source package just chooses one of the fields when this happens,
// so it isn't useful in counting duplicate fields.
type FieldCounter struct {
	counts  map[string]uint
	charLoc uint
	recent  []rune
	reader  *bytes.Reader
}

// NewFieldCounter returns a new FieldCounter created with the specified source JSON.
func NewFieldCounter(srcJSON []byte) *FieldCounter {
	return &FieldCounter{
		counts: make(map[string]uint),
		reader: bytes.NewReader(srcJSON),
	}
}

// Duplicates returns a map from field names to
// the number of times each field appears at the top level of the JSON.
// Fields will only appear in the map if they appear multiple (>1) times.
func (ctr *FieldCounter) Duplicates() map[string]uint {
	dupMap := make(map[string]uint)
	for field, count := range ctr.counts {
		if count > 1 {
			dupMap[field] = count
		}
	}
	return dupMap
}

// Fields returns an array of field names that appear at least once in the JSON.
// The list of field names is sorted alphabetically.
func (ctr *FieldCounter) Fields() []string {
	fldList := make([]string, 0, len(ctr.counts))
	for field := range ctr.counts {
		fldList = append(fldList, field)
	}
	sort.Strings(fldList)
	return fldList
}

// NumFields returns the number of fields that appear in the JSON.
func (ctr *FieldCounter) NumFields() uint {
	return uint(len(ctr.counts))
}

var plural = pluralize.NewClient()

// String returns the status of the counter as a string.
// The string format is:
//
//	@<characterLocation> '<sourceText>' <numFields> fields[, <dupMap>]
//
// where:
//
//	characterLocation     is current character count of the parse cursor in the source JSON
//	sourceText            shows characters before and after the cursor:
//	                        the cursor is indicated as <^>
//	                        if required the end of file is indicated by <EOF>
//	numFields             is the number of fields found (i.e. NumFields)
//	dupMap                is a map of duplicate field counts (i.e. Duplicates)
func (ctr *FieldCounter) String() string {
	var builder strings.Builder
	var r rune
	builder.WriteRune('@')
	builder.WriteString(strconv.Itoa(int(ctr.charLoc)))
	builder.WriteString(" '")
	for _, r = range ctr.recent {
		builder.WriteRune(r)
	}
	builder.WriteString("<^>")
	for i := 0; i < 15; i++ {
		if r, err := ctr.readRune(); err == nil {
			if r == '\n' {
				r = ' '
			}
			builder.WriteRune(r)
		} else {
			if errors.Is(err, io.EOF) {
				builder.WriteString("<EOF>")
			}
			break
		}
	}
	builder.WriteString("' ")
	builder.WriteString(strconv.Itoa(int(ctr.NumFields())))
	builder.WriteByte(' ')
	builder.WriteString(plural.Pluralize("field", int(ctr.NumFields()), false))
	duplicates := ctr.Duplicates()
	if len(duplicates) > 0 {
		builder.WriteRune(',')
		for field, count := range duplicates {
			builder.WriteByte(' ')
			builder.WriteString(field)
			builder.WriteByte(':')
			builder.WriteString(strconv.Itoa(int(count)))
		}
	} else {
		builder.WriteString(", no duplicates")
	}
	return builder.String()
}

const (
	msgParseArray      = "parse array"
	msgParseField      = "parse field"
	msgParseKeyword    = "parse keyword"
	msgParseNumber     = "parse number"
	msgParseObject     = "parse object"
	msgParseString     = "parse string"
	msgParseValue      = "parse value"
	msgUnreadRune      = "unread rune"
	fmtExpectedColon   = "expected colon, got '%c'"
	msgUnexpectedColon = "unexpected colon"
	fmtExpectedComma   = "expected comma, got '%c'"
	msgUnexpectedComma = "unexpected comma"
)

// Parse the JSON source to count field names.
func (ctr *FieldCounter) Parse() error {
	return ctr.readLoop(func(r rune) error {
		if unicode.IsSpace(r) {
			return nil
		} else if r == '{' {
			err := ctr.wrapCallError(ctr.parseObject(true), msgParseObject, true, true)
			if errors.Is(err, io.EOF) {
				return nil
			} else {
				return err
			}
		} else {
			return errUnexpected(r)
		}
	})
}

func (ctr *FieldCounter) parseArray() error {
	expectComma := false
	return ctr.readLoop(func(r rune) error {
		if unicode.IsSpace(r) {
			return nil
		} else if r == ',' {
			if expectComma {
				expectComma = false
				return nil
			} else {
				return fmt.Errorf(msgUnexpectedComma)
			}
		} else if r == ']' {
			return errFinished
		} else if expectComma {
			return fmt.Errorf(fmtExpectedComma, r)
		} else {
			if err := ctr.unreadRune(); err != nil {
				return ctr.wrapCallError(err, msgUnreadRune, false, false)
			}
			expectComma = true
			return ctr.wrapCallError(ctr.parseValue(), msgParseValue, false, false)
		}
	})
}

func (ctr *FieldCounter) parseField(field string) error {
	expectColon := true
	return ctr.readLoop(func(r rune) error {
		if unicode.IsSpace(r) {
			return nil
		} else if r == ':' {
			if expectColon {
				expectColon = false
				return nil
			} else {
				return fmt.Errorf(msgUnexpectedColon)
			}
		} else if expectColon {
			return fmt.Errorf(fmtExpectedColon, r)
		} else if r == '{' {
			if field == "fields" {
				// Special case for Apex wherein additional fields are in 'fields' object.
				delete(ctr.counts, field)
				return ctr.wrapCallError(ctr.parseObject(true), msgParseObject, true, false)
			} else {
				return ctr.wrapCallError(ctr.parseObject(false), msgParseObject, true, false)
			}
		} else if r == '[' {
			return ctr.wrapCallError(ctr.parseArray(), msgParseArray, true, false)
		} else if err := ctr.unreadRune(); err != nil {
			return ctr.wrapCallError(err, msgUnreadRune, false, false)
		} else {
			return ctr.wrapCallError(ctr.parseValue(), msgParseValue, true, false)
		}
	})
}

func (ctr *FieldCounter) parseKeyword(first rune) (string, error) {
	var builder strings.Builder
	builder.WriteRune(first)
	return builder.String(), ctr.readLoop(func(r rune) error {
		if syntax.IsWordChar(r) {
			builder.WriteRune(r)
		} else if err := ctr.unreadRune(); err != nil {
			return ctr.wrapCallError(err, msgUnreadRune, false, false)
		} else {
			return errFinished
		}
		return nil
	})
}

func (ctr *FieldCounter) parseNumber() error {
	foundDecimal := false
	foundExponent := false
	foundExponentSign := false
	return ctr.readLoop(func(r rune) error {
		if unicode.IsDigit(r) {
			return nil
		} else if r == '.' {
			if foundExponent {
				return fmt.Errorf("decimal in exponent")
			} else if foundDecimal {
				return fmt.Errorf("second decimal")
			} else {
				foundDecimal = true
				return nil
			}
		} else if r == 'e' || r == 'E' {
			if foundExponent {
				return fmt.Errorf("second exponent")
			} else {
				foundExponent = true
				return nil
			}
		} else if r == '-' || r == '+' {
			if !foundExponent {
				return fmt.Errorf("sign before exponent")
			} else if foundExponentSign {
				return fmt.Errorf("second exponent sign")
			} else {
				foundExponentSign = true
				return nil
			}
		} else if err := ctr.unreadRune(); err != nil {
			return ctr.wrapCallError(err, msgUnreadRune, false, false)
		} else {
			return errFinished
		}
	})
}

func (ctr *FieldCounter) parseObject(countFields bool) error {
	expectComma := false
	return ctr.readLoop(func(r rune) error {
		switch {
		case unicode.IsSpace(r):
			return nil
		case r == ',':
			if expectComma {
				expectComma = false
				return nil
			} else {
				return fmt.Errorf(msgUnexpectedComma)
			}
		case r == '"':
			if expectComma {
				return fmt.Errorf(fmtExpectedComma, r)
			}
			field, err := ctr.parseString()
			if err != nil {
				return ctr.wrapCallError(err, msgParseString, false, false)
			}
			if countFields {
				ctr.counts[field]++
			}
			expectComma = true
			return ctr.wrapCallError(ctr.parseField(field), msgParseField, false, false)
		case r == '}':
			return errFinished
		default:
			return errUnexpected(r)
		}
	})
}

func (ctr *FieldCounter) parseString() (string, error) {
	var builder strings.Builder
	err := ctr.readLoop(func(r rune) error {
		var err error
		switch r {
		case '\\':
			builder.WriteRune(r)
			if r, err = ctr.readRune(); err != nil {
				return fmt.Errorf("read escaped rune: %w", err)
			} else {
				builder.WriteRune(r)
				return nil
			}
		case '"':
			return errFinished
		default:
			builder.WriteRune(r)
			return nil
		}
	})
	return builder.String(), err
}

var goodKeywords = map[string]bool{
	"true":  true,
	"false": true,
	"null":  true,
}

func (ctr *FieldCounter) parseValue() error {
	return ctr.readLoop(func(r rune) error {
		if unicode.IsSpace(r) {
			return nil
		} else if unicode.IsLetter(r) {
			if keyword, err := ctr.parseKeyword(r); err != nil {
				return ctr.wrapCallError(err, msgParseKeyword, false, false)
			} else if !goodKeywords[strings.ToLower(keyword)] {
				return fmt.Errorf("bad keyword '%s'", keyword)
			}
		} else if r == '"' {
			_, err := ctr.parseString()
			return ctr.wrapCallError(err, msgParseString, true, false)
		} else if unicode.IsDigit(r) || r == '-' {
			return ctr.wrapCallError(ctr.parseNumber(), msgParseNumber, true, false)
		} else if r == '[' {
			return ctr.wrapCallError(ctr.parseArray(), msgParseArray, true, false)
		} else if r == '{' {
			return ctr.wrapCallError(ctr.parseObject(false), msgParseObject, true, false)
		} else {
			return errUnexpected(r)
		}
		return nil
	})
}

var errFinished = errors.New("finished readLoop")

func (ctr *FieldCounter) readLoop(fn func(r rune) error) error {
	for {
		r, err := ctr.readRune()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return fmt.Errorf("read rune: %w", err)
		}
		if err := fn(r); errors.Is(err, errFinished) {
			break
		} else if err != nil {
			return err
		}
	}
	return nil
}

func (ctr *FieldCounter) readRune() (rune, error) {
	r, _, err := ctr.reader.ReadRune()
	if err == nil {
		ctr.charLoc++
		if r == '\n' {
			r = ' '
		}
		ctr.recent = append(ctr.recent, r)
		if len(ctr.recent) > 15 {
			ctr.recent = ctr.recent[5:]
		}
	}
	return r, err
}

func (ctr *FieldCounter) unreadRune() error {
	err := ctr.reader.UnreadRune()
	if err == nil {
		if len(ctr.recent) > 0 {
			ctr.charLoc--
			ctr.recent = ctr.recent[:len(ctr.recent)-1]
		}
	}
	return err
}

func (ctr *FieldCounter) wrapCallError(err error, msg string, finished bool, status bool) error {
	if errors.Is(err, errFinished) {
		return err
	} else if err != nil {
		if status {
			return fmt.Errorf("%s %s: %w", msg, ctr, err)
		} else {
			return fmt.Errorf("%s: %w", msg, err)
		}
	} else if finished {
		return errFinished
	} else {
		return nil
	}
}

func errUnexpected(r rune) error {
	return fmt.Errorf("unexpected rune: %c", r)
}
