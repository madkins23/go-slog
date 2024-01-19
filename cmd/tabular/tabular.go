package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/phsym/console-slog"

	"github.com/madkins23/go-utils/text/table"

	"github.com/madkins23/go-slog/internal/bench"
	"github.com/madkins23/go-slog/internal/language"
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

	if err := language.Setup(); err != nil {
		slog.Error("Error during language setup", "err", err)
		return
	}

	var data bench.Data
	if err := data.LoadDataJSON(); err != nil {
		slog.Error("Error parsing benchmark JSON", "err", err)
		return
	}

	tableMgr := tableDefs()

	for _, test := range data.TestTags() {
		fmt.Printf("\nBenchmark %s\n", data.TestName(test))
		fmt.Println(tableMgr.BorderString(table.Top))
		fmt.Printf(tableMgr.HeaderFormat(), "Handler",
			bench.Runs.ShortName(), bench.Nanos.ShortName(),
			bench.MemAllocs.ShortName(), bench.MemBytes.ShortName(), bench.GbPerSec.ShortName())
		fmt.Println(tableMgr.SeparatorString(1))
		handlerRecords := data.HandlerRecords(test)
		for _, handler := range data.HandlerTags() {
			handlerRecord := handlerRecords[handler]
			if !handlerRecord.IsEmpty() {
				_, err := language.Printer().Printf(tableMgr.RowFormat(),
					data.HandlerName(handler), handlerRecord.Runs, handlerRecord.NanosPerOp,
					handlerRecord.MemAllocsPerOp, handlerRecord.MemBytesPerOp,
					handlerRecord.ItemValue(bench.GbPerSec))
				if err != nil {
					slog.Error("Unable to print data row", "err", err)
				}
			}
		}
		fmt.Println(tableMgr.BorderString(table.Bottom))
	}
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
				Width:  11,
				Format: "%11.2f",
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
