// Package replace defines functions that can be used as slog.HandlerOptions.ReplaceAttr values.
// For example:
//
//   - Remove attributes with empty key strings.
//   - Change level attributes named "lvl" to be named slog.LevelKey.
//   - Change message attributes named "message" to be named slog.MessageKey.
//   - Remove the "time" basic attribute.
//
// An early intention was to "fix" non-conformant handlers by providing
// ReplaceAttr functions that could make them conformant.
// Unfortunately, other issues prevent this usage in many cases:
//
//   - Attributes can't be directly removed, they can only be made empty
//     (defined as both an empty key and a nil value, as returned by [infra.EmptyAttr]),
//     but some handlers tested don't remove empty attributes as they should
//     so this fix doesn't work for them.
//   - Some handlers don't recognize `slog.HandlerOptions.ReplaceAttr`.
//   - Those that do don't always recognize them for the basic fields
//     (`time`, `level`, `message`, and `source`).
//
// # Function Types
//
// The ReplaceAttr option takes a function that will be applied to each attribute that is to be logged.
// The attribute may be returned as is or altered by the function.
// The function is defined inline in the HandlerOptions.ReplaceAttr field declaration
// but is not provided as a separate declaration.
// A convenience definition is found in [infra.AttrFn].
//
// Most of the functions in this package do not directly conform to this type,
// but instead return a function that conforms to this type.
// This provides more functionality in a smaller number of general purpose ReplaceAttr functions.
//
// # Arguments
//
// Commonly used ReplaceAttr arguments include:
//
//	caseInsensitive bool
//
// Most (if not all) ReplaceAttr in this package match against a field name.
// This match is much faster if it is exact (case-sensitive) so that is the default (false).
// Setting the caseInsensitive flag causes all matches to be done using strings.EqualFold.
//
//	grpChk GroupCheckFn
//
// A function that determines if the current stack of group names is acceptable.
// For example, replace.TopCheck returns true only if the stack is empty,
// indicating the attribute is not inside a group.
//
// [infra.AttrFn]: https://pkg.go.dev/github.com/madkins23/go-slog/infra#AttrFn
// [infra.EmptyAttr]: https://pkg.go.dev/github.com/madkins23/go-slog/infra#EmptyAttr
package replace
