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
	byTest        map[data.TestTag]*Average
	originalScore score.Value
	scores        map[score.Type]score.Value
	subScore      map[Weight]*Average
	rollup        map[RollOver]*Average
}

func NewHandlerData() *HandlerData {
	hd := &HandlerData{
		byTest:   make(map[data.TestTag]*Average),
		scores:   make(map[score.Type]score.Value),
		subScore: make(map[Weight]*Average),
		rollup:   make(map[RollOver]*Average),
	}
	for _, weight := range WeightOrder {
		hd.subScore[weight] = &Average{}
	}
	return hd
}

func (hd *HandlerData) ByTest(test data.TestTag) *Average {
	if hd.byTest[test] == nil {
		hd.byTest[test] = &Average{}
	}
	return hd.byTest[test]
}

func (hd *HandlerData) Rollup(over RollOver) *Average {
	if hd.rollup[over] == nil {
		hd.rollup[over] = &Average{}
	}
	return hd.rollup[over]
}

func (hd *HandlerData) Score(scoreType score.Type) score.Value {
	return hd.scores[scoreType]
}

func (hd *HandlerData) SetScore(scoreType score.Type, value score.Value) {
	hd.scores[scoreType] = value
}

func (hd *HandlerData) SubScore(weight Weight) *Average {
	if hd.subScore[weight] == nil {
		hd.subScore[weight] = &Average{}
	}
	return hd.subScore[weight]
}
