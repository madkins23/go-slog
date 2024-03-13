package data

import (
	_ "embed"
	"html/template"
	"math"

	"github.com/madkins23/go-slog/internal/markdown"
	"github.com/madkins23/go-slog/internal/warning"
)

var _ Scores = &ScoreDefault{}

type Scores interface {
	Initialize(bench *Benchmarks, warnings *Warnings) error
	HandlerBenchScores(handler HandlerTag) *TestScores
	HandlerWarningScore(handler HandlerTag) float64
	DocBench() template.HTML
	DocWarning() template.HTML
}

func NewScoreKeeper() Scores {
	return &ScoreDefault{}
}

// =============================================================================

type ScoreDefault struct {
	benchScores map[HandlerTag]*TestScores
	warnScores  map[HandlerTag]float64
}

func (sd *ScoreDefault) Initialize(bench *Benchmarks, warnings *Warnings) error {
	sd.initBenchmarkScores(bench)
	sd.initWarningScores(warnings)
	return nil
}

//go:embed scores/benchmarks.md
var benchDoc string

// DocBench returns HTML documentation on the benchmark scoring algorithm converted from markdown source.
func (sd *ScoreDefault) DocBench() template.HTML {
	return markdown.TemplateHTML(benchDoc)
}

//go:embed scores/warnings.md
var warningDoc string

// DocWarning returns HTML documentation on the warning scoring algorithm converted from markdown source.
func (sd *ScoreDefault) DocWarning() template.HTML {
	return markdown.TemplateHTML(warningDoc)
}

// HandlerBenchScores returns all the scores associated with the specified handler.
// There is a score for each test in which the handler participated and
// a single overall score averaged from the per-test scores.
func (sd *ScoreDefault) HandlerBenchScores(handler HandlerTag) *TestScores {
	return sd.benchScores[handler]
}

// -----------------------------------------------------------------------------

func (sd *ScoreDefault) initBenchmarkScores(bench *Benchmarks) {
	ranges := make(testRanges)
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
		sd.benchScores[handler] = scores
	}
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

// =============================================================================

func (sd *ScoreDefault) initWarningScores(w *Warnings) {
	var maxScore uint64
	for _, level := range warning.LevelOrder {
		maxScore += scoreWeight[level] * uint64(len(warning.WarningsForLevel(level)))
	}
	testScores := make(map[HandlerTag]uint64)
	for _, hdlr := range w.HandlerTags() {
		var score uint64
		for _, level := range w.byHandler[hdlr].Levels() {
			score += scoreWeight[level.level] * uint64(len(level.Warnings()))
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

// -----------------------------------------------------------------------------

// scoreWeight has the multipliers for different warning levels.
var scoreWeight = map[warning.Level]uint64{
	warning.LevelRequired:  8,
	warning.LevelImplied:   4,
	warning.LevelSuggested: 2,
	warning.LevelAdmin:     1,
}

// HandlerWarningScore returns a single floating point value in the range 0..100.0
// that is purported to be the score of the handler based on its warnings.
func (sd *ScoreDefault) HandlerWarningScore(handler HandlerTag) float64 {
	return sd.warnScores[handler]
}
