package warning

var (
	NoHandlerCreation = NewWarning(LevelAdmin, "NoHandlerCreation",
		"Test depends on unavailable handler creation")

	SkippingTest = NewWarning(LevelAdmin, "SkippingTest",
		"Skipping test")

	TestError = NewWarning(LevelAdmin, "TestError",
		"Test harness error")

	Undefined = NewWarning(LevelAdmin, "Undefined",
		"Undefined Warnings(s)")

	Unused = NewWarning(LevelAdmin, "Unused",
		"Unused Warnings(s)")
)

func init() {
	// Always update this number when adding or removing Warning objects.
	addTestCount(LevelAdmin, 5)
}

func Administrative() []*Warning {
	return WarningsForLevel(LevelAdmin)
}
