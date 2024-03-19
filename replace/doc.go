// Package replace defines functions that can be used as slog.HandlerOptions.ReplaceAttr values.
// For example:
// * Remove attributes with empty key strings.
// * Change level attributes named "lvl" to be named slog.LevelKey.
// * Change message attributes named "message" to be named slog.MessageKey.
// * Remove the "time" basic attribute.

// Package replace provides a number of functions that can provide appropriate functions for use with
// HandlerOptions.ReplaceAttr.
//
// An early intention was to "fix" non-conformant handlers by providing
// ReplaceAttr functions that could make them conformant.
// Sadly this has not thus far turned out to be the case (2024-03-18).
//
// # Function Types
//
// The ReplaceAttr option takes a function that will be applied to each attribute that is to be logged.
// The attribute may be returned as is or altered by the function.
// The function is defined in the HandlerOptions.ReplaceAttr field declaration,
// but a convenience definition is found in go-slog/infra.AttrFn:
//
//	type AttrFn func(groups []string, a slog.Attr) slog.Attr
//
// This package provides functions with two key function signatures.
//
// The original functions (e.g. replace.RemoveEmptyKey)
// implement the type described above as infra.AttrFn,
// which is described in the HandlerOptions.ReplaceAttr field declaration.
// These functions are now deprecated in favor of the newer functions.
//
// This may turn out to be a mistake but there's no V2 on the horizon;
// they can just be de-deprecated at a later time.
// However, be aware that they are no longer tested or supported and are
// likely to be really inefficient compared to the newer functions.
//
// The newer replace functions do not conform to this function signature,
// rather they return functions that conform to infra.AttrFn.
// Thus setting the value of the ReplaceAttr field requires the execution
// of one of the functions, not the function itself.
//
// ## Older Functions
//
// As mentioned before, these functions are currently marked as deprecated.
// The also conform to infra.AttrFn and so the function pointer is set
// as the value of ReplaceAttr, not the value returned by executing the function.
//
// Rather than go into detail here, the relevant functions will be added
// to the following sections.
//
// ## Newer Functions
//
// Some of the newer functions share similar arguments:
//
// ### caseInsensitive bool
//
// Most (if not all) replace functions will match against a field name.
// This match is much faster if it is exact (case-_sensitive_) so that is the default (false).
// Setting the caseInsensitive flag causes all matches to be done after converting
// both strings into the same case.
//
// ### grpChk GroupCheck
//
// A function that determines if the current stack of group names is acceptable.
// For example, replace.TopCheck returns true only if the stack is empty,
// indicating the attribute is not inside  a group.
//
//	type GroupCheck func(groups []string) bool
package replace
