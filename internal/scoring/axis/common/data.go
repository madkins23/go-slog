package common

import (
	"github.com/madkins23/go-slog/internal/data"
	"github.com/madkins23/go-slog/internal/scoring/score"
)

type HandlerData struct {
	byTest map[data.TestTag]*score.Average
	scores map[score.Type]score.Value
}

func NewHandlerData() *HandlerData {
	return &HandlerData{
		byTest: make(map[data.TestTag]*score.Average),
		scores: make(map[score.Type]score.Value),
	}
}

func (hd *HandlerData) ByTest(test data.TestTag) *score.Average {
	if hd.byTest[test] == nil {
		hd.byTest[test] = &score.Average{}
	}
	return hd.byTest[test]
}

func (hd *HandlerData) Score(scoreType score.Type) score.Value {
	return hd.scores[scoreType]
}

func (hd *HandlerData) SetScore(scoreType score.Type, value score.Value) {
	hd.scores[scoreType] = value
}
