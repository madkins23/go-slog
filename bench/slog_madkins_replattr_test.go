package bench

import (
	"testing"

	"github.com/madkins23/go-slog/bench/tests"
	"github.com/madkins23/go-slog/creator/madkinsreplattr"
)

// BenchmarkMadkinsReplAttr runs benchmarks for the madkins/replattr JSON handler.
func BenchmarkMadkinsReplAttr(b *testing.B) {
	slogSuite := tests.NewSlogBenchmarkSuite(madkinsflash.Creator())
	tests.Run(b, slogSuite)
}
