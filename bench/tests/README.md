# Technical Details for `SlogBenchmarkSuite`

## Benchmark Suite

The core code for verifying `slog` handlers is in `bench/tests/SlogBenchmarkSuite`,
which was generated from scratch for this repository in the absence of usable preexisting alternatives.

The `bench/tests` directory contains, at this time (2024-01-16),
nothing but the code for `SlogBenchmarkSuite`.
Since there are a lot of tests they have been broken up into separate files
by category or functionality.
The main files are:

* `suite.go`  
  Declares `SlogBenchmarkSuite` and a few top-level methods.
* `benchmarks.go`  
  Actual benchmark tests.
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
  Warnings specific to the benchmark suite.

Inherited:

* [`infra.WarningManager`](https://github.com/madkins23/go-slog/blob/main/infra/warnings.go)  
  The code that manages benchmark warnings is currently located in the `infra` package.
