package replace

import (
	"strings"
	"testing"
)

const (
	alpha = "Home, Home on the Range!"
	bravo = "home, home on the range!"
)

// BenchmarkCompareChangeCase benchmarks comparing strings with two case conversions.
// This seems to take a factor of ten longer than the strings.EqualFold version.
// It also results in 2 allocs/op with 48 bytes.
func BenchmarkCompareChangeCase(b *testing.B) {
	var count uint
	for i := 0; i < b.N; i++ {
		if strings.ToUpper(alpha) == strings.ToUpper(bravo) {
			count++
		}
	}
}

// BenchmarkCompareEqualFold benchmarks comparing strings using strings.EqualFold.
// This appears to be the winner at 10% of the dual case conversion version.
// There are no memory allocations.
func BenchmarkCompareEqualFold(b *testing.B) {
	var count uint
	for i := 0; i < b.N; i++ {
		if strings.EqualFold(alpha, bravo) {
			count++
		}
	}
}
