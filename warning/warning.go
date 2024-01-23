package warning

import "fmt"

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
)

func (l Level) String() string {
	if levelName, found := levelNames[l]; found {
		return levelName
	} else {
		return fmt.Sprintf("Unknown level %d", l)
	}
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
