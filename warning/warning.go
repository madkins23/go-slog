package warning

import (
	"fmt"
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

func ParseLevel(text string) (Level, error) {
	if level, found := levelParse[strings.ToLower(text)]; found {
		return level, nil
	}
	return levelUnused, fmt.Errorf("no warning level '%s'", text)
}

// Warning definition.
type Warning struct {
	// Level is the warning level.
	Level Level

	// Name of the warning.
	Name string

	// Description of the warning.
	Description string
}

var _ error = &Warning{}

func (w *Warning) Error() string {
	return fmt.Sprintf("%s [%s] %s", w.Level.String(), w.Name, w.Description)
}

func (w *Warning) ErrorExtra(extra string) error {
	return &warningError{
		Warning: w,
		extra:   extra,
	}
}

var _ error = &warningError{}

type warningError struct {
	*Warning
	extra string
}

func (we *warningError) Error() string {
	return fmt.Sprintf("%s: %s", we.Warning.Error(), we.extra)
}
