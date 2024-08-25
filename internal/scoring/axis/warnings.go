package axis

import (
	_ "embed"
	"html/template"
	"log/slog"
	"strconv"

	"github.com/madkins23/go-slog/infra/warning"
	"github.com/madkins23/go-slog/internal/data"
	"github.com/madkins23/go-slog/internal/markdown"
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

	warnScores map[data.HandlerTag]score.Value
}

func NewWarnings(levelWeight map[warning.Level]uint, summaryHTML template.HTML) score.Axis {
	w := &Warnings{
		levelWeight: levelWeight,
		warnScores:  make(map[data.HandlerTag]score.Value),
	}
	w.SetSummary(summaryHTML)
	return w
}

func (w *Warnings) Setup(_ *data.Benchmarks, warns *data.Warnings) error {
	var totalScore uint
	for _, level := range warning.LevelOrder {
		var count uint
		for _, wrn := range warning.WarningsForLevel(level) {
			if wx, found := warns.ByWarning[wrn.Name]; found {
				if len(wx.Count) > 0 {
					count++
				}
			}
		}
		totalScore += w.levelWeight[level] * count
	}
	testScores := make(map[data.HandlerTag]uint)
	for _, hdlr := range warns.HandlerTags() {
		var scoreWork uint
		for _, level := range warns.ByHandler[hdlr].Levels() {
			scoreWork += w.levelWeight[level.Level] * uint(len(level.Warnings()))
		}
		testScores[hdlr] = scoreWork
	}
	// The range for warning scores is zero to totalScore.
	for _, hdlr := range warns.HandlerTags() {
		if totalScore == 0 {
			// If we're all the same (the score range is essentially zero) we all get 100%.
			w.warnScores[hdlr] = 100.0
		} else {
			w.warnScores[hdlr] = 100.0 * score.Value(totalScore-testScores[hdlr]) / score.Value(totalScore)
		}
	}
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

func (w *Warnings) ScoreForTest(handler data.HandlerTag, test data.TestTag) score.Value {
	slog.Warn("made up data", "func", "ScoreForTest", "handler", handler, "test", test)
	return 0.0
}

func (w *Warnings) ScoreForType(handler data.HandlerTag, scoreType score.Type) score.Value {
	if scoreType == score.Default {
		scoreType = score.Original
	}
	return w.warnScores[handler]
}

func (w *Warnings) Documentation() template.HTML {
	return warnDocHTML
}
