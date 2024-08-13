// Package warning defines Warning objects as well as the warning Manager.
//
// # Warnings
//
// A "warning" facility is built into many of the tests to provide a way to:
//
//   - avoid scanning through a lot of `go test` error logs in detail over and over,
//   - get a compressed set of warnings about issues after testing is done, and
//   - provide a list of known issues in the test suite.
//
// # Manager
//
// The warning Manager provides all warning-related functionality and data for a test suite.
// Each test suite embeds a warning Manager for:
//
//   - Method calls to add warning records.
//   - Encapsulation and handling of warning data for a test suite.
//   - Display of warning data after completion of the test suite.
//
// # Usage
//
// Existing usage cases in bench and verify packages were implemented
// in more or less the following pattern:
//
//   - Define a `WarningManager`.
//   - Predefine various `Warning` objects that may be used in testing.
//   - When defining a test suite use `WarnOnly` to specify that
//     during testing warning code should be run instead of assertions.
//   - In test code check the `WarningManager` for applicable warnings.
//   - In warning-specific code use `AddWarning` to note the condition exists
//     or `UnusedWarning` to note that the warning is redundant.
//   - Run the tests with the `-useWarnings` flag set to true to invoke warning code.
//     The `useWarnings` flag is turned on by default.
//     Without this flag the `WarningManager` will never flag warning code
//     and test assertions will raise conventional errors.
package warning
