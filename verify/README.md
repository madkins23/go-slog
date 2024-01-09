# Verifying `log/slog` Handlers

The `verify` package provides various `log/slog` handler test suites.
This document discusses simple usage details.
Technical details for the test suite are provided in
the [`README.md`](test/README.md) file in
the [`test`](test) package subdirectory.

## Simple Example

Verification of a `log/slog` handler has been made fairly simple.
The following application runs the test suite on `log/slog.JSONHandler`:

```go
import (
    "io"
    "log/slog"
    "testing"

    "github.com/stretchr/testify/suite"

    "github.com/madkins23/go-slog/verify/test"
)

func TestMain(m *testing.M) {
    test.WithWarnings(m)
}

// Test_slog runs tests for the log/slog JSON handler.
func Test_slog(t *testing.T) {
    slogSuite := &test.SlogTestSuite{
        Creator: &SlogCreator{},
        Name:    "log/slog.JSONHandler",
    }
    if *test.UseWarnings {
        slogSuite.WarnOnly(test.WarnDuplicates)
    }
    suite.Run(t, slogSuite)
}

var _ test.LoggerCreator = &SlogCreator{}

type SlogCreator struct{}

func (creator *SlogCreator) SimpleLogger(w io.Writer) *slog.Logger {
    return slog.New(slog.NewJSONHandler(w, nil))
}

func (creator *SlogCreator) SourceLogger(w io.Writer) *slog.Logger {
    return slog.New(
    slog.NewJSONHandler(w, &slog.HandlerOptions{AddSource: true}))
}
```

The file itself must have the `_test.go` suffix in order to be executed as a test.

### More Examples

This package contains several examples, including the one above:
* `slog_test.go`
  Verifies the standard `log/slog.JSONHandler`.
* `slog_zerolog_phsym_test.go`
  Verifies the `zeroslog` handler.
* `slog_zerolog_samber_test.go`
  Verifies the `slog-zerolog` handler.

In addition, there is a `main_test.go` file which exists to provide
a global resource to the other tests.
