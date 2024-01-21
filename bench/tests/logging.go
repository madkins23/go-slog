package tests

import (
	"bufio"
	"bytes"
	_ "embed"
	"errors"
	"io"
	"log/slog"
	"regexp"
	"strconv"
	"strings"
)

// -----------------------------------------------------------------------------
// Acquire a bunch of log statements to use in Benchmark_Logging.
// The log data was generated running the server application.

//go:embed logging.txt
var logging []byte

var logDataMap [][]any

var (
	ptnCode  = regexp.MustCompile(`\s(\d+)\s*$`)
	ptnSplit = regexp.MustCompile(`\s+`)
)

// Read log data from embedded data, construct array of arg arrays for logging.
func init() {
	reader := bufio.NewReader(bytes.NewReader(logging))
	var line bytes.Buffer
	for {
		if chunk, isPrefix, err := reader.ReadLine(); errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			slog.Warn("Error reading logging data line", "err", err)
		} else {
			line.Write(chunk)
			if isPrefix {
				continue
			}
		}
		// Example line:
		//  07:55:52 INF 200 |    9.522199ms |             ::1 | GET      "/chart.svg?tag=samber_zap&item=MemAllocs" sys=gin
		if parts := strings.Split(string(line.Bytes()), "|"); len(parts) != 4 {
			slog.Warn("Wrong number of parts", "num", len(parts), "line", line, "func", "getLogData")
		} else {
			var args []any
			// Parse parts[0]:
			//  07:55:52 INF 200
			if matches := ptnCode.FindStringSubmatch(parts[0]); len(matches) != 2 {
				slog.Warn("Unable to parse code", "from", parts[0])
			} else if num, err := strconv.ParseInt(matches[1], 10, 64); err != nil {
				slog.Warn("Unable to parse int", "from", parts[0], "func", "getLogData")
			} else {
				args = append(args, "code", num)
			}
			args = append(args, "duration", strings.Trim(parts[1], " "))
			// Ignore parts[2] (::1) since I don't know what it is.
			// Parse parts[3]:
			//  GET      "/chart.svg?tag=samber_zap&item=MemAllocs" sys=gin
			parts = ptnSplit.Split(strings.Trim(parts[3], " "), -1)
			if len(parts) == 3 {
				args = append(args, "method", parts[0])
				args = append(args, "url", strings.Trim(parts[1], "\""))
			}
			args = append(args, "sys", "gin")
			logDataMap = append(logDataMap, args)
		}
		line.Reset()
	}
}

func logData() [][]any {
	return logDataMap
}
