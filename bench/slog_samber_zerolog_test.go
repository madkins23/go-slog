package bench

import (
	"testing"

	"github.com/madkins23/go-slog/bench/tests"
	"github.com/madkins23/go-slog/creator/samberzerolog"
)

// BenchmarkSamberZerolog runs benchmarks for the samber/slog-zerolog handler.
func BenchmarkSamberZerolog(b *testing.B) {
	slogSuite := tests.NewSlogBenchmarkSuite(samberzerolog.Creator(), "SamberZerolog")
	tests.Run(b, slogSuite)
}
