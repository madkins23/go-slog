package bench

import (
	"testing"

	"github.com/madkins23/go-slog/bench/tests"
	"github.com/madkins23/go-slog/creator/phsymzerolog"
)

// BenchmarkPhsymZerolog runs benchmarks for the phsym/zeroslog handler.
func BenchmarkPhsymZerolog(b *testing.B) {
	slogSuite := tests.NewSlogBenchmarkSuite(phsymzerolog.Creator())
	tests.Run(b, slogSuite)
}
