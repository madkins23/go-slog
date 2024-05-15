package keeper

import (
	"github.com/madkins23/go-slog/infra/warning"
	"github.com/madkins23/go-slog/internal/scoring/axis"
	"github.com/madkins23/go-slog/internal/scoring/score"
)

const simpleName = "Simple"

func setupSimple() error {
	return score.AddKeeper(
		score.NewKeeper(
			simpleName,
			axis.NewWarnings(simpleWarningScoreWeight),
			axis.NewBenchmarks(simpleBenchmarkScoreWeight),
			defaultDocHTML))
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
