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
	benchWeight map[BenchValue]uint
	benchScores map[data.HandlerTag]*testScores
	doc         template.HTML
	exhibits    []score.Exhibit
}

func NewBenchmarks(benchWeight map[BenchValue]uint) score.Axis {
	return &Benchmarks{
		benchWeight: benchWeight,
		doc:         benchDocHTML,
	}
}

func (b *Benchmarks) Setup(bench *data.Benchmarks, _ *data.Warnings) error {
	// Calculate test ranges used in calculating scores.
	ranges := make(map[data.TestTag]*testRange)
	for _, test := range bench.TestTags() {
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

	// Calculate scores using test ranges.
	b.benchScores = make(map[data.HandlerTag]*testScores)
	for _, handler := range bench.HandlerTags() {
		scores := &testScores{
			byTest: make(map[data.TestTag]score.Value),
		}
		for test, record := range bench.ByHandler[handler] {
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
			scores.byTest[test] = collect / score.Value(count)
		}
		var count uint
		for _, s := range scores.byTest {
			count++
			scores.Overall += s
		}
		scores.Overall /= score.Value(count)
		b.benchScores[handler] = scores
	}
	rows := make([][]string, 0, len(b.benchWeight))
	for name, value := range b.benchWeight {
		rows = append(rows, []string{string(name), strconv.Itoa(int(value))})
	}
	b.exhibits = []score.Exhibit{exhibit.NewTable("", []string{"Data", "Weight"}, rows)}
	return nil
}

func (b *Benchmarks) AxisTitle() string {
	return b.ColumnHeader() + " Score"
}

func (b *Benchmarks) ColumnHeader() string {
	return "Benchmark"
}

func (b *Benchmarks) ExhibitCount() uint {
	return uint(len(b.exhibits))
}

func (b *Benchmarks) Exhibits() []score.Exhibit {
	return b.exhibits
}

func (b *Benchmarks) HandlerScore(handler data.HandlerTag) score.Value {
	return b.benchScores[handler].Overall
}

func (b *Benchmarks) Documentation() template.HTML {
	return b.doc
}

// -----------------------------------------------------------------------------

type BenchValue string

const (
	Allocations BenchValue = "Allocations"
	AllocBytes  BenchValue = "Alloc Bytes"
	Nanoseconds BenchValue = "Nanoseconds"
)

var _ = []BenchValue{
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

// testScores maps test tags to scores for a handler.
// Deprecated: Get rid of testScores now that the per-test scores are not used.
// TODO: Get rid of this!
type testScores struct {
	// Overall score for a handler.
	Overall score.Value

	// Scores by test for a handler.
	byTest map[data.TestTag]score.Value
}
