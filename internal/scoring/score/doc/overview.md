The algorithms behind the scores shown on this page are somewhat arbitrary.
The original scoring algorithm (`Default`) was deemed "good enough",
but later work has focused on enabling multiple scoring algorithms.
These can be found on the [Home page](/go-slog/) or
in the `Scoring` drop-down in the upper right section of every page.

Algorithms are implemented by ["scorekeepers"](https://pkg.go.dev/github.com/madkins23/go-slog/internal/scoring/keeper).
Each scorekeeper is specified by the two axes shown in the scoring chart.
Each axis interprets test data according to its own algorithm.

The current scorekeeper and axis algorithms are described below: