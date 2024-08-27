package common

import (
	"math"

	"github.com/madkins23/go-slog/internal/scoring/score"
)

const epsilon = 0.00000001
const percent = 5.0

// -----------------------------------------------------------------------------

func FuzzyEqual(a, b score.Value) bool {
	return math.Abs(float64(a-b)) < epsilon
}

// -----------------------------------------------------------------------------

func PercentEqual(a, b score.Value) bool {
	return PercentDifference(a, b) < percent
}

func PercentDifference(a, b score.Value) score.Value {
	if a == b {
		return 0
	}
	return score.Value(200 * math.Abs(float64(a-b)) / math.Abs(float64(a+b)))
}
