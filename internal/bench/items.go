package bench

// -----------------------------------------------------------------------------

//go:generate go run github.com/dmarkham/enumer -type=TestItems
type TestItems uint

const (
	Runs TestItems = iota
	Nanos
	MemAllocs
	MemBytes
	MemMB
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
	MemMB: {
		short: "MB/Sec",
		long:  "Megabytes allocated per second",
	},
}

// -----------------------------------------------------------------------------

func (item TestItems) ShortName() string {
	return testItemData[item].short
}

func (item TestItems) LongName() string {
	return testItemData[item].long
}
