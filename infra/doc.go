// Package infra contains functionality shared between test and benchmark managers.
//
// The definitions in this package are publicly visible and may in some cases be usable elsewhere,
// as opposed to the internal package which also provides shared infrastructure.
//
// # Attribute Declarations
//
//   - AttrFn defines the HandlerOptions.ReplaceAttr function template.
//   - EmptyAttr() returns an empty attribute.
//
// # Creator Objects
//
// A [Creator] object is a factory used to generate slog.Logger objects for testing.
// A number of predefined `Creator` objects can be found in the [creator package].
//
// # Predefined Options
//
// Functions are provided to return various standard slog.HandlerOptions objects.
// These are used in testing and benchmarks and may have some utility elsewhere.
//
// [Creator]: https://pkg.go.dev/github.com/madkins23/go-slog/infra#Creator
// [creator package]: https://pkg.go.dev/github.com/madkins23/go-slog/creator
package infra
