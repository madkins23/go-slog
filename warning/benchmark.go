package warning

var (
	Mismatch = NewWarning(LevelRequired, "Mismatch",
		"Logged record does not match expected")

	NotDisabled = NewWarning(LevelRequired, "NotDisabled",
		"Logging was not properly disabled")
)

func Benchmark() []*Warning {
	return []*Warning{
		Mismatch,
		NotDisabled,
	}
}
