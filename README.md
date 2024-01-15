# go-slog
Tools and testing for `log/slog` handlers

## Contents

* Test harness for verifying `log/slog` handlers

## Verification

Verification of `log/slog` handlers can be done by creating
simple `_test.go` files that utilize the `verify/test.SlogTestSuite`
located in this repository.
Usage details for this facility are provided in
the [`Readme.MD`](verify/README.md) file
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

A small collection of functions in the `replace` package
can be used with `slog.HandlerOptions.ReplaceAttr`.

## Gin Logging Redirect

Package `gin` contains utilities for using `log/slog` with `gin-gonic/gin`.
In particular, this package provides `gin.Writer` which can be used to redirect Gin-internal logging:
```go
gin.DefaultWriter = gin.NewWriter(slog.LevelInfo)
gin.DefaultErrorWriter = gin.NewWriter(slog.LevelError)
```
Start this before starting Gin and all of the Gin-internal logging
should be redirected to the new `io.Writer` objects.
These objects will parse the Gin-internal logging formats and
use `log/slog` to do the actual logging so it will look the same.
