package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/phsym/console-slog"

	"github.com/madkins23/go-utils/text/table"

	"github.com/madkins23/go-slog/internal/data"
	"github.com/madkins23/go-slog/internal/language"
)

// Tabular reads output from `test -bench` and formats it into simple tables.
// See scripts/bench and scripts/tabulate for usage examples.

func main() {
	flag.Parse() // Necessary for -json=<file> argument defined in infra package.

	logger := slog.New(console.NewHandler(os.Stderr, &console.HandlerOptions{
		Level:      slog.LevelInfo,
		TimeFormat: time.TimeOnly,
	}))
	slog.SetDefault(logger)

	if err := language.Setup(); err != nil {
		slog.Error("Error during language setup", "err", err)
		return
	}

	bench := data.NewBenchmarks()
	if err := bench.ParseBenchmarkData(nil); err != nil {
		slog.Error("Error parsing -bench data", "err", err)
		return
	}

	tableMgr := tableDefs()

	for _, test := range bench.TestTags() {
		fmt.Printf("\nBenchmark %s\n", bench.TestName(test))
		fmt.Println(tableMgr.BorderString(table.Top))
		fmt.Printf(tableMgr.HeaderFormat(), "Handler",
			data.Runs.ShortName(), data.Nanos.ShortName(),
			data.MemAllocs.ShortName(), data.MemBytes.ShortName(), data.GbPerSec.ShortName())
		fmt.Println(tableMgr.SeparatorString(1))
		handlerRecords := bench.HandlerRecords(test)
		for _, handler := range bench.HandlerTags() {
			handlerRecord := handlerRecords[handler]
			if !handlerRecord.IsEmpty() {
				_, err := language.Printer().Printf(tableMgr.RowFormat(),
					bench.HandlerName(handler), handlerRecord.Runs, handlerRecord.NanosPerOp,
					handlerRecord.MemAllocsPerOp, handlerRecord.MemBytesPerOp,
					handlerRecord.GbPerSec)
				if err != nil {
					slog.Error("Unable to print data row", "err", err)
				}
			}
		}
		fmt.Println(tableMgr.BorderString(table.Bottom))
	}

	if bench.HasWarningText() {
		fmt.Println(string(bench.WarningText()))
	}

	fmt.Println()
}

func tableDefs() table.TableDef {
	return table.TableDef{
		Columns: []table.ColumnDef{
			{ // Handler
				Width:     20,
				AlignLeft: true,
			},
			{ // Runs
				Width:       11,
				Format:      "%11d",
				ColumnLines: 1,
			},
			{ // Nanoseconds/Op
				Width:  13,
				Format: "%13.2f",
			},
			{ // Allocs/Op
				Width:  11,
				Format: "%11d",
			},
			{ // Bytes/Op
				Width:  11,
				Format: "%11d",
			},
			{ // GB/Sec
				Width:  15,
				Format: "%15.2f",
			},
		},
		Prefix:      "  ",
		Border:      true,
		BorderLines: 1,
	}
}
