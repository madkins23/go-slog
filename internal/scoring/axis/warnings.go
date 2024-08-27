package axis

import (
	_ "embed"
	"html/template"
	"strconv"

	"github.com/madkins23/go-slog/infra/warning"
	"github.com/madkins23/go-slog/internal/data"
	"github.com/madkins23/go-slog/internal/markdown"
	"github.com/madkins23/go-slog/internal/scoring/axis/common"
	"github.com/madkins23/go-slog/internal/scoring/axis/warn"
	"github.com/madkins23/go-slog/internal/scoring/exhibit"
	"github.com/madkins23/go-slog/internal/scoring/score"
)

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
		var byData score.Average

		for _, level := range warning.LevelOrder {
			var ranged score.Value
			dataLevel := levels.Level(level)
			if dataLevel == nil {
				ranged = 100.0
			} else {
				count := dataLevel.Count()
				ranged = ranges[dataLevel.Level].RangedValue(float64(count))
			}
			hdlrData.ByLevel(level).Add(ranged)
			byData.AddMultiple(hdlrData.ByLevel(level).Average(), w.levelWeight[level])
		}
		// Set the scores.
		hdlrData.SetScore(score.ByData, byData.Average())
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

func (w *Warnings) LevelLog(handler data.HandlerTag, level warning.Level) []string {
	if hdlrData, found := w.handlerData[handler]; found {
		return hdlrData.LevelLog(level)
	}
	return []string{"No handler data for handler " + string(handler)}
}

func (w *Warnings) ScoreFor(handler data.HandlerTag) score.Value {
	return w.ScoreForType(handler, score.Default)
}

func (w *Warnings) ScoreForLevel(handler data.HandlerTag, level warning.Level) score.Value {
	result := w.handlerData[handler].ByLevel(level)
	if result.Count > 0 {
		return result.Average()
	}
	return 100
}

func (w *Warnings) ScoreForTest(handler data.HandlerTag, test data.TestTag) score.Value {
	result := w.handlerData[handler].ByTest(test)
	if result.Count > 0 {
		return result.Average()
	}
	return 100
}

func (w *Warnings) ScoreForType(handler data.HandlerTag, scoreType score.Type) score.Value {
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
