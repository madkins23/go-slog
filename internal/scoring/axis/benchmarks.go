package axis

import (
	_ "embed"
	"html/template"
	"strconv"

	"github.com/madkins23/go-slog/internal/data"
	"github.com/madkins23/go-slog/internal/markdown"
	"github.com/madkins23/go-slog/internal/scoring/axis/bench"
	"github.com/madkins23/go-slog/internal/scoring/exhibit"
	"github.com/madkins23/go-slog/internal/scoring/score"
)

const BenchmarksType = "Benchmarks"

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

var _ score.Axis = &Benchmarks{}

type Benchmarks struct {
	handlerData map[data.HandlerTag]*bench.HandlerData
	benchWeight map[bench.Weight]uint
	exhibits    []score.Exhibit
	summaryHTML template.HTML
	BenchOptions
}

type BenchOptions struct {
	Name         string
	IncludeTests []data.TestTag
	ExcludeTests []data.TestTag
}

func NewBenchmarks(benchWeight map[bench.Weight]uint, summaryHTML template.HTML, options *BenchOptions) score.Axis {
	b := &Benchmarks{
		benchWeight: benchWeight,
		handlerData: make(map[data.HandlerTag]*bench.HandlerData),
		summaryHTML: summaryHTML,
	}
	if options != nil {
		b.BenchOptions = *options
	}
	return b
}

func (b *Benchmarks) Setup(benchMarks *data.Benchmarks, _ *data.Warnings) error {
	testTags := b.testTagMap(benchMarks.TestTags())
	// Calculate data ranges
	// Score 0: Original
	original := bench.NewOriginal(benchMarks, testTags, b.benchWeight)
	// Calculate test ranges used in calculating scores.
	original.MakeRanges()
	// Score 1 & 2
	ranges := make(map[data.TestTag]map[bench.Weight]bench.Range)
	for _, test := range benchMarks.TestTags() {
		if testTags[test] {
			// Score 1 & 2
			ranges[test] = map[bench.Weight]bench.Range{
				bench.Nanoseconds: bench.NewRangeFloat64(),
				bench.Allocations: bench.NewRangeUint64(),
				bench.AllocBytes:  bench.NewRangeUint64(),
			}
			for _, records := range benchMarks.HandlerRecordsFor(test) {
				ranges[test][bench.Nanoseconds].AddValueFloat64(records.NanosPerOp)
				ranges[test][bench.Allocations].AddValueUint64(records.MemAllocsPerOp)
				ranges[test][bench.AllocBytes].AddValueUint64(records.MemBytesPerOp)
			}
		}
	}
	original.CheckRanges(ranges)

	// Calculate scores using test ranges.
	for _, handler := range benchMarks.HandlerTags() {
		// Score 0: Original
		original.ResetForHandler()
		// Score 1 & 2
		handlerData := b.handlerData[handler]
		if handlerData == nil {
			handlerData = bench.NewHandlerData()
			b.handlerData[handler] = handlerData
		}
		for test, record := range benchMarks.ByHandler[handler] {
			if testTags[test] {
				// Score 0: Original
				original.HandlerTest(test, record)
				// Score 1 & 2
				rngTest := ranges[test]
				for _, weight := range bench.WeightOrder {
					rngTestWeight := rngTest[weight]
					if length := rngTestWeight.Length(); length > 0 {
						ranged := rngTestWeight.RangedValue(record.ItemValue(weight.Item()))
						// Score 1: Refactored original algorithm
						handlerData.ByTest(test).AddMultiple(ranged, b.benchWeight[weight])
						// Score 2: Newer algorithm rollup over BenchWeight subs
						handlerData.SubScore(weight).Add(ranged)
					}
				}
				// Score 1: Refactored original algorithm
				handlerData.Rollup(bench.OverTests).Add(handlerData.ByTest(test).Average())
				original.CheckTest(handlerData, test)
			}
		}
		// Score 0: Original
		handlerData.SetScore(score.Original, original.Score())
		// Score 1: Refactored original algorithm
		handlerData.SetScore(score.ByTest, handlerData.Rollup(bench.OverTests).Average())
		// Score 2: Newer algorithm rollup over BenchWeight subs
		for _, weight := range bench.WeightOrder {
			handlerData.Rollup(bench.OverData).AddMultiple(handlerData.SubScore(weight).Average(), b.benchWeight[weight])
		}
		handlerData.SetScore(score.ByData, handlerData.Rollup(bench.OverData).Average())
		original.CheckTotal(handlerData)
	}
	rows := make([][]string, 0, len(b.benchWeight))
	for _, weight := range bench.WeightOrder {
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

func (b *Benchmarks) HasTest(test data.TestTag) bool {
	for _, hdlr := range b.handlerData {
		if average := hdlr.ByTest(test); average != nil {
			if average.Count > 0 {
				return true
			}
		}
	}
	return false
}

func (b *Benchmarks) ScoreFor(handler data.HandlerTag) score.Value {
	return b.ScoreForType(handler, score.Default)
}

func (b *Benchmarks) ScoreForTest(handler data.HandlerTag, test data.TestTag) score.Value {
	return b.handlerData[handler].ByTest(test).Average()
}

func (b *Benchmarks) ScoreForType(handler data.HandlerTag, scoreType score.Type) score.Value {
	if scoreType == score.Default {
		scoreType = score.Original
	}
	return b.handlerData[handler].Score(scoreType)
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

func (b *Benchmarks) Type() string {
	return BenchmarksType
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
