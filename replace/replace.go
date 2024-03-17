package replace

import (
	"log/slog"
	"strings"

	"github.com/madkins23/go-slog/infra"
)

type groupCheck func(groups []string) bool

var _ groupCheck = TopCheck

func TopCheck(groups []string) bool {
	return len(groups) == 0
}

// -----------------------------------------------------------------------------

// ChangeKey maps keys from one string to another ('from' to 'to').
// The noCase argument can be set to convert all strings to lower case before comparing them.
// The grpChk argument is a function that will be applied to the groups passed in
// to determine whether to change the key, returning bool if the key should be changed.
func ChangeKey(from, to string, caseInsensitive bool, grpChk groupCheck) infra.AttrFn {
	return func(groups []string, a slog.Attr) slog.Attr {
		var found bool
		if caseInsensitive {
			// TODO: Figure strings.EqualFold() is too slow but havent tested.
			found = strings.ToLower(a.Key) == strings.ToLower(from)
		} else {
			found = a.Key == from
		}
		if found && (grpChk == nil || grpChk(groups)) {
			a.Key = to
		}
		return a
	}
}
