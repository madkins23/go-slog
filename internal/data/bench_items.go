package data

// -----------------------------------------------------------------------------

// BenchItems are specific pieces of benchmark data available per test.
// The items are associated with both long and short display names.
//
//go:generate go run github.com/dmarkham/enumer -type=BenchItems
type BenchItems uint

const (
	Runs BenchItems = iota
	Nanos
	MemAllocs
	MemBytes
	MbPerSec
	GbPerSec
	TbPerSec
)

// -----------------------------------------------------------------------------

// testItemData associates BenchItems with both long and short display names.
var testItemData = map[BenchItems]struct {
	short string
	long  string
}{
	Runs: {
		short: "Runs",
		long:  "Test runs",
	},
	Nanos: {
		short: "Ns/Op",
		long:  "Nanoseconds per operation",
	},
	MemAllocs: {
		short: "Allocs/Op",
		long:  "Memory allocations per operation",
	},
	MemBytes: {
		short: "Bytes/Op",
		long:  "Bytes allocated per operation",
	},
	MbPerSec: {
		short: "MB/Sec",
		long:  "Megabytes processed per second",
	},
	GbPerSec: {
		short: "GB/Sec",
		long:  "Gigabytes processed per second",
	},
	TbPerSec: {
		short: "TB/Sec",
		long:  "Terabytes processed per second",
	},
}

// -----------------------------------------------------------------------------

// ShortName returns the short display name for a benchmark item.
func (item BenchItems) ShortName() string {
	return testItemData[item].short
}

// LongName returns the long display name for a benchmark item.
func (item BenchItems) LongName() string {
	return testItemData[item].long
}
