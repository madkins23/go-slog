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

The reason for two functions is the possibility that an `slog.Logger` is available
but a `slog.Handler` is not.
This handler is implemented directly as a `slog.Logger`,
without defining a `slog.Handler` interface.[^1]

## Options

The `options` package provides some predefined
[`slog.HandlerOptions`](https://pkg.go.dev/log/slog@master#HandlerOptions) objects.
These are used in testing and benchmarks and may have some utility elsewhere.

---

[^1]: This may or may not be a desirable thing.
      On the one hand, there is a lot of useful code in the `slog` package outside the handler.
      On the other hand, replacing that code might provide some advantage.
