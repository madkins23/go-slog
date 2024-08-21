The `Size` scoring algorithm is basically the same as the `Default` algorithm.
Only the handlers shown in the chart is changed.

The intention of this scorekeeper was to see if there were some handlers that worked
better for logging smaller items vs some that were optimized for logging larger items.

The `Size` score chart graphs various `slog` handlers by "small" vs. "large" tests.
Each handler is scored for a set of small tests and separately for a set of large tests.
Handlers with a preference for one or the other will show up on the chart
as off of the main diagonal (`0,0` → `100,100`).