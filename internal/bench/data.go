package bench

import (
	"flag"
	"log/slog"
	"regexp"
	"sort"
)

var benchFile = flag.String("bench", "", "Load benchmark data from path (optional)")

// -----------------------------------------------------------------------------

// TestTag is a unique name for a Benchmark test.
// The type is an alias for string so that types can't be confused.
type TestTag string

// HandlerTag is a unique name for a slog handler.
// The type is an alias for string so that types can't be confused.
type HandlerTag string

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

func (tr *TestRecord) IsEmpty() bool {
	return tr.Runs == 0
}

func (tr *TestRecord) ItemValue(item TestItems) float64 {
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

// Data encapsulates benchmark records by BenchmarkName and HandlerTag.
type Data struct {
	byTest       map[TestTag]HandlerRecords
	byHandler    map[HandlerTag]TestRecords
	tests        []TestTag
	handlers     []HandlerTag
	testNames    map[TestTag]string
	testCPUs     map[TestTag]uint64
	handlerNames map[HandlerTag]string
	warningText  []string
}

func NewData() *Data {
	return &Data{
		byTest:       make(map[TestTag]HandlerRecords),
		byHandler:    make(map[HandlerTag]TestRecords),
		testNames:    make(map[TestTag]string),
		testCPUs:     make(map[TestTag]uint64),
		handlerNames: make(map[HandlerTag]string),
	}
}

// -----------------------------------------------------------------------------

var (
	ptnWarnLine = regexp.MustCompile(`^# (.*)`)
	ptnDataLine = regexp.MustCompile(`^([^/]+)/Benchmark_([^-]+)-(\d+)\s+(\d+)\s+(\d+(?:\.\d+)?)\s+ns/op\b`)
	ptnAllocsOp = regexp.MustCompile(`\s(\d+)\s+allocs/op\b`)
	ptnBytesOp  = regexp.MustCompile(`\s(\d+)\s+B/op\b`)
	ptnMbSec    = regexp.MustCompile(`\s(\d+(?:\.\d+)?)\s+MB/s`)
)

// -----------------------------------------------------------------------------

// HandlerName returns the full name associated with a HandlerTag.
// If there is no full name the tag is returned.
func (d *Data) HandlerName(handler HandlerTag) string {
	if name, found := d.handlerNames[handler]; found {
		return name
	} else {
		return string(handler)
	}
}

// HandlerRecords returns a map of HandlerTag to TestRecord for the specified benchmark.
func (d *Data) HandlerRecords(test TestTag) HandlerRecords {
	return d.byTest[test]
}

// HandlerTags returns an array of all handler names sorted alphabetically.
func (d *Data) HandlerTags() []HandlerTag {
	if d.handlers == nil {
		for handler := range d.byHandler {
			d.handlers = append(d.handlers, handler)
		}
		sort.Slice(d.handlers, func(i, j int) bool {
			return d.HandlerName(d.handlers[i]) < d.HandlerName(d.handlers[j])
		})
	}
	return d.handlers
}

// TestName returns the full name associated with a TestTag.
// If there is no full name the tag is returned.
func (d *Data) TestName(test TestTag) string {
	if name, found := d.testNames[test]; found {
		return name
	} else {
		return string(test)
	}
}

// TestRecords returns a map of HandlerTag to TestRecord for the specified benchmark.
func (d *Data) TestRecords(handler HandlerTag) TestRecords {
	return d.byHandler[handler]
}

// TestTags returns an array of all test names sorted alphabetically.
func (d *Data) TestTags() []TestTag {
	if d.tests == nil {
		for test := range d.byTest {
			d.tests = append(d.tests, test)
		}
		sort.Slice(d.tests, func(i, j int) bool {
			return d.TestName(d.tests[i]) < d.TestName(d.tests[j])
		})
	}
	return d.tests
}

// HasWarningText from end of benchmark run.
func (d *Data) HasWarningText() bool {
	return len(d.warningText) > 0
}

// WarningText from end of benchmark run.
func (d *Data) WarningText() []string {
	return d.warningText
}
