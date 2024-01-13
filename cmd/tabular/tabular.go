package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

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

	data := make(map[string]map[string]testData)

	for _, b := range item.Suites[0].Benchmarks {
		parts := strings.Split(b.Name, "/")
		if len(parts) != 2 {
			fmt.Printf("* Name has %d parts: %s\n", len(parts), b.Name)
			continue
		}

		if data[parts[0]] == nil {
			data[parts[0]] = make(map[string]testData, 0)
		}

		data[parts[0]][parts[1]] = testData{
			iterations:     b.Runs,
			nanosPerOp:     b.NsPerOp,
			memBytesPerOp:  b.Mem.BytesPerOp,
			memAllocsPerOp: b.Mem.AllocsPerOp,
			memMbPerSec:    b.Mem.BytesPerOp,
		}
	}

	tests := make([]string, 0)
	for test := range data {
		tests = append(tests, test)
	}
	sort.Strings(tests)

	for _, test := range tests {
		fmt.Printf("\n%s\n  Handler                    Runs     Ns/Op  Bytes/Op Allocs/Op    MB/Sec\n", test)
		fmt.Println("  -----------------------------------------------------------------------")

		testData := data[test]
		hdlrs := make([]string, 0)
		for hdlr := range testData {
			hdlrs = append(hdlrs, hdlr)
		}
		sort.Strings(hdlrs)

		for _, hdlr := range hdlrs {
			hdlrData := testData[hdlr]
			fmt.Printf("  %-20s  %9d %9.3f %9d %9d %9d\n",
				hdlr, hdlrData.iterations, hdlrData.nanosPerOp,
				hdlrData.memBytesPerOp, hdlrData.memAllocsPerOp, hdlrData.memMbPerSec)
		}
	}
}
