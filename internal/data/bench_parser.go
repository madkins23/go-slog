package data

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
)

var (
	ptnHandlerDef = regexp.MustCompile(`^#\s*Handler\[(\S+)\]\s*=\s*"(\S+)"\s*$`)
	ptnWarnLine   = regexp.MustCompile(`^# (.*)`)
	ptnDataLine   = regexp.MustCompile(`^Benchmark([^/]+)/Benchmark([^-]+)-(\d+)\s+(\d+)\s+(\d+(?:\.\d+)?)\s+ns/op\b`)
	ptnAllocsOp   = regexp.MustCompile(`\s(\d+)\s+allocs/op\b`)
	ptnBytesOp    = regexp.MustCompile(`\s(\d+)\s+B/op\b`)
	ptnMbSec      = regexp.MustCompile(`\s(\d+(?:\.\d+)?)\s+MB/s`)
)

// -----------------------------------------------------------------------------

// ParseBenchmarkData parses benchmark data from the output of go -bench.
// If argument 'in' is nil then the data will be loaded from os.Stdin
// unless the -bench=<path> command line flag is set
// in which case the data will be loaded from the specified path.
// This can be overridden by passing in a non-nil io.Reader,
// in which that data will be parsed instead.
func (b *Benchmarks) ParseBenchmarkData(in io.Reader) error {
	var err error
	if in == nil {
		if *benchFile != "" {
			if in, err = os.Open(*benchFile); err != nil {
				return fmt.Errorf("open --bench=%s: %s\n", *benchFile, err)
			}
		} else {
			in = os.Stdin
		}
	}
	scanner := bufio.NewScanner(in)

	for scanner.Scan() {
		var ok bool
		var hdlrBytes, testBytes []byte
		var nsOps, mbSec float64
		var cpus, runs, allocsOp, bytesOp uint64
		line := scanner.Bytes()
		if matches := ptnHandlerDef.FindSubmatch(line); len(matches) == 3 {
			// Capture relationship between handler name in benchmark function vs. Creator.
			// Parse handler tag/name information written in the log by bench/suite.Run().
			// This way the handler name field is populated by the Creator name string.
			// The data will be parsed by internal/data.Benchmarks.ParseBenchmarkData() and
			// passed into Warnings.ParseWarningData().
			b.handlerNames[HandlerTag(matches[1])] = string(matches[2])
		} else if matches := ptnWarnLine.FindSubmatch(line); len(matches) == 2 {
			// Capture warning text marked with "# " at beginning of line.
			if len(b.warningText) > 0 {
				b.warningText = append(b.warningText, []byte("\n")...)
			}
			b.warningText = append(b.warningText, matches[1]...)
		} else if matches := ptnDataLine.FindSubmatch(line); len(matches) == 6 {
			// Process a data line.
			hdlrBytes = matches[1]
			testBytes = matches[2]
			if cpus, err = strconv.ParseUint(string(matches[3]), 10, 64); err != nil {
				return fmt.Errorf("parse cpus: %w", err)
			}
			if runs, err = strconv.ParseUint(string(matches[4]), 10, 64); err != nil {
				return fmt.Errorf("parse runs: %w", err)
			}
			if nsOps, err = strconv.ParseFloat(string(matches[5]), 64); err != nil {
				return fmt.Errorf("parse ns/op: %w", err)
			}
			if matches = ptnAllocsOp.FindSubmatch(line); len(matches) == 2 {
				if allocsOp, err = strconv.ParseUint(string(matches[1]), 10, 64); err != nil {
					return fmt.Errorf("parse allocs/op: %w", err)
				}
			}
			if matches = ptnBytesOp.FindSubmatch(line); len(matches) == 2 {
				if bytesOp, err = strconv.ParseUint(string(matches[1]), 10, 64); err != nil {
					return fmt.Errorf("parse bytes/op: %w", err)
				}
			}
			if matches = ptnMbSec.FindSubmatch(line); len(matches) == 2 {
				if mbSec, err = strconv.ParseFloat(string(matches[1]), 64); err != nil {
					return fmt.Errorf("parse mb/s: %w", err)
				}
			}
			ok = true
		}

		if ok {
			test := TestTag("Bench" + TagSeparator + string(testBytes))
			b.testNames[test] = test.Name()

			handler := HandlerTag(hdlrBytes)
			if _, found := b.handlerNames[handler]; !found {
				b.handlerNames[handler] = handler.Name()
			}
			b.testCPUs[test] = cpus

			if b.byTest[test] == nil {
				b.byTest[test] = make(HandlerRecords)
			}
			b.byTest[test][handler] = TestRecord{
				Runs:           runs,
				NanosPerOp:     nsOps,
				MemBytesPerOp:  bytesOp,
				MemAllocsPerOp: allocsOp,
				MbPerSec:       mbSec,
				GbPerSec:       mbSec / 1_000.0,
				TbPerSec:       mbSec / 1_000_000.0,
			}

			if b.ByHandler[handler] == nil {
				b.ByHandler[handler] = make(TestRecords)
			}
			b.ByHandler[handler][test] = TestRecord{
				Runs:           runs,
				NanosPerOp:     nsOps,
				MemBytesPerOp:  bytesOp,
				MemAllocsPerOp: allocsOp,
				MbPerSec:       mbSec,
			}
		}
	}
	if scanner.Err() != nil {
		return fmt.Errorf("scan input: %w", scanner.Err())
	}

	return nil
}
