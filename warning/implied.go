package warning

var (
	DefaultLevel = &Warning{
		Level:       LevelImplied,
		Name:        "DefaultLevel",
		Description: "Handler doesn't default to slog.LevelInfo",
	}

	MessageKey = &Warning{
		Level:       LevelImplied,
		Name:        "MessageKey",
		Description: "Wrong message key (should be 'msg')",
	}

	NoReplAttr = &Warning{
		Level:       LevelImplied,
		Name:        "NoReplAttr",
		Description: "HandlerOptions.ReplaceAttr not available",
	}

	NoReplAttrBasic = &Warning{
		Level:       LevelImplied,
		Name:        "NoReplAttrBasic",
		Description: "HandlerOptions.ReplaceAttr not available for basic fields",
	}

	SourceKey = &Warning{
		Level:       LevelImplied,
		Name:        "SourceKey",
		Description: "Source data not logged when AddSource flag set",
	}
)

func Implied() []*Warning {
	return []*Warning{
		DefaultLevel,
		MessageKey,
		NoReplAttr,
		NoReplAttrBasic,
		SourceKey,
	}
}
