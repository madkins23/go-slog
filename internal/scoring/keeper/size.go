package keeper

import (
	_ "embed"

	"github.com/madkins23/go-slog/internal/data"
	"github.com/madkins23/go-slog/internal/markdown"
	"github.com/madkins23/go-slog/internal/scoring/axis"
	"github.com/madkins23/go-slog/internal/scoring/score"
)

const sizeName = "Size"

var (
	//go:embed doc/simple.md
	sizeDocMD string
)

func setupSize() error {
	return score.AddKeeper(
		score.NewKeeper(
			sizeName,
			axis.NewBenchmarks("Large", defaultBenchmarkScoreWeight, largeTests, nil),
			axis.NewBenchmarks("Small", defaultBenchmarkScoreWeight, nil, largeTests),
			markdown.TemplateHTML(sizeDocMD, false)))
}

// -----------------------------------------------------------------------------

var (
	largeTests = []data.TestTag{
		"Bench.BigGroup",
		"Bench.Logging",
	}
)
