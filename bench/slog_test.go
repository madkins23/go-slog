package bench

import (
	"testing"

	"github.com/madkins23/go-slog/bench/tests"
	"github.com/madkins23/go-slog/creator/json/slog"
)

// Benchmark_slog runs benchmarks for the log/slog JSON handler.
func Benchmark_slog(b *testing.B) {
	slogSuite := tests.NewSlogBenchmarkSuite(slog.Slog())
	tests.Run(b, slogSuite)
}
