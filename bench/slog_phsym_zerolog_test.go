package bench

import (
	"testing"

	"github.com/madkins23/go-slog/bench/tests"
	"github.com/madkins23/go-slog/creator/phsym_zerolog"
)

// Benchmark_slog_phsym_zerolog runs benchmarks for the phsym/zeroslog handler.
func Benchmark_slog_phsym_zerolog(b *testing.B) {
	slogSuite := tests.NewSlogBenchmarkSuite(phsym_zerolog.SlogPhsymZerolog())
	tests.Run(b, slogSuite)
}
