package axis

import (
	_ "embed"
	"html/template"
	"strconv"

	"github.com/madkins23/go-slog/infra/warning"
	"github.com/madkins23/go-slog/internal/data"
	"github.com/madkins23/go-slog/internal/markdown"
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
	levelWeight map[warning.Level]uint
	warnScores  map[data.HandlerTag]score.Value
	exhibits    []score.Exhibit
	summaryHTML template.HTML
}

func NewWarnings(levelWeight map[warning.Level]uint, summaryHTML template.HTML) score.Axis {
	return &Warnings{
		levelWeight: levelWeight,
		summaryHTML: summaryHTML,
	}
}

func (w *Warnings) Setup(_ *data.Benchmarks, warns *data.Warnings) error {
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
	rows := make([][]string, 0, len(w.levelWeight))
	for _, level := range warning.LevelOrder {
		if value, found := w.levelWeight[level]; found {
			rows = append(rows, []string{level.String(), strconv.Itoa(int(value))})
		}
	}
	w.exhibits = []score.Exhibit{exhibit.NewTable("", []string{"Level", "Weight"}, rows)}
	return nil
}

func (w *Warnings) Name() string {
	return "Warnings"
}

func (w *Warnings) ScoreFor(handler data.HandlerTag) score.Value {
	return w.warnScores[handler]
}

func (w *Warnings) Summary() template.HTML {
	return w.summaryHTML
}

func (w *Warnings) Exhibits() []score.Exhibit {
	if w.exhibits == nil {
		// TODO: Need some exhibits!
	}
	return w.exhibits
}

func (w *Warnings) Documentation() template.HTML {
	return warnDocHTML
}
