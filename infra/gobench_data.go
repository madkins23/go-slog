package infra

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"
)

var jsonFile = flag.String("json", "", "Source JSON file (else stdin)")

// -----------------------------------------------------------------------------

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

// BenchTag is an alias for string so that types can't be confused.
type BenchTag string

// HandlerTag is an alias for string so that types can't be confused.
type HandlerTag string

// BenchmarkRecord represents a single benchmark/handler test result.
type BenchmarkRecord struct {
	Iterations     int
	NanosPerOp     float64
	MemBytesPerOp  int
	MemAllocsPerOp int
	MemMbPerSec    int
}

// BenchData encapsulates benchmark records by BenchmarkName and HandlerTag.
type BenchData struct {
	date         time.Time
	byBenchmark  map[BenchTag]map[HandlerTag]BenchmarkRecord
	byHandler    map[HandlerTag]map[BenchTag]BenchmarkRecord
	benches      []BenchTag
	handlers     []HandlerTag
	benchNames   map[BenchTag]string
	handlerNames map[HandlerTag]string
}

// -----------------------------------------------------------------------------

var ptnBenchName = regexp.MustCompile(`Benchmark(.+)-(\d+)`)
var ptnHandlerName = regexp.MustCompile(`Benchmark(?:_slog)?_(.*)`)

// LoadBenchJSON loads benchmark data from JSON emitted by gobenchdata.
// The data will be loaded from os.Stdin unless the -json=<path> flag is set
// in which case the data will be loaded from the specified path.
func (bd *BenchData) LoadBenchJSON() error {
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
		return fmt.Errorf("unmarshal source bench data: %w", err)
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
				fmt.Printf("* Name has %d parts: %s\n", len(parts), bm.Name)
				continue
			}
			handler := HandlerTag(parts[0])
			bench := BenchTag(parts[1])

			if matches := ptnBenchName.FindSubmatch([]byte(bench)); matches != nil && len(matches) > 2 {
				bench = BenchTag(strings.TrimLeft(string(matches[1]), "_"))
				if bd.benchNames == nil {
					bd.benchNames = make(map[BenchTag]string)
				}
				bd.benchNames[bench] =
					fmt.Sprintf("%s (%s CPU)",
						strings.Replace(string(bench), "_", " ", -1),
						string(matches[2]))
			}
			if bd.handlerNames == nil {
				bd.handlerNames = make(map[HandlerTag]string)
			}
			bd.handlerNames[handler] = strings.Replace(string(handler), "_", " ", -1)

			if matches := ptnHandlerName.FindSubmatch([]byte(handler)); matches != nil && len(matches) > 1 {
				handler = HandlerTag(matches[1])
			}

			if bd.byBenchmark == nil {
				bd.byBenchmark = make(map[BenchTag]map[HandlerTag]BenchmarkRecord)
			}
			if bd.byBenchmark[bench] == nil {
				bd.byBenchmark[bench] = make(map[HandlerTag]BenchmarkRecord)
			}
			bd.byBenchmark[bench][handler] = BenchmarkRecord{
				Iterations:     bm.Runs,
				NanosPerOp:     bm.NsPerOp,
				MemBytesPerOp:  bm.Mem.BytesPerOp,
				MemAllocsPerOp: bm.Mem.AllocsPerOp,
				MemMbPerSec:    bm.Mem.BytesPerOp,
			}

			if bd.byHandler == nil {
				bd.byHandler = make(map[HandlerTag]map[BenchTag]BenchmarkRecord)
			}
			if bd.byHandler[handler] == nil {
				bd.byHandler[handler] = make(map[BenchTag]BenchmarkRecord)
			}
			bd.byHandler[handler][bench] = BenchmarkRecord{
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

// BenchName returns the full name associated with a BenchTag.
// If there is no full name the tag is returned.
func (bd *BenchData) BenchName(bench BenchTag) string {
	if name, found := bd.benchNames[bench]; found {
		return name
	} else {
		return string(bench)
	}
}

// BenchRecords returns a map of HandlerTag to BenchmarkRecord for the specified benchmark.
func (bd *BenchData) BenchRecords(handler HandlerTag) map[BenchTag]BenchmarkRecord {
	return bd.byHandler[handler]
}

// BenchTags returns an array of all benchmark names sorted alphabetically.
func (bd *BenchData) BenchTags() []BenchTag {
	if bd.benches == nil {
		for bench := range bd.byBenchmark {
			bd.benches = append(bd.benches, bench)
		}
		sort.Slice(bd.benches, func(i, j int) bool {
			return bd.benches[i] > bd.benches[j]
		})
	}
	return bd.benches
}

// HandlerName returns the full name associated with a HandlerTag.
// If there is no full name the tag is returned.
func (bd *BenchData) HandlerName(handler HandlerTag) string {
	if name, found := bd.handlerNames[handler]; found {
		return name
	} else {
		return string(handler)
	}
}

// HandlerRecords returns a map of HandlerTag to BenchmarkRecord for the specified benchmark.
func (bd *BenchData) HandlerRecords(bench BenchTag) map[HandlerTag]BenchmarkRecord {
	return bd.byBenchmark[bench]
}

// HandlerTags returns an array of all handler names sorted alphabetically.
func (bd *BenchData) HandlerTags() []HandlerTag {
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

func (bd *BenchData) Date() time.Time {
	return bd.date
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
