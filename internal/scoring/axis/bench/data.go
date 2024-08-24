package bench

import (
	"github.com/madkins23/go-slog/internal/data"
	"github.com/madkins23/go-slog/internal/scoring/score"
)

// -----------------------------------------------------------------------------

type RollOver uint8

const (
	OverData = iota
	OverTests
)

type HandlerData struct {
	byTest        map[data.TestTag]*score.Average
	originalScore score.Value
	scores        map[score.Type]score.Value
	subScore      map[Weight]*score.Average
	rollup        map[RollOver]*score.Average
}

func NewHandlerData() *HandlerData {
	hd := &HandlerData{
		byTest:   make(map[data.TestTag]*score.Average),
		scores:   make(map[score.Type]score.Value),
		subScore: make(map[Weight]*score.Average),
		rollup:   make(map[RollOver]*score.Average),
	}
	for _, weight := range WeightOrder {
		hd.subScore[weight] = &score.Average{}
	}
	return hd
}

func (hd *HandlerData) ByTest(test data.TestTag) *score.Average {
	if hd.byTest[test] == nil {
		hd.byTest[test] = &score.Average{}
	}
	return hd.byTest[test]
}

func (hd *HandlerData) Rollup(over RollOver) *score.Average {
	if hd.rollup[over] == nil {
		hd.rollup[over] = &score.Average{}
	}
	return hd.rollup[over]
}

func (hd *HandlerData) Score(scoreType score.Type) score.Value {
	return hd.scores[scoreType]
}

func (hd *HandlerData) SetScore(scoreType score.Type, value score.Value) {
	hd.scores[scoreType] = value
}

func (hd *HandlerData) SubScore(weight Weight) *score.Average {
	if hd.subScore[weight] == nil {
		hd.subScore[weight] = &score.Average{}
	}
	return hd.subScore[weight]
}
