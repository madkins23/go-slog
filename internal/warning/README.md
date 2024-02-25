# `test`

This package contains testing support code.

## Warnings

A "warning" facility is built into many of the tests to provide a way to:
* avoid scanning through a lot of `go test` error logs in detail over and over,
* get a compressed set of warnings about issues after testing is done, and
* provide a list of known issues in the test suite.

Since this code is currently used primary in handler verification,
there are better examples in the [`verify` `README`](../../verify/README.md) file.

The actual warnings are defined in the [`warning`]() package.

### Usage

* Define a `WarningManager`.
* Predefine various `Warning` objects that may be used in testing.
* When defining a test suite use `WarnOnly` to specify that
  during testing warning code should be run instead of assertions.
* In test code check the `WarningManager` for applicable warnings.
* In warning-specific code use `AddWarning` to note the condition exists
  or `UnusedWarning` to note that the warning is redundant.
* Run the tests with the `-useWarnings` flag to invoke warning code.
  Without this flag the `WarningManager` will never flag warning code
  and test assertions will raise conventional errors.

Usage of the `WarningManager` is described in more detail
in the [`verify` `README`](../../verify/README.md).
