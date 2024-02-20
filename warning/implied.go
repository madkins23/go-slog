package warning

var (
	DefaultLevel = NewWarning(LevelImplied, "DefaultLevel",
		"Handler doesn't default to slog.LevelInfo")

	MessageKey = NewWarning(LevelImplied, "MessageKey",
		"Wrong message key (should be 'msg')")

	NoReplAttr = NewWarning(LevelImplied, "NoReplAttr",
		"HandlerOptions.ReplaceAttr not available")

	NoReplAttrBasic = NewWarning(LevelImplied, "NoReplAttrBasic",
		"HandlerOptions.ReplaceAttr not available for basic fields")

	SourceKey = NewWarning(LevelImplied, "SourceKey",
		"Source data not logged when AddSource flag set")
)

func Implied() []*Warning {
	return WarningsForLevel(LevelImplied)
}
