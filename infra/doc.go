// Package infra contains functionality shared between test and benchmark managers.
//
// The definitions in this package are publicly visible and may in some cases be usable elsewhere,
// as opposed to the [internal] package which also provides shared infrastructure.
//
// # Attribute Declarations
//
//   - AttrFn defines the HandlerOptions.ReplaceAttr function template.
//   - EmptyAttr() returns an empty attribute.
//
// # Creator Objects
//
// A [Creator] object is a factory used to generate slog.Logger objects for testing.
// A number of predefined Creator objects can be found in the [creator package].
//
// # Predefined Options
//
// Functions are provided to return various standard slog.HandlerOptions objects.
// These are used in testing and benchmarks and may have some utility elsewhere.
//
//   - SimpleOptions() provides a simple set of options for default usage.
//   - LevelOptions() provides a simple set of options with the specified slog.Leveler.
//   - SourceOptions() provides a simple set of options that adds source file data.
//   - ReplaceAttrOptions() provides a simple set of options with the specified AttrFn.
//
// # Warnings
//
// The [warning] sub-package provides the warning manager and predefined warnings.
//
// [Creator]: https://pkg.go.dev/github.com/madkins23/go-slog/infra#Creator
// [creator package]: https://pkg.go.dev/github.com/madkins23/go-slog/creator
// [internal]: https://pkg.go.dev/github.com/madkins23/go-slog/internal
// [warning]: https://pkg.go.dev/github.com/madkins23/go-slog/warning
package infra
