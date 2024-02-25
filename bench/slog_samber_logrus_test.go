package bench

import (
	"testing"

	"github.com/madkins23/go-slog/bench/tests"
	"github.com/madkins23/go-slog/creator/samberlogrus"
)

// BenchmarkSamberLogrus runs benchmarks for the samber/slog-logrus handler.
func BenchmarkSamberLogrus(b *testing.B) {
	slogSuite := tests.NewSlogBenchmarkSuite(samberlogrus.Creator())
	tests.Run(b, slogSuite)
}
