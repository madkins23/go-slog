Handlers are scored based on how few warnings are generated.
Warnings are worth different amounts depending on their warning level.
The weights applied during this process are shown on the right.

Scoring is done for all handlers at the same time:

```
for each handler
    score starts at zero
    for each warning level
        score = score + weight(level) * len(warnings)
        adjust score to range of zero to maximum possible number of warnings
```

Where the `weight(level)` comes from the predefined table shown above and to the right.

The scores for each handler are then divided by the maximum possible number
of warnings that any handler might receive (if it were really awful)
and that number is subtracted from `100.0`.
This results in a number from `0.0` (awful, all warnings logged) to `100.0` (no warnings at all).
That number is stored for use and displayed on this page and on each handler page.

Note that most scores are above 50 as it is difficult to throw _all_ the warnings.
