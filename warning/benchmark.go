package warning

var (
	Mismatch = NewWarning(LevelRequired, "Mismatch",
		"Logged record does not match expected")

	NotDisabled = NewWarning(LevelRequired, "NotDisabled",
		"Logging was not properly disabled")
)

func init() {
	// Always update this number when adding or removing Warning objects.
	addTestCount(LevelRequired, 2)
}

func Benchmark() []*Warning {
	return []*Warning{
		Mismatch,
		NotDisabled,
	}
}
