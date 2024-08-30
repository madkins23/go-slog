package score

import "math"

// -----------------------------------------------------------------------------

type Average struct {
	Value Value
	Count uint
}

func (ba *Average) Add(v Value) *Average {
	ba.Value += v
	ba.Count++
	return ba
}

func (ba *Average) AddMultiple(v Value, multiple uint) *Average {
	ba.Value += v * Value(multiple)
	ba.Count += multiple
	return ba
}

func (ba *Average) Average() Value {
	if ba.Count > 1 {
		return ba.Value / Value(ba.Count)
	}
	if ba.Count == 1 {
		return ba.Value
	}
	return Value(math.NaN())
}
