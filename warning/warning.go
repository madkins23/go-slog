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

	// Description of the warning.
	Description string

	// Summary of warning in Markdown
	summary string
}

var (
	allWarnings = make([]*Warning, 0, 25)
	byName      = make(map[string]*Warning)
	testCounts  = make(map[Level]int)
	warningTree = make(map[Level][]*Warning)
	warningLock sync.Mutex
)

func NewWarning(level Level, name string, description string, summary string) *Warning {
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
			summary:     fixSummary(summary),
		}
		allWarnings = append(allWarnings, warning)
		byName[name] = warning
	}
	return warning
}

// HasSummary returns true if there is summary data.
func (w *Warning) HasSummary() bool {
	return len(w.summary) > 0
}

// Summary converts the Markdown summary data into HTML and returns it.
func (w *Warning) Summary() template.HTML {
	return MD2HTML(w.summary)
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

func fixSummary(summary string) string {
	prefixSpaces := math.MaxInt
	scanner := bufio.NewScanner(strings.NewReader(summary))
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
	scanner = bufio.NewScanner(strings.NewReader(summary))
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
