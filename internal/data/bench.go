package data

import (
	"flag"
	"log/slog"
	"math"
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

func (tr *TestRecord) IsEmpty() bool {
	return tr.Runs == 0
}

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

// Benchmarks encapsulates benchmark records by BenchmarkName and HandlerTag.
type Benchmarks struct {
	byTest       map[TestTag]HandlerRecords
	byHandler    map[HandlerTag]TestRecords
	tests        []TestTag
	handlers     []HandlerTag
	testNames    map[TestTag]string
	testCPUs     map[TestTag]uint64
	handlerNames map[HandlerTag]string
	warningText  []byte
	lookup       map[string]HandlerTag
	ranges       testRanges
}

func NewBenchmarks() *Benchmarks {
	return &Benchmarks{
		byTest:       make(map[TestTag]HandlerRecords),
		byHandler:    make(map[HandlerTag]TestRecords),
		testNames:    make(map[TestTag]string),
		testCPUs:     make(map[TestTag]uint64),
		handlerNames: make(map[HandlerTag]string),
	}
}

// -----------------------------------------------------------------------------

func (b *Benchmarks) HasHandler(tag HandlerTag) bool {
	_, found := b.byHandler[tag]
	return found
}

func (b *Benchmarks) HasTest(tag TestTag) bool {
	_, found := b.byTest[tag]
	return found
}

// HandlerLookup returns a map from handler names to handler tags.
// Capture relationship between handler name in benchmark function vs. Creator.
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
		for handler := range b.byHandler {
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
	return b.byHandler[handler]
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

// -----------------------------------------------------------------------------

// testRange collects high and low values for a given handler/test combination.
type testRange struct {
	allocLow, allocHigh uint64
	bytesLow, bytesHigh uint64
	nanosLow, nanosHigh float64
}

// testRanges maps test tags to TestRange objects.
type testRanges map[TestTag]*testRange

// TestScores maps test tags to scores for a handler.
type TestScores struct {
	Overall float64
	byTest  map[TestTag]float64
}

// ForTest returns the floating point score for the test
// (for the implied handler for which the TestScores object was created).
func (ts *TestScores) ForTest(test TestTag) float64 {
	return ts.byTest[test]
}

// -----------------------------------------------------------------------------

// testRange returns a testRange object for the specified test tag.
// If the test tag is bad the result will be nil.
func (b *Benchmarks) testRange(test TestTag) *testRange {
	if b.ranges == nil {
		b.ranges = make(testRanges)
		for _, test := range b.TestTags() {
			aRange := &testRange{
				allocLow: math.MaxUint64,
				bytesLow: math.MaxUint64,
				nanosLow: math.MaxFloat64,
			}
			for _, records := range b.HandlerRecords(test) {
				if records.MemAllocsPerOp > aRange.allocHigh {
					aRange.allocHigh = records.MemAllocsPerOp
				}
				if records.MemAllocsPerOp < aRange.allocLow {
					aRange.allocLow = records.MemAllocsPerOp
				}
				if records.MemBytesPerOp > aRange.bytesHigh {
					aRange.bytesHigh = records.MemBytesPerOp
				}
				if records.MemBytesPerOp < aRange.bytesLow {
					aRange.bytesLow = records.MemBytesPerOp
				}
				if records.NanosPerOp > aRange.nanosHigh {
					aRange.nanosHigh = records.NanosPerOp
				}
				if records.NanosPerOp < aRange.nanosLow {
					aRange.nanosLow = records.NanosPerOp
				}
			}
			b.ranges[test] = aRange
		}
	}
	return b.ranges[test]
}

// HandlerScore returns all the scores associated with the specified handler.
// There is a score for each test in which the handler participated and
// a single overall score averaged from the per-test scores.
func (b *Benchmarks) HandlerScore(handler HandlerTag) *TestScores {
	scores := &TestScores{
		byTest: make(map[TestTag]float64),
	}
	for test, record := range b.byHandler[handler] {
		rng := b.testRange(test)
		var collect float64
		var count uint
		if scoreRange := float64(rng.allocHigh - rng.allocLow); scoreRange > 0 {
			collect += 100.0 * float64(rng.allocHigh-record.MemAllocsPerOp) / scoreRange
			count++
		}
		if scoreRange := float64(rng.bytesHigh - rng.bytesLow); scoreRange > 0 {
			collect += 200.0 * float64(rng.bytesHigh-record.MemBytesPerOp) / scoreRange
			count += 2
		}
		if scoreRange := rng.nanosHigh - rng.nanosLow; scoreRange > 0 {
			collect += 300.0 * (rng.nanosHigh - record.NanosPerOp) / scoreRange
			count += 3
		}
		scores.byTest[test] = collect / float64(count)
	}
	var count uint
	for _, s := range scores.byTest {
		count++
		scores.Overall += s
	}
	scores.Overall /= float64(count)
	return scores
}
