package tests

import (
	"bytes"
	"flag"
	"io"
	"log/slog"
	"os"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/madkins23/go-slog/infra"
	"github.com/madkins23/go-slog/internal/test"
	"github.com/madkins23/go-slog/warning"
)

var justTests = flag.Bool("justTests", false, "Don't run benchmarks, just tests")

type SlogBenchmarkSuite struct {
	infra.Creator
	*test.WarningManager

	b  *testing.B
	mu sync.RWMutex
}

func NewSlogBenchmarkSuite(creator infra.Creator) *SlogBenchmarkSuite {
	suite := &SlogBenchmarkSuite{
		Creator:        creator,
		WarningManager: NewWarningManager(creator.Name()),
	}
	suite.WarnOnly(warning.NoHandlerCreation)
	return suite
}

// B retrieves the current *testing.B context.
func (suite *SlogBenchmarkSuite) B() *testing.B {
	suite.mu.RLock()
	defer suite.mu.RUnlock()
	return suite.b
}

// SetB sets the current *testing.B context.
func (suite *SlogBenchmarkSuite) SetB(b *testing.B) {
	suite.mu.Lock()
	defer suite.mu.Unlock()
	suite.b = b
}

// logger for testing with handler tweaks if HandlerFn is specified in the benchmark.
// Assumes that if HandlerFn is present CanMakeHandler() is true.
func (suite *SlogBenchmarkSuite) logger(b *Benchmark, w io.Writer) *slog.Logger {
	if b.HandlerFn != nil {
		// Since we're here we know that CanMakeHandler() must be true.
		// Otherwise we would have hit the continue above.
		return slog.New(b.HandlerFn(suite.NewHandler(w, b.Options)))
	}
	return suite.NewLogger(w, b.Options)
}

// -----------------------------------------------------------------------------

const benchmarkMethodPrefix = "Benchmark"

func Run(b *testing.B, suite *SlogBenchmarkSuite) {
	defer recoverAndFailOnPanic(b)

	stdoutLogger := suite.NewLogger(os.Stdout, infra.SimpleOptions())
	suite.SetB(b)
	suiteType := reflect.TypeOf(suite)
	for i := 0; i < suiteType.NumMethod(); i++ {
		method := suiteType.Method(i)
		if strings.HasPrefix(method.Name, benchmarkMethodPrefix) {
			results := method.Func.Call([]reflect.Value{reflect.ValueOf(suite)})
			if len(results) < 1 {
				b.Fatalf("Unable to acquire benchmark definition")
			}
			benchmark, ok := results[0].Interface().(*Benchmark)
			if !ok {
				b.Fatalf("Could not convert benchmark definition %v", results[0].Interface())
			}

			if benchmark.BenchmarkFn == nil {
				slog.Error("No benchmark function", "method", method.Name)
				continue
			}

			if benchmark.HandlerFn != nil && !suite.CanMakeHandler() {
				// This test requires the handler to be adjusted before creating the logger
				// but the Creator object doesn't provide a handler so skip the test.
				test.Debugf(2, ">>>     Skip:   %s\n", method.Name)
				suite.AddWarningFn(warning.NoHandlerCreation, method.Name, "")
				// After this any benchmark with a non-nil HandlerFn must be able to make a handler.
				continue
			}

			if benchmark.VerifyFn != nil {
				var buffer bytes.Buffer
				logger := suite.logger(benchmark, &buffer)
				benchmark.BenchmarkFn(logger)
				if err := benchmark.VerifyFn(buffer.Bytes(), nil, suite.WarningManager); err != nil {
					slog.Warn("Verification Error", "err", err)
				}
			}

			if *justTests {
				continue
			}

			test.Debugf(2, ">>>     Method: %s\n", method.Name)
			// TODO: If I could call the following I could haz results now?
			//       testing.Benchmark(func(b *testing.B) {
			b.Run(method.Name, func(b *testing.B) {
				var count test.CountWriter
				function := benchmark.BenchmarkFn
				logger := suite.logger(benchmark, &count)
				if test.DebugLevel() > 0 {
					// Print the log record to STDOUT.
					function(stdoutLogger)
				}
				b.ReportAllocs()
				// TODO: This doesn't seem to make any difference?
				b.ResetTimer()
				b.RunParallel(func(pb *testing.PB) {
					for pb.Next() {
						function(logger)
					}
				})
				b.StopTimer()
				b.SetBytes(int64(count.Bytes()))
				if !benchmark.DontCount && b.N != int(count.Written()) {
					b.Fatalf("Mismatch in log write count. Expected: %d, Actual: %d",
						b.N, count.Written())
				}
			})
		}
	}
}
