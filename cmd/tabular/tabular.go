package main

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

// Tabular reads the JSON from gobenchdata and formats it into simple tables.
// See scripts/bench for usage example.

type benchmark struct {
	Name    string
	Runs    int
	NsPerOp float64
	Mem     struct {
		BytesPerOp  int
		AllocsPerOp int
		MBPerSec    int
	}
}

type suiteData struct {
	Benchmarks []benchmark
}

type jsonData struct {
	Suites []suiteData
}

type testData struct {
	iterations     int
	nanosPerOp     float64
	memBytesPerOp  int
	memAllocsPerOp int
	memMbPerSec    int
}

var ptnBenchName = regexp.MustCompile(`Benchmark(.+)-(\d+)`)
var ptnHandlerName = regexp.MustCompile(`Benchmark(?:_slog)?_(.*)`)

func main() {
	jsonFile := flag.String("json", "", "Source JSON file (else stdin)")
	flag.Parse()

	var err error
	out := os.Stdout
	if *jsonFile != "" {
		if out, err = os.Open(*jsonFile); err != nil {
			fmt.Printf("* Error opening JSON file %s: %s\n", *jsonFile, err)
			return
		}
	}

	var source []byte
	if source, err = io.ReadAll(out); err != nil {
		fmt.Printf("* Error reading JSON source: %s\n", err)
		return
	}

	var br []jsonData
	if err := json.Unmarshal(source, &br); err != nil {
		fmt.Printf("* Unable to unmarshal bench.json data: %s\n", err)
		return
	}

	if len(br) != 1 {
		fmt.Printf("* Top level array has %d items\n", len(br))
		return
	}

	item := br[0]
	if len(item.Suites) != 1 {
		fmt.Printf("* Suites array has %d items\n", len(item.Suites))
		return
	}

	type handlerName string
	type benchName string

	data := make(map[benchName]map[handlerName]testData)

	for _, bm := range item.Suites[0].Benchmarks {
		parts := strings.Split(bm.Name, "/")
		if len(parts) != 2 {
			fmt.Printf("* Name has %d parts: %s\n", len(parts), bm.Name)
			continue
		}
		handler := handlerName(parts[0])
		bench := benchName(parts[1])

		if matches := ptnBenchName.FindSubmatch([]byte(bench)); matches != nil && len(matches) > 2 {
			name := strings.TrimLeft(string(matches[1]), "_")
			name = strings.Replace(name, "_", " ", -1)
			bench = benchName(fmt.Sprintf("%s (%s CPU)", name, string(matches[2])))
		}

		if matches := ptnHandlerName.FindSubmatch([]byte(handler)); matches != nil && len(matches) > 1 {
			handler = handlerName(matches[1])
		}

		if data[bench] == nil {
			data[bench] = make(map[handlerName]testData, 0)
		}

		data[bench][handler] = testData{
			iterations:     bm.Runs,
			nanosPerOp:     bm.NsPerOp,
			memBytesPerOp:  bm.Mem.BytesPerOp,
			memAllocsPerOp: bm.Mem.AllocsPerOp,
			memMbPerSec:    bm.Mem.BytesPerOp,
		}
	}

	benches := make([]string, 0)
	for test := range data {
		benches = append(benches, string(test))
	}
	sort.Strings(benches)

	for _, test := range benches {
		fmt.Printf("\nBenchmark %s\n", test)
		fmt.Println("  Handler                    Runs     Ns/Op  Bytes/Op Allocs/Op    MB/Sec")
		fmt.Println("  -----------------------------------------------------------------------")

		testData := data[benchName(test)]
		hdlrs := make([]string, 0)
		for hdlr := range testData {
			hdlrs = append(hdlrs, string(hdlr))
		}
		sort.Strings(hdlrs)

		for _, hdlr := range hdlrs {
			hdlrData := testData[handlerName(hdlr)]
			fmt.Printf("  %-20s  %9d %9.3f %9d %9d %9d\n",
				hdlr, hdlrData.iterations, hdlrData.nanosPerOp,
				hdlrData.memBytesPerOp, hdlrData.memAllocsPerOp, hdlrData.memMbPerSec)
		}
	}
}
