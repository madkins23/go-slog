# Verifying `log/slog` Handlers

The `verify` package provides various `log/slog` handler test suites.
This document discusses simple usage details.
Technical details for the test suite are provided in
the [`README.md`](test/README.md) file in
the [`test`](test) package subdirectory.

The real benefit of `log/slog` is the ability to swap handlers without
rewriting all the log statements in existing code.
This only works if the various handlers behave in a similar manner.

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

## Caveats

The test harness that drives verification has some limitations.

_Verification only makes sense for JSON handlers_,
which are generally used to feed log records into downstream processing.
Text and console handlers don't have a consistent format.
While it might be useful to test those handlers as well,
the difficulty of parsing various output formats argues against it.

_The test suite interface to `log/slog` logging is via `slog.Handler`_.
This is a well-defined interface and the obvious place to swap in
a different logging backend.
This makes loggers that directly implement `log/slog`-like behavior
without instantiating a `slog.Handler`
(e.g. the [darvaza loggers](https://github.com/darvaza-proxy/slog))
inappropriate for this test suite.

## Warnings

A "warning" facility is built into many of the tests to provide a way to:
* avoid scanning through a lot of `go test` error logs in detail over and over,
* get a compressed set of warnings about issues after testing is done, and
* provide a list of known issues in the test suite.

Compare the simple example above with the current (2024-01-10) test suite for
[`samber/slog-zerolog`](https://github.com/samber/slog-zerolog):[^1]
```go
// Test_slog_samber_zerolog runs tests for the samber zerolog handler.
func Test_slog_samber_zerolog(t *testing.T) {
	sLogSuite := &test.SlogTestSuite{
		Creator: SlogSamberZerologHandlerCreator,
		Name:    "samber/slog-zerolog",
	}
	if *test.UseWarnings {
		sLogSuite.WarnOnly(test.WarnDefaultLevel)
		sLogSuite.WarnOnly(test.WarnMessageKey)
		sLogSuite.WarnOnly(test.WarnEmptyAttributes)
		sLogSuite.WarnOnly(test.WarnGroupInline)
		sLogSuite.WarnOnly(test.WarnLevelCase)
		sLogSuite.WarnOnly(test.WarnNanoDuration)
		sLogSuite.WarnOnly(test.WarnNanoTime)
		sLogSuite.WarnOnly(test.WarnNoReplAttrBasic)
		sLogSuite.WarnOnly(test.WarnResolver)
		sLogSuite.WarnOnly(test.WarnZeroPC)
		sLogSuite.WarnOnly(test.WarnZeroTime)
	}
	suite.Run(t, sLogSuite)
}
```
The various `WarnOnly()` calls configure a set of warnings that are recognized by the test suite.
Each warning is recognized by one or more tests,
which execute different code when the warning is configured.
When such a test is run and the relevant warning is set
the test executes code that "warns only" instead of running the usual Go testing assertions
(resulting in test failures).
The test suite will succeed (`PASS`) and the following result data will show at the end (2024-01-10):

```
Warnings for samber/slog-zerolog:
     2 Empty attribute(s) logged ("":null)
       TestAttributeWithEmpty
         {"level":"info","time":"2024-01-10T16:37:10-08:00","":null,"first":"one","pi":3.141592653589793,"message":"This is a message"}
       TestAttributesEmpty
         {"level":"info","time":"2024-01-10T16:37:10-08:00","":null,"first":"one","pi":3.141592653589793,"message":"This is a message"}
     1 Group with empty key does not inline subfields
       TestGroupInline
         {"level":"info","time":"2024-01-10T16:37:10-08:00","":{"fourth":"forth","second":2,"third":"3"},"first":"one","pi":3.141592653589793,"message":"This is a message"}
     2 Handler doesn't default to slog.LevelInfo
       TestDefaultLevel: defaultlevel is 'DEBUG'
       TestDefaultLevel: defaultlevel with AddSource is 'DEBUG'
     3 HandlerOptions.ReplAttr not available for basic fields
       TestReplaceAttrBasic: too many attributes: 4, time field still exists, message field still exists
       TestReplaceAttrFnLevelCase: level value not null
       TestReplaceAttrFnRemoveTime: time value not empty string
    10 Log level in lowercase
       TestCancelledContext: 'info'
       TestCancelledContext: 'info'
       TestKey: 'info'
       TestKeyCase: 'debug'
       TestKeyCase: 'info'
       TestKeyCase: 'warn'
       TestKeyCase: 'error'
       TestKeys: 'info'
       TestZeroPC: 'info'
       TestZeroTime: 'info'
     4 LogValuer objects are not resolved
       TestResolveGroup
         {"level":"info","time":"2024-01-10T16:37:10-08:00","group":{"hidden":{},"pi":3.141592653589793},"message":"This is a message"}
       TestResolveGroupWith
         {"level":"info","time":"2024-01-10T16:37:10-08:00","group":{"hidden":{},"pi":3.141592653589793},"message":"This is a message"}
       TestResolveValuer
         {"level":"info","time":"2024-01-10T16:37:10-08:00","hidden":{},"message":"This is a message"}
       TestResolveWith
         {"level":"info","time":"2024-01-10T16:37:10-08:00","hidden":{},"message":"This is a message"}
     1 SourceKey logged for zero PC
       TestZeroPC
         {"level":"info","time":"2024-01-10T16:37:10-08:00","source":{"function":"","file":"","line":0},"message":"This is a message"}
     6 Wrong message key (should be 'msg')
       TestCancelledContext: `message`
       TestCancelledContext: `message`
       TestKey: `message`
       TestKeys: `message`
       TestZeroPC: `message`
       TestZeroTime: `message`
     1 Zero time is logged
       TestZeroTime
         {"level":"info","time":"0001-01-01T00:00:00Z","message":"This is a message"}
```

### Warning Specifics

Each of the warnings is intended to represent a feature that might be expected.

#### Required

* **Empty attribute(s) logged ("":null)**  
test

#### Suggestions




[^1]: Respect to Samuel Berthe, I'm not picking on you, I just need an example here. :wink: