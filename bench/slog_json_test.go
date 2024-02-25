package bench

import (
	"testing"

	"github.com/madkins23/go-slog/bench/tests"
	"github.com/madkins23/go-slog/creator/slogjson"
)

// BenchmarkSlogJSON runs benchmarks for the log/slog JSON handler.
func BenchmarkSlogJSON(b *testing.B) {
	slogSuite := tests.NewSlogBenchmarkSuite(slogjson.Creator(), "SlogJSON")
	tests.Run(b, slogSuite)
}
