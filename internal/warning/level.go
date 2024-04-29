package warning

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/madkins23/go-slog/internal/markdown"
)

// -----------------------------------------------------------------------------

// Level for warnings, used mainly to organize warning on output.
type Level uint

const (
	levelUnused Level = iota

	// LevelAdmin contains administrative warnings.
	LevelAdmin

	// LevelSuggested contains "suggested" warnings.
	LevelSuggested

	// LevelImplied contains warnings implied by documentation.
	LevelImplied

	// LevelRequired contains warnings about conflicts with documentation.
	LevelRequired
)

var (
	// LevelOrder returns an array of Level items ordered from most to least important.
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

// String implements the string interface for Level objects.
func (l Level) String() string {
	if levelName, found := levelNames[l]; found {
		return levelName
	} else {
		return fmt.Sprintf("Unknown level %d", l)
	}
}

// levelSummaries maps levels to summary Markdown strings with short summaries.
var levelSummaries = map[Level]string{
	LevelRequired:  `Warnings that can be justified from requirements in the [^slog.Handler^](https://pkg.go.dev/log/slog@master#Handler) documentation.`,
	LevelImplied:   `Warnings that seem to be implied by documentation but can't be considered required.`,
	LevelSuggested: `Warnings not mandated by any documentation or requirements.`,
	LevelAdmin:     `Warnings that provide information about the tests or conflicts with other warnings.`,
}

// Summary returns template.HTML derived from the Level summary Markdown strings.
func (l Level) Summary() template.HTML {
	return markdown.TemplateHTML(levelSummaries[l], true)
}

// Warnings returns an array of warnings for the current Level.
func (l Level) Warnings() []*Warning {
	return WarningsForLevel(l)
}

// ParseLevel attempts to parse a string as a Level name.
// If found, the Level is returned, otherwise an error.
func ParseLevel(text string) (Level, error) {
	if level, found := levelParse[strings.ToLower(text)]; found {
		return level, nil
	}
	return levelUnused, fmt.Errorf("no warning level '%s'", text)
}
