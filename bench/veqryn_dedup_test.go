package bench

import (
	"testing"

	"github.com/madkins23/go-slog/bench/tests"
	"github.com/madkins23/go-slog/creator/veqryndedup"
)

// BenchmarkVeqrynDedupGroup runs benchmarks for the veqryn/dedup JSON handler in Group mode.
func BenchmarkVeqrynDedupGroup(b *testing.B) {
	slogSuite := tests.NewSlogBenchmarkSuite(veqryndedup.Creator(veqryndedup.Group))
	tests.Run(b, slogSuite)
}

// BenchmarkVeqrynDedupIgnore runs benchmarks for the veqryn/dedup JSON handler in Ignore mode.
func BenchmarkVeqrynDedupIgnore(b *testing.B) {
	slogSuite := tests.NewSlogBenchmarkSuite(veqryndedup.Creator(veqryndedup.Ignore))
	tests.Run(b, slogSuite)
}

// BenchmarkVeqrynDedupIncr runs benchmarks for the veqryn/dedup JSON handler in Incr mode.
func BenchmarkVeqrynDedupIncr(b *testing.B) {
	slogSuite := tests.NewSlogBenchmarkSuite(veqryndedup.Creator(veqryndedup.Incr))
	tests.Run(b, slogSuite)
}

// BenchmarkVeqrynDedupOverwrite runs benchmarks for the veqryn/dedup JSON handler in Overwrite mode.
func BenchmarkVeqrynDedupOverwrite(b *testing.B) {
	slogSuite := tests.NewSlogBenchmarkSuite(veqryndedup.Creator(veqryndedup.Over))
	tests.Run(b, slogSuite)
}
