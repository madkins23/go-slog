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

	var data bench.Data
	if err := data.LoadDataJSON(); err != nil {
		slog.Error("Error parsing benchmark JSON", "err", err)
		return
	}

	tableMgr := tableDefs()

	for _, test := range data.TestTags() {
		fmt.Printf("\nBenchmark %s\n", test)
		fmt.Println(tableMgr.BorderString(table.Top))
		fmt.Printf(tableMgr.HeaderFormat(), "Handler",
			bench.Runs.ShortName(), bench.Nanos.ShortName(),
			bench.MemAllocs.ShortName(), bench.MemBytes.ShortName(), bench.GBperSec.ShortName())
		fmt.Println(tableMgr.SeparatorString(1))
		handlerRecords := data.HandlerRecords(test)
		for _, handler := range data.HandlerTags() {
			handlerRecord := handlerRecords[handler]
			if !handlerRecord.IsEmpty() {
				fmt.Printf(tableMgr.RowFormat(),
					handler, handlerRecord.Iterations, handlerRecord.NanosPerOp,
					handlerRecord.MemAllocsPerOp, handlerRecord.MemBytesPerOp, handlerRecord.GbPerSec)
			}
		}
		fmt.Println(tableMgr.BorderString(table.Bottom))
	}
}

func tableDefs() table.TableDef {
	return table.TableDef{
		Columns: []table.ColumnDef{
			{
				Width:     20,
				AlignLeft: true,
			},
			{
				Width:       9,
				Format:      "%9d",
				ColumnLines: 1,
			},
			{
				Width:  11,
				Format: "%11.3f",
			},
			{
				Width:  9,
				Format: "%9d",
			},
			{
				Width:  9,
				Format: "%9d",
			},
			{
				Width:  13,
				Format: "%13.3f",
			},
		},
		Prefix:      "  ",
		Border:      true,
		BorderLines: 1,
	}
}
