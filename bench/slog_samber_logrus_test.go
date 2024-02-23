package bench

import (
	"testing"

	"github.com/madkins23/go-slog/bench/tests"
	"github.com/madkins23/go-slog/creator/logrus/samber"
)

// Benchmark_slog_samber_logrus runs benchmarks for the samber/slog-logrus handler.
func Benchmark_slog_samber_logrus(b *testing.B) {
	slogSuite := tests.NewSlogBenchmarkSuite(samber.SlogSamberLogrus())
	tests.Run(b, slogSuite)
}
