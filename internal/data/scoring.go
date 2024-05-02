package data

import (
	_ "embed"
	"html/template"
	"math"

	warning2 "github.com/madkins23/go-slog/infra/warning"
	"github.com/madkins23/go-slog/internal/markdown"
)

// Scores defines the interface for scoring objects.
// This supports replacing the default scorekeeper with a better one.
type Scores interface {
	Initialize(bench *Benchmarks, warnings *Warnings) error
	HandlerBenchScores(handler HandlerTag) *TestScores
	HandlerWarningScore(handler HandlerTag) float64
	DocOverview() template.HTML
	DocBench() template.HTML
	DocWarning() template.HTML
	WeightBench() map[benchValue]uint
	WeightBenchOrder() []benchValue
	WeightWarning() map[warning2.Level]uint64
	// WeightWarningOrder unnecessary, use warnings.LevelOrder().
}

// NewScoreKeeper returns a default Scores object.
func NewScoreKeeper() Scores {
	return &ScoreDefault{}
}

// =============================================================================

var _ Scores = &ScoreDefault{}

// ScoreDefault is the default scoring object.
type ScoreDefault struct {
	benchScores map[HandlerTag]*TestScores
	warnScores  map[HandlerTag]float64
}

// Initialize a ScoreDefault object.
func (sd *ScoreDefault) Initialize(bench *Benchmarks, warnings *Warnings) error {
	sd.initBenchmarkScores(bench)
	sd.initWarningScores(warnings)
	return nil
}

//go:embed scores/overview.md
var overviewDoc string

// DocOverview returns the scoring overview document converted from Markdown source to template.HTML.
func (sd *ScoreDefault) DocOverview() template.HTML {
	return markdown.TemplateHTML(overviewDoc, false)
}

//go:embed scores/benchmarks.md
var benchDoc string

// DocBench returns HTML documentation on the benchmark scoring algorithm converted from Markdown source to template.HTML.
func (sd *ScoreDefault) DocBench() template.HTML {
	return markdown.TemplateHTML(benchDoc, false)
}

//go:embed scores/warnings.md
var warningDoc string

// DocWarning returns HTML documentation on the warning scoring algorithm converted from markdown source to template.HTML.
func (sd *ScoreDefault) DocWarning() template.HTML {
	return markdown.TemplateHTML(warningDoc, false)
}

// HandlerBenchScores returns all the scores associated with the specified handler.
// There is a score for each test in which the handler participated and
// a single overall score averaged from the per-test scores.
func (sd *ScoreDefault) HandlerBenchScores(handler HandlerTag) *TestScores {
	return sd.benchScores[handler]
}

// WeightBench returns a map of algorithm weights by bench "values"
// (i.e. allocations, bytes allocated, and nanoseconds per operation).
func (sd *ScoreDefault) WeightBench() map[benchValue]uint {
	return benchScoreWeight
}

// WeightBenchOrder returns an array of bench "value" names in the order
// they should be referenced in a tabular format.
func (sd *ScoreDefault) WeightBenchOrder() []benchValue {
	return benchScoreWeightOrder
}

// WeightWarning returns a map of algorithm weights by warning levels.
func (sd *ScoreDefault) WeightWarning() map[warning2.Level]uint64 {
	return warningScoreWeight
}

// =============================================================================

func (sd *ScoreDefault) initBenchmarkScores(bench *Benchmarks) {
	// Calculate test ranges used in calculating scores.
	ranges := make(map[TestTag]*testRange)
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
	sd.benchScores = make(map[HandlerTag]*TestScores)
	for _, handler := range bench.HandlerTags() {
		scores := &TestScores{
			byTest: make(map[TestTag]float64),
		}
		for test, record := range bench.byHandler[handler] {
			rng := ranges[test]
			var collect float64
			var count uint
			if scoreRange := float64(rng.allocHigh - rng.allocLow); scoreRange > 0 {
				collect += float64(benchScoreWeight[allocations]) * 100.0 * float64(rng.allocHigh-record.MemAllocsPerOp) / scoreRange
				count += benchScoreWeight[allocations]
			}
			if scoreRange := float64(rng.bytesHigh - rng.bytesLow); scoreRange > 0 {
				collect += float64(benchScoreWeight[allocBytes]) * 100.0 * float64(rng.bytesHigh-record.MemBytesPerOp) / scoreRange
				count += benchScoreWeight[allocBytes]
			}
			if scoreRange := rng.nanosHigh - rng.nanosLow; scoreRange > 0 {
				collect += float64(benchScoreWeight[nanoseconds]) * 100.0 * (rng.nanosHigh - record.NanosPerOp) / scoreRange
				count += benchScoreWeight[nanoseconds]
			}
			scores.byTest[test] = collect / float64(count)
		}
		var count uint
		for _, s := range scores.byTest {
			count++
			scores.Overall += s
		}
		scores.Overall /= float64(count)
		sd.benchScores[handler] = scores
	}
}

// -----------------------------------------------------------------------------

type benchValue string

const (
	allocations benchValue = "Allocations"
	allocBytes  benchValue = "Alloc Bytes"
	nanoseconds benchValue = "Nanoseconds"
)

// benchScoreWeight has the multipliers for different benchmark values.
var benchScoreWeight = map[benchValue]uint{
	allocations: 1,
	allocBytes:  2,
	nanoseconds: 3,
}

var benchScoreWeightOrder = []benchValue{
	nanoseconds,
	allocBytes,
	allocations,
}

// testRange collects high and low values for a given handler/test combination.
type testRange struct {
	allocLow, allocHigh uint64
	bytesLow, bytesHigh uint64
	nanosLow, nanosHigh float64
}

// TestScores maps test tags to scores for a handler.
type TestScores struct {
	// Overall score for a handler.
	Overall float64

	// Scores by test for a handler.
	byTest map[TestTag]float64
}

// ForTest returns the floating point score for the test
// (for the implied handler for which the TestScores object was created).
func (ts *TestScores) ForTest(test TestTag) float64 {
	return ts.byTest[test]
}

// =============================================================================

func (sd *ScoreDefault) initWarningScores(w *Warnings) {
	var maxScore uint64
	for _, level := range warning2.LevelOrder {
		var count uint64
		for _, warn := range warning2.WarningsForLevel(level) {
			if wx, found := w.byWarning[warn.Name]; found {
				if len(wx.count) > 0 {
					count++
				}
			}
		}
		maxScore += warningScoreWeight[level] * count
	}
	testScores := make(map[HandlerTag]uint64)
	for _, hdlr := range w.HandlerTags() {
		var score uint64
		for _, level := range w.byHandler[hdlr].Levels() {
			score += warningScoreWeight[level.level] * uint64(len(level.Warnings()))
		}
		testScores[hdlr] = score
	}
	// The range for warning scores is zero to maxScore.
	sd.warnScores = make(map[HandlerTag]float64)
	for _, hdlr := range w.HandlerTags() {
		if maxScore == 0 {
			// If we're all the same (the score range is essentially zero) we all get 100%.
			sd.warnScores[hdlr] = 100.0
		} else {
			sd.warnScores[hdlr] = 100.0 * float64(maxScore-testScores[hdlr]) / float64(maxScore)
		}
	}
}

// HandlerWarningScore returns a single floating point value in the range 0..100.0
// that is purported to be the score of the handler based on its warnings.
func (sd *ScoreDefault) HandlerWarningScore(handler HandlerTag) float64 {
	return sd.warnScores[handler]
}

// -----------------------------------------------------------------------------

// warningScoreWeight has the multipliers for different warning levels.
var warningScoreWeight = map[warning2.Level]uint64{
	warning2.LevelRequired:  8,
	warning2.LevelImplied:   4,
	warning2.LevelSuggested: 2,
	warning2.LevelAdmin:     1,
}
