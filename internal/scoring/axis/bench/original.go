package bench

import (
	"fmt"
	"log/slog"
	"math"

	"github.com/madkins23/go-slog/internal/data"
	"github.com/madkins23/go-slog/internal/scoring/score"
)

type testRange struct {
	allocLow, allocHigh uint64
	bytesLow, bytesHigh uint64
	nanosLow, nanosHigh float64
}

func (tr *testRange) String(bv Weight) string {
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

type Original struct {
	bench          *data.Benchmarks
	count, tests   uint
	collect, total score.Value
	ranges         map[data.TestTag]*testRange
	testTags       map[data.TestTag]bool
	weight         map[Weight]uint
}

func NewOriginal(bench *data.Benchmarks, tagMap map[data.TestTag]bool, weights map[Weight]uint) *Original {
	return &Original{
		bench:    bench,
		ranges:   make(map[data.TestTag]*testRange),
		testTags: tagMap,
		weight:   weights,
	}
}

func (o *Original) HandlerTest(test data.TestTag, record data.TestRecord) {
	o.collect = 0
	o.count = 0
	rngTest := o.ranges[test]
	if scoreRange := float64(rngTest.allocHigh - rngTest.allocLow); scoreRange > 0 {
		o.collect += score.Value(float64(o.weight[Allocations]) * 100.0 * float64(rngTest.allocHigh-record.MemAllocsPerOp) / scoreRange)
		o.count += o.weight[Allocations]
	}
	if scoreRange := float64(rngTest.bytesHigh - rngTest.bytesLow); scoreRange > 0 {
		o.collect += score.Value(float64(o.weight[AllocBytes]) * 100.0 * float64(rngTest.bytesHigh-record.MemBytesPerOp) / scoreRange)
		o.count += o.weight[AllocBytes]
	}
	if scoreRange := rngTest.nanosHigh - rngTest.nanosLow; scoreRange > 0 {
		o.collect += score.Value(float64(o.weight[Nanoseconds]) * 100.0 * (rngTest.nanosHigh - record.NanosPerOp) / scoreRange)
		o.count += o.weight[Nanoseconds]
	}
	o.total += o.collect / score.Value(o.count)
	o.tests++
}

func (o *Original) CheckRanges(ranges map[data.TestTag]map[Weight]Range) {
	for test := range ranges {
		if o.testTags[test] {
			for _, weight := range WeightOrder {
				original := o.ranges[test].String(weight)
				byOthers := ranges[test][weight].String()
				if byOthers != original {
					slog.Error("range comparison", "weight", weight,
						"Original", original,
						"ByOthers", byOthers)
				}
			}
		}
	}
}

func (o *Original) CheckTest(handlerData *HandlerData, test data.TestTag) {
	if !fuzzyEqual(o.collect, handlerData.byTest[test].Value) {
		slog.Error("collect comparison", "Original", o.collect, "by Test", handlerData.byTest[test].Value)
	}
	if o.count != handlerData.byTest[test].Count {
		slog.Error("count comparison", "Original", o.count, "by Test", handlerData.byTest[test].Count)
	}
}

func (o *Original) CheckTotal(handlerData *HandlerData) {
	if !fuzzyEqual(o.total.Round(), handlerData.Rollup(OverTests).Value) {
		slog.Error("total comparison",
			"Original", o.total.Round(),
			"by Test", handlerData.Rollup(OverTests).Value)
	}
	if o.tests != handlerData.Rollup(OverTests).Count {
		slog.Warn("count comparison",
			"Original", o.tests,
			"by Test", handlerData.Rollup(OverTests).Count)
	}
}

func (o *Original) MakeRanges() {
	for _, test := range o.bench.TestTags() {
		if o.testTags[test] {
			aRange := &testRange{
				allocLow: math.MaxUint64,
				bytesLow: math.MaxUint64,
				nanosLow: math.MaxFloat64,
			}
			for _, records := range o.bench.HandlerRecordsFor(test) {
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
			o.ranges[test] = aRange
		}
	}
}

func (o *Original) ResetForHandler() {
	o.tests = 0
	o.total = 0
}

func (o *Original) Score() score.Value {
	return o.total.Round() / score.Value(o.tests)
}

// -----------------------------------------------------------------------------

func fuzzyEqual(a, b score.Value) bool {
	const epsilon = 0.00000001
	return math.Abs(float64(a-b)) < epsilon
}
