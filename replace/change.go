package replace

import (
	"log/slog"
	"strings"

	"github.com/madkins23/go-slog/infra"
)

// -----------------------------------------------------------------------------

// ChangeFn converts a value to another value without reference to a containing slog.Attr or key.
type ChangeFn func(value slog.Value) slog.Value

// SetValueTo returns a ChangeFn that will set the value of an attribute to the specified value.
// The specified value is fixed at the time SetValueTo is executed.
func SetValueTo(value slog.Value) ChangeFn {
	return func(_ slog.Value) slog.Value {
		return value
	}
}

// ChangeCases is an enumerated type for type case settings.
type ChangeCases uint8

const (
	CaseNone ChangeCases = iota
	CaseLower
	CaseUpper
)

var caseFn = map[ChangeCases]func(string) string{
	CaseLower: strings.ToLower,
	CaseUpper: strings.ToUpper,
}

// -----------------------------------------------------------------------------

// ChangeCase changes the case of the string value of any attribute
// matching the specified key to upper-, or lower-case as required. For example:
//
//	options := &slog.HandlerOptions{
//		ReplaceAttr: replace.ChangeCase(
//			slog.LevelKey, replace.CaseLower, false, replace.TopCheck)
//	}
//
// returns an infra.AttrFn that will match on attributes with
// the key "level" matched precisely (not case-insensitive)
// at the top level of the log record (not within a group)
// and change the string value of that attribute to lower case.
//
// The second argument is of type replace.ChangeCases with values
// replace.CaseNone, replace.CaseLower, and replace.CaseUpper.
// The default value is replace.CaseNone which means no change.
func ChangeCase(key string, chgCase ChangeCases, caseInsensitive bool, grpChk GroupCheckFn) infra.AttrFn {
	if chgCase != CaseNone {
		if fn, found := caseFn[chgCase]; found {
			// Return a function that changes the case of any string per chgCase.
			return ChangeValue(key, func(val slog.Value) slog.Value {
				return slog.StringValue(fn(val.String()))
			}, caseInsensitive, grpChk)
		}
		slog.Error("unknown chgCase value", "value", chgCase)
	}
	// Return a function that just returns whatever value is passed in.
	return func(groups []string, a slog.Attr) slog.Attr {
		return a
	}
}

// ChangeKey changes the key of the attribute with the specified new value.
//
//	options := &slog.HandlerOptions{
//		ReplaceAttr: replace.ChangeKey(
//			"message", slog.MessageKey, false, replace.TopCheck)
//	}
//
// returns an infra.AttrFn that will match on attributes with the key "message" and
// change it to the value of slog.MessageKey (which is "msg").
func ChangeKey(from, to string, caseInsensitive bool, grpChk GroupCheckFn) infra.AttrFn {
	return func(groups []string, a slog.Attr) slog.Attr {
		if checkFieldGroups(groups, a, from, caseInsensitive, grpChk) {
			a.Key = to
		}
		return a
	}
}

// ChangeValue changes the value of the string value of any attribute matching the specified key.
// The new value is generated by executing a ChangeFn:
//
//	type ChangeFn func(value slog.Value) slog.Value
//
// provided by the caller.
func ChangeValue(key string, chgFn ChangeFn, caseInsensitive bool, grpChk GroupCheckFn) infra.AttrFn {
	return func(groups []string, a slog.Attr) slog.Attr {
		if checkFieldGroups(groups, a, key, caseInsensitive, grpChk) {
			a.Value = chgFn(a.Value)
		}
		return a
	}
}

// RemoveKey removes a specified attribute by changing it to
// the empty attribute (`slog.Attr{}`) which is supposed to be ignored.
//
//	options := &slog.HandlerOptions{
//		ReplaceAttr: replace.RemoveKey(slog.TimeKey, false, TopCheck)
//	}
//
// removes the slog.TimeKey ("time") key so that the time will not be shown.
func RemoveKey(key string, caseInsensitive bool, grpChk GroupCheckFn) infra.AttrFn {
	return func(groups []string, a slog.Attr) slog.Attr {
		if checkFieldGroups(groups, a, key, caseInsensitive, grpChk) {
			return infra.EmptyAttr()
		}
		return a
	}
}

// -----------------------------------------------------------------------------

// checkFieldGroups runs basic tests common to change functions above.
// The key of the provided attribute is checked per the caseInsensitive argument.
// If grpChk is non-nil it will be applied to the groups array.
// The result is true if the calling change function applies to this attribute.
func checkFieldGroups(groups []string, a slog.Attr, key string, caseInsensitive bool, grpChk GroupCheckFn) bool {
	var found bool
	if caseInsensitive {
		// Faster than converting two strings to the same case.
		// See BenchmarkCompareChangeCase and BenchmarkCompareEqualFold in change_test.go.
		found = strings.EqualFold(a.Key, key)
	} else {
		found = a.Key == key
	}
	return found && (grpChk == nil || grpChk(groups))
}
