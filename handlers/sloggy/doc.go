// Package sloggy provides a fairly straightforward slog.Handler implementation.
// It is feature-complete per observed slog.JSONHandler behavior.
//
// This was an initial, naive attempt to write a "better" handler (because [hubris]).
// It is the second "feature complete" handler after slog.JSONHandler
// (which is admittedly used as a "default behavior" model throughout the verification suite).
// A user could switch between the two handlers and be reasonably confident that
// the log output would be the same.
//
// Performance not as good as slog.JSONHandler.
// The only planned performance enhancement was to use prefix/postfix byte arrays,
// otherwise this was a green field build with performance left until later.
// The [flash] handler, originally a copy of this one, has been tweaked for performance.
//
// [flash]: https://pkg.go.dev/github.com/madkins23/go-slog/handlers/flash
// [hubris]: https://wiki.c2.com/?LazinessImpatienceHubris
package sloggy
