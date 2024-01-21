package warning

var (
	Duplicates = &Warning{
		Level:       LevelSuggested,
		Name:        "Duplicates",
		Description: "Duplicate field(s) found",
	}
	DurationSeconds = &Warning{
		Level:       LevelSuggested,
		Name:        "Duplicates",
		Description: "slog.Duration() logs seconds instead of nanoseconds",
	}
	DurationMillis = &Warning{
		Level:       LevelSuggested,
		Name:        "DurationMillis",
		Description: "slog.Duration() logs milliseconds instead of nanoseconds",
	}

	LevelCase = &Warning{
		Level:       LevelSuggested,
		Name:        "LevelCase",
		Description: "Log level in lowercase",
	}

	TimeMillis = &Warning{
		Level:       LevelSuggested,
		Name:        "TimeMillis",
		Description: "slog.Time() logs milliseconds instead of nanoseconds",
	}
)

func Suggested() []*Warning {
	return []*Warning{
		Duplicates,
		DurationSeconds,
		DurationMillis,
		LevelCase,
		TimeMillis,
	}
}
