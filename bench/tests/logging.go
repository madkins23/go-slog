package tests

import (
	"bufio"
	"bytes"
	_ "embed"
	"errors"
	"io"
	"log/slog"
	"regexp"

	"github.com/madkins23/go-slog/gin"
)

// -----------------------------------------------------------------------------
// Acquire a bunch of log statements to use in Benchmark_Logging.
// The log data was generated running the server application.

//go:embed logging.txt
var logging []byte

var logDataMap [][]any

// Example line:
//
//	07:55:52 INF 200 |    9.522199ms |             ::1 | GET      "/chart.svg?tag=samber_zap&item=MemAllocs" sys=gin
var ptnTrimTimeLevel = regexp.MustCompile(`^\s*[\d:]{8}\s+\w+\s+(\d+.*?)\s*$`)

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
		msg := line.String()
		matches := ptnTrimTimeLevel.FindStringSubmatch(msg)
		if len(matches) == 2 {
			msg = matches[1]
		}
		if args, err := gin.Parse(line.String()); err != nil {
			panic("Error parsing logging traffic: " + err.Error())
		} else {
			logDataMap = append(logDataMap, args)
		}
		line.Reset()
	}
}

func logData() [][]any {
	return logDataMap
}
