# Benchmarking `log/slog` Handlers

The `bench` package provides various `log/slog` (henceforth just `slog`) handler benchmark tests.
This document discusses simple usage details.
Technical details for the test suite are provided in
the [`README.md`](https://pkg.go.dev/github.com/madkins23/go-slog/bench/tests#section-readme) file in
the [`tests`](tests) package subdirectory.

## Simple Example

Benchmarking a `slog` handler using the `benchmark` test suite is fairly simple.
The following [code](https://github.com/madkins23/go-slog/blob/main/bench/slog_json_test.go)
runs the test suite on `slog.JSONHandler`:

```go
package bench

import (
  "testing"

  "github.com/madkins23/go-slog/bench/tests"
  "github.com/madkins23/go-slog/creator/slogjson"
)

// BenchmarkSlogJSON runs benchmarks for the slog/JSONHandler JSON handler.
func BenchmarkSlogJSON(b *testing.B) {
  slogSuite := tests.NewSlogBenchmarkSuite(slogjson.Creator())
  tests.Run(b, slogSuite)
}
```

The file itself must have the `_test.go` suffix and
contain a function with a name of the pattern `Benchmark<tag_name>`
where `<tag_name>` will likely be something like `PhsymZerolog` or `SlogJSON`.

The first line in `BenchmarkSlogJSON` creates a new test suite.
The argument to the `NewSlogBenchmarkSuite` function is
an [`infra.Creator`](https://pkg.go.dev/github.com/madkins23/go-slog/infra#Creator) object,
which is responsible for creating new `slog.Logger`
(and optionally `slog.Handler`) objects for benchmarks.

In this case an appropriate factory is created by the pre-existing
[`slogjson.Creator`](https://pkg.go.dev/github.com/madkins23/go-slog/creator/slogjson#Creator) function.
In order to test a new handler instance
(one that has not been tested in this repository)
it is necessary to [create a new `infra.Creator`](https://pkg.go.dev/github.com/madkins23/go-slog/infra#readme-creator) for it.
Existing examples can be found in the `creator` package.

Finally, the suite is executed via the
[`bench/tests.Run`](https://pkg.go.dev/github.com/madkins23/go-slog/bench/tests#Run) function,
passing in the test suite object.

In short:
* The `BenchmarkXxx` function is executed by the [Go test harness](https://pkg.go.dev/testing).
* The test function configures a `SlogBenchmarkSuite` using an `infra.Creator` factory object.
* The test function executes the test suite `test.Run()`.

More examples are available in this package.

In addition, there is a [`main_test.go`](https://github.com/madkins23/go-slog/blob/main/bench/main_test.go) file which exists to provide
a global resource to the other tests ([described below](#testmain)).

### Running Benchmarks

Run the handler verification tests installed in this repository with:
```shell
go test -bench=. bench/*.go
```

It is not necessary to use the `-benchmem` argument to generate memory statistics.
This feature is turned on within the benchmark harness (nor can it be turned off).

Due to the way Go benchmark testing is configured
it is not possible to gather results internally.
Processing of results must be done using external tools
such as the commands
[`tabular`](https://pkg.go.dev/github.com/madkins23/go-slog/cmd/tabular) and
[`server`](https://pkg.go.dev/github.com/madkins23/go-slog/cmd/server) in this repository.

On an operating system that supports `bash` scripts you can use
the [`scripts/bench`](https://github.com/madkins23/go-slog/blob/main/scripts/bench) script which is configured
with appropriate post-processing via `scripts/tabulate`.

#### Test Flags

There are two flags defined for testing the verification code:
* `-debug=<level>`  
  Sets an integer level for showing any `Debugf()` statements in the code.
* `-justTests`
  Just run benchmark verification tests, not the actual benchmarks (see [below](#supporting-tests)).

### Supporting Tests

In addition to the benchmarks there are tests that verify the benchmarks.
The goal of these tests is to make sure that the benchmark is actually testing something.

The supporting tests are not the same as normal Go test harness tests:
* they don't use the standard test assertions and
* they report issues via the [`WarningManager`](https://pkg.go.dev/github.com/madkins23/go-slog/internal/warning#Manager).

When running benchmarks the warning data from supporting tests is specified at the end of the output (as of 2024-02-22):
```
# Warnings for chanchal/ZapHandler:
#   Implied
#      1 [SourceKey] Source data not logged when AddSource flag set
#          SimpleSource: no source key
#            {"level":"info","time":"2024-02-21T12:21:40-08:00","msg":"This is a message"}
#
# Warnings for phsym/zeroslog:
#   Implied
#      1 [SourceKey] Source data not logged when AddSource flag set
#          SimpleSource: no source key
#            {"level":"info","caller":"/home/marc/work/go/src/github.com/madkins23/go-slog/bench/tests/benchmarks.go:70","time":"2024-02-21T12:21:58-08:00","message":"This is a message"}
#
# Warnings for samber/slog-zap:
#   Implied
#      1 [SourceKey] Source data not logged when AddSource flag set
#          SimpleSource: no source key
#            {"level":"info","time":"2024-02-21T12:22:33-08:00","caller":"tests/benchmarks.go:70","msg":"This is a message"}
#
#  Handlers by warning:
#   Implied
#     [SourceKey] Source data not logged when AddSource flag set
#       chanchal/ZapHandler
#       phsym/zeroslog
#       samber/slog-zap
```

The prefixed [octothorpe](https://en.wiktionary.org/wiki/octothorpe)
characters (`#`, often referred to as "pound signs")
are used to mark the warning output for later consumption by
result display commands (i.e.
[`tabular`](https://pkg.go.dev/github.com/madkins23/go-slog/cmd/tabular) and
[`server`](https://pkg.go.dev/github.com/madkins23/go-slog/cmd/server)).

### Making a Benchmark Test

Benchmark tests can live in any repository,
though it may not make as much sense to run benchmarks for a single handler.
Handler authors may want to do this when making changes to the code.

* Build an appropriate `Creator` object  
  When testing a single handler it will be necessary to point to the local handler code,
  whereas the provided `Creator` object within `go-slog` will point to released code.
* Build a Benchmark test function
* Run the benchmark tests
* Process the data for consumption:
  - `tabular` generates text output in tabular form
  - `server` provides tabular and chart data plus warnings

## Creators

`Creator` objects are factories for generating new `slog.Logger` objects.
Detailed documentation on defining and using `Creator` objects is provided in
the [`infra` package](https://pkg.go.dev/github.com/madkins23/go-slog/infra#readme-creator).

## Caveats

* Actual testing is done by calling through a `slog.Logger` object.
* Benchmark tests only operate against the final call (e.g. `Info()`).
  The initial creation of a `slog.Logger` object,
  which may include the use of `Handler.WithAttrs()` and/or `Handler.WithGroup()` calls,
  is not measured as it is generally an initialization step and (in theory) only called once,
  whereas the logging calls may be executed many times.
* Text and console handlers don't have a consistent format.
  While it might be useful to test those handlers as well,
  the difficulty of parsing various output formats argues against it.[^1]
* The `-useWarnings` flag tends to result in the results being buried in the normal `go test` output.
  This can be fixed by implementing a global [`TestMain()`](#testmain) function.

### `TestMain`

Normally the warning results will show up in a block in the middle of `go test` output.
This is due to the way the default test harness works.

It is possible to override the default test harness by defining a global function
[`TestMain()`](https://pkg.go.dev/testing#hdr-Main).
The `bench/tests` package provides a convenient function to support this.
Define the following `TestMain()` function:
```go
func TestMain(m *testing.M) {
    tests.WithWarnings(m)
}
```

This function may be defined in the same `_test.go` file as the handler test.
If multiple handler tests are in the same directory:

* It will be necessary to move the `TestMain()` definition to a separate file,
  such as the [`bench/main_test.go`](main_test.go).
* An addition listing of which handlers throw each warning
  will be added after the normal output.
