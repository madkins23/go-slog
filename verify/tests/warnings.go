package tests

import (
	warning2 "github.com/madkins23/go-slog/infra/warning"
)

// NewWarningManager generates an infra.Manager configured for SlogTestSuite.
func NewWarningManager(name string) *warning2.Manager {
	mgr := warning2.NewWarningManager(name, "Test", "")
	mgr.Predefine(warning2.Required()...)
	mgr.Predefine(warning2.Implied()...)
	mgr.Predefine(warning2.Suggested()...)
	return mgr
}
