# `sloggy` Handler

The `sloggy` package provides a fairly straightforward implementation
of a `slog.Handler`.

The `sloggy` handler does use prefix/suffix byte arrays similar to the `zerolog` context.
This doesn't help for simple loggers, but when using `WithGroup` or `WithAttrs`
these fields can store pre-formatted JSON fragments that will later speed up
eventual log statements.

This was an initial, naive attempt to write a "better" handler,
because [*hubris*](https://wiki.c2.com/?LazinessImpatienceHubris).

In one sense this was a success:
this is the second "feature complete" handler after `slog.JSONHandler`.[^1]
A user could switch between the two handlers and be reasonably confident that
the log output would be the same.

Nor is it terribly slow, but it isn't in the top ranks in terms of performance.
Given that, `slog.JSONHandler` is the better choice. :cry:

## Example

```go
logger := slog.New(flash.NewHandler(os.Stdout, nil))
logger.Info("hello", "count", 3)
```

[^1]: Feature complete in this case being according to the verification test suite
defined in `go-slog/verify` and the tests and warnings defined therein.
