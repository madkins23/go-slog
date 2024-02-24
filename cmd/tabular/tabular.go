// Reads benchmark output and displays prints it as text tables.
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

/*
tabular parses benchmark test and verification test output and
formats it into simple tables and warning listings.

Usage:

	go run cmd/tabular/tabular.go [flags]

The flags are:

	-bench string
	    Load benchmark data from path (optional)
	-language value
	    One or more language tags to be tried, defaults to US English.
	-useWarnings
	    Show warning instead of known errors
	-verify string
	    Load verification data from path (optional)

See scripts/bench, scripts/verify and scripts/server for usage examples.
*/
func main() {
	flag.Parse() // Necessary for -bench=<file> argument defined in infra package.

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
	warns := data.NewWarningData()
	if err := data.Setup(bench, warns); err != nil {
		slog.Error("Setup error", "err", err)
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

	for _, tag := range warns.HandlerTags() {
		fmt.Printf("\nWarnings for %s:\n", warns.HandlerName(tag))
		levels := warns.ForHandler(tag)
		for _, level := range levels.Levels() {
			fmt.Printf("  %s\n", level.Name())
			for _, warn := range level.Warnings() {
				fmt.Printf("  %4d [%s] %s\n", len(warn.Instances()), warn.Name(), warn.Summary())
				for _, instance := range warn.Instances() {
					line := instance.Name()
					if extra := instance.Extra(); extra != "" {
						line += ": " + extra
					}
					fmt.Printf("         %s\n", line)
					if line := instance.Line(); line != "" {
						fmt.Printf("           %s\n", line)
					}
				}
			}
		}
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
