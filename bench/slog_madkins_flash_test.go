package bench

import (
	"testing"

	"github.com/madkins23/go-slog/bench/tests"
	"github.com/madkins23/go-slog/creator/madkinsflash"
)

// BenchmarkMadkinsFlash runs benchmarks for the madkins/flash JSON handler.
func BenchmarkMadkinsFlash(b *testing.B) {
	slogSuite := tests.NewSlogBenchmarkSuite(madkinsflash.Creator())
	tests.Run(b, slogSuite)
}
