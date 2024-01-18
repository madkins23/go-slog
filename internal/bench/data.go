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
		MBPerSec    int
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
	MemMbPerSec    int
}

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
func (bd *Data) LoadDataJSON() error {
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

	bd.date = time.UnixMilli(int64(testData.Date))
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
				if bd.testNames == nil {
					bd.testNames = make(map[TestTag]string)
				}
				if bd.testCPUs == nil {
					bd.testCPUs = make(map[TestTag]uint64)
				}
				bd.testNames[test] = strings.Replace(string(test), "_", " ", -1)
				if cpuCount, err := strconv.ParseUint(string(matches[2]), 10, 64); err != nil {
					slog.Warn("Unable to parse CPU count", "from", matches[2], "err", err)
				} else {
					bd.testCPUs[test] = cpuCount
				}
			}

			if handler == "Benchmark_slog" {
				// Fix this so the handler name doesn't get edited down to nothing.
				handler = "Benchmark_slog_slog_JSONHandler"
			}
			handler = HandlerTag(strings.TrimLeft(strings.TrimPrefix(string(handler), "Benchmark_slog"), "_"))
			parts = strings.Split(strings.TrimLeft(string(handler), "_"), "_")
			if bd.handlerNames == nil {
				bd.handlerNames = make(map[HandlerTag]string)
			}
			for i, part := range parts {
				if len(part) > 0 {
					parts[i] = strings.ToUpper(part[:1]) + part[1:]
				}
			}
			bd.handlerNames[handler] = strings.Join(parts, " ")

			if matches := ptnHandlerName.FindSubmatch([]byte(handler)); matches != nil && len(matches) > 1 {
				handler = HandlerTag(matches[1])
			}

			if bd.byTest == nil {
				bd.byTest = make(map[TestTag]HandlerRecords)
			}
			if bd.byTest[test] == nil {
				bd.byTest[test] = make(HandlerRecords)
			}
			bd.byTest[test][handler] = TestRecord{
				Iterations:     bm.Runs,
				NanosPerOp:     bm.NsPerOp,
				MemBytesPerOp:  bm.Mem.BytesPerOp,
				MemAllocsPerOp: bm.Mem.AllocsPerOp,
				MemMbPerSec:    bm.Mem.BytesPerOp,
			}

			if bd.byHandler == nil {
				bd.byHandler = make(map[HandlerTag]TestRecords)
			}
			if bd.byHandler[handler] == nil {
				bd.byHandler[handler] = make(TestRecords)
			}
			bd.byHandler[handler][test] = TestRecord{
				Iterations:     bm.Runs,
				NanosPerOp:     bm.NsPerOp,
				MemBytesPerOp:  bm.Mem.BytesPerOp,
				MemAllocsPerOp: bm.Mem.AllocsPerOp,
				MemMbPerSec:    bm.Mem.BytesPerOp,
			}
		}
	}
	return nil
}

// -----------------------------------------------------------------------------

// TestName returns the full name associated with a TestTag.
// If there is no full name the tag is returned.
func (bd *Data) TestName(test TestTag) string {
	if name, found := bd.testNames[test]; found {
		return name
	} else {
		return string(test)
	}
}

// TestRecords returns a map of HandlerTag to TestRecord for the specified benchmark.
func (bd *Data) TestRecords(handler HandlerTag) TestRecords {
	return bd.byHandler[handler]
}

// TestTags returns an array of all test names sorted alphabetically.
func (bd *Data) TestTags() []TestTag {
	if bd.tests == nil {
		for test := range bd.byTest {
			bd.tests = append(bd.tests, test)
		}
		sort.Slice(bd.tests, func(i, j int) bool {
			return bd.TestName(bd.tests[i]) > bd.TestName(bd.tests[j])
		})
	}
	return bd.tests
}

// HandlerName returns the full name associated with a HandlerTag.
// If there is no full name the tag is returned.
func (bd *Data) HandlerName(handler HandlerTag) string {
	if name, found := bd.handlerNames[handler]; found {
		return name
	} else {
		return string(handler)
	}
}

// HandlerRecords returns a map of HandlerTag to TestRecord for the specified benchmark.
func (bd *Data) HandlerRecords(test TestTag) HandlerRecords {
	return bd.byTest[test]
}

// HandlerTags returns an array of all handler names sorted alphabetically.
func (bd *Data) HandlerTags() []HandlerTag {
	if bd.handlers == nil {
		for handler := range bd.byHandler {
			bd.handlers = append(bd.handlers, handler)
		}
		sort.Slice(bd.handlers, func(i, j int) bool {
			return bd.handlers[i] > bd.handlers[j]
		})
	}
	return bd.handlers
}

func (bd *Data) Date() time.Time {
	return bd.date
}

// -----------------------------------------------------------------------------

func (br *TestRecord) IsEmpty() bool {
	return br.Iterations == 0
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
