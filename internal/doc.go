// Package internal contains various private packages.
//
// # Data
//
// Object definitions and parsing functionality for acquiring
// data from the output of benchmark and verification tests
// for use by [display commands].
//
// # JSON
//
// JSON functionality used by various tests.
//
// # Language
//
// Provide language-appropriate formatting of numbers.
//
// # Markdown
//
// Utility for parsing Markdown format and converting it into HTML.
//
// # Utility
//
// Common utilities:
//
//   - Acquire current function name off of stack.
//
// # Scoring
//
// Data structures for handler "scores" and functionality for
// converting benchmark and verification data into score data.
//
// # Test
//
// Utilities for constructing test files.
//
// # Support for go:generate tools
//
// Keep 'go tidy' from removing go:generate modules from go.mod file.
//
// [display commands]: https://pkg.go.dev/github.com/madkins23/go-slog/cmd
package internal
