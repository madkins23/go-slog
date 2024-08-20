package bench

import (
	"fmt"
	"math"

	"github.com/madkins23/go-slog/internal/scoring/score"
)

// -----------------------------------------------------------------------------

type Range interface {
	AddValueUint64(val uint64)
	AddValueFloat64(val float64)
	Length() float64
	RangedValue(from float64) score.Value
	String() string
}

// -----------------------------------------------------------------------------

var _ Range = &RangeFloat64{}

type RangeFloat64 struct {
	low, high float64
}

func NewRangeFloat64() *RangeFloat64 {
	return &RangeFloat64{
		low:  math.MaxFloat64,
		high: 0.0,
	}
}

func (r *RangeFloat64) AddValueUint64(val uint64) {
	r.AddValueFloat64(float64(val))
}

func (r *RangeFloat64) AddValueFloat64(val float64) {
	if val < r.low {
		r.low = val
	}
	if val > r.high {
		r.high = val
	}
}

func (r *RangeFloat64) Length() float64 {
	return r.high - r.low
}

func (r *RangeFloat64) RangedValue(from float64) score.Value {
	return score.Value(100.0 * (r.high - from) / r.Length())
}

func (r *RangeFloat64) String() string {
	return fmt.Sprintf("%0.2f -> %0.2f", r.low, r.high)
}

// -----------------------------------------------------------------------------

var _ Range = &RangeUint64{}

type RangeUint64 struct {
	low, high uint64
}

func NewRangeUint64() *RangeUint64 {
	return &RangeUint64{
		low:  math.MaxUint64,
		high: 0,
	}
}

func (r *RangeUint64) AddValueFloat64(val float64) {
	r.AddValueUint64(uint64(val))
}

func (r *RangeUint64) AddValueUint64(val uint64) {
	if val < r.low {
		r.low = val
	}
	if val > r.high {
		r.high = val
	}
}

func (r *RangeUint64) Length() float64 {
	return float64(r.high - r.low)
}

func (r *RangeUint64) RangedValue(from float64) score.Value {
	return score.Value(100.0 * (float64(r.high) - from) / r.Length())
}

func (r *RangeUint64) String() string {
	return fmt.Sprintf("%0d -> %0d", r.low, r.high)
}
