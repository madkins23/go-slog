package keeper

import (
	_ "embed"

	"github.com/madkins23/go-slog/internal/markdown"
	"github.com/madkins23/go-slog/internal/scoring/axis"
	"github.com/madkins23/go-slog/internal/scoring/filter"
	"github.com/madkins23/go-slog/internal/scoring/score"
)

const ReplAttrName = "~ReplAttr"

var (
	//go:embed doc/repl-attr-doc.md
	replAttrDocMD string
)

func setupReplAttr() error {
	return score.AddKeeper(
		score.NewKeeper(
			ReplAttrName,
			axis.NewWarnings(
				defaultWarningScoreWeight,
				markdown.TemplateHTML(defaultXSumMD, false)),
			axis.NewBenchmarks(
				defaultBenchmarkScoreWeight,
				markdown.TemplateHTML(defaultYSumMD, false),
				nil),
			markdown.TemplateHTML(replAttrDocMD, false),
			defaultOptions,
			filter.ReplAttr()))
}
