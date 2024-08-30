Handlers are scored based on how few warnings are generated.
Warnings are worth different amounts depending on their warning level.
The weights applied during this process are shown on the right.

#### Source Data

Each handler results in a lot of verification test output:

```
Warnings for slog/JSONHandler:
  Suggested
     2 [Duplicates] Duplicate field(s) found
         TestAttributeDuplicate: map[alpha:2 charlie:3]
         TestAttributeWithDuplicate: map[alpha:2 charlie:3]
```

The `Suggested` line is an example of a warning "level"
(warnings are grouped into levels on the [warning page](https://madkins23.github.io/go-slog/warnings.html)).
In this example there are two instances of the `Duplicates` warning.

#### Warnings Algorithm

Scoring is done for all handlers at the same time:

```
for each handler
    score starts at zero
    for each warning level
        for each warning in level
            if warning shows up during testing
                score = score + weight(level) * len(warnings)
                adjust score to range of
                    zero to maximum possible number of warnings
```

Where the `weight(level)` comes from the predefined table shown above and to the right.

The scores for each handler are then divided by the maximum possible number
of warnings that any handler might receive (if it were really awful)
and that number is subtracted from `100.0`.
This results in a number from `0.0` (awful, all warnings logged) to `100.0` (no warnings at all).
That number is stored for use and displayed on this page and on each handler page.

Note that most scores are above `~40` as it is difficult to throw _all_ the warnings.

#### Scores

Multiple scores are generated for each handler.
The _main_ (or "default") score is shown in the data tables
with the column header **Score** with an associated checkbox.
The checkbox can be used to show several other "score" columns, as follows:

* `Default` (**Score**)  
  This is the score that is shown in the overall chart
  at the top of the page in the column labeled **Warnings**.
  The default score is the same as the `By Data` score.
* `By Data`  
  This score is calculated by rolling up scores calculated per warning level.
* `Original`  
  This is the "original" score which has been overtaken by newer code.
  The `Original` score is within 5% of `by Data` value.
