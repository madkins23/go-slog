package score

import (
	"math/rand/v2"
	"testing"

	"github.com/stretchr/testify/require"
)

var expectedTotals = map[string]Value{
	"Bench.Attributes":          98.547561494122533,
	"Bench.BigGroup":            96.703480622510213,
	"Bench.Disabled":            0.000000000000000,
	"Bench.KeyValues":           98.841919191919189,
	"Bench.Logging":             98.844005415612685,
	"Bench.Simple":              97.294605741059854,
	"Bench.SimpleSource":        100.000000000000000,
	"Bench.WithAttrsAttributes": 98.819085849983182,
	"Bench.WithAttrsKeyValues":  98.845662454799069,
	"Bench.WithAttrsSimple":     99.360456122804365,
	"Bench.WithGroupAttributes": 98.750441779116628,
	"Bench.WithGroupKeyValues":  98.853746377672564,
}

// TestAdditionOrder checks the score.Value.Round function's performance.
// The problem is that adding floating point numbers can be order dependent
// due to [accuracy] issues in floating point arithmetic.
// The score.Value.Round function is used to get rid of minor accuracy issue
//
// [accuracy] https://en.wikipedia.org/wiki/Floating-point_arithmetic#Accuracy_problems
func TestAdditionOrder(t *testing.T) {
	tests := make([]string, 0, len(expectedTotals))
	for test := range expectedTotals {
		tests = append(tests, test)
	}
	var last Value
	for i := 0; i < 25; i++ {
		rand.Shuffle(len(tests), func(i, j int) {
			tests[i], tests[j] = tests[j], tests[i]
		})
		var total Value
		for _, test := range tests {
			total += expectedTotals[test]
		}
		total = total.Round()
		if i > 0 {
			require.Equal(t, last, total)
		}
		last = total
	}
}
