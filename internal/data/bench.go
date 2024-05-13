package data

import (
	"flag"
	"log/slog"
	"sort"
)

var benchFile = flag.String("bench", "", "Load benchmark data from path (optional)")

// -----------------------------------------------------------------------------

// TestRecords is a map of test records by test tag.
type TestRecords map[TestTag]TestRecord

// HandlerRecords is a map of test records by handler tag.
type HandlerRecords map[HandlerTag]TestRecord

// TestRecord represents a single benchmark/handler test result.
type TestRecord struct {
	Runs           uint64
	NanosPerOp     float64
	MemBytesPerOp  uint64
	MemAllocsPerOp uint64
	MbPerSec       float64
	GbPerSec       float64
	TbPerSec       float64
}

// IsEmpty returns true if the TestRecord has no data.
func (tr *TestRecord) IsEmpty() bool {
	return tr.Runs == 0
}

// ItemValue returns the numeric value for the specified item.
func (tr *TestRecord) ItemValue(item BenchItems) float64 {
	switch item {
	case Runs:
		return float64(tr.Runs)
	case Nanos:
		return tr.NanosPerOp
	case MemAllocs:
		return float64(tr.MemAllocsPerOp)
	case MemBytes:
		return float64(tr.MemBytesPerOp)
	case MbPerSec:
		return tr.MbPerSec
	case GbPerSec:
		return tr.GbPerSec
	case TbPerSec:
		return tr.TbPerSec
	default:
		slog.Warn("Unknown bench.TestItem", "item", item)
		return 0
	}
}

// -----------------------------------------------------------------------------

// Benchmarks encapsulates benchmark records by TestTag and HandlerTag.
type Benchmarks struct {
	byTest       map[TestTag]HandlerRecords
	ByHandler    map[HandlerTag]TestRecords
	tests        []TestTag
	handlers     []HandlerTag
	testNames    map[TestTag]string
	testCPUs     map[TestTag]uint64
	handlerNames map[HandlerTag]string
	warningText  []byte
	lookup       map[string]HandlerTag
}

func NewBenchmarks() *Benchmarks {
	return &Benchmarks{
		byTest:       make(map[TestTag]HandlerRecords),
		ByHandler:    make(map[HandlerTag]TestRecords),
		testNames:    make(map[TestTag]string),
		testCPUs:     make(map[TestTag]uint64),
		handlerNames: make(map[HandlerTag]string),
	}
}

// -----------------------------------------------------------------------------

// HasHandler returns true if a handler is defined with the specified tag.
func (b *Benchmarks) HasHandler(tag HandlerTag) bool {
	_, found := b.ByHandler[tag]
	return found
}

// HasTest returns true if a test is defined with the specified tag.
func (b *Benchmarks) HasTest(tag TestTag) bool {
	_, found := b.byTest[tag]
	return found
}

// HandlerLookup returns a map from handler names to handler tags,
// capturing the relationship between handler name in benchmark function vs. Creator.
// The result will be passed into Warnings.ParseWarningData(),
// where it will be used to convert handler names to tags.
// This makes all handler tags the same between Benchmarks and Warnings.
func (b *Benchmarks) HandlerLookup() map[string]HandlerTag {
	if b.lookup == nil {
		b.lookup = make(map[string]HandlerTag, len(b.handlerNames))
		for tag, name := range b.handlerNames {
			b.lookup[name] = tag
		}
	}
	return b.lookup
}

// HandlerName returns the full name associated with a HandlerTag.
// If there is no full name the tag is returned.
func (b *Benchmarks) HandlerName(handler HandlerTag) string {
	if name, found := b.handlerNames[handler]; found {
		return name
	} else {
		return string(handler)
	}
}

// HandlerRecords returns a map of HandlerTag to TestRecord for the specified benchmark.
func (b *Benchmarks) HandlerRecords(test TestTag) HandlerRecords {
	return b.byTest[test]
}

// HandlerTags returns an array of all handler names sorted alphabetically.
func (b *Benchmarks) HandlerTags() []HandlerTag {
	if b.handlers == nil {
		for handler := range b.ByHandler {
			b.handlers = append(b.handlers, handler)
		}
		sort.Slice(b.handlers, func(i, j int) bool {
			return b.HandlerName(b.handlers[i]) < b.HandlerName(b.handlers[j])
		})
	}
	return b.handlers
}

// TestName returns the full name associated with a TestTag.
// If there is no full name the tag is returned.
func (b *Benchmarks) TestName(test TestTag) string {
	if name, found := b.testNames[test]; found {
		return name
	} else {
		return string(test)
	}
}

// TestRecords returns a map of HandlerTag to TestRecord for the specified benchmark.
func (b *Benchmarks) TestRecords(handler HandlerTag) TestRecords {
	return b.ByHandler[handler]
}

// TestTags returns an array of all test names sorted alphabetically.
func (b *Benchmarks) TestTags() []TestTag {
	if b.tests == nil {
		for test := range b.byTest {
			b.tests = append(b.tests, test)
		}
		sort.Slice(b.tests, func(i, j int) bool {
			return b.TestName(b.tests[i]) < b.TestName(b.tests[j])
		})
	}
	return b.tests
}

// HasWarningText from end of benchmark run.
func (b *Benchmarks) HasWarningText() bool {
	return len(b.warningText) > 0
}

// WarningText from end of benchmark run.
func (b *Benchmarks) WarningText() []byte {
	return b.warningText
}
