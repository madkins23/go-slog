package bench

import (
	"github.com/madkins23/go-slog/internal/scoring/axis/common"
	"github.com/madkins23/go-slog/internal/scoring/score"
)

// -----------------------------------------------------------------------------

type RollOver uint8

const (
	OverData = iota
	OverTests
)

type HandlerData struct {
	*common.HandlerData
	subScore map[Weight]*score.Average
	rollup   map[RollOver]*score.Average
}

func NewHandlerData() *HandlerData {
	hd := &HandlerData{
		HandlerData: common.NewHandlerData(),
		subScore:    make(map[Weight]*score.Average),
		rollup:      make(map[RollOver]*score.Average),
	}
	for _, weight := range WeightOrder {
		hd.subScore[weight] = &score.Average{}
	}
	return hd
}

func (hd *HandlerData) Rollup(over RollOver) *score.Average {
	if hd.rollup[over] == nil {
		hd.rollup[over] = &score.Average{}
	}
	return hd.rollup[over]
}

func (hd *HandlerData) SubScore(weight Weight) *score.Average {
	if hd.subScore[weight] == nil {
		hd.subScore[weight] = &score.Average{}
	}
	return hd.subScore[weight]
}
