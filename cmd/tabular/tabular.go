package main

import (
	"flag"
	"fmt"

	"github.com/madkins23/go-slog/infra"
)

// Tabular reads the JSON from gobenchdata and formats it into simple tables.
// See scripts/bench for usage example.

func main() {
	flag.Parse() // Necessary for -json=<file> argument defined in infra package.

	var data infra.BenchData
	if err := data.LoadBenchJSON(); err != nil {
		fmt.Printf("* Error parsing benchmark JSON: %s\n", err)
	}

	for _, bench := range data.Benches() {
		fmt.Printf("\nBenchmark %s\n", bench)
		fmt.Println("  Handler                    Runs     Ns/Op  Bytes/Op Allocs/Op    MB/Sec")
		fmt.Println("  -----------------------------------------------------------------------")

		handlerRecords := data.HandlerRecords(bench)
		for _, handler := range data.Handlers() {
			handlerRecord := handlerRecords[handler]
			fmt.Printf("  %-20s  %9d %9.3f %9d %9d %9d\n",
				handler, handlerRecord.Iterations, handlerRecord.NanosPerOp,
				handlerRecord.MemBytesPerOp, handlerRecord.MemAllocsPerOp, handlerRecord.MemMbPerSec)
		}
	}
}
