The `Size` scoring algorithm is intended to show any difference between performance of
`slog` handlers separating "small" benchmarks from "large" ones.
The scorekeeper for this algorithm compares two sets of benchmarks according to this separation.

The `Size` score chart graphs various `slog` handlers by "small" vs. "large" tests.
The intention of this scorekeeper was to see if there were some handlers that worked
better for logging smaller items vs some that were optimized for logging larger items.
