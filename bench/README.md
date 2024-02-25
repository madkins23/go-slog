
# Benchmarking `log/slog` Handlers

The `bench` package provides various `log/slog` (henceforth just `slog`) handler benchmark suites.
This document discusses simple usage details.
Technical details for the test suite are provided in
the [`README.md`](tests/README.m4) file in
the [`tests`](https://pkg.go.dev/github.com/madkins23/go-slog/bench/tests) package subdirectory.

## Making a Benchmark Test

Benchmark tests can live in any repository,
though it may not make as much sense to run benchmarks for a single handler.
Handler authors may want to do this when making changes to the code.

:construction: **TBD** :construction:

* Build a `Creator` object
* Build a Benchmark test function
* Run the benchmark tests
* Process the data for consumption using
  - tabular generates text output in tabular form
  - server provides tabular and chart data plus warnings

### Simple Example

Benchmarking a `slog` handler using the `benchmark` test suite is fairly simple.
The following [code](https://github.com/madkins23/go-slog/blob/main/bench/slog_test.go)
runs the test suite on `slog.JSONHandler`:

```go
package bench

import (
    "testing"

    "github.com/madkins23/go-slog/bench/tests"
    "github.com/madkins23/go-slog/creator/slogjson"
)

// BenchmarkSlogJSON runs benchmarks for the log/slog JSON handler.
func BenchmarkSlogJSON(b *testing.B) {
    slogSuite := tests.NewSlogBenchmarkSuite(slogjson.Creator(), "SlogJSON")
    tests.Run(b, slogSuite)
}
```

The file itself must have the `_test.go` suffix and
contain a function with a name of the pattern `Benchmark<tag_name>`
where `<tag_name>` will likely be something like `PhsymZerolog` or `SlogJSON`.

The first line in `BenchmarkSlogJSON` creates a new test suite.
The argument to the `NewSlogBenchmarkSuite` function is an [`infra.Creator`](../infra/creator.go) object,
which is responsible for creating new `slog.Logger`
(and optionally `slog.Handler`) objects for benchmarks.

In this case an appropriate factory is created by the `creator.Slog` function
that is already defined in the `creator` package.
In order to test a new handler instance
(one that has not been tested in this repository)
it is necessary to [create a new `infra.Creator`](#creators) for it.
Existing examples can be found in the `creator` package.

Finally, the suite is run via its `Run` method.
In short:
* The `BenchmarkXxx` function is executed by the [Go test harness](https://pkg.go.dev/testing).
* The test function configures a `SlogTestSuite` using an `infra.Creator` factory object.
* The test function executes the test suite via its `Run` method.

#### More Examples

This package contains several examples, including the one above:
* [`slog_test.go`](https://github.com/madkins23/go-slog/blob/main/verify/slog_test.go)
  Verifies the [standard `slog.JSONHandler`](https://pkg.go.dev/log/slog@master#JSONHandler).
* [`slog_phsym_zerolog_test.go`](https://github.com/madkins23/go-slog/blob/main/verify/slog_phsym_zerolog_test.go)
  Verifies the [`phsym zeroslog` handler](https://github.com/phsym/zeroslog/tree/2bf737d6422a5de048845cd3bdd2db6363555eb4).
* [`slog_samber_zap_test.go`](https://github.com/madkins23/go-slog/blob/main/verify/slog_samber_zerolog_test.go)
  Verifies the [`samber slog-zap` handler](https://github.com/samber/slog-zap).
* [`slog_samber_zerolog_test.go`](https://github.com/madkins23/go-slog/blob/main/verify/slog_samber_zerolog_test.go)
  Verifies the [`samber slog-zerolog` handler](https://github.com/samber/slog-zerolog).

In addition to the test files for individual handlers,
there is a [`main_test.go`](https://github.com/madkins23/go-slog/blob/main/verify/main_test.go) file which exists to provide
a global resource to the other tests ([described below](#testmain)).

## Running Benchmarks

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
with appropriate post-processing via
[`scripts/tabulate`](https://github.com/madkins23/go-slog/blob/main/scripts/tabulate) or
[`scripts/server`](https://github.com/madkins23/go-slog/blob/main/scripts/server).

#### Test Flags

* `-debug=<level>`
  Sets an integer level for showing any `Debugf()` statements in the code.
* `-justTests`
  Just run benchmark verification tests, not the actual benchmarks.


### Supporting Tests

In addition to the benchmarks there are tests that verify the benchmarks.
The goal of these tests is to make sure that the benchmark is actually testing something.

The supporting tests are not the same as normal Go test harness tests:
* they don't use the standard test assertions and
* they report issues via the [`WarningManager`](../infra/warnings.go).

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
result display commands (e.g.
[`tabular`](https://pkg.go.dev/github.com/madkins23/go-slog/cmd/tabular) and
[`server`](https://pkg.go.dev/github.com/madkins23/go-slog/cmd/server)).

## Caveats

* Actual testing is done by calling through a `slog.Logger` object.
* Documentation for functions in `_test.go` files in this directory
  is not included in [`pkg.go.dev`](https://pkg.go.dev/github.com/madkins23/go-slog/bench)
