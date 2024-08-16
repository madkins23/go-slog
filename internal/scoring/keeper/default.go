package keeper

import (
	_ "embed"

	"github.com/madkins23/go-slog/infra/warning"
	"github.com/madkins23/go-slog/internal/markdown"
	"github.com/madkins23/go-slog/internal/scoring/axis"
	"github.com/madkins23/go-slog/internal/scoring/filter"
	"github.com/madkins23/go-slog/internal/scoring/score"
)

const DefaultName = "Default"

var (
	//go:embed doc/default-doc.md
	defaultDocMD string

	//go:embed doc/default-sum-x.md
	defaultXSumMD string

	//go:embed doc/default-sum-y.md
	defaultYSumMD string
)

var defaultOptions = &score.KeeperOptions{
	Title: "Speed vs. Functionality",
	ChartCaption: `
		Higher numbers are better on both axes. The "good" zone is the upper right and the "bad" zone is the lower left.<br/>
		The top is fast, the bottom is slow. Left is more warnings, right is less.`,
}

func setupDefault() error {
	return score.AddKeeper(
		score.NewKeeper(
			DefaultName,
			axis.NewWarnings(
				defaultWarningScoreWeight,
				markdown.TemplateHTML(defaultXSumMD, false)),
			axis.NewBenchmarks(
				defaultBenchmarkScoreWeight,
				markdown.TemplateHTML(defaultYSumMD, false),
				nil),
			markdown.TemplateHTML(defaultDocMD, false),
			defaultOptions,
			filter.Basic()))
}

// -----------------------------------------------------------------------------

// defaultWarningScoreWeight has the multipliers for different warning levels.
var defaultWarningScoreWeight = map[warning.Level]uint{
	warning.LevelRequired:  8,
	warning.LevelImplied:   4,
	warning.LevelSuggested: 2,
	warning.LevelAdmin:     1,
}

// benchScoreWeight has the multipliers for different benchmark values.
var defaultBenchmarkScoreWeight = map[axis.BenchValue]uint{
	axis.Allocations: 1,
	axis.AllocBytes:  2,
	axis.Nanoseconds: 3,
}
