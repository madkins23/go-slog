package common

import (
	"math"

	"github.com/madkins23/go-slog/internal/scoring/score"
)

// -----------------------------------------------------------------------------

func FuzzyEqual(a, b score.Value) bool {
	const epsilon = 0.00000001
	return math.Abs(float64(a-b)) < epsilon
}
