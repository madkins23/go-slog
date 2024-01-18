package bench

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

var jsonFile = flag.String("json", "", "Source JSON file (else stdin)")

// -----------------------------------------------------------------------------
// Records matching gobenchdata JSON output.

type testData struct {
	Date   uint64
	Suites []suiteData
}

type suiteData struct {
	GoOS       string `json:"Goos"`
	GoArch     string `json:"Goarch"`
	Package    string `json:"Pkg"`
	Benchmarks []benchmarkData
}

type benchmarkData struct {
	Name    string
	Runs    int
	NsPerOp float64
	Mem     struct {
		BytesPerOp  int
		AllocsPerOp int
		MBPerSec    float64
	}
}

// -----------------------------------------------------------------------------

// TestTag is a unique name for a Benchmark test.
// The type is an alias for string so that types can't be confused.
type TestTag string

// HandlerTag is a unique name for a slog handler.
// The type is an alias for string so that types can't be confused.
type HandlerTag string

// TestRecords is a map of test records by test tag.
type TestRecords map[TestTag]TestRecord

// HandlerRecords is a map of test records by handler tag.
type HandlerRecords map[HandlerTag]TestRecord

// TestRecord represents a single benchmark/handler test result.
type TestRecord struct {
	Iterations     int
	NanosPerOp     float64
	MemBytesPerOp  int
	MemAllocsPerOp int
	GbPerSec       float64
}

func (tr *TestRecord) IsEmpty() bool {
	return tr.Iterations == 0
}

func (tr *TestRecord) ItemValue(item TestItems) float64 {
	switch item {
	case Runs:
		return float64(tr.Iterations)
	case Nanos:
		return tr.NanosPerOp
	case MemAllocs:
		return float64(tr.MemAllocsPerOp)
	case MemBytes:
		return float64(tr.MemBytesPerOp)
	case GBperSec:
		return float64(tr.GbPerSec)
	default:
		slog.Warn("Unknown bench.TestItem", "item", item)
		return 0
	}
}

// -----------------------------------------------------------------------------

// Data encapsulates benchmark records by BenchmarkName and HandlerTag.
type Data struct {
	date         time.Time
	byTest       map[TestTag]HandlerRecords
	byHandler    map[HandlerTag]TestRecords
	tests        []TestTag
	handlers     []HandlerTag
	testNames    map[TestTag]string
	testCPUs     map[TestTag]uint64
	handlerNames map[HandlerTag]string
}

// -----------------------------------------------------------------------------

var ptrTestName = regexp.MustCompile(`Benchmark(.+)-(\d+)`)
var ptnHandlerName = regexp.MustCompile(`Benchmark(?:_slog)?_(.*)`)

// LoadDataJSON loads benchmark data from JSON emitted by gobenchdata.
// The data will be loaded from os.Stdin unless the -json=<path> flag is set
// in which case the data will be loaded from the specified path.
func (d *Data) LoadDataJSON() error {
	var err error
	var in io.Reader = os.Stdin
	if *jsonFile != "" {
		if in, err = os.Open(*jsonFile); err != nil {
			return fmt.Errorf("opening JSON file %s: %s\n", *jsonFile, err)
		}
	}

	source, err := io.ReadAll(in)
	if err != nil {
		return fmt.Errorf("read all source data: %w\n", err)
	}
	// Skip over non-JSON first line from gobenchdata.
	source = skipTextBeforeJSON(source)

	// Unmarshal the data:
	var data []testData
	if err = json.Unmarshal(source, &data); err != nil {
		return fmt.Errorf("unmarshal source gobenchdata: %w", err)
	}
	if len(data) != 1 {
		return fmt.Errorf("top level array has %d items\n", len(data))
	}
	testData := data[0]

	d.date = time.UnixMilli(int64(testData.Date))
	if len(testData.Suites) != 1 {
		return fmt.Errorf("suites array has %d items\n", len(testData.Suites))
	}
	for _, suiteData := range testData.Suites {
		// TODO: What about other suite data?
		for _, bm := range suiteData.Benchmarks {
			parts := strings.Split(bm.Name, "/")
			if len(parts) != 2 {
				slog.Warn("Name has wrong number of parts", "name", bm.Name, "parts", len(parts))
				continue
			}
			handler := HandlerTag(parts[0])
			test := TestTag(parts[1])

			if matches := ptrTestName.FindSubmatch([]byte(test)); matches != nil && len(matches) > 2 {
				test = TestTag(strings.TrimLeft(string(matches[1]), "_"))
				if d.testNames == nil {
					d.testNames = make(map[TestTag]string)
				}
				if d.testCPUs == nil {
					d.testCPUs = make(map[TestTag]uint64)
				}
				d.testNames[test] = strings.Replace(string(test), "_", " ", -1)
				if cpuCount, err := strconv.ParseUint(string(matches[2]), 10, 64); err != nil {
					slog.Warn("Unable to parse CPU count", "from", matches[2], "err", err)
				} else {
					d.testCPUs[test] = cpuCount
				}
			}

			if handler == "Benchmark_slog" {
				// Fix this so the handler name doesn't get edited down to nothing.
				handler = "Benchmark_slog_slog_JSONHandler"
			}
			handler = HandlerTag(strings.TrimLeft(strings.TrimPrefix(string(handler), "Benchmark_slog"), "_"))
			parts = strings.Split(strings.TrimLeft(string(handler), "_"), "_")
			if d.handlerNames == nil {
				d.handlerNames = make(map[HandlerTag]string)
			}
			for i, part := range parts {
				if len(part) > 0 {
					parts[i] = strings.ToUpper(part[:1]) + part[1:]
				}
			}
			d.handlerNames[handler] = strings.Join(parts, " ")

			if matches := ptnHandlerName.FindSubmatch([]byte(handler)); matches != nil && len(matches) > 1 {
				handler = HandlerTag(matches[1])
			}

			if d.byTest == nil {
				d.byTest = make(map[TestTag]HandlerRecords)
			}
			if d.byTest[test] == nil {
				d.byTest[test] = make(HandlerRecords)
			}
			d.byTest[test][handler] = TestRecord{
				Iterations:     bm.Runs,
				NanosPerOp:     bm.NsPerOp,
				MemBytesPerOp:  bm.Mem.BytesPerOp,
				MemAllocsPerOp: bm.Mem.AllocsPerOp,
				GbPerSec:       bm.Mem.MBPerSec / 1000,
			}

			if d.byHandler == nil {
				d.byHandler = make(map[HandlerTag]TestRecords)
			}
			if d.byHandler[handler] == nil {
				d.byHandler[handler] = make(TestRecords)
			}
			d.byHandler[handler][test] = TestRecord{
				Iterations:     bm.Runs,
				NanosPerOp:     bm.NsPerOp,
				MemBytesPerOp:  bm.Mem.BytesPerOp,
				MemAllocsPerOp: bm.Mem.AllocsPerOp,
				GbPerSec:       bm.Mem.MBPerSec,
			}
		}
	}
	return nil
}

// -----------------------------------------------------------------------------

// TestName returns the full name associated with a TestTag.
// If there is no full name the tag is returned.
func (d *Data) TestName(test TestTag) string {
	if name, found := d.testNames[test]; found {
		return name
	} else {
		return string(test)
	}
}

// TestRecords returns a map of HandlerTag to TestRecord for the specified benchmark.
func (d *Data) TestRecords(handler HandlerTag) TestRecords {
	return d.byHandler[handler]
}

// TestTags returns an array of all test names sorted alphabetically.
func (d *Data) TestTags() []TestTag {
	if d.tests == nil {
		for test := range d.byTest {
			d.tests = append(d.tests, test)
		}
		sort.Slice(d.tests, func(i, j int) bool {
			return d.TestName(d.tests[i]) > d.TestName(d.tests[j])
		})
	}
	return d.tests
}

// HandlerName returns the full name associated with a HandlerTag.
// If there is no full name the tag is returned.
func (d *Data) HandlerName(handler HandlerTag) string {
	if name, found := d.handlerNames[handler]; found {
		return name
	} else {
		return string(handler)
	}
}

// HandlerRecords returns a map of HandlerTag to TestRecord for the specified benchmark.
func (d *Data) HandlerRecords(test TestTag) HandlerRecords {
	return d.byTest[test]
}

// HandlerTags returns an array of all handler names sorted alphabetically.
func (d *Data) HandlerTags() []HandlerTag {
	if d.handlers == nil {
		for handler := range d.byHandler {
			d.handlers = append(d.handlers, handler)
		}
		sort.Slice(d.handlers, func(i, j int) bool {
			return d.HandlerName(d.handlers[i]) > d.HandlerName(d.handlers[j])
		})
	}
	return d.handlers
}

func (d *Data) Date() time.Time {
	return d.date
}

// -----------------------------------------------------------------------------

// skipTextBeforeJSON skips over any non-JSON lines until some JSON is found
// then returns the remainder of the source data starting with the JSON.
// This support using the gobenchdata application which places a line of text
// ahead of the JSON output unless the output is redirected to a file.
// This supports reading from gobenchdata standard output via a shell pipe.
func skipTextBeforeJSON(source []byte) []byte {
	newLine := true
	for i, b := range source {
		if newLine && b == '[' || b == '{' {
			return source[i:]
		} else if b == '\n' {
			newLine = true
		} else if newLine {
			newLine = false
		}
	}
	return []byte{}
}
