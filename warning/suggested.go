package warning

var (
	Duplicates = NewWarning(LevelSuggested, "Duplicates",
		"Duplicate field(s) found")

	DurationSeconds = NewWarning(LevelSuggested, "DurationSeconds",
		"slog.Duration() logs seconds instead of nanoseconds")

	DurationMillis = NewWarning(LevelSuggested, "DurationMillis",
		"slog.Duration() logs milliseconds instead of nanoseconds")

	LevelCase = NewWarning(LevelSuggested, "LevelCase",
		"Log level in lowercase")

	TimeMillis = NewWarning(LevelSuggested, "TimeMillis",
		"slog.Time() logs milliseconds instead of nanoseconds")
)

func init() {
	// Always update this number when adding or removing Warning objects.
	addTestCount(LevelSuggested, 5)
}

func Suggested() []*Warning {
	return WarningsForLevel(LevelSuggested)
}
