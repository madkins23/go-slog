package data

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	ptnWarnLine = regexp.MustCompile(`^# (.*)`)
	ptnDataLine = regexp.MustCompile(`^([^/]+)/(Benchmark_[^-]+)-(\d+)\s+(\d+)\s+(\d+(?:\.\d+)?)\s+ns/op\b`)
	ptnAllocsOp = regexp.MustCompile(`\s(\d+)\s+allocs/op\b`)
	ptnBytesOp  = regexp.MustCompile(`\s(\d+)\s+B/op\b`)
	ptnMbSec    = regexp.MustCompile(`\s(\d+(?:\.\d+)?)\s+MB/s`)
)

// -----------------------------------------------------------------------------

// ParseBenchmarkData parses benchmark data from the output of go -bench.
// The data will be loaded from os.Stdin unless the -bench=<path> flag is set
// in which case the data will be loaded from the specified path.
func (d *Benchmarks) ParseBenchmarkData(in io.Reader) error {
	var err error
	if *benchFile != "" {
		if in, err = os.Open(*benchFile); err != nil {
			return fmt.Errorf("open --bench=%s: %s\n", *benchFile, err)
		}
	} else {
		in = os.Stdin
	}
	scanner := bufio.NewScanner(in)

	for scanner.Scan() {
		var ok bool
		var hdlrBytes, testBytes []byte
		var nsOps, mbSec float64
		var cpus, runs, allocsOp, bytesOp uint64
		line := scanner.Bytes()
		if matches := ptnWarnLine.FindSubmatch(line); len(matches) == 2 {
			// Capture warning text marked with "# " at beginning of line.
			if len(d.warningText) > 0 {
				d.warningText = append(d.warningText, []byte("\n")...)
			}
			d.warningText = append(d.warningText, matches[1]...)
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
			tagName := strings.TrimPrefix(string(testBytes), "Benchmark_")
			test := TestTag("Bench:" + tagName)
			d.testNames[test] = strings.Replace(tagName, "_", " ", -1)

			if string(hdlrBytes) == "Benchmark_slog" {
				// Fix this so the handler name doesn't get edited down to nothing.
				hdlrBytes = []byte("Benchmark_slog_slog_JSONHandler")
			}
			handler := HandlerTag(strings.TrimPrefix(string(hdlrBytes), "Benchmark_slog_"))
			parts := strings.Split(string(handler), "_")
			for i, part := range parts {
				parts[i] = strings.ToUpper(part[:1]) + strings.ToLower(part[1:])
			}
			d.handlerNames[handler] = strings.Join(parts, " ")

			d.testCPUs[test] = cpus

			if d.byTest[test] == nil {
				d.byTest[test] = make(HandlerRecords)
			}
			d.byTest[test][handler] = TestRecord{
				Runs:           runs,
				NanosPerOp:     nsOps,
				MemBytesPerOp:  bytesOp,
				MemAllocsPerOp: allocsOp,
				MbPerSec:       mbSec,
				GbPerSec:       mbSec / 1_000.0,
				TbPerSec:       mbSec / 1_000_000.0,
			}

			if d.byHandler[handler] == nil {
				d.byHandler[handler] = make(TestRecords)
			}
			d.byHandler[handler][test] = TestRecord{
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