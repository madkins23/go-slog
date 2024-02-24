package bench

import (
	"testing"

	"github.com/madkins23/go-slog/bench/tests"
	"github.com/madkins23/go-slog/creator/slog_json"
)

// Benchmark_slog runs benchmarks for the log/slog JSON handler.
func Benchmark_slog(b *testing.B) {
	slogSuite := tests.NewSlogBenchmarkSuite(slog_json.SlogJson())
	tests.Run(b, slogSuite)
}
