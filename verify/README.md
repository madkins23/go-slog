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

## Simple Example

Verification of a `slog` handler using the `verify` test suite is fairly simple.
The following [code](https://github.com/madkins23/go-slog/blob/main/verify/slog_test.go)
runs the test suite on `slog.JSONHandler`:

```go
package verify

import (
  "testing"

    "github.com/stretchr/testify/suite"

    "github.com/madkins23/go-slog/creator"
    "github.com/madkins23/go-slog/infra"
    "github.com/madkins23/go-slog/verify/tests"
)

// Test_slog runs tests for the log/slog JSON handler.
func Test_slog(t *testing.T) {
    slogSuite := tests.NewSlogTestSuite(creator.Slog())
    slogSuite.WarnOnly(infra.WarnDuplicates)
    suite.Run(t, slogSuite)
}
```
The file itself must have the `_test.go` suffix and
contain a function with a name beginning with 'Test' in order to be executed as a test.
In this instance `slog_test.go` contains function `Test_slog`.

The first line in `Test_slog` creates a new test suite.
The argument to the `NewSlogTestSuite` function is an `infra.Creator` object,
which is a object for creating new `slog.Handler` objects during testing.
In this case a factory for creating `slog.JSONHandler` objects
is created by the `creator.Slog` function that is already defined.
In order to test a new `slog.Handler` (one that has not been tested in this repository)
it is necessary to create a new `infra.Creator` for it.
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
  Verifies the standard [`slog.JSONHandler`](https://pkg.go.dev/log/slog@master#JSONHandler).
* [`slog_zerolog_phsym_test.go`](https://github.com/madkins23/go-slog/blob/main/verify/slog_phsym_zerolog_test.go)
  Verifies the [`zeroslog` handler](https://github.com/phsym/zeroslog/tree/2bf737d6422a5de048845cd3bdd2db6363555eb4).
* [`slog_zerolog_samber_test.go`](https://github.com/madkins23/go-slog/blob/main/verify/slog_samber_zerolog_test.go)
  Verifies the [`slog-zerolog` handler](https://github.com/samber/slog-zerolog).

In addition, there is a [`main_test.go`](https://github.com/madkins23/go-slog/blob/main/verify/main_test.go) file which exists to provide
a global resource to the other tests ([described below](#testmain)).

## Running Tests

Run the handler verification tests installed in this repository with:
```shell
go test -v ./verify -args -useWarnings
```

On an operating system that supports `bash` scripts you can use
the [`scripts/verify`](https://github.com/madkins23/go-slog/blob/main/scripts/verify) script.

**Note**:  running `go test ./... -args -useWarnings` will fail as
the other tests in the repository don't recognize the `-useWarnings` flag.

### Test Flags

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

## Caveats

The test harness that drives verification has some limitations.

_Verification only makes sense for JSON handlers_,
which are generally used to feed log records into downstream processing.
Text and console handlers don't have a consistent format.
While it might be useful to test those handlers as well,
the difficulty of parsing various output formats argues against it.[^2]

_The test suite interface to `slog` logging is via `slog.Handler`_.
This is a well-defined interface and the obvious place to swap in
a different logging backend.
This makes loggers that directly implement `slog`-like behavior
without instantiating a `slog.Handler`
(e.g. the [darvaza loggers](https://github.com/darvaza-proxy/slog))
inappropriate for this test suite.

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
    slogSuite.WarnOnly(tests.WarnDuplicates)
    slogSuite.WarnOnly(tests.WarnDurationMillis)
    slogSuite.WarnOnly(tests.WarnEmptyAttributes)
    slogSuite.WarnOnly(tests.WarnGroupEmpty)
    slogSuite.WarnOnly(tests.WarnGroupInline)
    slogSuite.WarnOnly(tests.WarnLevelCase)
    slogSuite.WarnOnly(tests.WarnMessageKey)
    slogSuite.WarnOnly(tests.WarnTimeMillis)
    slogSuite.WarnOnly(tests.WarnNoReplAttr)
    slogSuite.WarnOnly(tests.WarnSourceKey)
    slogSuite.WarnOnly(tests.WarnZeroTime)
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
     1 Empty (sub)group(s) logged
       TestGroupWithMultiSubEmpty
         {"level":"info","first":"one","group":{"second":2,"third":"3","subGroup":{}},"time":"2024-01-15T17:36:33-08:00","message":"This is a message"}
     2 Empty attribute(s) logged ("":null)
       TestAttributeWithEmpty
         {"level":"info","":null,"first":"one","pi":3.141592653589793,"time":"2024-01-15T17:36:33-08:00","message":"This is a message"}
       TestAttributesEmpty
         {"level":"info","first":"one","":null,"pi":3.141592653589793,"time":"2024-01-15T17:36:33-08:00","message":"This is a message"}
     1 Group with empty key does not inline subfields
       TestGroupInline
         {"level":"info","first":"one","":{"second":2,"third":"3","fourth":"forth"},"pi":3.141592653589793,"time":"2024-01-15T17:36:33-08:00","message":"This is a message"}
     1 Zero time is logged
       TestZeroTime: 0001-01-01T00:00:00Z
         {"level":"info","time":"0001-01-01T00:00:00Z","message":"This is a message"}
  Implied
     4 HandlerOptions.ReplaceAttr not available
       TestReplaceAttr: too many attributes: 6, alpha == beta, change still exists, remove still exists
       TestReplaceAttrBasic: too many attributes: 4, time field still exists, message field still exists, source == <nil>
       TestReplaceAttrFnLevelCase: level value not null
       TestReplaceAttrFnRemoveTime: time value not empty string
     1 Source data not logged when AddSource flag set
       TestKey: no 'source' key
         {"level":"info","caller":"/snap/go/current/src/reflect/value.go:596","time":"2024-01-15T17:36:33-08:00","message":"This is a message"}
     6 Wrong message key (should be 'msg')
       TestCancelledContext: `message`
       TestCancelledContext: `message`
       TestKey: `message`
       TestKeys: `message`
       TestZeroPC: `message`
       TestZeroTime: `message`
  Suggested
     2 Duplicate field(s) found
       TestAttributeDuplicate: map[alpha:2 charlie:3]
       TestAttributeWithDuplicate: map[alpha:2 charlie:3]
    10 Log level in lowercase
       TestCancelledContext: 'info'
       TestCancelledContext: 'info'
       TestKey: 'info'
       TestKeyCase: 'warn'
       TestKeyCase: 'error'
       TestKeyCase: 'debug'
       TestKeyCase: 'info'
       TestKeys: 'info'
       TestZeroPC: 'info'
       TestZeroTime: 'info'
     1 slog.Duration() logs milliseconds instead of nanoseconds
       TestLogAttributes
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

### Warning Specifics

Each of the warnings is intended to represent a feature that is required,
might be expected, or provides administrative information.
For each warning described below:

* The first line shows the name of the warning (for calls to `WarnOnly`)
  followed by the warning description which would show up in the warnings section during testing.
* Documentation of the warning follows as required.
* Optional bullet items with references that support the test/warning.
  Generally these will be linked to supporting documentation of some sort.

#### Required

The following warnings relate to tests that I can justify from requirements in the
[`slog.Handler`](https://pkg.go.dev/log/slog@master#Handler) documentation.

* `WarnZeroTime`: 'Zero time is logged'  
  Handlers should not log the basic `time` field if it is zero.  
  * ['- If r.Time is the zero time, ignore the time.'](https://pkg.go.dev/log/slog@master#Handler)
* `WarnZeroPC`: 'SourceKey logged for zero PC'  
  The `slog.Record.PC` field can be loaded with a program counter (PC).
  This is normally done by the `slog.Logger` code.
  If the PC is non-zero and the `slog.HandlerOptions.AddSource` flag is set
  the `source` field will contain a [`slog.Source`](https://pkg.go.dev/log/slog@master#Source) record
  containing the function name, file name, and file line at which the log record was generated.
  If the PC is zero then this field and its associated group should not be logged.  
  * ['- If r.PC is zero, ignore it.'](https://pkg.go.dev/log/slog@master#Handler)
* `WarnResolver`: 'LogValuer objects are not resolved'  
  Handlers should resolve all objects implementing the
  [`LogValuer`](https://pkg.go.dev/log/slog@master#LogValuer) or
  [`Stringer`](https://pkg.go.dev/fmt#Stringer) interfaces.  
  This is a powerful feature which can customize logging of objects and
  [speed up logging by delaying argument resolution until logging time](https://pkg.go.dev/log/slog@master#hdr-Performance_considerations).  
  * ['- Attr's values should be resolved.'](https://pkg.go.dev/log/slog@master#Handler)
* `WarnEmptyAttributes`: 'Empty attribute(s) logged "":null'  
  Handlers are supposed to avoid logging empty attributes.  
  * ['- If an Attr's key and value are both the zero value, ignore the Attr.'](https://pkg.go.dev/log/slog@master#Handler)
* `WarnGroupInline`: 'Group with empty key does not inline subfields'  
  Handlers should expand groups named "" (the empty string) into the enclosing log record.  
  * ['- If a group's key is empty, inline the group's Attrs.'](https://pkg.go.dev/log/slog@master#Handler)
* `WarnGroupEmpty`: 'Empty (sub)group(s) logged'  
  Handlers should not log groups (or subgroups) without fields,
  whether or not the have non-empty names.  
  * ['- If a group has no Attrs (even if it has a non-empty key), ignore it.'](https://pkg.go.dev/log/slog@master#Handler)

#### Implied

Warnings that seem to be implied by documentation but can't be considered required.

* `WarnDefaultLevel`: 'Handler doesn't default to slog.LevelInfo'  
  A new `slog.Handler` should default to `slog.LevelInfo`.  
  * ['First, we wanted the default level to be Info, Since Levels are ints, Info is the default value for int, zero.'](https://pkg.go.dev/log/slog@master#Level)
* `WarnMessageKey`: 'Wrong message key (should be 'msg')'  
  The field name of the "message" key should be `msg`.  
  * [Constant values are defined for `slog/log`](https://pkg.go.dev/log/slog@master#pkg-constants)  
  * [Field values are defined for the `JSONHandler.Handle()` implementation](https://pkg.go.dev/log/slog@master#JSONHandler.Handle)
* `WarnSourceKey`: 'Source data not logged when AddSource flag set'  
  Handlers should log source data when the `slog.HandlerOptions.AddSource` flag is set.  
  * [Flag declaration as `slog.HandlerOptions` field](https://pkg.go.dev/log/slog@master#HandlerOptions)  
  * [Behavior defined for `JSONHandler.Handle()`](https://pkg.go.dev/log/slog@master#JSONHandler.Handle)  
  * [Definition of source data record](https://pkg.go.dev/log/slog@master#Source)
* `WarnNoReplAttr`: 'HandlerOptions.ReplaceAttr not available'  
  If `HandlerOptions.ReplaceAttr` is provided it should be honored by the handler.
  However, documentation on implementing handler methods seems to suggest it is optional.  
  * [Behavior defined for `slog.HandlerOptions`](https://pkg.go.dev/log/slog@master#HandlerOptions)  
  * ['You might also consider adding a ReplaceAttr option to your handler, like the one for the built-in handlers.'](https://github.com/golang/example/tree/master/slog-handler-guide#implementing-handler-methods)
* `WarnNoReplAttrBasic`: 'HandlerOptions.ReplaceAttr not available for basic field'  
  Some handlers (e.g. `phsym/zeroslog`) support `HandlerOptions.ReplaceAttr`
  except for the four main fields `time`, `level`, `msg`, and `source`.
  When that is the case it is better to use this (`WarnNoReplAttrBasic`) warning.

#### Suggested

These warnings are not AFAIK mandated by any documentation or requirements.[^4]

* `WarnDuplicates`: 'Duplicate field(s) found'  
  Some handlers (e.g. `slog.JSONHandler`)
  will output multiple occurrences of the same field name
  if the logger is called with multiple instances of the same field.
  This behavior is currently [under debate](https://github.com/golang/go/issues/59365)
  with no resolution at this time (2024-01-15) and a
  [release milestone of (currently unscheduled) Go 1.23](https://github.com/golang/go/milestone/212),
  (whereas [Go Release 1.22](https://tip.golang.org/doc/go1.22)
  is currently expected in February 2024).
* `WarnDurationSeconds`: 'slog.Duration() logs seconds instead of nanoseconds'
  The `slog.JSONHandler` uses nanoseconds for `time.Duration` but some other handlers use seconds.
  * [Go issue 59345: Nanoseconds is a recent change with Go 1.21](https://github.com/golang/go/issues/59345)
* `WarnDurationMillis`: 'slog.Duration() logs milliseconds instead of nanoseconds'  
  The `slog.JSONHandler` uses nanoseconds for `time.Duration` but some other handlers use seconds.
  * [Go issue 59345: Nanoseconds is a recent change with Go 1.21](https://github.com/golang/go/issues/59345)
* `WarnLevelCaseLog`: 'level in lowercase'  
  Each JSON log record contains the logging level of the log statement as a string.
  Different handlers provide that string in uppercase or lowercase.
  Documentation for [`slog.Level`](https://pkg.go.dev/log/slog@master#Level)
  says that its `String()` and `MarshalJSON()` methods will return uppercase
  but `UnmarshalJSON()` will parse in a case-insensitive manner.
* `WarnTimeMillis`: 'slog.Time() logs milliseconds instead of nanoseconds'  
  The `slog.JSONHandler` uses nanoseconds for `time.Time` but some other handlers use seconds.
  This does _not_ apply to the basic `time` field, only attribute fields.
  I can't find any supporting documentation or bug on this but
  [Go issue 59345](https://github.com/golang/go/issues/59345) (see previous warning)
  may have fixed this as well in Go 1.21.

#### Administrative

The last warnings provide information about the tests or conflicts with other warnings.

* `WarnSkippingTest`: 'Skipping test'  
  A test has been skipped, likely due to the specification of some other warning.
  Not currently used (2024-01-15).
* `WarnUndefined`: 'Undefined Warning(s)'  
  An attempt to call WarnOnly() with an undefined warning.
  Warnings must be predefined to the `WarningManager` prior to use.
* `WarnUnused`: 'Unused Warning(s)'  
  If a warning is specified but the condition is not actually present
  one of these warnings will be issued with the specified warning.
  These are intended to help clean out unnecessary `WarnOnly` settings
  from a test suite as issues are fixed in the tested handler.

### Caveats

Warnings have been defined for cases that I have seen thus far for the rather
limited number of handlers for which I have configured tests.
I can think of other possible warnings, but I am loath to configure them unless they are needed.[^5]
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
If multiple handler tests are in the same directory it will be necessary to
move the `TestMain()` definition to a separate file,
such as the [`verify/main_test.go`](main_test.go).

[^1]: The `--debug` flag and `Debugf` function are defined in the `test` package in this repository.

[^2]: An additional argument is that using a non-JSON output is generally only done
when writing/testing/debugging code manually and doesn't require the verification level
as JSON output which is generally done to support downstream processing of logs.

[^3]: Respect to the handler's author, I'm not picking on you, I just need an example here. :wink:

[^4]: I favor more rigorous guidelines and handlers that require fewer warnings.
Most JSON logging will be done to feed downstream log consumers.
The looser the guidelines the greater the chance that swapping `slog` handlers
will necessitate changes to downstream processes.

[^5]: I am lazy.
