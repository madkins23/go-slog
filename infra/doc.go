// Package infra contains functionality shared between test and benchmark managers.
//
// The definitions in this package are publicly visible and may in some cases be usable elsewhere,
// as opposed to the `internal` package which also provides shared infrastructure.
//
// Convenient definitions:
//   - AttrFn defines the slog.HandlerOptions.ReplaceAttr function template.
//   - EmptyAttr() returns an empty attribute.
//   - Creator struct used to generate slog.Handler and/or slog.Logger objects.
//   - Functions to return slog.HandlerOptions of general utility.
package infra
