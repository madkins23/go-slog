package axis

import (
	_ "embed"
	"fmt"
	"html/template"
	"log/slog"
	"math"
	"strconv"

	"github.com/madkins23/go-slog/internal/data"
	"github.com/madkins23/go-slog/internal/markdown"
	"github.com/madkins23/go-slog/internal/scoring/exhibit"
	"github.com/madkins23/go-slog/internal/scoring/score"
)

var (
	//go:embed doc/bench-doc.md
	benchDocMD   string
	benchDocHTML template.HTML
)

func setupBenchmarks() error {
	benchDocHTML = markdown.TemplateHTML(benchDocMD, false)
	return nil
}

// -----------------------------------------------------------------------------

type HandlerData struct {
	scores          map[score.Type]score.Value
	originalScore   score.Value
	byTest          map[data.TestTag]*BenchAverage
	rollupOverTests BenchAverage // original
	subScore        map[BenchWeight]*BenchAverage
	rollupOverData  BenchAverage
}

// -----------------------------------------------------------------------------

var _ score.Axis = &Benchmarks{}

type Benchmarks struct {
	handlerData map[data.HandlerTag]*HandlerData
	benchWeight map[BenchWeight]uint
	exhibits    []score.Exhibit
	summaryHTML template.HTML
	BenchOptions
}

type BenchOptions struct {
	Name         string
	IncludeTests []data.TestTag
	ExcludeTests []data.TestTag
}

func NewBenchmarks(benchWeight map[BenchWeight]uint, summaryHTML template.HTML, options *BenchOptions) score.Axis {
	b := &Benchmarks{
		benchWeight: benchWeight,
		handlerData: make(map[data.HandlerTag]*HandlerData),
		summaryHTML: summaryHTML,
	}
	if options != nil {
		b.BenchOptions = *options
	}
	return b
}

func (b *Benchmarks) Setup(bench *data.Benchmarks, _ *data.Warnings) error {
	// Calculate test ranges used in calculating scores.
	ranges := make(map[data.TestTag]map[BenchWeight]iRange)
	testTags := b.testTagMap(bench.TestTags())
	// Score 0: original algorithm
	xRanges := make(map[data.TestTag]*testRange)
	for _, test := range bench.TestTags() {
		if testTags[test] {
			// Score 0: original algorithm
			aRange := &testRange{
				allocLow: math.MaxUint64,
				bytesLow: math.MaxUint64,
				nanosLow: math.MaxFloat64,
			}
			for _, records := range bench.HandlerRecordsFor(test) {
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
			xRanges[test] = aRange

			// Score 1 & 2
			ranges[test] = map[BenchWeight]iRange{
				Nanoseconds: newRangeFloat64(),
				Allocations: newRangeUint64(),
				AllocBytes:  newRangeUint64(),
			}
			for _, records := range bench.HandlerRecordsFor(test) {
				ranges[test][Nanoseconds].addValueFloat64(records.NanosPerOp)
				ranges[test][Allocations].addValueUint64(records.MemAllocsPerOp)
				ranges[test][AllocBytes].addValueUint64(records.MemBytesPerOp)
			}
		}
	}

	// Calculate scores using test ranges.
	for _, handler := range bench.HandlerTags() {
		// Score 0: original algorithm
		var tests uint
		var total score.Value

		handlerData := b.handlerData[handler]
		if handlerData == nil {
			handlerData = &HandlerData{
				byTest:   make(map[data.TestTag]*BenchAverage),
				scores:   make(map[score.Type]score.Value),
				subScore: make(map[BenchWeight]*BenchAverage),
			}
			b.handlerData[handler] = handlerData
			for _, weight := range benchWeightOrder {
				handlerData.subScore[weight] = &BenchAverage{}
			}
		}
		for test, record := range bench.ByHandler[handler] {
			if testTags[test] {
				// Score 0: original algorithm
				var collect score.Value
				var count uint
				xRngTest := xRanges[test]
				if scoreRange := float64(xRngTest.allocHigh - xRngTest.allocLow); scoreRange > 0 {
					collect += score.Value(float64(b.benchWeight[Allocations]) * 100.0 * float64(xRngTest.allocHigh-record.MemAllocsPerOp) / scoreRange)
					count += b.benchWeight[Allocations]
				}
				if scoreRange := float64(xRngTest.bytesHigh - xRngTest.bytesLow); scoreRange > 0 {
					collect += score.Value(float64(b.benchWeight[AllocBytes]) * 100.0 * float64(xRngTest.bytesHigh-record.MemBytesPerOp) / scoreRange)
					count += b.benchWeight[AllocBytes]
				}
				if scoreRange := xRngTest.nanosHigh - xRngTest.nanosLow; scoreRange > 0 {
					collect += score.Value(float64(b.benchWeight[Nanoseconds]) * 100.0 * (xRngTest.nanosHigh - record.NanosPerOp) / scoreRange)
					count += b.benchWeight[Nanoseconds]
				}
				total += collect / score.Value(count)
				tests++

				// Score 1 & 2
				rngTest := ranges[test]
				handlerData.byTest[test] = &BenchAverage{}
				for _, weight := range benchWeightOrder {
					rngTestWeight := rngTest[weight]
					if length := rngTestWeight.length(); length > 0 {
						ranged := rngTestWeight.rangedValue(record.ItemValue(weight.Item()))
						// Score 1: Refactored original algorithm
						handlerData.byTest[test].addMultiple(ranged, b.benchWeight[weight])
						// Score 2: Newer algorithm rollup over BenchWeight subs
						handlerData.subScore[weight].add(ranged)
					}
				}
				// Score 1: Refactored original algorithm
				handlerData.rollupOverTests.add(handlerData.byTest[test].average())

				// TODO: remove?
				if !fuzzyEqual(collect, handlerData.byTest[test].Value) {
					slog.Error("collect comparison", "Original", collect, "by Test", handlerData.byTest[test].Value)
				}
				if count != handlerData.byTest[test].Count {
					slog.Error("collect comparison", "Original", count, "by Test", handlerData.byTest[test].Count)
				}
				for _, weight := range benchWeightOrder {
					original := xRngTest.String(weight)
					byOthers := rngTest[weight].String()
					if byOthers != original {
						slog.Error("range comparison", "weight", weight,
							"Original", original,
							"ByOthers", byOthers)
					}
				}
			}
		}
		// Score 0: original algorithm
		handlerData.scores[score.Original] = total.Round() / score.Value(tests)
		// Score 1: Refactored original algorithm
		handlerData.scores[score.ByTest] = handlerData.rollupOverTests.average()
		// Score 2: Newer algorithm rollup over BenchWeight subs
		for _, weight := range benchWeightOrder {
			handlerData.rollupOverData.addMultiple(handlerData.subScore[weight].average(), b.benchWeight[weight])
		}
		handlerData.scores[score.ByData] = handlerData.rollupOverData.average()

		// TODO: remove?
		if !fuzzyEqual(total.Round(), handlerData.rollupOverTests.Value) {
			slog.Error("total comparison",
				"Original", total.Round(),
				"by Test", handlerData.rollupOverTests.Value)
		}
		if tests != handlerData.rollupOverTests.Count {
			slog.Warn("count comparison",
				"Original", tests,
				"by Test", handlerData.rollupOverTests.Count)
		}
	}
	rows := make([][]string, 0, len(b.benchWeight))
	for _, weight := range benchWeightOrder {
		if value, found := b.benchWeight[weight]; found {
			rows = append(rows, []string{string(weight), strconv.Itoa(int(value))})
		}
	}
	b.exhibits = []score.Exhibit{exhibit.NewTable("", []string{"Data", "Weight"}, rows)}
	if b.IncludeTests != nil {
		b.exhibits = append(b.exhibits, exhibit.NewList("Included", testTagsToStrings(b.IncludeTests)))
	}
	if b.ExcludeTests != nil {
		b.exhibits = append(b.exhibits, exhibit.NewList("Excluded", testTagsToStrings(b.ExcludeTests)))
	}
	return nil
}

func (b *Benchmarks) Name() string {
	if b.BenchOptions.Name != "" {
		return b.BenchOptions.Name
	}
	return "Benchmarks"
}

func (b *Benchmarks) ScoreFor(handler data.HandlerTag) score.Value {
	return b.ScoreForType(handler, score.Original)
}

func (b *Benchmarks) ScoreForTest(handler data.HandlerTag, test data.TestTag) score.Value {
	return b.handlerData[handler].byTest[test].average()
}

func (b *Benchmarks) ScoreForType(handler data.HandlerTag, scoreType score.Type) score.Value {
	return b.handlerData[handler].scores[scoreType]
}

func (b *Benchmarks) Summary() template.HTML {
	return b.summaryHTML
}

func (b *Benchmarks) Exhibits() []score.Exhibit {
	return b.exhibits
}

func (b *Benchmarks) Documentation() template.HTML {
	return benchDocHTML
}

// -----------------------------------------------------------------------------

func (b *Benchmarks) testTagMap(allTags []data.TestTag) map[data.TestTag]bool {
	include := b.IncludeTests
	if include == nil {
		include = allTags
	}
	ttm := make(map[data.TestTag]bool)
	if len(include) > 0 {
		for _, test := range include {
			ttm[test] = true
		}
	}
	if len(b.ExcludeTests) > 0 {
		for _, test := range b.ExcludeTests {
			ttm[test] = false
		}
	}
	return ttm
}

// -----------------------------------------------------------------------------

func testTagsToStrings(tags []data.TestTag) []string {
	result := make([]string, len(tags))
	for i, tag := range tags {
		result[i] = tag.Name()
	}
	return result
}

// -----------------------------------------------------------------------------

type BenchWeight string

const (
	Allocations BenchWeight = "Allocations"
	AllocBytes  BenchWeight = "Alloc Bytes"
	Nanoseconds BenchWeight = "Nanoseconds"
)

var benchWeightOrder = []BenchWeight{
	Nanoseconds,
	AllocBytes,
	Allocations,
}

func (bw BenchWeight) Item() data.BenchItems {
	switch bw {
	case Allocations:
		return data.MemAllocs
	case AllocBytes:
		return data.MemBytes
	case Nanoseconds:
		return data.Nanos
	default:
		return 0.0
	}
}

// -----------------------------------------------------------------------------

type BenchAverage struct {
	Value score.Value
	Count uint
}

func (ba *BenchAverage) add(v score.Value) *BenchAverage {
	ba.Value += v
	ba.Count++
	return ba
}

func (ba *BenchAverage) addMultiple(v score.Value, multiple uint) *BenchAverage {
	ba.Value += v * score.Value(multiple)
	ba.Count += multiple
	return ba
}

func (ba *BenchAverage) average() score.Value {
	return ba.Value.Round() / score.Value(ba.Count)
}

// -----------------------------------------------------------------------------
// Score 0: original algorithm

type testRange struct {
	allocLow, allocHigh uint64
	bytesLow, bytesHigh uint64
	nanosLow, nanosHigh float64
}

func (tr *testRange) String(bv BenchWeight) string {
	switch bv {
	case Allocations:
		return fmt.Sprintf("%0d -> %0d", tr.allocLow, tr.allocHigh)
	case AllocBytes:
		return fmt.Sprintf("%0d -> %0d", tr.bytesLow, tr.bytesHigh)
	case Nanoseconds:
		return fmt.Sprintf("%0.2f -> %0.2f", tr.nanosLow, tr.nanosHigh)
	default:
		return "<unknown:" + string(bv) + ">"
	}
}

// -----------------------------------------------------------------------------

type iRange interface {
	addValueUint64(val uint64)
	addValueFloat64(val float64)
	length() float64
	rangedValue(from float64) score.Value
	String() string
}

// -----------------------------------------------------------------------------

var _ iRange = &RangeFloat64{}

type RangeFloat64 struct {
	low, high float64
}

func newRangeFloat64() *RangeFloat64 {
	return &RangeFloat64{
		low:  math.MaxFloat64,
		high: 0.0,
	}
}

func (r *RangeFloat64) addValueUint64(val uint64) {
	r.addValueFloat64(float64(val))
}

func (r *RangeFloat64) addValueFloat64(val float64) {
	if val < r.low {
		r.low = val
	}
	if val > r.high {
		r.high = val
	}
}

func (r *RangeFloat64) length() float64 {
	return r.high - r.low
}

func (r *RangeFloat64) rangedValue(from float64) score.Value {
	return score.Value(100.0 * (r.high - from) / r.length())
}

func (r *RangeFloat64) String() string {
	return fmt.Sprintf("%0.2f -> %0.2f", r.low, r.high)
}

// -----------------------------------------------------------------------------

var _ iRange = &RangeUint64{}

type RangeUint64 struct {
	low, high uint64
}

func newRangeUint64() *RangeUint64 {
	return &RangeUint64{
		low:  math.MaxUint64,
		high: 0,
	}
}

func (r *RangeUint64) addValueFloat64(val float64) {
	r.addValueUint64(uint64(val))
}

func (r *RangeUint64) addValueUint64(val uint64) {
	if val < r.low {
		r.low = val
	}
	if val > r.high {
		r.high = val
	}
}

func (r *RangeUint64) length() float64 {
	return float64(r.high - r.low)
}

func (r *RangeUint64) rangedValue(from float64) score.Value {
	return score.Value(100.0 * (float64(r.high) - from) / r.length())
}

func (r *RangeUint64) String() string {
	return fmt.Sprintf("%0d -> %0d", r.low, r.high)
}

// -----------------------------------------------------------------------------

func fuzzyEqual(a, b score.Value) bool {
	const epsilon = 0.000000001
	return math.Abs(float64(a-b)) < epsilon
}
