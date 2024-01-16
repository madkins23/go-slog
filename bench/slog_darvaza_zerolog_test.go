package bench

import (
	"testing"

	"github.com/madkins23/go-slog/bench/tests"
	"github.com/madkins23/go-slog/creator"
)

// Benchmark_slog_darvaza_zerolog runs benchmarks for the darvaza/zerolog handler.
func Benchmark_slog_darvaza_zerolog(b *testing.B) {
	slogSuite := tests.NewSlogBenchmarkSuite(creator.SlogDarvazaZerolog())
	tests.Run(b, slogSuite)
}
