package bench

import (
	"testing"

	"github.com/madkins23/go-slog/bench/tests"
	"github.com/madkins23/go-slog/creator/chanchalzap"
)

// BenchmarkChanchalZap runs benchmarks for the chanchal/zap handler.
func BenchmarkChanchalZap(b *testing.B) {
	slogSuite := tests.NewSlogBenchmarkSuite(chanchalzap.Creator(), "ChanchalZap")
	tests.Run(b, slogSuite)
}
