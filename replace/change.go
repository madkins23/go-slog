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

// ChangeCase returns a ChangeFn that will change the case of a string value.
func ChangeCase(key string, chgCase ChangeCases, caseInsensitive bool, grpChk GroupCheck) infra.AttrFn {
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

// ChangeKey maps keys from one string to another ('from' to 'to').
// The caseInsensitive argument can be set to convert all strings to lower case before comparing them.
// The grpChk argument is a function that will be applied to the groups passed in
// to determine whether to change the key, returning bool if the key should be changed.
func ChangeKey(from, to string, caseInsensitive bool, grpChk GroupCheck) infra.AttrFn {
	return func(groups []string, a slog.Attr) slog.Attr {
		if checkFieldGroups(groups, a, from, caseInsensitive, grpChk) {
			a.Key = to
		}
		return a
	}
}

// ChangeValue changes the value for the specified key to the result of executing
// the value of argument chgFn against the current value associated with the specified key.
// The caseInsensitive argument can be set to convert all strings to lower case before comparing them.
// The grpChk argument is a function that will be applied to the groups passed in
// to determine whether to change the key, returning bool if the key should be changed.
func ChangeValue(key string, chgFn ChangeFn, caseInsensitive bool, grpChk GroupCheck) infra.AttrFn {
	return func(groups []string, a slog.Attr) slog.Attr {
		if checkFieldGroups(groups, a, key, caseInsensitive, grpChk) {
			a.Value = chgFn(a.Value)
		}
		return a
	}
}

// RemoveKey removes an attribute by field name.
// It is intended to be used as a ReplaceAttr function, to make example output deterministic.
func RemoveKey(key string, caseInsensitive bool, grpChk GroupCheck) infra.AttrFn {
	return func(groups []string, a slog.Attr) slog.Attr {
		if checkFieldGroups(groups, a, key, caseInsensitive, grpChk) {
			return infra.EmptyAttr()
		}
		return a
	}
}

// -----------------------------------------------------------------------------

func checkFieldGroups(groups []string, a slog.Attr, key string, caseInsensitive bool, grpChk GroupCheck) bool {
	var found bool
	if caseInsensitive {
		// TODO: Figure strings.EqualFold() is too slow but havent tested.
		found = strings.ToLower(a.Key) == strings.ToLower(key)
	} else {
		found = a.Key == key
	}
	return found && (grpChk == nil || grpChk(groups))
}
