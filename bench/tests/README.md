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
  Warnings specific to the benchmark suite.

Inherited:

* [`test.WarningManager`](https://github.com/madkins23/go-slog/blob/main/internal/test/warnings.go)  
  The code that manages benchmark warnings is currently located in the `internal/test` package.

## Benchmark Tests

Benchmark tests are defined by the `Benchmark` structure.

Each test _must_ have:
* a pointer to `slog.HandlerOptions` to be used in generating the `slog.Logger`,
* a pointer to a `BenchmarkFn` which executes the actual benchmark test,

and _may_ have:
* an optional pointer to a `HandlerFn` which is used to adjust
  the `slog.Handler` object (if available) before constructing the `slog.Logger`, and
* an optional pointer to a `VerifyFn` which is used to verify the test.

For example:

```Go
//goland:noinspection GoSnakeCaseUsage
func (suite *SlogBenchmarkSuite) Benchmark_Simple() *Benchmark {
    return &Benchmark{
        Options: infra.SimpleOptions(),
        BenchmarkFn: func(logger *slog.Logger) {
            logger.Info(message)
        },
        VerifyFn: matcher("Simple", expectedBasic()),
    }
}```
