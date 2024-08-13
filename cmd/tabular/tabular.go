/*
tabular parses benchmark test and verification test output and
formats it into simple tables and warning listings.

# Usage

	go run cmd/tabular/tabular.go [flags]

The flags are:

	-bench string
	    Load benchmark data from path (optional)
	-language value
	    One or more language tags to be tried, defaults to US English.
	-useWarnings=<bool>
	    Show warning instead of known errors, defaults true
	-verify string
	    Load verification data from path (optional)

The scripts/tabulate script will run cmd/tabular,
taking input from temporary files created by scripts/verify and scripts/bench.

The -language flag is used to enable proper formatting of displayed numbers.

# Output

	Benchmark Attributes
	╔══════════════════════╦═════════════╤═══════════════╤═════════════╤═════════════╤═════════════════╗
	║ Handler              ║        Runs │         Ns/Op │   Allocs/Op │    Bytes/Op │          GB/Sec ║
	╠══════════════════════╬═════════════╪═══════════════╪═════════════╪═════════════╪═════════════════╣
	║ chanchal/zap         ║   1,228,219 │      1,041.00 │           5 │         418 │      470,809.52 ║
	║ phsym/zerolog        ║   1,690,441 │        725.10 │           2 │         272 │      946,553.07 ║
	║ samber/logrus        ║      50,173 │     25,886.00 │          90 │       8,519 │          810.17 ║
	║ samber/zap           ║     215,511 │      6,005.00 │          46 │       6,649 │       14,318.54 ║
	║ samber/zerolog       ║     224,841 │      5,123.00 │          54 │       4,837 │       17,820.35 ║
	║ slog/json            ║     751,104 │      1,369.00 │           6 │         473 │      234,700.62 ║
	╚══════════════════════╩═════════════╧═══════════════╧═════════════╧═════════════╧═════════════════╝

	...[tables for other benchmark tests]...

	Warnings for slog/json:
	  Suggested
	   2 [Duplicates] Duplicate field(s) found
	       Verify$AttributeDuplicate: map[alpha:2 charlie:3]
	       Verify$AttributeWithDuplicate: map[alpha:2 charlie:3]

	...[warnings for other handlers]...
*/
package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/phsym/console-slog"

	"github.com/madkins23/go-utils/text/table"

	"github.com/madkins23/go-slog/internal/data"
	"github.com/madkins23/go-slog/internal/language"
)

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
	warns := data.NewWarnings()
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
			data.MemAllocs.ShortName(), data.MemBytes.ShortName(), data.MbPerSec.ShortName())
		fmt.Println(tableMgr.SeparatorString(1))
		handlerRecords := bench.HandlerRecordsFor(test)
		for _, handler := range bench.HandlerTags() {
			handlerRecord := handlerRecords[handler]
			if !handlerRecord.IsEmpty() {
				_, err := language.Printer().Printf(tableMgr.RowFormat(),
					bench.HandlerName(handler), handlerRecord.Runs, handlerRecord.NanosPerOp,
					handlerRecord.MemAllocsPerOp, handlerRecord.MemBytesPerOp,
					handlerRecord.MbPerSec)
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
					line := instance.Tag()
					if instance.HasSource() {
						line = instance.Source() + ":" + line
					}
					if extra := instance.Extra(); extra != "" {
						const newLine = "\n"
						const nextLine = "\n           "
						if strings.Index(extra, newLine) >= 0 {
							extra = nextLine + strings.Replace(extra, newLine, nextLine, -1)
						}
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
				Width:     23,
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
			{ // MB/Op
				Width:  11,
				Format: "%11.2f",
			},
		},
		Prefix:      "  ",
		Border:      true,
		BorderLines: 1,
	}
}
