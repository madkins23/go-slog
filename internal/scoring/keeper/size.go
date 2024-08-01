package keeper

import (
	_ "embed"

	"github.com/madkins23/go-slog/internal/data"
	"github.com/madkins23/go-slog/internal/markdown"
	"github.com/madkins23/go-slog/internal/scoring/axis"
	"github.com/madkins23/go-slog/internal/scoring/filter"
	"github.com/madkins23/go-slog/internal/scoring/score"
)

const sizeName = "Size"

var (
	//go:embed doc/size-doc.md
	sizeDocMD string

	//go:embed doc/size-sum-x.md
	sizeXSumMD string

	//go:embed doc/size-sum-y.md
	sizeYSumMD string
)

var sizeOptions = &score.KeeperOptions{
	ChartTitle: "Large vs Small",
	ChartCaption: `
		Higher numbers are better on both axes. The "good" zone is the upper right and the "bad" zone is the lower left.`,
}

func setupSize() error {
	return score.AddKeeper(
		score.NewKeeper(
			sizeName,
			axis.NewBenchmarks(
				defaultBenchmarkScoreWeight,
				markdown.TemplateHTML(sizeXSumMD, false),
				&axis.BenchOptions{Name: "Large", IncludeTests: largeTests}),
			axis.NewBenchmarks(
				defaultBenchmarkScoreWeight,
				markdown.TemplateHTML(sizeYSumMD, false),
				&axis.BenchOptions{Name: "Small", ExcludeTests: largeTests}),
			markdown.TemplateHTML(sizeDocMD, false),
			sizeOptions,
			filter.Basic()))
}

// -----------------------------------------------------------------------------

var (
	largeTests = []data.TestTag{
		"Bench.BigGroup",
		"Bench.Logging",
	}
)
