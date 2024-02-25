package warning

var (
	EmptyAttributes = NewWarning(LevelRequired, "EmptyAttributes", "Empty attribute(s) logged (\"\":null)", `
		Handlers are supposed to avoid logging empty attributes.  
		* ["- If an Attr's key and value are both the zero value, ignore the Attr."](https://pkg.go.dev/log/slog@master#Handler)`)

	GroupEmpty = NewWarning(LevelRequired, "GroupEmpty", "Empty (sub)group(s) logged", `
		Handlers should not log groups (or subgroups) without fields,
		whether or not the have non-empty names.
		* ["- If a group has no Attrs (even if it has a non-empty key), ignore it."](https://pkg.go.dev/log/slog@master#Handler)`)

	GroupInline = NewWarning(LevelRequired, "GroupInline", "Group with empty key does not inline subfields", `
		Handlers should expand groups named "" (the empty string) into the enclosing log record.  
		* ["- If a group's key is empty, inline the group's Attrs."](https://pkg.go.dev/log/slog@master#Handler)`)

	Resolver = NewWarning(LevelRequired, "Resolver", "LogValuer objects are not resolved", `
		Handlers should resolve all objects implementing the
		[^LogValuer^](https://pkg.go.dev/log/slog@master#LogValuer) or
		[^Stringer^](https://pkg.go.dev/fmt#Stringer) interfaces.
		This is a powerful feature which can customize logging of objects and
		[speed up logging by delaying argument resolution until logging time](https://pkg.go.dev/log/slog@master#hdr-Performance_considerations).
		* ["- Attr's values should be resolved."](https://pkg.go.dev/log/slog@master#Handler)`)

	ZeroPC = NewWarning(LevelRequired, "ZeroPC", "SourceKey logged for zero PC", `
		The ^slog.Record.PC^ field can be loaded with a program counter (PC).
		This is normally done by the ^slog.Logger^ code.
		If the PC is non-zero and the ^slog.HandlerOptions.AddSource^ flag is set
		the ^source^ field will contain a [^slog.Source^](https://pkg.go.dev/log/slog@master#Source) record
		containing the function name, file name, and file line at which the log record was generated.
		If the PC is zero then this field and its associated group should not be logged.
		* ["- If r.PC is zero, ignore it."](https://pkg.go.dev/log/slog@master#Handler)`)

	ZeroTime = NewWarning(LevelRequired, "ZeroTime", "Zero time is logged", `
		Handlers should not log the basic ^time^ field if it is zero.
		* ["- If r.Time is the zero time, ignore the time."](https://pkg.go.dev/log/slog@master#Handler)`)
)

func init() {
	// Always update this number when adding or removing Warning objects.
	addTestCount(LevelRequired, 6)
}

func Required() []*Warning {
	return WarningsForLevel(LevelRequired)
}
