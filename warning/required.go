package warning

var (
	EmptyAttributes = &Warning{
		Level:       LevelRequired,
		Name:        "EmptyAttributes",
		Description: "Empty attribute(s) logged (\"\":null)",
	}

	GroupEmpty = &Warning{
		Level:       LevelRequired,
		Name:        "GroupEmpty",
		Description: "Empty (sub)group(s) logged",
	}

	GroupInline = &Warning{
		Level:       LevelRequired,
		Name:        "Duplicates",
		Description: "Group with empty key does not inline subfields",
	}

	Resolver = &Warning{
		Level:       LevelRequired,
		Name:        "Resolver",
		Description: "LogValuer objects are not resolved",
	}

	ZeroPC = &Warning{
		Level:       LevelRequired,
		Name:        "ZeroPC",
		Description: "SourceKey logged for zero PC",
	}

	ZeroTime = &Warning{
		Level:       LevelRequired,
		Name:        "ZeroTime",
		Description: "Zero time is logged",
	}
)

func Required() []*Warning {
	return []*Warning{
		EmptyAttributes,
		GroupEmpty,
		GroupInline,
		Resolver,
		ZeroPC,
		ZeroTime,
	}
}
