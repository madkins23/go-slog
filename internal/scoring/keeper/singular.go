package keeper

import (
	_ "embed"

	"github.com/madkins23/go-slog/internal/markdown"
	"github.com/madkins23/go-slog/internal/scoring/axis"
	"github.com/madkins23/go-slog/internal/scoring/filter"
	"github.com/madkins23/go-slog/internal/scoring/score"
)

const SingleName = "~Singular"

var (
	//go:embed doc/singular-doc.md
	singularDocMD string
)

func setupSingular() error {
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
			markdown.TemplateHTML(singularDocMD, false),
			defaultOptions,
			filter.Singular()))
}
