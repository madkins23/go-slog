package warning

var (
	NoHandlerCreation = &Warning{
		Level:       LevelAdmin,
		Name:        "NoHandlerCreation",
		Description: "Test depends on unavailable handler creation",
	}
	SkippingTest = &Warning{
		Level:       LevelAdmin,
		Name:        "SkippingTest",
		Description: "Skipping test",
	}
	Undefined = &Warning{
		Level:       LevelAdmin,
		Name:        "Undefined",
		Description: "Undefined Warnings(s)",
	}
	Unused = &Warning{
		Level:       LevelAdmin,
		Name:        "Unused",
		Description: "Unused Warnings(s)",
	}
)

func Administrative() []*Warning {
	return []*Warning{
		NoHandlerCreation,
		SkippingTest,
		Undefined,
		Unused,
	}
}
