package keeper

import (
	_ "embed"

	"github.com/madkins23/go-slog/infra/warning"
	"github.com/madkins23/go-slog/internal/markdown"
	"github.com/madkins23/go-slog/internal/scoring/axis"
	"github.com/madkins23/go-slog/internal/scoring/score"
)

const simpleName = "Simple"

var (
	//go:embed doc/simple-doc.md
	simpleDocMD string

	//go:embed doc/simple-sum-x.md
	simpleXSumMD string

	//go:embed doc/simple-sum-y.md
	simpleYSumMD string
)

func setupSimple() error {
	return score.AddKeeper(
		score.NewKeeper(
			simpleName,
			axis.NewWarnings(
				simpleWarningScoreWeight,
				markdown.TemplateHTML(simpleXSumMD, false)),
			axis.NewBenchmarks(
				simpleBenchmarkScoreWeight,
				markdown.TemplateHTML(simpleYSumMD, false),
				nil),
			markdown.TemplateHTML(simpleDocMD, false),
			defaultOptions))
}

// -----------------------------------------------------------------------------

// simpleWarningScoreWeight has the multipliers for different warning levels.
var simpleWarningScoreWeight = map[warning.Level]uint{
	warning.LevelRequired: 2,
	warning.LevelImplied:  1,
}

// simpleScoreWeight has the multipliers for different benchmark values.
var simpleBenchmarkScoreWeight = map[axis.BenchValue]uint{
	axis.Allocations: 1,
	axis.Nanoseconds: 2,
}
