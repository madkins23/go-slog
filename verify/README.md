# Verifying `log/slog` Handlers

The `verify` package provides various `log/slog` (henceforth just `slog`) handler test suites.
This document discusses simple usage details.
Technical details for the test suite are provided in
the [`README.md`](tests/README.md) file in
the [`tests`](tests) package subdirectory.

There are two main benefits to using `slog`:
1. If everyone uses it then logging from a package imported into an application
   will go through whatever `slog` handler (`slog.Handler`) is configured for that application.
   Logs will all look the same and be amenable to downstream processing.
2. The ability to swap `slog` handlers without rewriting all the log statements in existing code.
   This might be done for CPU and/or memory efficiency.
   This only works if the various handlers behave in a similar manner.

The second benefit justifies a strong verification suite.

### Simple Example

Verification of a `slog` handler using the `verify` test suite is fairly simple.
The following [code](https://github.com/madkins23/go-slog/blob/main/verify/slog_test.go)
runs the test suite on `slog.JSONHandler`:

```go
package verify

import (
    "testing"

    "github.com/stretchr/testify/suite"

    "github.com/madkins23/go-slog/creator/slogjson"
    "github.com/madkins23/go-slog/verify/tests"
    "github.com/madkins23/go-slog/warning"
)

// TestVerifySlogJSON runs tests for the slog/JSONHandler JSON handler.
func TestVerifySlogJSON(t *testing.T) {
    slogSuite := tests.NewSlogTestSuite(slogjson.Creator())
    slogSuite.WarnOnly(warning.Duplicates)
    suite.Run(t, slogSuite)
}
```

The file itself must have the `_test.go` suffix and
contain a function with a name of the pattern `TestVerify<tag_name>`
where `<tag_name>` will likely be something like `PhsymZerolog` or `SlogJSON`.

The first line in `Test_slog` creates a new test suite.
The argument to the `NewSlogTestSuite` function is an [`infra.Creator`](../infra/creator.go) object,
which is responsible for creating new `slog.Logger`
(and optionally `slog.Handler`) objects during testing.
In this case an appropriate factory is created by the `creator.Slog` function
that is already defined in the `creator` package.
In order to test a new handler instance
(one that has not been tested in this repository)
it is necessary to [create a new `infra.Creator`](#creators) for it.
Existing examples can be found in the `creator` package.

Once the test suite exists the second line configures a warning to be tracked.
The meaning of `WarnOnly` is to only warn about an error condition, not fail the test.
The warning mechanism [documented below](#warnings) describes this in more detail.

Finally the suite is run via its `Run` method.
In short:
* The `TestXxxxxx` function is executed by the [Go test harness](https://pkg.go.dev/testing).
* The test function configures a `SlogTestSuite` using an `infra.Creator` factory object.
* The test function executes the test suite via its `Run` method.

### More Examples

This package contains several examples, including the one above:
* [`slog_test.go`](https://github.com/madkins23/go-slog/blob/main/verify/slog_test.go)
  Verifies the [standard `slog.JSONHandler`](https://pkg.go.dev/log/slog@master#JSONHandler).
* [`slog_phsym_zerolog_test.go`](https://github.com/madkins23/go-slog/blob/main/verify/slog_phsym_zerolog_test.go)
  Verifies the [`phsym zeroslog` handler](https://github.com/phsym/zeroslog/tree/2bf737d6422a5de048845cd3bdd2db6363555eb4).
* [`slog_samber_zap_test.go`](https://github.com/madkins23/go-slog/blob/main/verify/slog_samber_zerolog_test.go)
  Verifies the [`samber slog-zap` handler](https://github.com/samber/slog-zap).
* [`slog_samber_zerolog_test.go`](https://github.com/madkins23/go-slog/blob/main/verify/slog_samber_zerolog_test.go)
  Verifies the [`samber slog-zerolog` handler](https://github.com/samber/slog-zerolog).

In addition, there is a [`main_test.go`](https://github.com/madkins23/go-slog/blob/main/verify/main_test.go) file which exists to provide
a global resource to the other tests ([described below](#testmain)).

### Running Tests

Run the handler verification tests installed in this repository with:
```shell
go test -v ./verify -args -useWarnings
```

On an operating system that supports `bash` scripts you can use
the [`scripts/verify`](https://github.com/madkins23/go-slog/blob/main/scripts/verify) script.

**Note**:  running `go test ./... -args -useWarnings` will fail as
the other tests in the repository don't recognize the `-useWarnings` flag.

#### Test Flags

There are two flags defined for testing the verification code:
* `-debug=<level>`  
  Sets an integer level for showing any `Debugf()` statements in the code.
  As of 2024-01-11 there is only one in the test suite code at a level of `1`.
  This statement dumps the current JSON log record.[^1]
* `-useWarnings`  
  Activates the warning system (see [**Warnings**](#warnings)) below.
  Without this flag the tests fail on errors in the usual manner.
  When this flag is present tests succeed and warnings are presented
  in the `go test` output.

### Caveats

The test harness that drives verification has some limitations.

_Verification only makes sense for JSON handlers_,
which are generally used to feed log records into downstream processing.
Text and console handlers don't have a consistent format.
While it might be useful to test those handlers as well,
the difficulty of parsing various output formats argues against it.[^2]

## Creators

`Creator` objects are factories for generating new `slog.Logger` objects.


Detailed documentation on this is provided in the [`README`](../infra/README.md)
for the [`infra` package](../infra).

## Warnings

A "warning" facility is built into many of the tests to provide a way to:
* avoid scanning through a lot of `go test` error logs in detail over and over,
* get a compressed set of warnings about issues after testing is done, and
* provide a list of known issues in the test suite.

### Example

Compare the simple example [above](#simple-example) with the following excerpt from
the current (2024-01-15) test suite for
[`phsym/zeroslog`](https://github.com/phsym/zeroslog):[^3]
```go
// Test_slog_zerolog_phsym runs tests for the physym zerolog handler.
func Test_slog_zerolog_phsym(t *testing.T) {
    slogSuite := tests.NewSlogTestSuite(creator.SlogPhsymZerolog())
    slogSuite.WarnOnly(warning.Duplicates)
    slogSuite.WarnOnly(warning.EmptyAttributes)
    slogSuite.WarnOnly(warning.GroupEmpty)
    slogSuite.WarnOnly(warning.GroupInline)
    slogSuite.WarnOnly(warning.LevelCase)
    slogSuite.WarnOnly(warning.MessageKey)
    slogSuite.WarnOnly(warning.TimeMillis)
    slogSuite.WarnOnly(warning.NoReplAttr)
    slogSuite.WarnOnly(warning.SourceKey)
    slogSuite.WarnOnly(warning.ZeroTime)
    suite.Run(t, slogSuite)
}
```

The various `WarnOnly()` calls configure a set of warnings that are recognized by the test suite.
Each warning is recognized by one or more tests,
which execute different code when the warning is configured.
When such a test is run and the relevant warning is set
the test executes code that "warns only" instead of running the usual Go testing assertions
(resulting in test failures).
The test suite will succeed (`PASS`) and the following result data will show at the end (2024-01-15):

```
Warnings for phsym/zeroslog:
  Required
     2 [EmptyAttributes] Empty attribute(s) logged ("":null)
         TestAttributeWithEmpty
           {"level":"info","":null,"first":"one","pi":3.141592653589793,"time":"2024-01-21T08:57:18-08:00","message":"This is a message"}
         TestAttributesEmpty
           {"level":"info","first":"one","":null,"pi":3.141592653589793,"time":"2024-01-21T08:57:18-08:00","message":"This is a message"}
     1 [GroupEmpty] Empty (sub)group(s) logged
         TestGroupWithMultiSubEmpty
           {"level":"info","first":"one","group":{"second":2,"third":"3","subGroup":{}},"time":"2024-01-21T08:57:18-08:00","message":"This is a message"}
     1 [ZeroTime] Zero time is logged
         TestZeroTime: 0001-01-01T00:00:00Z
           {"level":"info","time":"0001-01-01T00:00:00Z","message":"This is a message"}
  Implied
     6 [MessageKey] Wrong message key (should be 'msg')
         TestCancelledContext: `message`
         TestCancelledContext: `message`
         TestKey: `message`
         TestKeys: `message`
         TestZeroPC: `message`
         TestZeroTime: `message`
     4 [NoReplAttr] HandlerOptions.ReplaceAttr not available
         TestReplaceAttr: too many attributes: 6, alpha == beta, change still exists, remove still exists
         TestReplaceAttrBasic: too many attributes: 4, time field still exists, message field still exists, source == <nil>
         TestReplaceAttrFnLevelCase: level value not null
         TestReplaceAttrFnRemoveTime: time value not empty string
     1 [SourceKey] Source data not logged when AddSource flag set
         TestKey: no 'source' key
           {"level":"info","caller":"/snap/go/10489/src/reflect/value.go:596","time":"2024-01-21T08:57:18-08:00","message":"This is a message"}
  Suggested
     3 [Duplicates] Duplicate field(s) found
         TestAttributeDuplicate: map[alpha:2 charlie:3]
         TestAttributeWithDuplicate: map[alpha:2 charlie:3]
         TestGroupInline
           {"level":"info","first":"one","":{"second":2,"third":"3","fourth":"forth"},"pi":3.141592653589793,"time":"2024-01-21T08:57:18-08:00","message":"This is a message"}
    10 [LevelCase] Log level in lowercase
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
     1 [TimeMillis] slog.Time() logs milliseconds instead of nanoseconds
         TestLogAttributes: 2024-01-21T08:57:18-08:00
```

### Warning Result Format

Each warning entry in the results has the following format:
```
  <count> <warning-text>
            <test-function-name>[: <optional-text>]
              [<optional-log-record>]
            ...<more instances>...
```

The first line has the `<count>` of times the warning was raised
followed by the `<warning-text>` of the warning.
After that there are `<count>` blocks of data,
one for each instance of the warning.

For each instance there are one or two lines.
The first line shows the `<test-function-name>` in which the warning was raised.
There may also be a colon followed by `<optional-text>`
further describing the specific instance of the warning.

The optional second line shows the `<optional-log-record>`
which is the JSON log record generated by the `slog` handler for the test.
It is often, but not always, the case that `<optional-text>`
makes an `<optional-log-record>` redundant.

### Warning Details

Each of the warnings is intended to represent a feature that is required,
might be expected, or provides administrative information.
[Details on individual warnings](https://madkins23.github.io/go-slog/warnings.html).

### Caveats

Warnings have been defined for cases that I have seen thus far for the rather
limited number of handlers for which I have configured tests.
I can think of other possible warnings, but I am loath to configure them unless they are needed.[^4]
If your handler comes up with a new error condition for which there are tests but no warning
you can either fix your handler or file a ticket.

The `-useWarnings` flag tends to result in the results being buried in the normal `go test` output.
This can be fixed by implementing a global [`TestMain()`](#testmain) function.

Warnings will only be visible when running `go test` if the `-v` flag is used.

### `TestMain`

Normally the warning results will show up in a block in the middle of `go test` output.
This is due to the way the default test harness works.

It is possible to override the default test harness by defining a global function
[`TestMain()`](https://pkg.go.dev/testing#hdr-Main).
The `verify/tests` package provides a convenient function to support this.
Define the following `TestMain()` function:
```go
func TestMain(m *testing.M) {
    tests.WithWarnings(m)
}
```

This function may be defined in the same `_test.go` file as the handler test.
If multiple handler tests are in the same directory:

* It will be necessary to move the `TestMain()` definition to a separate file,
  such as the [`verify/main_test.go`](main_test.go).
* An addition listing of which handlers throw each warning
  will be added after the normal output:

```
 Handlers by warning:
  Required
    [EmptyAttributes] Empty attribute(s) logged ("":null)
      phsym/zeroslog
      samber/slog-zap
      samber/slog-zerolog
    [GroupEmpty] Empty (sub)group(s) logged
      phsym/zeroslog
    [Resolver] LogValuer objects are not resolved
      samber/slog-zap
      samber/slog-zerolog
    [ZeroPC] SourceKey logged for zero PC
      samber/slog-zap
      samber/slog-zerolog
    [ZeroTime] Zero time is logged
      phsym/zeroslog
      samber/slog-zap
      samber/slog-zerolog
  Implied
    [DefaultLevel] Handler doesn't default to slog.LevelInfo
      samber/slog-zerolog
    [MessageKey] Wrong message key (should be 'msg')
      phsym/zeroslog
      samber/slog-zerolog
    [NoReplAttr] HandlerOptions.ReplaceAttr not available
      phsym/zeroslog
    [NoReplAttrBasic] HandlerOptions.ReplaceAttr not available for basic fields
      samber/slog-zap
      samber/slog-zerolog
    [SourceKey] Source data not logged when AddSource flag set
      phsym/zeroslog
  Suggested
    [Duplicates] slog.Duration() logs seconds instead of nanoseconds
      log/slog.JSONHandler
      phsym/zeroslog
      samber/slog-zap
      samber/slog-zerolog
    [LevelCase] Log level in lowercase
      phsym/zeroslog
      samber/slog-zap
      samber/slog-zerolog
    [TimeMillis] slog.Time() logs milliseconds instead of nanoseconds
      phsym/zeroslog
      samber/slog-zap
      samber/slog-zerolog
```
This listing isn't shown unless there are multiple test suites run in the same call to `go test`.

---

In the absence of a `TestMain` definition the warnings will not be shown.
Add the `tests.WithWarnings(m)` line to the end of the main test,
after the call to `Run` the test suite.

---

[^1]: The `--debug` flag and `Debugf` function are defined in the `test` package in this repository.

[^2]: An additional argument is that using a non-JSON output is generally only done
when writing/testing/debugging code manually and doesn't require the verification level
as JSON output which is generally done to support downstream processing of logs.

[^3]: Respect to the handler's author, I'm not picking on you, I just need an example here. :wink:

[^4]: I am lazy.
