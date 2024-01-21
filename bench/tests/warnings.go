package tests

import (
	"github.com/madkins23/go-slog/infra"
	"github.com/madkins23/go-slog/warning"
)

var (
	Mismatch = &warning.Warning{
		Level:       warning.LevelRequired,
		Name:        "Mismatch",
		Description: "Logged record does not match expected",
	}
	NotDisabled = &warning.Warning{
		Level:       warning.LevelRequired,
		Name:        "NotDisabled",
		Description: "Logging was not properly disabled",
	}
)

var benchmarkWarnings = []*warning.Warning{
	Mismatch,
	NotDisabled,
}

// NewWarningManager generates an infra.WarningManager configured for SlogBenchmarkSuite.
func NewWarningManager(name string) *infra.WarningManager {
	mgr := infra.NewWarningManager(name, benchmarkMethodPrefix, "# ")
	mgr.Predefine(benchmarkWarnings...)
	return mgr
}
