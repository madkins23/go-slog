package bench

import (
	"testing"

	"github.com/madkins23/go-slog/bench/tests"
	"github.com/madkins23/go-slog/creator/samberzap"
)

// BenchmarkSamberZap runs benchmarks for the darvaza/slog-zap handler.
func BenchmarkSamberZap(b *testing.B) {
	slogSuite := tests.NewSlogBenchmarkSuite(samberzap.Creator(), "SamberZap")
	tests.Run(b, slogSuite)
}
