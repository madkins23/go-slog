package tests

import "github.com/madkins23/go-slog/infra"

// NewWarningManager generates an infra.WarningManager configured for SlogBenchmarkSuite.
func NewWarningManager(name string) *infra.WarningManager {
	mgr := infra.NewWarningManager(name, benchmarkMethodPrefix, "# ")
	// No extra predefined warnings (yet?).
	return mgr
}
