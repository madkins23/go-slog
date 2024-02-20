package warning

var (
	EmptyAttributes = NewWarning(LevelRequired, "EmptyAttributes",
		"Empty attribute(s) logged (\"\":null)")

	GroupEmpty = NewWarning(LevelRequired, "GroupEmpty",
		"Empty (sub)group(s) logged")

	GroupInline = NewWarning(LevelRequired, "GroupInline",
		"Group with empty key does not inline subfields")

	Resolver = NewWarning(LevelRequired, "Resolver",
		"LogValuer objects are not resolved")

	ZeroPC = NewWarning(LevelRequired, "ZeroPC",
		"SourceKey logged for zero PC")

	ZeroTime = NewWarning(LevelRequired, "ZeroTime",
		"Zero time is logged")
)

func Required() []*Warning {
	return WarningsForLevel(LevelRequired)
}
