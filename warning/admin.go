package warning

var (
	NoHandlerCreation = NewWarning(LevelAdmin, "NoHandlerCreation", "Test depends on unavailable handler creation", `
		Some benchmark tests depend on access to a ^slog.Handler^ object.
		Some ^slog^ implementations create a ^slog.Logger^ but no ^slog.Handler^.
		In this case the relevant benchmark tests can't be run.`)

	SkippingTest = NewWarning(LevelAdmin, "SkippingTest", "Skipping test", `
		A test has been skipped, likely due to the specification of some other warning.`)

	TestError = NewWarning(LevelAdmin, "TestError", "Test harness error", `
		Some sort of error has occurred during testing.
		This will generally require a programming fix.`)

	Undefined = NewWarning(LevelAdmin, "Undefined", "Undefined Warnings(s)", `
		An attempt to call ^WarnOnly^ with an undefined warning.
		Warnings must be predefined to the ^Manager^ prior to use.`)

	Unused = NewWarning(LevelAdmin, "Unused", "Unused Warnings(s)", `
		If a warning is specified but the condition is not actually present
		one of these warnings will be issued with the specified warning.
		These are intended to help clean out unnecessary ^WarnOnly^ settings
		from a test suite as issues are fixed in the tested handler.`)
)

func init() {
	// Always update this number when adding or removing Warning objects.
	addTestCount(LevelAdmin, 5)
}

func Administrative() []*Warning {
	return WarningsForLevel(LevelAdmin)
}
