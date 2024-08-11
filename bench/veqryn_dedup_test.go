package bench

import (
	"testing"

	"github.com/madkins23/go-slog/bench/tests"
	"github.com/madkins23/go-slog/creator/veqryndedup"
)

// BenchmarkVeqrynDedupAppend runs benchmarks for the veqryn/dedup JSON handler in Append mode.
func BenchmarkVeqrynDedupAppend(b *testing.B) {
	slogSuite := tests.NewSlogBenchmarkSuite(veqryndedup.Creator(veqryndedup.Append))
	tests.Run(b, slogSuite)
}

// BenchmarkVeqrynDedupIgnore runs benchmarks for the veqryn/dedup JSON handler in Ignore mode.
func BenchmarkVeqrynDedupIgnore(b *testing.B) {
	slogSuite := tests.NewSlogBenchmarkSuite(veqryndedup.Creator(veqryndedup.Ignore))
	tests.Run(b, slogSuite)
}

// BenchmarkVeqrynDedupIncrement runs benchmarks for the veqryn/dedup JSON handler in Increment mode.
func BenchmarkVeqrynDedupIncrement(b *testing.B) {
	slogSuite := tests.NewSlogBenchmarkSuite(veqryndedup.Creator(veqryndedup.Increment))
	tests.Run(b, slogSuite)
}

// BenchmarkVeqrynDedupOverwrite runs benchmarks for the veqryn/dedup JSON handler in Overwrite mode.
func BenchmarkVeqrynDedupOverwrite(b *testing.B) {
	slogSuite := tests.NewSlogBenchmarkSuite(veqryndedup.Creator(veqryndedup.Overwrite))
	tests.Run(b, slogSuite)
}
