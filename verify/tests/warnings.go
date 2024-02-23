package tests

import (
	"github.com/madkins23/go-slog/warning"
)

// NewWarningManager generates an infra.Manager configured for SlogTestSuite.
func NewWarningManager(name string) *warning.Manager {
	mgr := warning.NewWarningManager(name, "Test", "")
	mgr.Predefine(warning.Required()...)
	mgr.Predefine(warning.Implied()...)
	mgr.Predefine(warning.Suggested()...)
	return mgr
}
