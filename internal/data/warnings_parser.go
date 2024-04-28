package data

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"regexp"
	"strings"

	"github.com/madkins23/go-slog/internal/warning"
)

// -----------------------------------------------------------------------------

var (
	ptnWarningsFor  = regexp.MustCompile(`^\s*Warnings\s+for\s+(.*):\s*$`)
	ptnLevel        = regexp.MustCompile(`^\s*(\S+)\s*$`)
	ptnWarning      = regexp.MustCompile(`^\s*\d+\s+\[(.*)]\s*(.*?)\s*$`)
	ptnInstance     = regexp.MustCompile(`^\s*(\S+)(?::\s*(.*?))?\s*$`)
	ptnExtra        = regexp.MustCompile(`^\s*\+(.*?)\s*$`)
	ptnLogLine      = regexp.MustCompile(`^\s*\{`)
	ptnByWarning    = regexp.MustCompile(`^\s*Handlers\s+by\s+warning:\s*$`)
	ptnSummaryStart = regexp.MustCompile(`^:\[\s*(\S.*?)\s*$`)
	ptnSummaryLine  = regexp.MustCompile(`^::\s*(\S.*?)\s*$`)
	ptnSummaryLink  = regexp.MustCompile(`^:>\s*(\S.*?)\s*-->\s*(\S.*?)\s*$`)
	ptnSummaryEnd   = regexp.MustCompile(`^:]\s*$`)
)

// ParseWarningData parses warning data from the output of benchmark and verification testing.
// The data will be loaded from os.Stdin unless the -bench=<path> flag is set
// in which case the data will be loaded from the specified path.
//
// TODO: Refactor this method into a series of smaller ones?
func (w *Warnings) ParseWarningData(in io.Reader, source string, lookup map[string]HandlerTag) error {
	var err error
	if in == nil {
		if *verifyFile != "" {
			if in, err = os.Open(*verifyFile); err != nil {
				return fmt.Errorf("open --verify=%s: %s\n", *verifyFile, err)
			}
		} else {
			slog.Warn("unable to parse verification data without -verify flag")
			return nil
		}
	}
	scanner := bufio.NewScanner(in)

	var handler HandlerTag
	var test TestTag
	var level warning.Level
	var dWarning *dataWarning
	var instance *dataInstance
	saveInstance := func(line []byte) {
		if instance != nil {
			if dWarning == nil {
				slog.Warn("Nil dataWarning", "line", line, "instance", instance)
			} else {
				dWarning.AddInstance(instance)
				tWarning := w.findTest(test, level, dWarning.warning.name)
				tWarning.warning.summary = dWarning.warning.summary
				tWarning.AddInstance(
					&dataInstance{
						name:  handler.Name(),
						extra: instance.extra,
						log:   instance.log,
					})
				wd := w.FindWarning(dWarning.warning.name)
				wd.count[handler]++
				wd.hdlrMap[handler] = true
				wd.testMap[test] = true
			}
			instance = nil
		}
	}
	for scanner.Scan() {
		// Remove prefix output during benchmark testing to mark warning data.
		// TODO: Is there ever any prefix octothorpe? Remove TrimPrefix()?
		line := bytes.Trim(bytes.TrimPrefix(scanner.Bytes(), []byte{'#'}), " ")
		if len(line) == 0 {
			continue
		}

		if matches := ptnWarningsFor.FindSubmatch(line); len(matches) == 2 {
			saveInstance(line)
			handler = HandlerTag(matches[1])
			// Capture relationship between handler name in benchmark function vs. Creator.
			// The handler string here is the Creator name,
			// converting it through the lookup map makes it into the Benchmarks variant,
			// which makes all handler tags the same between Benchmarks and Warnings.
			// The Creator name can't be used because they all contain slashes
			// which breaks up the URL pattern matching in the server.
			if h, found := lookup[string(handler)]; found {
				w.getHandlerData(h).name = string(handler)
				handler = h
			} else {
				slog.Warn("Default handler name", "handler", handler)
				parts := strings.Split(string(handler), "/")
				for i, part := range parts {
					if len(part) > 0 {
						parts[i] = strings.ToUpper(part[:1]) + strings.ToLower(part[1:])
					}
				}
				w.getHandlerData(handler).name = strings.Join(parts, " ")
			}
			continue
		}
		if ptnByWarning.Match(line) {
			// End of data.
			break
		}
		if matches := ptnSummaryStart.FindSubmatch(line); len(matches) == 2 {
			hdlrName := string(matches[1])
			hdlrSummary := make([]string, 0, 5)
			hdlrLinks := make(map[string]string, 5)
			for scanner.Scan() {
				sumLine := scanner.Bytes()
				if sumMatches := ptnSummaryLine.FindSubmatch(sumLine); len(sumMatches) == 2 {
					hdlrSummary = append(hdlrSummary, string(sumMatches[1]))
				} else if sumMatches = ptnSummaryLink.FindSubmatch(sumLine); len(sumMatches) == 3 {
					hdlrLinks[string(sumMatches[1])] = string(sumMatches[2])
				} else {
					if !ptnSummaryEnd.Match(sumLine) {
						slog.Warn("Unexpected line during parse of handler summary", "line", sumLine)
					}
					if lookup != nil {
						if h, found := lookup[hdlrName]; found {
							if w.getHandlerData(h).name != "" && w.getHandlerData(h).name != hdlrName {
								slog.Warn("Handler name mismatch", "current", hdlrName, "previous", w.getHandlerData(h).name)
							}
							w.getHandlerData(h).summary = strings.Join(hdlrSummary, "\n")
							w.getHandlerData(h).links = hdlrLinks
						}
					}
					line = sumLine
					break
				}
			}
		}
		if handler == "" {
			// Can't do anything until we recognize something.
			continue
		}
		if matches := ptnLevel.FindSubmatch(line); len(matches) == 2 {
			if lvl, err := warning.ParseLevel(string(matches[1])); err == nil {
				saveInstance(line)
				level = lvl
				dWarning = nil
				continue
			}
			// Else not a level, keep looking.
		}
		if matches := ptnWarning.FindSubmatch(line); len(matches) == 3 {
			warningName := string(matches[1])
			saveInstance(line)
			dWarning = w.findHandler(handler, level, warningName)
			dWarning.warning.summary = string(matches[2])
			instance = nil
			continue
		}
		// Do this before ptnInstance as they can otherwise get confused.
		if ptnLogLine.Match(line) {
			instance.line = string(line)
			// Attempt to pretty-print the log line.
			var jm map[string]any
			if json.Unmarshal(line, &jm) == nil {
				if indented, err := json.MarshalIndent(jm, "", "\t"); err == nil {
					instance.log = string(indented)
				}
			}
			continue
		}
		if matches := ptnExtra.FindSubmatch(line); len(matches) == 2 {
			if instance != nil {
				instance.extra = instance.extra + "\n" + string(matches[1])
			}
			continue
		}
		if matches := ptnInstance.FindSubmatch(line); len(matches) == 3 {
			saveInstance(line)
			testTagStr := string(matches[1])
			for {
				changed := false
				for _, src := range []string{
					"Benchmark_", "Benchmark",
					"Test_", "Test",
				} {
					if strings.HasPrefix(testTagStr, src) {
						testTagStr = strings.TrimPrefix(testTagStr, src)
						changed = true
					}
				}
				if !changed {
					break
				}
			}
			instance = &dataInstance{
				tag:    testTagStr,
				source: source,
				name:   TestTag(testTagStr).Name(),
				extra:  string(matches[2]),
			}
			if source != "" {
				testTagStr = source + ":" + testTagStr
			}
			test = TestTag(testTagStr)
			if _, found := w.testNames[test]; !found {
				w.testNames[test] = instance.name
			}
			continue
		}
		if handler != "" {
			slog.Warn("Unprocessed line", "line", string(line))
		}
	}
	saveInstance([]byte(""))

	if err := scanner.Err(); err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("scan input: %w", scanner.Err())
	}

	return nil
}
