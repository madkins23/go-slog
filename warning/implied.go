package warning

var (
	CanceledContext = NewWarning(LevelImplied, "Canceled Context",
		"Canceled context blocks logging")

	DefaultLevel = NewWarning(LevelImplied, "DefaultLevel",
		"Handler doesn't default to slog.LevelInfo")

	LevelMath = NewWarning(LevelImplied, "LevelMath",
		"Log levels are not properly treated as integers")

	MessageKey = NewWarning(LevelImplied, "MessageKey",
		"Wrong message key (should be 'msg')")

	NoReplAttr = NewWarning(LevelImplied, "NoReplAttr",
		"HandlerOptions.ReplaceAttr not available")

	NoReplAttrBasic = NewWarning(LevelImplied, "NoReplAttrBasic",
		"HandlerOptions.ReplaceAttr not available for basic fields")

	SourceKey = NewWarning(LevelImplied, "SourceKey",
		"Source data not logged when AddSource flag set")

	WithGroup = NewWarning(LevelImplied, "WithGroup",
		"WithGroup doesn't embed following attributes into group")
)

func init() {
	// Always update this number when adding or removing Warning objects.
	addTestCount(LevelImplied, 8)
}

func Implied() []*Warning {
	return WarningsForLevel(LevelImplied)
}
