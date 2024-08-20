package bench

import "github.com/madkins23/go-slog/internal/data"

// -----------------------------------------------------------------------------

type Weight string

const (
	Allocations Weight = "Allocations"
	AllocBytes  Weight = "Alloc Bytes"
	Nanoseconds Weight = "Nanoseconds"
)

var WeightOrder = []Weight{
	Nanoseconds,
	AllocBytes,
	Allocations,
}

func (bw Weight) Item() data.BenchItems {
	switch bw {
	case Allocations:
		return data.MemAllocs
	case AllocBytes:
		return data.MemBytes
	case Nanoseconds:
		return data.Nanos
	default:
		return 0.0
	}
}
