package warning

var (
	CanceledContext = NewWarning(LevelImplied, "CanceledContext", "Canceled context blocks logging", `
		["The context is provided to support applications that provide logging information along the call chain. In a break with usual Go practice,
		the Handle method should not treat a canceled context as a signal to stop work."](https://github.com/golang/example/tree/master/slog-handler-guide#the-handle-method)`)

	DefaultLevel = NewWarning(LevelImplied, "DefaultLevel", "Handler doesn't default to slog.LevelInfo", `
		A new ^slog.Handler^ should default to ^slog.LevelInfo^.  
		* ["First, we wanted the default level to be Info, Since Levels are ints, Info is the default value for int, zero."](https://pkg.go.dev/log/slog@master#Level)`)

	LevelMath = NewWarning(LevelImplied, "LevelMath", "Log levels are not properly treated as integers", `
		[Log levels are actually numbers](https://pkg.go.dev/log/slog@master#Level), with space between them for user-defined levels.
		Handlers should properly handle numeric levels and math applied to level values.`)

	MessageKey = NewWarning(LevelImplied, "MessageKey", "Wrong message key (should be 'msg')", `
		The field name of the "message" key should be ^msg^.  
		* [Constant values are defined for ^slog/log^](https://pkg.go.dev/log/slog@master#pkg-constants)  
		* [Field values are defined for the ^JSONHandler.Handle()^ implementation](https://pkg.go.dev/log/slog@master#JSONHandler.Handle)`)

	NoReplAttr = NewWarning(LevelImplied, "NoReplAttr", "HandlerOptions.ReplaceAttr not available", `
		If [^HandlerOptions.ReplaceAttr^](https://pkg.go.dev/log/slog@master#HandlerOptions)
		is provided it should be honored by the handler.
		However, documentation on implementing handler methods seems to suggest it is optional.  
		* [Behavior defined for ^slog.HandlerOptions^](https://pkg.go.dev/log/slog@master#HandlerOptions)  
		* ["You might also consider adding a ReplaceAttr option to your handler, like the one for the built-in handlers."](https://github.com/golang/example/tree/master/slog-handler-guide#implementing-handler-methods)`)

	NoReplAttrBasic = NewWarning(LevelImplied, "NoReplAttrBasic", "HandlerOptions.ReplaceAttr not available for basic fields", `
		Some handlers (e.g. ^phsym/zeroslog^) support
		[^HandlerOptions.ReplaceAttr^](https://pkg.go.dev/log/slog@master#HandlerOptions)
		except for the four main fields ^time^, ^level^, ^msg^, and ^source^.
		When that is the case it is better to use this (^WarnNoReplAttrBasic^) warning.`)

	SourceKey = NewWarning(LevelImplied, "SourceKey", "Source data not logged when AddSource flag set", `
		Handlers should log source data when the ^slog.HandlerOptions.AddSource^ flag is set.
		* [Flag declaration as ^slog.HandlerOptions^ field](https://pkg.go.dev/log/slog@master#HandlerOptions)
		* [Behavior defined for ^JSONHandler.Handle()^](https://pkg.go.dev/log/slog@master#JSONHandler.Handle)
		* [Definition of source data record](https://pkg.go.dev/log/slog@master#Source)`)

	WithGroup = NewWarning(LevelImplied, "WithGroup", "WithGroup doesn't embed following attributes into group", `
		Complex log statements involving ^WithGroup^ require attributes to be attached to groups.
		This warning represents situations where the attributes are attached to the wrong log group.`)
)

func init() {
	// Always update this number when adding or removing Warning objects.
	addTestCount(LevelImplied, 8)
}

func Implied() []*Warning {
	return WarningsForLevel(LevelImplied)
}
