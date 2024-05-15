package keeper

import (
	_ "embed"
	"html/template"

	"github.com/madkins23/go-slog/infra/warning"
	"github.com/madkins23/go-slog/internal/markdown"
	"github.com/madkins23/go-slog/internal/scoring/axis"
	"github.com/madkins23/go-slog/internal/scoring/score"
)

const DefaultName = "Default"

var (
	//go:embed doc/default.md
	defaultDocMD   string
	defaultDocHTML template.HTML
)

func setupDefault() error {
	defaultDocHTML = markdown.TemplateHTML(defaultDocMD, false)
	return score.AddKeeper(
		score.NewKeeper(
			DefaultName,
			axis.NewWarnings(defaultWarningScoreWeight),
			axis.NewBenchmarks(defaultBenchmarkScoreWeight),
			defaultDocHTML))
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
