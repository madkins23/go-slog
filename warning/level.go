package warning

import (
	"fmt"
	"html/template"
	"strings"
)

// -----------------------------------------------------------------------------

// Level for warnings, used mainly to organize warning on output.
type Level uint

const (
	levelUnused Level = iota
	LevelAdmin
	LevelSuggested
	LevelImplied
	LevelRequired
)

var (
	LevelOrder = []Level{
		LevelRequired,
		LevelImplied,
		LevelSuggested,
		LevelAdmin,
	}
	levelNames = map[Level]string{
		LevelAdmin:     "Administrative",
		LevelSuggested: "Suggested",
		LevelImplied:   "Implied",
		LevelRequired:  "Required",
	}
	levelParse = map[string]Level{
		"admin":          LevelAdmin,
		"administrative": LevelAdmin,
		"suggested":      LevelSuggested,
		"implied":        LevelImplied,
		"required":       LevelRequired,
	}
)

func (l Level) String() string {
	if levelName, found := levelNames[l]; found {
		return levelName
	} else {
		return fmt.Sprintf("Unknown level %d", l)
	}
}

var levelDescriptions = map[Level]string{
	LevelRequired:  "The following warnings can be justified from requirements in the [`slog.Handler`](https://pkg.go.dev/log/slog@master#Handler) documentation.",
	LevelImplied:   "Warnings that seem to be implied by documentation but can't be considered required.",
	LevelSuggested: "These warnings are not AFAIK mandated by any documentation or requirements.",
	LevelAdmin:     "Warnings that provide information about the tests or conflicts with other warnings.",
}

func (l Level) Description() template.HTML {
	return MD2HTML(levelDescriptions[l])
}

func (l Level) Warnings() []*Warning {
	return WarningsForLevel(l)
}

func ParseLevel(text string) (Level, error) {
	if level, found := levelParse[strings.ToLower(text)]; found {
		return level, nil
	}
	return levelUnused, fmt.Errorf("no warning level '%s'", text)
}
