package bench

import (
	"testing"

	"github.com/madkins23/go-slog/bench/tests"
	"github.com/madkins23/go-slog/creator/svcrunnerjsonlog"
)

// BenchmarkSvcrunnerJsonlog runs benchmarks for the svcrunner/jsonlog handler.
func BenchmarkSvcrunnerJsonlog(b *testing.B) {
	slogSuite := tests.NewSlogBenchmarkSuite(svcrunnerjsonlog.Creator())
	tests.Run(b, slogSuite)
}
