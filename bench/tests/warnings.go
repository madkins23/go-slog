package tests

import (
	"github.com/madkins23/go-slog/internal/warning"
)

// NewWarningManager generates an infra.Manager configured for SlogBenchmarkSuite.
func NewWarningManager(name string) *warning.Manager {
	mgr := warning.NewWarningManager(name, benchmarkMethodPrefix, "# ")
	mgr.Predefine(warning.Benchmark()...)
	mgr.Predefine(warning.Implied()...)
	mgr.Predefine(warning.Suggested()...)
	return mgr
}
