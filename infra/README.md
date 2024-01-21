# Infrastructure

This package contains infrastructure shared by various packages.
The main packages sharing this code are [`verify`](../verify) and [`bench`](../bench).
The definitions in this package are publicly visible and may in some cases be usable elsewhere.

## Attributes

A few useful items relating to [`slog.Attr`](https://pkg.go.dev/log/slog@master#Attr):

* `AttrFn` defines a type useful when working with
  [`HandlerOptions.ReplaceAttr`](https://pkg.go.dev/log/slog@master#HandlerOptions)
  functionality.
* The `EmptyAttr` function returns a new empty `slog.Attr` object.

## Creator

A [`Creator`](../infra/creator.go) object is a factory used to generate
`slog.Logger` objects for testing.
A number of predefined `Creator` objects can be found in the [`creator` package](../creator).
The simplest of these returns loggers for the
[`slog.NewJSONHandler`](https://pkg.go.dev/log/slog@master#JSONHandler)
handler.

Creation of a new `infra.Creator` object is fairly simple.
The `infra.NewCreator` function takes the name of the handler package
and one or two functions as appropriate:
* A `CreateHandlerFn` creates a new `slog.Handler`.
* A `CreateLoggerFn` creates a new `slog.Logger`.

At least one of the two is required, though the `CreateHandlerFn` is preferred,
as in the following example:

```Go
package creator

import (
	"io"
	"log/slog"

	"github.com/madkins23/go-slog/infra"
)

// Slog returns a Creator object for the log/slog.JSONHandler.
func Slog() infra.Creator {
	return infra.NewCreator("log/slog.JSONHandler", SlogHandlerFn, nil)
}

func SlogHandlerFn(w io.Writer, options *slog.HandlerOptions) slog.Handler {
	return slog.NewJSONHandler(w, options)
}
```

`Creator` factories can generate both `slog.Handler` and `slog.Logger` objects.
Most tests use the latter, but a few tests require the former.

* If only a `CreateLoggerFn` is provided the `Creator.NewHandler` method returns `nil`
  and the `Creator.CanMakeHandler` method returns `false`.
  Tests requiring handler creation use the latter method to skip the test.
* If only a `CreateHandlerFn` is provided the `Creator.NewHandler` method
  uses that function to return a new handler and the `Creator.NewLogger` method
  also uses that function, then wraps the handler in `slog.New`.
* If both functions are provided, they will each be used for the appropriate method.
  This has yet to be required.

The reason for these two functions is the existence of the [`darvaza zerolog` handler](https://pkg.go.dev/darvaza.org/slog/handlers/zerolog).
This handler is implemented directly as a `slog.Logger`,
without defining a `slog.Handler` interface.[^1]

## Options

The `options` package provides some predefined
[`slog.HandlerOptions`](https://pkg.go.dev/log/slog@master#HandlerOptions) objects.
These are used in testing and benchmarks and may have some utility elsewhere.

## Utility

There is a single utility function: `CurrentFunctionName`.
This function will look up the stack for the first function with the specified prefix
(e.g. `Test` or `Benchmark`) and return that function name.
It is currently (2024-01-19) only used in the `WarningManager`,
also in this package, so it may not need to be public.

## Warnings

A "warning" facility is built into many of the tests to provide a way to:
* avoid scanning through a lot of `go test` error logs in detail over and over,
* get a compressed set of warnings about issues after testing is done, and
* provide a list of known issues in the test suite.

Since this code is currently used primary in handler verification,
there are better examples in the [`verify` `README`](../verify/README.md) file.

The actual warnings are defined in the [`warning`](../warning) package.

### Usage

* Define a `WarningManager`.
* Predefine various `Warning` objects that may be used in testing.
* When defining a test suite use `WarnOnly` to specify that
  during testing warning code should be run instead of assertions.
* In test code check the `WarningManager` for applicable warnings.
* In warning-specific code use `AddWarning` to note the condition exists
  or `UnusedWarning` to note that the warning is redundant.
* Run the tests with the `-useWarnings` flag to invoke warning code.
  Without this flag the `WarningManager` will never flag warning code
  and test assertions will raise conventional errors.

Usage of the `WarningManager` is described in more detail
in the [`verify` `README`](../verify/README.md).

---

[^1]: This may or may not be a desirable thing.
      On the one hand, there is a lot of useful code in the `slog` package outside the handler.
      On the other hand, replacing that code might provide some advantage.
