package warning

var (
	// LevelSuggested warnings.
	//
	// Note: Update the number of warnings in the init function below.

	Duplicates = NewWarning(LevelSuggested, "Duplicates", "Duplicate field(s) found", `
		Some handlers (e.g. ^slog.JSONHandler^)
		will output multiple occurrences of the same field name
		if the logger is called with multiple instances of the same field,
		generally by using WithAttrs and then the same fields in the eventual log call (e.g. Info).
		This behavior is currently [under debate](https://github.com/golang/go/issues/59365)
		with no resolution at this time (2024-01-15) and a
		[release milestone of (currently unscheduled) Go 1.23](https://github.com/golang/go/milestone/212),
		(whereas [Go Release 1.22](https://tip.golang.org/doc/go1.22)
		is currently expected in February 2024).`)

	DurationSeconds = NewWarning(LevelSuggested, "DurationSeconds", "slog.Duration() logs seconds instead of nanoseconds", `
		The ^slog.JSONHandler^ uses nanoseconds for ^time.Duration^ but some other handlers use seconds.
		* [Go issue 59345: Nanoseconds is a recent change with Go 1.21](https://github.com/golang/go/issues/59345)`)

	DurationMillis = NewWarning(LevelSuggested, "DurationMillis", "slog.Duration() logs milliseconds instead of nanoseconds", `
		The ^slog.JSONHandler^ uses nanoseconds for ^time.Duration^ but some other handlers use seconds.
		* [Go issue 59345: Nanoseconds is a recent change with Go 1.21](https://github.com/golang/go/issues/59345)`)

	GroupWithTop = NewWarning(LevelSuggested, "GroupWithTop",
		"^WithGroup().With()^ ends up at top level of log record instead of in the group", `
		Almost all handlers treat ^logger.WithGroup(<name>).With(<attrs>)^ as writing ^<attrs>^ to the group ^<name>^.
		Some handlers write ^<attrs>^ to the top level of the log record.`)

	GroupDuration = NewWarning(LevelSuggested, "GroupDuration", "", `
		Some handlers that change the way ^time.Duration^ objects are logged (see warnings ^DurationMillis^ and ^DurationSeconds^)
		only manage to make the change at the top level of the logged record, duration objects in groups are still in nanoseconds.`)

	LevelCase = NewWarning(LevelSuggested, "LevelCase", "Log level in lowercase", `
		Each JSON log record contains the logging level of the log statement as a string.
		Different handlers provide that string in uppercase or lowercase.
		Documentation for [^slog.Level^](https://pkg.go.dev/log/slog@master#Level)
		says that its ^String()^ and ^MarshalJSON()^ methods will return uppercase
		but ^UnmarshalJSON()^ will parse in a case-insensitive manner.`)

	LevelWrong = NewWarning(LevelSuggested, "LevelWrong", "Log level is incorrect", `
		The log level name is not what was expected (e.g. "WARNING" instead of "WARN").
		This is different from the LevelCase warning which is from the right level name but the wrong character case.`)

	TimeMillis = NewWarning(LevelSuggested, "TimeMillis", "slog.Time() logs milliseconds instead of nanoseconds", `
		The ^slog.JSONHandler^ uses nanoseconds for ^time.Time^ but some other handlers use milliseconds.
		This does _not_ apply to the basic ^time^ field, only attribute fields.
		I can't find any supporting documentation or bug on this but
		[Go issue 59345](https://github.com/golang/go/issues/59345) (see previous warning)
		may have fixed this as well in Go 1.21.`)

	TimeSeconds = NewWarning(LevelSuggested, "TimeSeconds", "slog.Time() logs seconds instead of nanoseconds", `
		The ^slog.JSONHandler^ uses nanoseconds for ^time.Time^ but some other handlers use seconds.
		This does _not_ apply to the basic ^time^ field, only attribute fields.
		I can't find any supporting documentation or bug on this but
		[Go issue 59345](https://github.com/golang/go/issues/59345) (see previous warning)
		may have fixed this as well in Go 1.21.`)
)

func init() {
	// Always update this number when adding or removing Warning objects.
	addTestCount(LevelSuggested, 9)
}

// Suggested returns an array of all LevelSuggested warnings.
func Suggested() []*Warning {
	return WarningsForLevel(LevelSuggested)
}
