Handlers are scored based on how few warnings are generated.
Warnings are worth different amounts depending on their warning level.

Scoring is done for all handlers at the same time:
```
for each handler
    score starts at zero
    for each warning level
        score = score + weight(level) * len(warnings)
```

The scores for each handler are then divided by the maximum possible number
of warnings that any handler might receive (if it were really awful)
and that number subtracted from ^100.0^.
This results in a number from ^0.0^ (awful, all warnings logged) to ^100.0^ (no warnings at all).
