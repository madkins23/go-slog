# Technical Details for `SlogBenchmarkSuite`

## Benchmark Suite

The core code for verifying `slog` handlers is in `bench/tests/SlogBenchmarkSuite`,
which was generated from scratch for this repository in the absence of
usable preexisting alternatives.

The `bench/tests` directory contains, at this time (2024-01-16),
nothing but the code for `SlogBenchmarkSuite`.
Since there are a lot of tests they have been broken up into separate files
by category or functionality.

The main files are:

* `suite.go`  
  Declares `SlogBenchmarkSuite` and a few top-level methods.
* `benchmarks.go`  
  Actual benchmark tests.
* `checks.go`  
  Common checks used by various benchmark tests.
* `logging.go`  
  Code to load test data cases from `logging.txt`
* `utility.go`  
  `SlogTestSuite` utility methods used in multiple places in the test suite.

Supporting files:

* `README.md`  
  The file you are currently reading.
* `doc.go`  
  Source code documentation for the `bench/tests` package.
* `data.go`  
  A few data items that are used in multiple places in the test suite.
* `warnings.go`  
  Custom `internal/warning.Manager` for the benchmark test suite.

Inherited:

* [`infra/warning.Manager`](https://pkg.go.dev/github.com/madkins23/go-slog/infra/warning#Manager)  
  The code that manages benchmark warnings is currently located in the `internal/test` package.

## Benchmark Tests

Benchmark tests are defined by the
[`Benchmark` structure](https://pkg.go.dev/github.com/madkins23/go-slog/bench/tests#Benchmark).

## Test Execution

The main part of the test harness is in the
[`bench/tests.Run`](https://pkg.go.dev/github.com/madkins23/go-slog/bench/tests#Run) function.

* For each method name beginning with `Benchmark`:
  * Execute the method, returning an pointer to an object of class `Benchmark`.
  * If the `Benchmark` has a handler function[^1]  
    then the `Creator` must be able to provide a `Handler` (some can't),  
    or else a `Warning` is logged and the test is skipped.
  * If the `Benchmark` has a verify function to test the log output
    (or more accurately, to test the test itself) then:
    * Get a logger, using the handler function if present.
    * Run a single test using that logger.
    * Verify the output with the function.
  * If the `-justTests` flag is false (not set):
    * Get a logger, using the handler function if present.
    * The Go test harness is used to run the `Benchmark` test function
      in parallel in ever-larger batches until enough testing has been done.
    * The test harness emits a line of data with results of the test.

[//]: # (Remove the following --- if github footnotes are ever implemented in pkg.go.dev per https://github.com/golang/go/issues/65922)

---

[^1]: A `HandlerFn` is used to adjust a newly created `Handler`
      prior to using it to create a custom `Logger`,
      instead of the normal mechanism which returns a generic `Logger`.
      Some `Logger` customisations must be done by
      manipulating the `Handler` in this fashion.