package axis

import (
	_ "embed"
	"html/template"
	"log/slog"
	"strconv"

	"github.com/madkins23/go-slog/infra/warning"
	"github.com/madkins23/go-slog/internal/data"
	"github.com/madkins23/go-slog/internal/markdown"
	"github.com/madkins23/go-slog/internal/scoring/axis/common"
	"github.com/madkins23/go-slog/internal/scoring/axis/warn"
	"github.com/madkins23/go-slog/internal/scoring/exhibit"
	"github.com/madkins23/go-slog/internal/scoring/score"
)

const defaultScoreType = score.Original

var (
	//go:embed doc/warn-doc.md
	warnDocMD   string
	warnDocHTML template.HTML
)

func setupWarnings() error {
	warnDocHTML = markdown.TemplateHTML(warnDocMD, false)
	return nil
}

var _ score.Axis = &Warnings{}

type Warnings struct {
	score.AxisCore
	handlerData map[data.HandlerTag]*warn.HandlerData
	levelWeight map[warning.Level]uint
}

func NewWarnings(levelWeight map[warning.Level]uint, summaryHTML template.HTML) score.Axis {
	w := &Warnings{
		levelWeight: levelWeight,
		handlerData: make(map[data.HandlerTag]*warn.HandlerData),
	}
	w.SetSummary(summaryHTML)
	return w
}

func (w *Warnings) Setup(_ *data.Benchmarks, warns *data.Warnings) error {
	testTags := w.testTagMap(warns.TestTags())
	// Calculate data ranges
	// Do original calculation and newer calculations and compare them for error.
	original := warn.NewOriginal(warns, testTags, w.levelWeight)
	// Calculate test ranges used in calculating scores.
	original.MakeMaxScore()

	// Ranges for warning numbers are simple,
	// they go from 0 to the maximum number of unique warnings per level.
	var maxScore uint
	ranges := make(map[warning.Level]common.Range)
	for _, level := range warning.LevelOrder {
		var count uint
		for _, wrn := range warning.WarningsForLevel(level) {
			if wx, found := warns.ByWarning[wrn.Name]; found {
				if len(wx.Count) > 0 {
					count++
				}
			}
		}
		ranges[level] = &common.RangeUint64{}
		ranges[level].AddValueUint64(uint64(count))
		maxScore += w.levelWeight[level] * count
	}
	original.CheckMaxScore(maxScore)
	original.MakeWarnScores()
	for _, hdlr := range warns.HandlerTags() {
		if w.handlerData[hdlr] == nil {
			w.handlerData[hdlr] = warn.NewHandlerData()
		}
		hdlrData := w.handlerData[hdlr]
		// Get handler/level data.
		levels := warns.ForHandler(hdlr)
		var a score.Average
		for _, level := range levels.Levels() {
			count := level.Count()
			ranged := ranges[level.Level].RangedValue(float64(count))
			hdlrData.ByLevel(level.Level).Add(ranged)
			slog.Info("by Level", "hdlr", hdlr, "level", level, "count", count, "ranged", ranged, "after", hdlrData.ByLevel(level.Level).Average())
			a.AddMultiple(ranged, w.levelWeight[level.Level])
		}
		hdlrData.SetScore(score.ByData, a.Average())
		hdlrData.SetScore(score.Default, hdlrData.Score(score.ByData))
		hdlrData.SetScore(score.Original, original.Score(hdlr))
	}
	original.CheckByDataScores(w.handlerData)
	// Create Exhibits.
	rows := make([][]string, 0, len(w.levelWeight))
	for _, level := range warning.LevelOrder {
		if value, found := w.levelWeight[level]; found {
			rows = append(rows, []string{level.String(), strconv.Itoa(int(value))})
		}
	}
	w.AddExhibit(exhibit.NewTable("", []string{"Level", "Weight"}, rows))
	return nil
}

func (w *Warnings) Name() string {
	return "Warnings"
}

func (w *Warnings) HasTest(_ data.TestTag) bool {
	return true
}

func (w *Warnings) ScoreFor(handler data.HandlerTag) score.Value {
	return w.ScoreForType(handler, score.Default)
}

func (w *Warnings) ScoreForLevel(handler data.HandlerTag, level warning.Level) *score.Average {
	return w.handlerData[handler].ByLevel(level)
}

func (w *Warnings) ScoreForTest(handler data.HandlerTag, test data.TestTag) score.Value {
	return w.handlerData[handler].ByTest(test).Average()
}

func (w *Warnings) ScoreForType(handler data.HandlerTag, scoreType score.Type) score.Value {
	if scoreType == score.Default {
		scoreType = defaultScoreType
	}
	return w.handlerData[handler].Score(scoreType)
}

func (w *Warnings) Documentation() template.HTML {
	return warnDocHTML
}

// -----------------------------------------------------------------------------

// testTagMap returns a map from data.TestTag to bool
// in order to track which tests are in the scoring for this axis.
// The allTags argument specifies the list of all possible tests.
// Currently (2024-08-26) the map includes allTags tests with no modification.
func (w *Warnings) testTagMap(allTags []data.TestTag) map[data.TestTag]bool {
	ttm := make(map[data.TestTag]bool)
	for _, test := range allTags {
		ttm[test] = true
	}
	return ttm
}
