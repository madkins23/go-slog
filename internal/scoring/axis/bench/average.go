package bench

import "github.com/madkins23/go-slog/internal/scoring/score"

// -----------------------------------------------------------------------------

type Average struct {
	Value score.Value
	Count uint
}

func (ba *Average) Add(v score.Value) *Average {
	ba.Value += v
	ba.Count++
	return ba
}

func (ba *Average) AddMultiple(v score.Value, multiple uint) *Average {
	ba.Value += v * score.Value(multiple)
	ba.Count += multiple
	return ba
}

func (ba *Average) Average() score.Value {
	return ba.Value.Round() / score.Value(ba.Count)
}
