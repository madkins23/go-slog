# go-slog
Tools and testing for `log/slog` (hereafter just `slog`) handlers.

[Recent benchmark data](https://madkins23.github.io/go-slog/index.html)
is available via [GitHub Pages](https://pages.github.com/).

See the [source](https://github.com/madkins23/go-slog)
or [documentation](https://pkg.go.dev/github.com/madkins23/go-slog)
for more detailed documentation.

[![Go Report Card](https://goreportcard.com/badge/github.com/madkins23/go-slog)](https://goreportcard.com/report/github.com/madkins23/go-slog)
![GitHub](https://img.shields.io/github/license/madkins23/go-slog)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/madkins23/go-slog)
[![Go Reference](https://pkg.go.dev/badge/github.com/madkins23/go-slog.svg)](https://pkg.go.dev/github.com/madkins23/go-slog)

## Repository Contents

* [Test harness for benchmarking `slog` handler performance](#handler-benchmarks)
* [Test harness for verifying `slog` handler functionality](#handler-verification)
* Web Server that processes output of previous two steps and provides access to result.
* Static copies of server pages provided via [GitHub Pages](https://pages.github.com/).
* [Functions for use with `slog.HandlerOptions.ReplaceAttr`](#replace-attributes-functions)
* [Utility to redirect internal `gin` logging to `slog`](#gin-logging-redirect)

## Handler Benchmarks

Benchmarks of `slog` handlers can be done by creating
simple `_test.go` files that utilize the `bench/test.SlogBenchmarkSuite`
located in this repository.
Usage details for this facility are provided in
the [`README`](bench/README.md) file located in the `bench` package directory.

Benchmarks are intended to compare multiple handlers.
The benchmark data generated can be processed by two applications:
* [`tabular`](cmd/tabular/tabular.go)  
  generates a set of tables, each of which compares handlers for a given benchmark test.
* [`server`](cmd/server/server.go)  
  runs a simple web server showing the same tables plus bar charts and
  handler verification warnings.

## Handler Verification

Verification of `slog` handlers can be done by creating
simple `_test.go` files that utilize the `verify/test.SlogTestSuite`
located in this repository.
Usage details for this facility are provided in
the [`README`](verify/README.md) file
located in the [`verify`](https://pkg.go.dev/github.com/madkins23/go-slog@v0.7.1-alpha-gin/verify) package directory.

The tests implemented herein were inspired by:
* the [`slogtest`](https://pkg.go.dev/golang.org/x/exp/slog/slogtest) application,
* rules specified in
  the [`log/slog.Hander`](https://pkg.go.dev/log/slog@master#Handler) and
  [handler writing guide](https://github.com/golang/example/tree/master/slog-handler-guide)
  documentation,
* issues I noticed while exploring
  [`go-logging-benchmark`](https://github.com/betterstack-community/go-logging-benchmarks)
* as well as some other test cases that seemed useful.

## Web Server

:construction: **TBD** :construction:

### Generating GitHub Pages

:construction: **TBD** :construction:

[Recent benchmark data](https://madkins23.github.io/go-slog/index.html).
This set of web pages is generated on a weekly basis using github actions and pages.
The pages are generate using the `server` command defined herein and
then scraped into the repository from which the [GitHub Pages](https://pages.github.com/) are served.

## Replace Attributes Functions

A small collection of functions in the [`replace`](replace) package
can be used with `slog.HandlerOptions.ReplaceAttr`.

These were intended to "fix" some of the verification issues with various handlers.
Unfortunately, other issues prevent them from being fixed:
* Attributes can't be directly removed, they can only be made empty,
  but some handlers tested don't remove empty attributes as they should
  so this fix doesn't work for them.
* Some handlers don't recognize `slog.HandlerOptions.ReplaceAttr`.
* Those that do don't always recognize them for the basic fields
  (`time`, `level`, `message`, and `source`).

## Gin Logging Redirect

Package `gin` contains utilities for using `slog` with `gin-gonic/gin`.
In particular, this package provides `gin.Writer` which can be used to redirect Gin-internal logging:
```go
import (
    "github.com/gin-gonic/gin"
    ginslog "github.com/madkins23/go-slog/gin"
)

gin.DefaultWriter = ginslog.NewWriter(&ginslog.Options{})
gin.DefaultErrorWriter = ginslog.NewWriter(&ginslog.Options{Level: slog.LevelError})
```
Configure this before starting Gin and all the Gin-internal logging
should be redirected to the new `io.Writer` objects.
These objects will parse the Gin-internal logging formats and
use `slog` to do the actual logging, so the log output will all look the same.

The `gin.Writer` objects can further parse the "standard" Gin traffic lines containing:
```
200 |  5.529751605s |             ::1 | GET      "/chart.svg?tag=With_Attrs_Attributes&item=MemBytes"
```
To embed the traffic data at the top level of the log messages:
```go
gin.DefaultWriter = ginslog.NewWriter(&ginslog.Options{
	Traffic: ginslog.Traffic{Parse: true, Embed: true},
})
```
To aggregate the traffic data into a group named by `ginslog.DefaultTrafficGroup`:
```go
gin.DefaultWriter = ginslog.NewWriter(&ginslog.Options{
	Traffic: ginslog.Traffic{Parse: true},
})
```
To aggregate the traffic data into a group named `"bob"`:
```go
gin.DefaultWriter = ginslog.NewWriter(&ginslog.Options{
	Traffic: ginslog.Traffic{Parse: true, Group: "bob"},
})
```
Further options can be found in the code documentation of `go-slog/gin.Options`.

## Caveats

### JSON Only

The tests in this repository only apply to `slog` JSON output.
Console output can come in a variety of formats and
generally doesn't have a performance issue as only humans will look at it.

### Who Am I?

Your response to this repository, especially if you are a `slog` handler author,
may well be "who are you to make these rules?"
This is a reasonable question.

The author has no authority to dictate `slog` handler behavior.
The tests and warnings contained herein more or less defined themselves
based on the behavior of various `slog` handlers under test.

Benchmark and verification tests come from `slog` documentation,
the `slogtest` test harness, and tests embedded in
[betterstack-community/go-logging-benchmarks](https://github.com/betterstack-community/go-logging-benchmarks).
Warnings generated by these tests pretty much defined themselves (with the author's help, of course).

Each test and/or warning is based on some sort of justification.
This justification is reflected in documentation and comments throughout the code.
It should be possible to follow this trail of bread crumbs to justify each test or warning.
Whether the reader agrees with this justification is subjective.

The several levels of verification tests are defined based on the strength of justification:

* **Required**  
  Justified from requirements in the `slog.Handler` documentation.
* **Implied**  
  Implied by documentation but can't be considered required.
* **Suggested**  
  Not mandated by any documentation or requirements.
  These are the ones that the author just made up because they seemed appropriate.
* **Administrative**  
  Information about the tests or conflicts with other warnings.

#### As a `slog` Author

* You don't have to pay attention to any of this. Really.
* If you _do_ pay attention, you should probably work down from `Required` warnings.
* Consider your users' viewpoint...

#### As a `slog` User

* Are you confident that the `slog` handler you choose will be fine forever?
* Consider the trade-off between performance and functionality.  
  - Do you need faster logging without support for picky warnings?
  - Do you need full support of all verification tests at the cost of performance?
* If you need to swap out `slog` handlers will the new one support your usage of the old one?
  - Should your code use only generally supported features even if more useful ones are available?
  - Are you prepared to change logging statements to use a less functional handler?

#### IMHO

The author considers the interoperability of logging via `slog` to be very important,
and possibly the best aspect of `slog` logging.
Editing every log statement in a large project can be a real pain.

## Links

### Slog Documentation

* [Documentation of `log/slog`](https://pkg.go.dev/log/slog@master)
* [Guide to Implementation of `slog` Handlers](https://github.com/golang/example/tree/master/slog-handler-guide)
* [Test Harness `slogtest`](https://pkg.go.dev/golang.org/x/exp/slog/slogtest)

### Slog Handlers

The following handlers are currently under test in this repository:

* [`chanchal/zaphandler`](https://github.com/chanchal1987/zaphandler)
* [`log/slog.JSONHandler`](https://pkg.go.dev/log/slog#JSONHandler)
* [`phsym/zeroslog`](https://github.com/phsym/zeroslog)
* [`samber/slog-zap`](https://github.com/samber/slog-zap)
* [`samber/slog-zerolog`](https://github.com/samber/slog-zerolog)
* [`samber/slog-logrus`](https://github.com/samber/slog-logrus)

Handlers that have been investigated and found wanting:

* `darvaza` handlers are based on a different definition of `log/slog`
  as an interface that is not compatible with the "real" `log/slog/Logger`.
  Since the latter is _not_ an interface there is no way to build a shim.
  In addition, there is no separate 'Handler' object.
  * [`darvaza/logrus`](https://pkg.go.dev/darvaza.org/slog/handlers/logrus)
  * [`darvaza/zap`](https://pkg.go.dev/darvaza.org/slog/handlers/zap)
  * [`darvaza/zerolog`](https://pkg.go.dev/darvaza.org/slog/handlers/zerolog)

* Some handlers are still using `golang.org/x/exp/slog`
  instead of the standard library `log/slog`:
  * [`galecore/xslog`](https://github.com/galecore/xslog)
  * [`evanphx/go-hclog-slog`](https://github.com/evanphx/go-hclog-slog)

Console handlers are not tested in this repository, but the author likes this one:

* [`phsym/console-slog`](https://github.com/phsym/console-slog)

### Miscellaneous

* [Awesome `slog`](https://github.com/go-slog/awesome-slog)
* [Go Logging Benchmarks](https://github.com/betterstack-community/go-logging-benchmarks)
  - Benchmarks of various Go logging packages, not just `slog` loggers.
  - Used github action in this project as template for
    generating github [pages](https://pages.github.com/).
