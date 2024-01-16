package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/phsym/console-slog"

	"github.com/madkins23/go-slog/infra"
)

// Tabular reads the JSON from gobenchdata and formats it into simple tables.
// See scripts/bench for usage example.

func main() {
	flag.Parse() // Necessary for -json=<file> argument defined in infra package.

	logger := slog.New(console.NewHandler(os.Stderr, &console.HandlerOptions{
		Level:      slog.LevelInfo,
		TimeFormat: time.TimeOnly,
	}))
	slog.SetDefault(logger)

	var data infra.BenchData
	if err := data.LoadBenchJSON(); err != nil {
		slog.Error("Error parsing benchmark JSON", "err", err)
		return
	}

	for _, bench := range data.BenchTags() {
		fmt.Printf("\nBenchmark %s\n", bench)
		fmt.Println("  Handler                    Runs     Ns/Op  Bytes/Op Allocs/Op    MB/Sec")
		fmt.Println("  -----------------------------------------------------------------------")
		handlerRecords := data.HandlerRecords(bench)
		for _, handler := range data.HandlerTags() {
			handlerRecord := handlerRecords[handler]
			if !handlerRecord.IsEmpty() {
				fmt.Printf("  %-20s  %9d %9.3f %9d %9d %9d\n",
					handler, handlerRecord.Iterations, handlerRecord.NanosPerOp,
					handlerRecord.MemBytesPerOp, handlerRecord.MemAllocsPerOp, handlerRecord.MemMbPerSec)
			}
		}
	}
}
