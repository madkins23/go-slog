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

// BenchName is an alias for string so that types can't be confused.
type BenchName string

// HandlerName is an alias for string so that types can't be confused.
type HandlerName string

// BenchmarkRecord represents a single benchmark/handler test result.
type BenchmarkRecord struct {
	Iterations     int
	NanosPerOp     float64
	MemBytesPerOp  int
	MemAllocsPerOp int
	MemMbPerSec    int
}

// BenchData encapsulates benchmark records by BenchmarkName and HandlerName.
type BenchData struct {
	byBenchmark map[BenchName]map[HandlerName]BenchmarkRecord
	byHandler   map[HandlerName]map[BenchName]BenchmarkRecord
	benches     []BenchName
	handlers    []HandlerName
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

	// TODO: What about other test data?
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
			handler := HandlerName(parts[0])
			bench := BenchName(parts[1])

			if matches := ptnBenchName.FindSubmatch([]byte(bench)); matches != nil && len(matches) > 2 {
				name := strings.TrimLeft(string(matches[1]), "_")
				name = strings.Replace(name, "_", " ", -1)
				bench = BenchName(fmt.Sprintf("%s (%s CPU)", name, string(matches[2])))
			}

			if matches := ptnHandlerName.FindSubmatch([]byte(handler)); matches != nil && len(matches) > 1 {
				handler = HandlerName(matches[1])
			}

			if bd.byBenchmark == nil {
				bd.byBenchmark = make(map[BenchName]map[HandlerName]BenchmarkRecord)
			}
			if bd.byBenchmark[bench] == nil {
				bd.byBenchmark[bench] = make(map[HandlerName]BenchmarkRecord)
			}
			bd.byBenchmark[bench][handler] = BenchmarkRecord{
				Iterations:     bm.Runs,
				NanosPerOp:     bm.NsPerOp,
				MemBytesPerOp:  bm.Mem.BytesPerOp,
				MemAllocsPerOp: bm.Mem.AllocsPerOp,
				MemMbPerSec:    bm.Mem.BytesPerOp,
			}

			if bd.byHandler == nil {
				bd.byHandler = make(map[HandlerName]map[BenchName]BenchmarkRecord)
			}
			if bd.byHandler[handler] == nil {
				bd.byHandler[handler] = make(map[BenchName]BenchmarkRecord)
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

// Benches returns an array of all benchmark names sorted alphabetically.
func (bd *BenchData) Benches() []BenchName {
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

// HandlerRecords returns a map of HandlerName to BenchmarkRecord for the specified benchmark.
func (bd *BenchData) HandlerRecords(bench BenchName) map[HandlerName]BenchmarkRecord {
	return bd.byBenchmark[bench]
}

// Handlers returns an array of all handler names sorted alphabetically.
func (bd *BenchData) Handlers() []HandlerName {
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
