package tests

import (
	"github.com/madkins23/go-slog/infra"
	"github.com/madkins23/go-slog/warning"
)

// NewWarningManager generates an infra.WarningManager configured for SlogTestSuite.
func NewWarningManager(name string) *infra.WarningManager {
	mgr := infra.NewWarningManager(name, "Test", "")
	mgr.Predefine(warning.Required()...)
	mgr.Predefine(warning.Implied()...)
	mgr.Predefine(warning.Suggested()...)
	return mgr
}
