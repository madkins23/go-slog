package tests

import "github.com/madkins23/go-slog/infra"

var (
	WarnDefaultLevel = &infra.Warning{
		Level: infra.WarnLevelImplied,
		Name:  "Handler doesn't default to slog.LevelInfo",
	}
	WarnDuplicates = &infra.Warning{
		Level: infra.WarnLevelSuggested,
		Name:  "Duplicate field(s) found",
	}
	WarnDurationSeconds = &infra.Warning{
		Level: infra.WarnLevelSuggested,
		Name:  "slog.Duration() logs seconds instead of nanoseconds",
	}
	WarnDurationMillis = &infra.Warning{
		Level: infra.WarnLevelSuggested,
		Name:  "slog.Duration() logs milliseconds instead of nanoseconds",
	}
	WarnEmptyAttributes = &infra.Warning{
		Level: infra.WarnLevelRequired,
		Name:  "Empty attribute(s) logged (\"\":null)",
	}
	WarnGroupEmpty = &infra.Warning{
		Level: infra.WarnLevelRequired,
		Name:  "Empty (sub)group(s) logged",
	}
	WarnGroupInline = &infra.Warning{
		Level: infra.WarnLevelRequired,
		Name:  "Group with empty key does not inline subfields",
	}
	WarnLevelCase = &infra.Warning{
		Level: infra.WarnLevelSuggested,
		Name:  "Log level in lowercase",
	}
	WarnMessageKey = &infra.Warning{
		Level: infra.WarnLevelImplied,
		Name:  "Wrong message key (should be 'msg')",
	}
	WarnNoReplAttr = &infra.Warning{
		Level: infra.WarnLevelImplied,
		Name:  "HandlerOptions.ReplaceAttr not available",
	}
	WarnNoReplAttrBasic = &infra.Warning{
		Level: infra.WarnLevelImplied,
		Name:  "HandlerOptions.ReplaceAttr not available for basic fields",
	}
	WarnResolver = &infra.Warning{
		Level: infra.WarnLevelRequired,
		Name:  "LogValuer objects are not resolved",
	}
	WarnSourceKey = &infra.Warning{
		Level: infra.WarnLevelImplied,
		Name:  "Source data not logged when AddSource flag set",
	}
	WarnTimeMillis = &infra.Warning{
		Level: infra.WarnLevelSuggested,
		Name:  "slog.Time() logs milliseconds instead of nanoseconds",
	}
	WarnZeroPC = &infra.Warning{
		Level: infra.WarnLevelRequired,
		Name:  "SourceKey logged for zero PC",
	}
	WarnZeroTime = &infra.Warning{
		Level: infra.WarnLevelRequired,
		Name:  "Zero time is logged",
	}
)

var warnings = []*infra.Warning{
	WarnDefaultLevel,
	WarnDuplicates,
	WarnEmptyAttributes,
	WarnGroupInline,
	WarnLevelCase,
	WarnMessageKey,
	WarnDurationMillis,
	WarnDurationSeconds,
	WarnTimeMillis,
	WarnNoReplAttr,
	WarnNoReplAttrBasic,
	WarnResolver,
	WarnSourceKey,
	WarnGroupEmpty,
	WarnZeroPC,
	WarnZeroTime,
}

// NewWarningManager generates an infra.WarningManager configured for SlogTestSuite.
func NewWarningManager(name string) *infra.WarningManager {
	mgr := infra.NewWarningManager(name)
	mgr.Predefine(warnings...)
	return mgr
}
