package axis

import (
	"html/template"

	"github.com/madkins23/go-slog/infra/warning"
	"github.com/madkins23/go-slog/internal/data"
	"github.com/madkins23/go-slog/internal/scoring/score"
)

var _ score.Axis = &Warnings{}

type Warnings struct {
	levelWeight map[warning.Level]uint
	warnScores  map[data.HandlerTag]score.Value
}

func NewWarnings(levelWeight map[warning.Level]uint) score.Axis {
	return &Warnings{
		levelWeight: levelWeight,
	}
}

func (w *Warnings) Initialize(_ *data.Benchmarks, warns *data.Warnings) error {
	var maxScore uint
	for _, level := range warning.LevelOrder {
		var count uint
		for _, warn := range warning.WarningsForLevel(level) {
			if wx, found := warns.ByWarning[warn.Name]; found {
				if len(wx.Count) > 0 {
					count++
				}
			}
		}
		maxScore += w.levelWeight[level] * count
	}
	testScores := make(map[data.HandlerTag]uint)
	for _, hdlr := range warns.HandlerTags() {
		var scoreWork uint
		for _, level := range warns.ByHandler[hdlr].Levels() {
			scoreWork += w.levelWeight[level.Level] * uint(len(level.Warnings()))
		}
		testScores[hdlr] = scoreWork
	}
	// The range for warning scores is zero to maxScore.
	w.warnScores = make(map[data.HandlerTag]score.Value)
	for _, hdlr := range warns.HandlerTags() {
		if maxScore == 0 {
			// If we're all the same (the score range is essentially zero) we all get 100%.
			w.warnScores[hdlr] = 100.0
		} else {
			w.warnScores[hdlr] = 100.0 * score.Value(maxScore-testScores[hdlr]) / score.Value(maxScore)
		}
	}
	return nil
}

func (w *Warnings) ColumnHeader() string {
	return "Warnings"
}

func (w *Warnings) HandlerScore(handler data.HandlerTag) score.Value {
	return w.warnScores[handler]
}

func (w *Warnings) Documentation() template.HTML {
	return "<strong>TBD: Warnings.Documentation</strong>"
}
