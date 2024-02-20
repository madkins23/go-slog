package warning

import (
	"log/slog"
	"sort"
	"sync"
)

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
	allWarnings = make([]*Warning, 0, 25)
	byName      = make(map[string]*Warning)
	warningTree = make(map[Level][]*Warning)
	warningLock sync.Mutex
)

func NewWarning(level Level, name string, description string) *Warning {
	warningLock.Lock()
	defer warningLock.Unlock()
	var found bool
	var warning *Warning
	if warning, found = byName[name]; found {
		// Warnings must have unique names despite having different levels.
		slog.Warn("Duplicate warning name", "name", name, "warning", warning)
	} else {
		warning = &Warning{
			Level:       level,
			Name:        name,
			Description: description,
		}
		allWarnings = append(allWarnings, warning)
		byName[name] = warning
	}
	return warning
}

// WarningsForLevel returns a list of warnings for the specified level.
func WarningsForLevel(level Level) []*Warning {
	warningLock.Lock()
	defer warningLock.Unlock()
	if len(warningTree) < 1 {
		buildTree()
	}
	return warningTree[level]
}

// buildTree constructs the warningTree global variable if it is empty.
// This function doesn't invoke warningLock so make sure it is locked prior to calling.
func buildTree() {
	var array []*Warning
	for _, warning := range allWarnings {
		warningTree[warning.Level] = append(warningTree[warning.Level], warning)
	}
	for _, array = range warningTree {
		sort.Slice(array, func(i, j int) bool {
			return array[i].Name < array[j].Name
		})
	}
}
