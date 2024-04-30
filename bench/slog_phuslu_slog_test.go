package bench

import (
	"testing"

	"github.com/madkins23/go-slog/bench/tests"
	"github.com/madkins23/go-slog/creator/phusluslog"
)

// BenchmarkPhusluSlog runs benchmarks for the phuslu/slog handler.
func BenchmarkPhusluSlog(b *testing.B) {
	slogSuite := tests.NewSlogBenchmarkSuite(phusluslog.Creator())
	tests.Run(b, slogSuite)
}
