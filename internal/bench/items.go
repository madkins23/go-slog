package bench

// -----------------------------------------------------------------------------

//go:generate go run github.com/dmarkham/enumer -type=TestItems
type TestItems uint

const (
	Runs TestItems = iota
	Nanos
	MemAllocs
	MemBytes
	MbPerSec
	GbPerSec
	TbPerSec
)

// -----------------------------------------------------------------------------

var testItemData = map[TestItems]struct {
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

func (item TestItems) ShortName() string {
	return testItemData[item].short
}

func (item TestItems) LongName() string {
	return testItemData[item].long
}
