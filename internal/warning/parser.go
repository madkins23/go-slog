package warning

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

	"github.com/madkins23/go-slog/warning"
)

// -----------------------------------------------------------------------------

var (
	ptnWarningsFor = regexp.MustCompile(`^\s*Warnings\s+for\s+(.*):\s*$`)
	ptnLevel       = regexp.MustCompile(`^\s*(\S+)\s*$`)
	ptnWarning     = regexp.MustCompile(`^\s*\d+\s+\[(.*)]\s+(.*?)\s*$`)
	ptnInstance    = regexp.MustCompile(`^\s*(\S+)(?::\s*(.*?))?\s*$`)
	ptnLogLine     = regexp.MustCompile(`^\s*\{`)
	ptnByWarning   = regexp.MustCompile(`^\s*Handlers\s+by\s+warning:\s*$`)
)

// ParseWarningData parses warning data from the output of benchmark and verification testing.
// The data will be loaded from os.Stdin unless the -bench=<path> flag is set
// in which case the data will be loaded from the specified path.
func (d *Data) ParseWarningData(in io.Reader) error {
	var err error
	if *verifyFile != "" {
		if in, err = os.Open(*verifyFile); err != nil {
			return fmt.Errorf("open --verify=%s: %s\n", *verifyFile, err)
		}
	}
	scanner := bufio.NewScanner(in)

	var handler HandlerTag
	var level warning.Level
	var dWarning *dataWarning
	var instance *dataInstance
	saveInstance := func(line []byte) {
		if instance != nil {
			if dWarning == nil {
				slog.Warn("Nil dWarning", "line", line, "instance", instance)
			} else {
				dWarning.AddInstance(instance)
				d.findTest(TestTag(instance.name), level, dWarning.warning.name).AddInstance(
					&dataInstance{
						name:    string(handler),
						extra:   instance.extra,
						logLine: instance.logLine,
					})
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
			if d.handlerNames == nil {
				d.handlerNames = make(map[HandlerTag]string)
			}
			parts := strings.Split(string(handler), "/")
			for i, part := range parts {
				if len(part) > 0 {
					parts[i] = strings.ToUpper(part[:1]) + part[1:]
				}
			}
			d.handlerNames[handler] = strings.Join(parts, " ")
			continue
		}
		if ptnByWarning.Match(line) {
			// End of data.
			break
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
			dWarning = d.findHandler(handler, level, warningName)
			dWarning.warning.description = string(matches[2])
			instance = nil
			continue
		}
		if ptnLogLine.Match(line) {
			instance.logLine = string(line)
			// Attempt to pretty-print the log line.
			var jm map[string]any
			if json.Unmarshal(line, &jm) == nil {
				if indented, err := json.MarshalIndent(jm, "", "\t"); err == nil {
					instance.logLine = string(indented)
				}
			}
			continue
		}
		if matches := ptnInstance.FindSubmatch(line); len(matches) == 3 {
			saveInstance(line)
			instance = &dataInstance{
				name:  string(matches[1]),
				extra: string(matches[2]),
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
