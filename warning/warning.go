package warning

// Warning definition.
type Warning struct {
	// Level is the warning level.
	Level Level

	// Name of the warning.
	Name string

	// Description of the warning.
	Description string
}

var (
	allWarnings []*Warning
	usedNames   map[string]bool
)

func NewWarning(level Level, name string, description string) *Warning {
	return &Warning{
		Level:       level,
		Name:        name,
		Description: description,
	}
}
