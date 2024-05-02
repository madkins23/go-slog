package tests

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/madkins23/go-slog/infra"
	warning2 "github.com/madkins23/go-slog/infra/warning"
	"github.com/madkins23/go-slog/internal/data"
	"github.com/madkins23/go-slog/internal/misc"
	"github.com/madkins23/go-slog/internal/test"
)

var justTests = flag.Bool("justTests", false, "Don't run benchmarks, just tests")

// SlogBenchmarkSuite implements the benchmark test harness.
type SlogBenchmarkSuite struct {
	infra.Creator
	*warning2.Manager

	b  *testing.B
	mu sync.RWMutex
}

// NewSlogBenchmarkSuite creates a new benchmark test suite for the specified Creator.
// The handler string must match the suffix of the name of the enclosing test.
// This is:
//   - the string that will be matched from benchmark test output data,
//   - the string that will be conveyed to the parser to match against warning data, and
//   - the string that will be used as a tag for displaying server pages.
func NewSlogBenchmarkSuite(creator infra.Creator) *SlogBenchmarkSuite {
	suite := &SlogBenchmarkSuite{
		Creator: creator,
		Manager: NewWarningManager(creator.Name()),
	}
	suite.WarnOnly(warning2.NoHandlerCreation)
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
func (suite *SlogBenchmarkSuite) logger(b *Benchmark, w io.Writer) *slog.Logger {
	if b.HandlerFn != nil {
		if suite.CanMakeHandler() {
			return slog.New(b.HandlerFn(suite.NewHandler(w, b.Options)))
		} else {
			slog.Warn("b.handlerFn without suite.CanMakeHandler()")
		}
	}
	return suite.NewLogger(w, b.Options)
}

// -----------------------------------------------------------------------------

const benchmarkMethodPrefix = "Benchmark"

// Run all benchmark tests in the suite for the specified suite.
//
// This is the core algorithm for the benchmark test harness.
// It handles running all tests using reflection.
func Run(b *testing.B, suite *SlogBenchmarkSuite) {
	defer recoverAndFailOnPanic(b)

	// Capture relationship between handler name in benchmark function vs. Creator.
	// This way the handler name field is populated by the Creator name string.
	// The data will be parsed by internal/data.Benchmarks.ParseBenchmarkData() and
	// passed into Warnings.ParseWarningData().
	functionName := misc.CurrentFunctionName(benchmarkMethodPrefix)
	handler := data.HandlerTag(strings.TrimPrefix(functionName, benchmarkMethodPrefix))
	fmt.Printf("# Handler[%s]=\"%s\"\n", handler, suite.Creator.Name())

	stdoutLogger := suite.NewLogger(os.Stdout, infra.SimpleOptions())
	suite.SetB(b)
	suiteType := reflect.TypeOf(suite)
	// For each method name...
	for i := 0; i < suiteType.NumMethod(); i++ {
		method := suiteType.Method(i)
		// ...beginning with `Benchmark`:
		if strings.HasPrefix(method.Name, benchmarkMethodPrefix) {
			// Execute the method, returning a pointer to an object of class `Benchmark`.
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

			// If the `Benchmark` has a handler function
			// then the `Creator` must be able to provide a `Handler`:
			if benchmark.HandlerFn != nil && !suite.CanMakeHandler() {
				// This test requires the handler to be adjusted before creating the logger
				// but the Creator object doesn't provide a handler so skip the test.
				test.Debugf(2, ">>>     Skip:   %s\n", method.Name)
				suite.AddWarningFn(warning2.NoHandlerCreation, method.Name, "")
				// After this any benchmark with a non-nil HandlerFn must be able to make a handler.
				continue
			}

			var buffer bytes.Buffer
			// Get a logger, using the handler function if present.
			logger := suite.logger(benchmark, &buffer)
			// Run a single test using that logger.
			benchmark.BenchmarkFn(logger)
			// Track the size of the output line.
			bytesPerOp := int64(buffer.Len())

			// If the `Benchmark` has a verify function to test the log output:
			if benchmark.VerifyFn != nil {
				// Verify the output with the function.
				if err := benchmark.VerifyFn(buffer.Bytes(), nil, suite.Manager); err != nil {
					slog.Warn("Verification Error", "err", err)
				}
			}

			if *justTests {
				// The -justTests flag is set, don't do the actual benchmarks.
				continue
			}

			// TODO: If I could call the following I could haz results now?
			//       testing.Benchmark(func(b *testing.B) {
			b.Run(method.Name, func(b *testing.B) {
				var count test.CountWriter
				function := benchmark.BenchmarkFn
				// Capture warnings from a single run.
				if test.DebugLevel() > 0 {
					// Print the log record to STDOUT.
					function(stdoutLogger)
				}
				// Get a logger, using the handler function if present.
				// NOTE: the creation of the logger,
				//       which may involve Handler.WithAttrs() and/or Handler.WithGroup(),
				//       is NOT counted towards results.
				logger := suite.logger(benchmark, &count)
				// Now move on to the actual test.
				b.ReportAllocs()
				b.SetBytes(bytesPerOp)
				b.ResetTimer()
				// The Go test harness is used to run the `Benchmark` test function
				// in parallel in ever-larger batches until enough testing has been done.
				// The test harness emits a line of data with results of the test.
				b.RunParallel(func(pb *testing.PB) {
					for pb.Next() {
						function(logger)
					}
				})
				b.StopTimer()
				if !benchmark.DontCount && b.N != int(count.Written()) {
					b.Fatalf("Mismatch in log write count. Expected: %d, Actual: %d",
						b.N, count.Written())
				}
			})
		}
	}
}
