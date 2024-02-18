package data

// -----------------------------------------------------------------------------

// BenchItems are specific pieces of benchmark data available per test.
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

func (item BenchItems) ShortName() string {
	return testItemData[item].short
}

func (item BenchItems) LongName() string {
	return testItemData[item].long
}
