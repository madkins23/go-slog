package tests

import "github.com/madkins23/go-slog/infra"

var (
	WarnNoHandlerCreation = &infra.Warning{
		Level: infra.WarnLevelAdmin,
		Name:  "Test depends on unavailable handler creation",
	}
)

var warnings = []*infra.Warning{
	WarnNoHandlerCreation,
}

// NewWarningManager generates an infra.WarningManager configured for SlogBenchmarkSuite.
func NewWarningManager(name string) *infra.WarningManager {
	mgr := infra.NewWarningManager(name, benchmarkMethodPrefix)
	mgr.Predefine(warnings...)
	return mgr
}
