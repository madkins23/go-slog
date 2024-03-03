package warning

import (
	"bufio"
	"bytes"
	"html/template"
	"log/slog"
	"math"
	"sort"
	"strings"
	"sync"
	"unicode"
)

// Warning definition.
type Warning struct {
	// Level is the warning level.
	Level Level

	// Name of the warning.
	Name string

	// Summary of the warning.
	Summary string

	// Description of warning in Markdown
	description string
}

var (
	allWarnings = make([]*Warning, 0, 25)
	byName      = make(map[string]*Warning)
	testCounts  = make(map[Level]int)
	warningTree = make(map[Level][]*Warning)
	warningLock sync.Mutex
)

func NewWarning(level Level, name, summary, description string) *Warning {
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
			Summary:     summary,
			description: fixDescription(description),
		}
		allWarnings = append(allWarnings, warning)
		byName[name] = warning
	}
	return warning
}

// ByName returns the warning with the specified name, if any, else nil.
func ByName(name string) *Warning {
	return byName[name]
}

// HasDescription returns true if there is description data.
func (w *Warning) HasDescription() bool {
	return w != nil && len(w.description) > 0
}

// Description converts the Markdown description data into HTML and returns it.
func (w *Warning) Description() template.HTML {
	return MD2HTML(w.description)
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

// addTestCount supports unit test TestWarnings.
// Whenever creating a new Warning object make sure to update this count.
func addTestCount(level Level, increment uint) {
	testCounts[level] += int(increment)
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

func fixDescription(description string) string {
	prefixSpaces := math.MaxInt
	scanner := bufio.NewScanner(strings.NewReader(description))
	for scanner.Scan() {
		line := scanner.Text()
		onlySpaces := true
		var numSpaces int
		for i, c := range line {
			if !unicode.IsSpace(c) {
				onlySpaces = false
				numSpaces = i
				break
			}
		}
		if !onlySpaces && numSpaces < prefixSpaces {
			prefixSpaces = numSpaces
		}
	}
	var result bytes.Buffer
	scanner = bufio.NewScanner(strings.NewReader(description))
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > prefixSpaces {
			line = line[prefixSpaces:]
		}
		result.WriteString(line)
		result.WriteByte('\n')
	}
	return result.String()
}
