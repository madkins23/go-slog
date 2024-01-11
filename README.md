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
* as well as some other stuff I thought would be useful.

