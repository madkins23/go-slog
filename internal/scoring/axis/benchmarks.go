package axis

import (
	_ "embed"
	"html/template"
	"strconv"

	"github.com/madkins23/go-slog/internal/data"
	"github.com/madkins23/go-slog/internal/markdown"
	"github.com/madkins23/go-slog/internal/scoring/axis/bench"
	"github.com/madkins23/go-slog/internal/scoring/axis/common"
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
	// Do original calculation and newer calculations and compare them for error.
	original := bench.NewOriginal(benchMarks, testTags, b.benchWeight)
	original.MakeRanges()
	// Data ranges for newer calculations.
	ranges := make(map[data.TestTag]map[bench.Weight]common.Range)
	for _, test := range benchMarks.TestTags() {
		if testTags[test] {
			// Data ranges for newer calculations.
			ranges[test] = map[bench.Weight]common.Range{
				bench.Nanoseconds: common.NewRangeFloat64(),
				bench.Allocations: common.NewRangeUint64(),
				bench.AllocBytes:  common.NewRangeUint64(),
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
		// Original calculation.
		original.ResetForHandler()
		// Data ranges for newer calculations.
		handlerData := b.handlerData[handler]
		if handlerData == nil {
			b.handlerData[handler] = bench.NewHandlerData()
			handlerData = b.handlerData[handler]
		}
		for test, record := range benchMarks.ByHandler[handler] {
			if testTags[test] {
				// Original calculation.
				original.HandlerTest(test, record)
				// Data ranges for newer calculations.
				rngTest := ranges[test]
				for _, weight := range bench.WeightOrder {
					rngTestWeight := rngTest[weight]
					if length := rngTestWeight.Length(); length > 0 {
						ranged := rngTestWeight.RangedValue(record.ItemValue(weight.Item()))
						// Refactored original algorithm for byTest
						handlerData.ByTest(test).AddMultiple(ranged, b.benchWeight[weight])
						// Newer algorithm rollup over BenchWeight subs
						handlerData.SubScore(weight).Add(ranged)
					}
				}
				// Newer algorithm.
				handlerData.Rollup(bench.OverTests).Add(handlerData.ByTest(test).Average())
				original.CheckTest(handlerData, test)
			}
		}
		// Original calculation.
		handlerData.SetScore(score.Original, original.Score())
		// Refactored original algorithm for byTest
		handlerData.SetScore(score.ByTest, handlerData.Rollup(bench.OverTests).Average())
		// Newer algorithm rollup over BenchWeight subs
		for _, weight := range bench.WeightOrder {
			handlerData.Rollup(bench.OverData).AddMultiple(handlerData.SubScore(weight).Average(), b.benchWeight[weight])
		}
		handlerData.SetScore(score.ByData, handlerData.Rollup(bench.OverData).Average())
		original.CheckTotal(handlerData)
		handlerData.SetScore(score.Default,
			(handlerData.Score(score.ByData)+handlerData.Score(score.ByTest))/2.0)
	}
	// Create Exhibits.
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

// -----------------------------------------------------------------------------

// testTagMap returns a map from data.TestTag to bool
// in order to track which tests are in the scoring for this axis.
// The allTags argument specifies the list of all possible tests.
// Options IncludeTests and ExcludeTests modify this list.
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
