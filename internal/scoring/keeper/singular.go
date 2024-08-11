package keeper

import (
	_ "embed"

	"github.com/madkins23/go-slog/internal/markdown"
	"github.com/madkins23/go-slog/internal/scoring/axis"
	"github.com/madkins23/go-slog/internal/scoring/filter"
	"github.com/madkins23/go-slog/internal/scoring/score"
)

const SingleName = "~Dedup"

var (
	//go:embed doc/dedup-doc.md
	dedupDocMD string
)

func setupDedup() error {
	return score.AddKeeper(
		score.NewKeeper(
			SingleName,
			axis.NewWarnings(
				defaultWarningScoreWeight,
				markdown.TemplateHTML(defaultXSumMD, false)),
			axis.NewBenchmarks(
				defaultBenchmarkScoreWeight,
				markdown.TemplateHTML(defaultYSumMD, false),
				nil),
			markdown.TemplateHTML(dedupDocMD, false),
			defaultOptions,
			filter.Dedup()))
}
