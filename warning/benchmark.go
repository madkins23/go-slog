package warning

var (
	Mismatch = NewWarning(LevelRequired, "Mismatch", "Logged record does not match expected", `
		During benchmark testing the test is run once to see if the logged record matches expectations.`)

	NotDisabled = NewWarning(LevelRequired, "NotDisabled", "Logging was not properly disabled", `
		During benchmark testing a Debug log line is made with the current level set to Info.
		This warning is thrown if there is any logged output.`)
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
