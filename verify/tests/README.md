# Technical Details for `SlogTestSuite`

## Testing Suite

The core code for verifying `slog` handlers is in `verify/tests/SlogTestSuite`.
This type uses (and is composed with)
[`stretchr/testify/suite.Suite`](https://pkg.go.dev/github.com/stretchr/testify/suite).
All of the features of `testify/suite` are accessible to any `SlogTestSuite` or derivative.

The `verify/tests` directory contains, at this time (2024-01-15),
nothing but the code for `SlogTestSuite`.
Since there are a lot of tests they have been broken up into separate files
by category or functionality.
The main files are:

* `suite.go`  
  Declares `SlogTestSuite` and a few top-level methods.
* `documented.go`  
  Tests based on `slog` documentation.
  Comments for the tests include links to referenced statements therein.
* `slogtest.go`  
  Tests based on [`slogtest/slogtest`](https://pkg.go.dev/golang.org/x/exp/slog/slogtest),
  a test harness (presumably) developed along with `log/slog`.
  Comments for the tests reference the original `slogtest` tests
  as well as a few links to documentation.
* `duplicate.go`  
  Duplicate testing, which isn't currently regarded as an error.
  The status of this issue is currently
  [under discussion](https://github.com/golang/go/issues/59365).
* `replace.go`
  * Tests of
    [`slog.HandlerOptions.ReplaceAttr`](https://pkg.go.dev/golang.org/x/exp/slog#HandlerOptions)
    functionality, which appears to be
    [optional](https://github.com/golang/example/tree/master/slog-handler-guide#implementing-handler-methods).
  * Tests of replace functions defined in the `replace` package.
* `other.go`  
  Tests that don't seem to fit into any other category.
  These include log level functionality and log record time format.
* `checks.go`  
  Subtest methods that are called from multiple tests or are really long.
  Think of these as complex assertions.
* `utility.go`  
  `SlogTestSuite` utility methods used in multiple places in the test suite.

Supporting files:

* `README.md`  
  The file you are currently reading.
* `doc.go`  
  Source code documentation for the `verify/tests` package.
* `data.go`  
  A few data items that are used in multiple places in the test suite.
* `valuer.go`  
  A `struct` that implements the
  [`slog.LogValuer`](https://pkg.go.dev/log/slog@master#LogValuer)
  interface for testing.

Inherited:

* [`infra.WarningManager`](https://github.com/madkins23/go-slog/blob/main/infra/warnings.go)  
  The code that manages verification warnings is currently located in the `infra` package.
