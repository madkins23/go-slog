package warning

import (
	"flag"
	"sort"

	"github.com/madkins23/go-slog/warning"
)

var verifyFile = flag.String("verify", "", "Load verification data from path (optional)")

// -----------------------------------------------------------------------------

// TestTag is a unique name for a Benchmark or Verification test.
// The type is an alias for string so that types can't be confused.
type TestTag string

// HandlerTag is a unique name for a slog handler.
// The type is an alias for string so that types can't be confused.
type HandlerTag string

// -----------------------------------------------------------------------------

// Data encapsulates benchmark records by BenchmarkName and HandlerTag.
type Data struct {
	byTest       map[TestTag]*Levels
	byHandler    map[HandlerTag]*Levels
	tests        []TestTag
	handlers     []HandlerTag
	handlerNames map[HandlerTag]string
	testNames    map[TestTag]string
}

func NewData() *Data {
	return &Data{
		byTest:       make(map[TestTag]*Levels),
		byHandler:    make(map[HandlerTag]*Levels),
		tests:        make([]TestTag, 0),
		handlers:     make([]HandlerTag, 0),
		handlerNames: make(map[HandlerTag]string),
		testNames:    make(map[TestTag]string),
	}
}

func (d *Data) HasTest(test TestTag) bool {
	_, found := d.byTest[test]
	return found
}

func (d *Data) HasHandler(handler HandlerTag) bool {
	_, found := d.byHandler[handler]
	return found
}

func (d *Data) ForTest(test TestTag) *Levels {
	return d.byTest[test]
}

func (d *Data) ForHandler(handler HandlerTag) *Levels {
	return d.byHandler[handler]
}

// HandlerName returns the full name associated with a HandlerTag.
// If there is no full name the tag is returned.
func (d *Data) HandlerName(handler HandlerTag) string {
	if name, found := d.handlerNames[handler]; found {
		return name
	} else {
		return string(handler)
	}
}

// HandlerTags returns an array of all handler names sorted alphabetically.
func (d *Data) HandlerTags() []HandlerTag {
	if len(d.handlers) < 1 {
		for handler := range d.byHandler {
			d.handlers = append(d.handlers, handler)
		}
		sort.Slice(d.handlers, func(i, j int) bool {
			return d.HandlerName(d.handlers[i]) < d.HandlerName(d.handlers[j])
		})
	}
	return d.handlers
}

func (d *Data) findHandler(handler HandlerTag, level warning.Level, warningName string) *dataWarning {
	levels, ok := d.byHandler[handler]
	if !ok {
		levels = &Levels{
			lookup: make(map[string]*dataLevel),
			levels: make([]*dataLevel, 0),
		}
		d.byHandler[handler] = levels
	}
	return levels.findLevel(level, warningName)
}

func (d *Data) findTest(test TestTag, level warning.Level, warningName string) *dataWarning {
	levels, ok := d.byTest[test]
	if !ok {
		levels = &Levels{
			lookup: make(map[string]*dataLevel),
			levels: make([]*dataLevel, 0),
		}
		d.byTest[test] = levels
	}
	return levels.findLevel(level, warningName)
}

// -----------------------------------------------------------------------------

type Levels struct {
	lookup map[string]*dataLevel
	levels []*dataLevel
}

func (l *Levels) Levels() []*dataLevel {
	if len(l.levels) < 1 {
		l.levels = make([]*dataLevel, 0, len(l.lookup))
		for _, lv := range l.lookup {
			l.levels = append(l.levels, lv)
		}
		sort.Slice(l.levels, func(i, j int) bool {
			return l.levels[i].name < l.levels[j].name
		})
	}
	return l.levels
}

func (l *Levels) findLevel(lvl warning.Level, warningName string) *dataWarning {
	lv, ok := l.lookup[lvl.String()]
	if !ok {
		lv = &dataLevel{
			name:   lvl.String(),
			lookup: make(map[string]*dataWarning),
		}
		l.lookup[lvl.String()] = lv
	}
	return lv.findWarningGroup(warningName)
}

// -----------------------------------------------------------------------------

type dataLevel struct {
	name     string
	lookup   map[string]*dataWarning
	warnings []*dataWarning
}

func (l *dataLevel) Name() string {
	return l.name
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
		name, description string
	}
	instances []*dataInstance
	sorted    bool
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
	name    string
	extra   string
	logLine string
}

// -----------------------------------------------------------------------------
