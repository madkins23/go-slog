package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/madkins23/go-slog/creator/slogjson"
	"github.com/madkins23/go-slog/infra"
)

func ExampleBenchmark() {
	bm := &Benchmark{
		Options: infra.SimpleOptions(),
		BenchmarkFn: func(logger *slog.Logger) {
			logger.Info(message)
		},
		VerifyFn: matcher("Simple", expectedBasic()),
	}
	cr := slogjson.Creator()
	var buffer bytes.Buffer
	logger := slog.New(cr.NewHandler(&buffer, bm.Options))
	bm.BenchmarkFn(logger)
	var logMap map[string]any
	_ = json.Unmarshal(buffer.Bytes(), &logMap)
	fmt.Println(logMap["msg"])
	// Output: This is a message
}
