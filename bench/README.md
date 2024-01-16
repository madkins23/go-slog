# Benchmarking `log/slog` Handlers

The `bench` package provides various `log/slog` (henceforth just `slog`) handler benchmark suites.
This document discusses simple usage details.
Technical details for the test suite are provided in
the [`README.md`](tests/README.md) file in
the [`tests`](tests) package subdirectory.

### Simple Example

Benchmarking a `slog` handler using the `benchmark` test suite is fairly simple.
The following [code](https://github.com/madkins23/go-slog/blob/main/bench/slog_test.go)
runs the test suite on `slog.JSONHandler`:

```go
package bench

import (
    "testing"

    "github.com/madkins23/go-slog/bench/tests"
    "github.com/madkins23/go-slog/creator"
)

// Benchmark_slog runs benchmarks for the log/slog JSON handler.
func Benchmark_slog(b *testing.B) {
    slogSuite := tests.NewSlogBenchmarkSuite(creator.Slog())
    tests.Run(b, slogSuite)
}
```
The file itself must have the `_test.go` suffix and
contain a function with a name beginning with `Benchmark` in order to be executed as a benchmark.
In this instance `slog_test.go` contains function `Benchmark_slog`.

The first line in `Benchmark_slog` creates a new test suite.
The argument to the `NewSlogBenchmarkSuite` function is an [`infra.Creator`](../infra/creator.go) object,
which is responsible for creating new `slog.Logger`
(and optionally `slog.Handler`) objects for benchmarks.
In this case an appropriate factory is created by the `creator.Slog` function
that is already defined in the `creator` package.
In order to test a new handler instance
(one that has not been tested in this repository)
it is necessary to [create a new `infra.Creator`](#creators) for it.
Existing examples can be found in the `creator` package.

Finally the suite is run via its `Run` method.
In short:
* The `TestXxxxxx` function is executed by the [Go test harness](https://pkg.go.dev/testing).
* The test function configures a `SlogTestSuite` using an `infra.Creator` factory object.
* The test function executes the test suite via its `Run` method.

### More Examples

This package contains several examples, including the one above:
* [`slog_test.go`](https://github.com/madkins23/go-slog/blob/main/verify/slog_test.go)
  Verifies the standard [`slog.JSONHandler`](https://pkg.go.dev/log/slog@master#JSONHandler).
* [`slog_darvaza_zerolog_test.go`](https://github.com/madkins23/go-slog/blob/main/verify/slog_darvaza_zerolog_test.go)
  Verifies the [`zerolog` handler](https://pkg.go.dev/darvaza.org/slog/handlers/zerolog).
* [`slog_phsym_zerolog_test.go`](https://github.com/madkins23/go-slog/blob/main/verify/slog_phsym_zerolog_test.go)
  Verifies the [`zeroslog` handler](https://github.com/phsym/zeroslog/tree/2bf737d6422a5de048845cd3bdd2db6363555eb4).
* [`slog_samber_zap_test.go`](https://github.com/madkins23/go-slog/blob/main/verify/slog_samber_zerolog_test.go)
  Verifies the [`slog-zap` handler](https://github.com/samber/slog-zap).
* [`slog_samber_zerolog_test.go`](https://github.com/madkins23/go-slog/blob/main/verify/slog_samber_zerolog_test.go)
  Verifies the [`slog-zerolog` handler](https://github.com/samber/slog-zerolog).

In addition, there is a [`main_test.go`](https://github.com/madkins23/go-slog/blob/main/verify/main_test.go) file which exists to provide
a global resource to the other tests ([described below](#testmain)).

### Running Tests

Run the handler verification tests installed in this repository with:
```shell
go test -v -bench=. bench/*.go
```

Due to the way Go benchmark testing is configured
it is not possible to gather results internally.
Processing of results must be done using external tools
such as [`gobenchdata`](https://github.com/bobheadxi/gobenchdata)
and the command [`tabular`](../cmd/tabular/tabular.go)
in this repository.

On an operating system that supports `bash` scripts you can use
the [`scripts/bench`](https://github.com/madkins23/go-slog/blob/main/scripts/verify) script which is configured
with appropriate post-processing.

#### Test Flags

There is one flag defined for testing the verification code:
* `-debug=<level>`  
  Sets an integer level for showing any `Debugf()` statements in the code.
  As of 2024-01-11 there is only one in the test suite code at a level of `1`.
  This statement dumps the current JSON log record.[^1]

### Caveats

* Actually testing by calling through a `slog.Logger` object.
* Some tests are skipped because the require a `slog.Handler` object
  which is not available for some handler instances
  (e.g. [`darvaza`](https://github.com/darvaza-proxy/slog)) handlers.
