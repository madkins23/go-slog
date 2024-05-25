package axis

import (
	_ "embed"
	"html/template"
	"math"
	"strconv"

	"github.com/madkins23/go-slog/internal/data"
	"github.com/madkins23/go-slog/internal/markdown"
	"github.com/madkins23/go-slog/internal/scoring/exhibit"
	"github.com/madkins23/go-slog/internal/scoring/score"
)

var (
	//go:embed doc/benchmarks.md
	benchDocMD   string
	benchDocHTML template.HTML
)

func setupBenchmarks() error {
	benchDocHTML = markdown.TemplateHTML(benchDocMD, false)
	return nil
}

var _ score.Axis = &Benchmarks{}

type Benchmarks struct {
	name        string
	benchWeight map[BenchValue]uint
	benchScores map[data.HandlerTag]score.Value
	doc         template.HTML
	exhibits    []score.Exhibit
	includeTags []data.TestTag
	excludeTags []data.TestTag
}

func NewBenchmarks(name string, benchWeight map[BenchValue]uint, includeTags, excludeTags []data.TestTag) score.Axis {
	return &Benchmarks{
		name:        name,
		benchWeight: benchWeight,
		doc:         benchDocHTML,
		includeTags: includeTags,
		excludeTags: excludeTags,
	}
}

func (b *Benchmarks) Setup(bench *data.Benchmarks, _ *data.Warnings) error {
	// Calculate test ranges used in calculating scores.
	ranges := make(map[data.TestTag]*testRange)
	testTags := b.testTagMap(bench.TestTags())
	for _, test := range bench.TestTags() {
		if testTags[test] {
			aRange := &testRange{
				allocLow: math.MaxUint64,
				bytesLow: math.MaxUint64,
				nanosLow: math.MaxFloat64,
			}
			for _, records := range bench.HandlerRecords(test) {
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
			ranges[test] = aRange
		}
	}

	// Calculate scores using test ranges.
	b.benchScores = make(map[data.HandlerTag]score.Value)
	for _, handler := range bench.HandlerTags() {
		var tests uint
		var total score.Value
		for test, record := range bench.ByHandler[handler] {
			if testTags[test] {
				rng := ranges[test]
				var collect score.Value
				var count uint
				if scoreRange := float64(rng.allocHigh - rng.allocLow); scoreRange > 0 {
					collect += score.Value(float64(b.benchWeight[Allocations]) * 100.0 * float64(rng.allocHigh-record.MemAllocsPerOp) / scoreRange)
					count += b.benchWeight[Allocations]
				}
				if scoreRange := float64(rng.bytesHigh - rng.bytesLow); scoreRange > 0 {
					collect += score.Value(float64(b.benchWeight[AllocBytes]) * 100.0 * float64(rng.bytesHigh-record.MemBytesPerOp) / scoreRange)
					count += b.benchWeight[AllocBytes]
				}
				if scoreRange := rng.nanosHigh - rng.nanosLow; scoreRange > 0 {
					collect += score.Value(float64(b.benchWeight[Nanoseconds]) * 100.0 * (rng.nanosHigh - record.NanosPerOp) / scoreRange)
					count += b.benchWeight[Nanoseconds]
				}
				total += collect / score.Value(count)
				tests++
			}
		}
		b.benchScores[handler] = total.Round() / score.Value(tests)
	}
	rows := make([][]string, 0, len(b.benchWeight))
	for _, weight := range benchWeightOrder {
		if value, found := b.benchWeight[weight]; found {
			rows = append(rows, []string{string(weight), strconv.Itoa(int(value))})
		}
	}
	b.exhibits = []score.Exhibit{exhibit.NewTable("", []string{"Data", "Weight"}, rows)}
	if b.includeTags != nil {
		b.exhibits = append(b.exhibits, exhibit.NewList("Included", testTagsToStrings(b.includeTags)))
	}
	if b.excludeTags != nil {
		b.exhibits = append(b.exhibits, exhibit.NewList("Excluded", testTagsToStrings(b.excludeTags)))
	}
	return nil
}

func (b *Benchmarks) AxisTitle() string {
	return b.Name() + " Score"
}

func (b *Benchmarks) Name() string {
	if b.name != "" {
		return b.name
	}
	return "Benchmark"
}

func (b *Benchmarks) ExhibitCount() uint {
	return uint(len(b.exhibits))
}

func (b *Benchmarks) Exhibits() []score.Exhibit {
	return b.exhibits
}

func (b *Benchmarks) HandlerScore(handler data.HandlerTag) score.Value {
	return b.benchScores[handler]
}

func (b *Benchmarks) Documentation() template.HTML {
	return b.doc
}

// -----------------------------------------------------------------------------

func (b *Benchmarks) testTagMap(allTags []data.TestTag) map[data.TestTag]bool {
	include := b.includeTags
	if include == nil {
		include = allTags
	}
	ttm := make(map[data.TestTag]bool)
	if len(include) > 0 {
		for _, test := range include {
			ttm[test] = true
		}
	}
	if len(b.excludeTags) > 0 {
		for _, test := range b.excludeTags {
			ttm[test] = false
		}
	}
	return ttm
}

// -----------------------------------------------------------------------------

func testTagsToStrings(tags []data.TestTag) []string {
	result := make([]string, len(tags))
	for i, tag := range tags {
		result[i] = string(tag)
	}
	return result
}

// -----------------------------------------------------------------------------

type BenchValue string

const (
	Allocations BenchValue = "Allocations"
	AllocBytes  BenchValue = "Alloc Bytes"
	Nanoseconds BenchValue = "Nanoseconds"
)

var benchWeightOrder = []BenchValue{
	Nanoseconds,
	AllocBytes,
	Allocations,
}

// -----------------------------------------------------------------------------

// testRange collects high and low values for a given handler/test combination.
type testRange struct {
	allocLow, allocHigh uint64
	bytesLow, bytesHigh uint64
	nanosLow, nanosHigh float64
}
