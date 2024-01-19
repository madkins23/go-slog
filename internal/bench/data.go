package bench

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var fromFile = flag.String("from", "", "Load data from path (optional)")

// -----------------------------------------------------------------------------
// Records matching gobenchdata JSON output.

// -----------------------------------------------------------------------------

// TestTag is a unique name for a Benchmark test.
// The type is an alias for string so that types can't be confused.
type TestTag string

// HandlerTag is a unique name for a slog handler.
// The type is an alias for string so that types can't be confused.
type HandlerTag string

// TestRecords is a map of test records by test tag.
type TestRecords map[TestTag]TestRecord

// HandlerRecords is a map of test records by handler tag.
type HandlerRecords map[HandlerTag]TestRecord

// TestRecord represents a single benchmark/handler test result.
type TestRecord struct {
	Runs           uint64
	NanosPerOp     float64
	MemBytesPerOp  uint64
	MemAllocsPerOp uint64
	MbPerSec       float64
	GbPerSec       float64
	TbPerSec       float64
}

func (tr *TestRecord) IsEmpty() bool {
	return tr.Runs == 0
}

func (tr *TestRecord) ItemValue(item TestItems) float64 {
	switch item {
	case Runs:
		return float64(tr.Runs)
	case Nanos:
		return tr.NanosPerOp
	case MemAllocs:
		return float64(tr.MemAllocsPerOp)
	case MemBytes:
		return float64(tr.MemBytesPerOp)
	case MbPerSec:
		return tr.MbPerSec
	case GbPerSec:
		return tr.GbPerSec
	case TbPerSec:
		return tr.TbPerSec
	default:
		slog.Warn("Unknown bench.TestItem", "item", item)
		return 0
	}
}

// -----------------------------------------------------------------------------

// Data encapsulates benchmark records by BenchmarkName and HandlerTag.
type Data struct {
	byTest       map[TestTag]HandlerRecords
	byHandler    map[HandlerTag]TestRecords
	tests        []TestTag
	handlers     []HandlerTag
	testNames    map[TestTag]string
	testCPUs     map[TestTag]uint64
	handlerNames map[HandlerTag]string
}

// -----------------------------------------------------------------------------

var (
	ptnDataLine = regexp.MustCompile(`^([^/]+)/Benchmark_([^-]+)-(\d+)\s+(\d+)\s+(\d+(?:\.\d+)?)\s+ns/op\b`)
	ptnAllocsOp = regexp.MustCompile(`\s(\d+)\s+allocs/op\b`)
	ptnBytesOp  = regexp.MustCompile(`\s(\d+)\s+B/op\b`)
	ptnMbSec    = regexp.MustCompile(`\s(\d+(?:\.\d+)?)\s+MB/s`)
)

// LoadDataJSON loads benchmark data from JSON emitted by gobenchdata.
// The data will be loaded from os.Stdin unless the -json=<path> flag is set
// in which case the data will be loaded from the specified path.
func (d *Data) LoadDataJSON() error {
	var err error
	var in io.Reader = os.Stdin
	if *fromFile != "" {
		if in, err = os.Open(*fromFile); err != nil {
			return fmt.Errorf("opening input file %s: %s\n", *fromFile, err)
		}
	}
	reader := bufio.NewReader(in)

	var line bytes.Buffer
	for {
		chunk, isPrefix, err := reader.ReadLine()
		if errors.Is(err, io.EOF) {
			if line.Len() < 1 {
				break
			}
		} else if err != nil {
			return fmt.Errorf("reading line: %w", err)
		}
		line.Write(chunk)
		if isPrefix {
			// Line not yet complete, go around again.
			continue
		}

		// Process the line here.
		var ok bool
		var hdlrBytes, testBytes []byte
		var nsOps, mbSec float64
		var cpus, runs, allocsOp, bytesOp uint64
		if matches := ptnDataLine.FindSubmatch(line.Bytes()); matches != nil && len(matches) == 6 {
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
			if matches = ptnAllocsOp.FindSubmatch(line.Bytes()); matches != nil && len(matches) == 2 {
				if allocsOp, err = strconv.ParseUint(string(matches[1]), 10, 64); err != nil {
					return fmt.Errorf("parse allocs/op: %w", err)
				}
			}
			if matches = ptnBytesOp.FindSubmatch(line.Bytes()); matches != nil && len(matches) == 2 {
				if bytesOp, err = strconv.ParseUint(string(matches[1]), 10, 64); err != nil {
					return fmt.Errorf("parse bytes/op: %w", err)
				}
			}
			if matches = ptnMbSec.FindSubmatch(line.Bytes()); matches != nil && len(matches) == 2 {
				if mbSec, err = strconv.ParseFloat(string(matches[1]), 64); err != nil {
					return fmt.Errorf("parse mb/s: %w", err)
				}
			}
			ok = true
		}

		if ok {
			test := TestTag(strings.TrimLeft(string(testBytes), "_"))
			if d.testNames == nil {
				d.testNames = make(map[TestTag]string)
			}
			d.testNames[test] = strings.Replace(string(test), "_", " ", -1)

			if string(hdlrBytes) == "Benchmark_slog" {
				// Fix this so the handler name doesn't get edited down to nothing.
				hdlrBytes = []byte("Benchmark_slog_slog_JSONHandler")
			}
			handler := HandlerTag(
				strings.TrimLeft(
					strings.TrimPrefix(string(hdlrBytes), "Benchmark_slog"),
					"_"))
			if d.handlerNames == nil {
				d.handlerNames = make(map[HandlerTag]string)
			}
			parts := strings.Split(strings.TrimLeft(string(handler), "_"), "_")
			for i, part := range parts {
				if len(part) > 0 {
					parts[i] = strings.ToUpper(part[:1]) + part[1:]
				}
			}
			d.handlerNames[handler] = strings.Join(parts, " ")

			if d.testCPUs == nil {
				d.testCPUs = make(map[TestTag]uint64)
			}
			d.testCPUs[test] = cpus

			if d.byTest == nil {
				d.byTest = make(map[TestTag]HandlerRecords)
			}
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

			if d.byHandler == nil {
				d.byHandler = make(map[HandlerTag]TestRecords)
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

		line.Reset()
	}

	return nil
}

// -----------------------------------------------------------------------------

// TestName returns the full name associated with a TestTag.
// If there is no full name the tag is returned.
func (d *Data) TestName(test TestTag) string {
	if name, found := d.testNames[test]; found {
		return name
	} else {
		return string(test)
	}
}

// TestRecords returns a map of HandlerTag to TestRecord for the specified benchmark.
func (d *Data) TestRecords(handler HandlerTag) TestRecords {
	return d.byHandler[handler]
}

// TestTags returns an array of all test names sorted alphabetically.
func (d *Data) TestTags() []TestTag {
	if d.tests == nil {
		for test := range d.byTest {
			d.tests = append(d.tests, test)
		}
		sort.Slice(d.tests, func(i, j int) bool {
			return d.TestName(d.tests[i]) < d.TestName(d.tests[j])
		})
	}
	return d.tests
}

// HandlerName returns the full name associated with a HandlerTag.
// If there is no full name the tag is returned.
func (d *Data) HandlerName(handler HandlerTag) string {
	if name, found := d.handlerNames[handler]; found {
		return name
	} else {
		return string(handler)
	}
}

// HandlerRecords returns a map of HandlerTag to TestRecord for the specified benchmark.
func (d *Data) HandlerRecords(test TestTag) HandlerRecords {
	return d.byTest[test]
}

// HandlerTags returns an array of all handler names sorted alphabetically.
func (d *Data) HandlerTags() []HandlerTag {
	if d.handlers == nil {
		for handler := range d.byHandler {
			d.handlers = append(d.handlers, handler)
		}
		sort.Slice(d.handlers, func(i, j int) bool {
			return d.HandlerName(d.handlers[i]) < d.HandlerName(d.handlers[j])
		})
	}
	return d.handlers
}
