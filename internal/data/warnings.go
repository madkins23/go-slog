package data

import (
	"flag"
	"fmt"
	"html/template"
	"regexp"
	"sort"
	"strings"

	"github.com/madkins23/go-slog/infra/warning"
	"github.com/madkins23/go-slog/internal/markdown"
)

var verifyFile = flag.String("verify", "", "Load verification data from path (optional)")

// -----------------------------------------------------------------------------

type handlerData struct {
	name    string
	summary string
	links   map[string]string
}

// Warnings encapsulates benchmark records by TestTag and HandlerTag.
type Warnings struct {
	byTest      map[TestTag]*Levels
	byHandler   map[HandlerTag]*Levels
	byWarning   map[string]*WarningData
	tests       []TestTag
	handlers    []HandlerTag
	handlerData map[HandlerTag]*handlerData
	testNames   map[TestTag]string
}

func NewWarnings() *Warnings {
	return &Warnings{
		byTest:      make(map[TestTag]*Levels),
		byHandler:   make(map[HandlerTag]*Levels),
		byWarning:   make(map[string]*WarningData),
		tests:       make([]TestTag, 0),
		handlers:    make([]HandlerTag, 0),
		handlerData: make(map[HandlerTag]*handlerData),
		testNames:   make(map[TestTag]string),
	}
}

// HasHandler returns true if a handler is defined with the specified tag.
func (w *Warnings) HasHandler(handler HandlerTag) bool {
	_, found := w.byHandler[handler]
	return found
}

// ForHandler returns warnings for the specified handler grouped by level.
func (w *Warnings) ForHandler(handler HandlerTag) *Levels {
	return w.byHandler[handler]
}

// HandlerName returns the full name associated with a HandlerTag.
// If there is no full name the tag is returned.
func (w *Warnings) HandlerName(handler HandlerTag) string {
	if data, found := w.handlerData[handler]; found {
		return data.name
	} else {
		return string(handler)
	}
}

// HasHandlerSummary returns true if the specified handler has a summary text.
func (w *Warnings) HasHandlerSummary(handler HandlerTag) bool {
	if data, found := w.handlerData[handler]; found {
		return found && data.summary != ""
	}
	return false
}

// HandlerSummaryHTML returns summary text for the specified handler
// converted from Markdown to template.HTML.
func (w *Warnings) HandlerSummaryHTML(handler HandlerTag) template.HTML {
	if data, found := w.handlerData[handler]; found {
		return markdown.TemplateHTML(data.summary, true)
	} else {
		return ""
	}
}

// HasHandlerLinks returns true if the specified handler has reference links.
func (w *Warnings) HasHandlerLinks(handler HandlerTag) bool {
	if data, found := w.handlerData[handler]; found {
		return found && len(data.links) > 0
	}
	return false
}

// HandlerLinks returns handler reference links.
func (w *Warnings) HandlerLinks(handler HandlerTag) map[string]string {
	if data, found := w.handlerData[handler]; found {
		return data.links
	} else {
		return nil
	}
}

// HandlerTags returns an array of all handler names sorted alphabetically.
func (w *Warnings) HandlerTags() []HandlerTag {
	if len(w.handlers) < 1 {
		for handler := range w.byHandler {
			w.handlers = append(w.handlers, handler)
		}
		sort.Slice(w.handlers, func(i, j int) bool {
			return w.HandlerName(w.handlers[i]) < w.HandlerName(w.handlers[j])
		})
	}
	return w.handlers
}

// HandlerWarningCount returns the number of the specified warning associated with the specified handler.
func (w *Warnings) HandlerWarningCount(handler HandlerTag, warning *warning.Warning) uint {
	return w.byWarning[warning.Name].count[handler]
}

// HasTest returns true if a test is defined with the specified tag.
func (w *Warnings) HasTest(test TestTag) bool {
	_, found := w.byTest[test]
	return found
}

// ForTest returns warnings for the specified test grouped by level.
func (w *Warnings) ForTest(test TestTag) *Levels {
	return w.byTest[test]
}

// TestName returns the full name associated with a TestTag.
// If there is no full name the tag is returned.
func (w *Warnings) TestName(test TestTag) string {
	if name, found := w.testNames[test]; found {
		return name
	} else {
		return string(test)
	}
}

// TestTags returns an array of all handler names sorted alphabetically.
func (w *Warnings) TestTags() []TestTag {
	if len(w.tests) < 1 {
		for test := range w.byTest {
			w.tests = append(w.tests, test)
		}
		sort.Slice(w.tests, func(i, j int) bool {
			return w.TestName(w.tests[i]) < w.TestName(w.tests[j])
		})
	}
	return w.tests
}

// TestTagsForSource returns an array of all test names for the specified source sorted alphabetically.
func (w *Warnings) TestTagsForSource(source string) []TestTag {
	if !strings.HasSuffix(source, ":") {
		source += ":"
	}
	result := make([]TestTag, 0, len(w.TestTags())) // Causes w.tests to be loaded
	for _, test := range w.tests {
		if strings.HasPrefix(string(test), source) {
			result = append(result, test)
		}
	}
	return result
}

// -----------------------------------------------------------------------------

func (w *Warnings) findHandler(handler HandlerTag, level warning.Level, warningName string) *dataWarning {
	levels, ok := w.byHandler[handler]
	if !ok {
		levels = &Levels{
			lookup: make(map[string]*dataLevel),
			levels: make([]*dataLevel, 0),
		}
		w.byHandler[handler] = levels
	}
	return levels.findLevel(level, warningName)
}

func (w *Warnings) findTest(test TestTag, level warning.Level, warningName string) *dataWarning {
	levels, ok := w.byTest[test]
	if !ok {
		levels = &Levels{
			lookup: make(map[string]*dataLevel),
			levels: make([]*dataLevel, 0),
		}
		w.byTest[test] = levels
	}
	return levels.findLevel(level, warningName)
}

// FindWarning returns warning data for the specified warning name.
// If there is no warning data for that name a new WarningData record is created and cached
// which supports data collection during parsing.
//
// This method is public so that it can be used during template expansion in cmd/server.
func (w *Warnings) FindWarning(name string) *WarningData {
	data, found := w.byWarning[name]
	if !found {
		data = &WarningData{
			count:   make(map[HandlerTag]uint),
			hdlrMap: make(map[HandlerTag]bool),
			testMap: make(map[TestTag]bool),
		}
		w.byWarning[name] = data
	}
	return data
}

// -----------------------------------------------------------------------------

func (w *Warnings) getHandlerData(handler HandlerTag) *handlerData {
	data, found := w.handlerData[handler]
	if !found || data == nil {
		w.handlerData[handler] = &handlerData{}
	}
	return w.handlerData[handler]
}

// -----------------------------------------------------------------------------

// WarningData collects handler and test records for a specific warning.
//
// This type is public so that it can be used during template expansion in cmd/server.
// Methods for WarningData are largely public as well for the same reason.
type WarningData struct {
	count    map[HandlerTag]uint
	hdlrMap  map[HandlerTag]bool
	handlers []HandlerTag
	testMap  map[TestTag]bool
	tests    []TestTag
}

func (wd *WarningData) HasHandlers() bool {
	return len(wd.hdlrMap) > 0
}

func (wd *WarningData) Handlers() []HandlerTag {
	if wd.handlers == nil {
		wd.handlers = make([]HandlerTag, 0, len(wd.hdlrMap))
		for key := range wd.hdlrMap {
			wd.handlers = append(wd.handlers, key)
		}
		sort.Slice(wd.handlers, func(i, j int) bool { return wd.handlers[i].Name() < wd.handlers[j].Name() })
	}

	return wd.handlers
}

func (wd *WarningData) HasTests() bool {
	return len(wd.testMap) > 0
}

func (wd *WarningData) Tests() []TestTag {
	if wd.tests == nil {
		wd.tests = make([]TestTag, 0, len(wd.testMap))
		for key := range wd.testMap {
			wd.tests = append(wd.tests, key)
		}
		sort.Slice(wd.tests, func(i, j int) bool { return wd.tests[i].Name() < wd.tests[j].Name() })
	}

	return wd.tests
}

func (wd *WarningData) HasUsage() bool {
	return len(wd.hdlrMap) > 0 || len(wd.testMap) > 0
}

// -----------------------------------------------------------------------------

type Levels struct {
	lookup map[string]*dataLevel
	levels []*dataLevel
}

func (l *Levels) Levels() []*dataLevel {
	if len(l.levels) < 1 {
		l.levels = make([]*dataLevel, 0, len(warning.LevelOrder))
		for _, lvl := range warning.LevelOrder {
			if lv, ok := l.lookup[lvl.String()]; ok {
				l.levels = append(l.levels, lv)
			}
		}
	}
	return l.levels
}

func (l *Levels) findLevel(lvl warning.Level, warningName string) *dataWarning {
	lv, ok := l.lookup[lvl.String()]
	if !ok {
		lv = &dataLevel{
			level:  lvl,
			lookup: make(map[string]*dataWarning),
		}
		l.lookup[lvl.String()] = lv
	}
	return lv.findWarningGroup(warningName)
}

// -----------------------------------------------------------------------------

type dataLevel struct {
	level    warning.Level
	lookup   map[string]*dataWarning
	warnings []*dataWarning
}

func (l *dataLevel) Name() string {
	return l.level.String()
}

func (l *dataLevel) Warnings() []*dataWarning {
	if l.warnings == nil {
		l.warnings = make([]*dataWarning, 0, len(l.lookup))
		for _, w := range l.lookup {
			l.warnings = append(l.warnings, w)
		}
		sort.Slice(l.warnings, func(i, j int) bool {
			return l.warnings[i].warning.name < l.warnings[j].warning.name
		})
	}
	return l.warnings
}

func (l *dataLevel) findWarningGroup(warningName string) *dataWarning {
	grp, ok := l.lookup[warningName]
	if !ok {
		grp = &dataWarning{
			instances: make([]*dataInstance, 0),
		}
		grp.warning.name = warningName
		l.lookup[warningName] = grp
	}
	return grp
}

// -----------------------------------------------------------------------------

type dataWarning struct {
	warning struct {
		name, summary string
		instance      *warning.Warning
	}
	instances []*dataInstance
	sorted    bool
}

func (w *dataWarning) Name() string {
	return w.warning.name
}

func (w *dataWarning) HasDescription() bool {
	if w.warning.instance == nil {
		w.warning.instance = warning.ByName(w.warning.name)
	}
	return w.warning.instance.HasDescription()
}

func (w *dataWarning) Description() template.HTML {
	if w.warning.instance == nil {
		w.warning.instance = warning.ByName(w.warning.name)
	}
	return w.warning.instance.Description()
}

func (w *dataWarning) Summary() string {
	return w.warning.summary
}

func (w *dataWarning) AddInstance(instance *dataInstance) {
	w.instances = append(w.instances, instance)
}

func (w *dataWarning) Instances() []*dataInstance {
	if !w.sorted {
		sort.Slice(w.instances, func(i, j int) bool {
			return w.instances[i].name < w.instances[j].name
		})
		w.sorted = true
	}
	return w.instances
}

// -----------------------------------------------------------------------------

type dataInstance struct {
	tag    string
	source string
	name   string
	extra  string
	line   string // log in one line
	log    string // log indented
}

func (di *dataInstance) Extra() string {
	return di.extra
}

const srcRoot = "https://cs.opensource.google/go/go/+/master:src"

// test message: longer description (/snap/go/10535/src/testing/slogtest/slogtest.go:61)
var ptnCaller = regexp.MustCompile(`^\s*(.+?):(.+?)\s*\(.*/(testing/.*\.go):(\d+)\)`)

func (di *dataInstance) ExtraHTML() template.HTML {
	fix := func(lines []string) []string {
		for i, line := range lines {
			if matches := ptnCaller.FindStringSubmatch(line); len(matches) == 5 {
				// Error from slog/slogtest.
				line = fmt.Sprintf(
					"%s<br/><a href=\"%s/%s;l=%s\">%s</a>",
					matches[1], srcRoot, matches[3], matches[4], matches[2])
			}
			lines[i] = "<p>" + line + "</p>"
		}
		return lines
	}
	x := template.HTML(strings.Join(fix(strings.Split(di.extra, "\n")), "\n"))
	return x
}

func (di *dataInstance) Line() string {
	return di.line
}

func (di *dataInstance) HasLog() bool {
	return len(strings.Trim(di.log, " ")) > 0
}

func (di *dataInstance) Log() string {
	if di.log != "" {
		return di.log
	}
	return di.line
}

func (di *dataInstance) Name() string {
	return di.name
}

func (di *dataInstance) HasSource() bool {
	return len(strings.Trim(di.source, " ")) > 0
}

func (di *dataInstance) Source() string {
	return di.source
}

func (di *dataInstance) Tag() string {
	return di.tag
}
