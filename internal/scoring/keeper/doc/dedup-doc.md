The `Dedup` scoring algorithm is basically the same as the `Default` algorithm.
Only the handlers shown in the chart is changed.

The chart compares the behavior of the `slog/JSONHandler` handler
to available wrappers that de-duplicate fields in the result.
As of 2024-08-21 both such handlers are wrappers around preexisting
`log/slog` handlers, in this case `slog/JSONHandler`.
