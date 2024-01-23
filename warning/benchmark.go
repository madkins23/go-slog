package warning

var (
	Mismatch = &Warning{
		Level:       LevelRequired,
		Name:        "Mismatch",
		Description: "Logged record does not match expected",
	}
	NotDisabled = &Warning{
		Level:       LevelRequired,
		Name:        "NotDisabled",
		Description: "Logging was not properly disabled",
	}
)

func Benchmark() []*Warning {
	return []*Warning{
		Mismatch,
		NotDisabled,
	}
}
