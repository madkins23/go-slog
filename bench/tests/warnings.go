package tests

import (
	"github.com/madkins23/go-slog/internal/test"
	"github.com/madkins23/go-slog/warning"
)

// NewWarningManager generates an infra.WarningManager configured for SlogBenchmarkSuite.
func NewWarningManager(name string) *test.WarningManager {
	mgr := test.NewWarningManager(name, benchmarkMethodPrefix, "# ")
	mgr.Predefine(warning.Benchmark()...)
	mgr.Predefine(warning.Implied()...)
	mgr.Predefine(warning.Suggested()...)
	return mgr
}
