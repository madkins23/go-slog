package bench

import (
	"testing"

	"github.com/madkins23/go-slog/bench/tests"
	"github.com/madkins23/go-slog/creator/madkinssloggy"
)

// BenchmarkMadkinsSloggy runs benchmarks for the slog/JSONHandler JSON handler.
func BenchmarkMadkinsSloggy(b *testing.B) {
	slogSuite := tests.NewSlogBenchmarkSuite(madkinssloggy.Creator())
	tests.Run(b, slogSuite)
}
