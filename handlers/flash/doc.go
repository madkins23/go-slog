// Package flash provides a feature-complete, performant slog.Handler implementation.
//
// This handler started as a clone of the [sloggy] handler and was then tweaked for performance.
//
// # Extra Options
//
// When creating a flash handler via NewHandler
// basic options can be specified by the slog.HandlerOptions options argument.
// In addition, the handler also supports optional [flash.Extras] options
// which can be used to adjust basic logging behavior slightly.
// This can be used to test the behavior of ReplaceAttr functionality or to
// match the behavior of another logging library.
//
// # Performance Edits
//
// After flash was cloned from sloggy it went through a number of performance-related [edits].
//
// [flash.Extras]: https://pkg.go.dev/github.com/madkins23/go-slog/handlers/flash#Extras
// [sloggy]: https://pkg.go.dev/github.com/madkins23/go-slog/handlers/sloggy
// [edits]: https://github.com/madkins23/go-slog/blob/main/handlers/flash/EDITS.md
package flash
