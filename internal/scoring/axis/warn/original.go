package warn

import (
	"log/slog"

	"github.com/madkins23/go-slog/infra/warning"
	"github.com/madkins23/go-slog/internal/data"
	"github.com/madkins23/go-slog/internal/scoring/axis/common"
	"github.com/madkins23/go-slog/internal/scoring/score"
)

// -----------------------------------------------------------------------------

// Original contains the original warning score calculations.
// The newer calculations work differently (more efficiently).
// This code was kept so that the calculations could be compared.
// In the fullness of time it may be removed.
type Original struct {
	warns          *data.Warnings
	count, tests   uint
	collect, total score.Value
	maxScore       uint
	testTags       map[data.TestTag]bool
	weight         map[warning.Level]uint
	testScores     map[data.HandlerTag]uint
	warnScores     map[data.HandlerTag]score.Value
}

// NewOriginal returns a new Original object.
func NewOriginal(warns *data.Warnings, tagMap map[data.TestTag]bool, levels map[warning.Level]uint) *Original {
	return &Original{
		warns: warns,
		//ranges:   make(map[data.TestTag]*testRange),
		testTags:   tagMap,
		weight:     levels,
		testScores: make(map[data.HandlerTag]uint),
	}
}

func (o *Original) CheckMaxScore(maxScore uint) {
	if o.maxScore != maxScore {
		slog.Error("maxScore comparison", "original", o.maxScore, "new", maxScore)
	}
}

func (o *Original) MakeMaxScore() {
	for _, level := range warning.LevelOrder {
		var count uint
		for _, warn := range warning.WarningsForLevel(level) {
			if wx, found := o.warns.ByWarning[warn.Name]; found {
				if len(wx.Count) > 0 {
					count++
				}
			}
		}
		o.maxScore += o.weight[level] * count
	}
}

func (o *Original) MakeWarnScores() {
	for _, hdlr := range o.warns.HandlerTags() {
		var scoreWork uint
		for _, level := range o.warns.ByHandler[hdlr].Levels() {
			scoreWork += o.weight[level.Level] * uint(len(level.Warnings()))
		}
		o.testScores[hdlr] = scoreWork
	}

	o.warnScores = make(map[data.HandlerTag]score.Value)
	for _, hdlr := range o.warns.HandlerTags() {
		if o.maxScore == 0 {
			// If we're all the same (the score range is essentially zero) we all get 100%.
			o.warnScores[hdlr] = 100.0
		} else {
			o.warnScores[hdlr] = 100.0 * score.Value(o.maxScore-o.testScores[hdlr]) / score.Value(o.maxScore)
		}
	}
}

func (o *Original) Score(hdlr data.HandlerTag) score.Value {
	return o.warnScores[hdlr]
}

// CheckByDataScores checks the calculated score value for a handler
// against the value from the original calculation.
// If the values differ by too much an error is logged.
func (o *Original) CheckByDataScores(handlerData map[data.HandlerTag]*HandlerData) {
	for _, hdlr := range o.warns.HandlerTags() {
		original := o.warnScores[hdlr]
		byData := handlerData[hdlr].Score(score.ByData)
		if !common.FuzzyEqual(original, byData) {
			slog.Error("warn score comparison", "Original", original, "by Data", byData)
		}
	}
}
