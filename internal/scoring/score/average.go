package score

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
	return ba.Value.Round() / Value(ba.Count)
}
