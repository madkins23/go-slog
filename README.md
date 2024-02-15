# go-slog
Tools and testing for `log/slog` (hereafter just `slog`) handlers.

## Contents

* Test harness for verifying `slog` handlers
* Functions for use with `slog.HandlerOptions.ReplaceAttr`
* Test harness for benchmarking `slog` handlers
* Utility to redirect internal `gin` logging to `slog`

## Handler Verification

Verification of `slog` handlers can be done by creating
simple `_test.go` files that utilize the `verify/test.SlogTestSuite`
located in this repository.
Usage details for this facility are provided in
the [`README`](verify/README.md) file
located in the [`verify`](verify) package directory.

The tests implemented herein were inspired by:
* the [`slogtest`](https://pkg.go.dev/golang.org/x/exp/slog/slogtest) application,
* rules specified in
  the [`log/slog.Hander`](https://pkg.go.dev/log/slog@master#Handler) and
  [handler writing guide](https://github.com/golang/example/tree/master/slog-handler-guide)
  documentation,
* issues I noticed while exploring
  [`go-logging-benchmark`](https://github.com/betterstack-community/go-logging-benchmarks)
* as well as some other test cases that seemed useful.

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
  (time, level, message, and source).

## Handler Benchmarks

Benchmarks of `slog` handlers can be done by creating
simple `_test.go` files that utilize the `bench/test.SlogBenchmarkSuite`
located in this repository.
Use details for this facility are provided in
the [`README`](bench/README.md) file
located in the [`bench`](bench) package directory.

Benchmarks are intended to compare multiple handlers.
The benchmark data generated can be processed by two applications:
* [`tabular`](cmd/tabular/tabular.go)  
  generates a set of tables, each of which compares handlers for a given benchmark test.
* [`server`](cmd/server/server.go)  
  runs a simple web server showing the same tables plus bar charts.

[Recent benchmark data](https://madkins23.github.io/go-slog/index.html).

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
