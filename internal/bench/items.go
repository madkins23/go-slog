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

var shortNames = map[TestItems]string{
	Runs:      "Runs",
	Nanos:     "Ns/Op",
	MemAllocs: "Allocs/Op",
	MemBytes:  "Bytes/Op",
	MemMB:     "MB/Sec",
}

func (item TestItems) ShortName() string {
	return shortNames[item]
}

var longNames = map[TestItems]string{
	Runs:      "Test runs",
	Nanos:     "Nanoseconds per test",
	MemAllocs: "Memory allocations per test",
	MemBytes:  "Memory bytes allocated per test",
	MemMB:     "Memory Megabytes allocated per second",
}

func (item TestItems) LongName() string {
	return longNames[item]
}
