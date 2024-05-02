package tests

import (
	warning2 "github.com/madkins23/go-slog/infra/warning"
)

// NewWarningManager generates an infra.Manager configured for SlogBenchmarkSuite.
func NewWarningManager(name string) *warning2.Manager {
	mgr := warning2.NewWarningManager(name, benchmarkMethodPrefix, "# ")
	mgr.Predefine(warning2.Benchmark()...)
	mgr.Predefine(warning2.Implied()...)
	mgr.Predefine(warning2.Suggested()...)
	return mgr
}
