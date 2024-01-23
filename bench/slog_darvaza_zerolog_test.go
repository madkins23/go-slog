package bench

import (
	"testing"

	"github.com/madkins23/go-slog/bench/tests"
	"github.com/madkins23/go-slog/creator"
	"github.com/madkins23/go-slog/warning"
)

// Benchmark_slog_darvaza_zerolog runs benchmarks for the darvaza/zerolog handler.
func Benchmark_slog_darvaza_zerolog(b *testing.B) {
	slogSuite := tests.NewSlogBenchmarkSuite(creator.SlogDarvazaZerolog())
	slogSuite.WarnOnly(warning.DurationMillis)
	slogSuite.WarnOnly(warning.TimeMillis)
	tests.Run(b, slogSuite)
}
