package tests

import (
	"log/slog"
	"os"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/madkins23/go-slog/infra"
	"github.com/madkins23/go-slog/internal/test"
)

type SlogBenchmarkSuite struct {
	infra.Creator
	*infra.WarningManager

	b  *testing.B
	mu sync.RWMutex
}

func NewSlogBenchmarkSuite(creator infra.Creator) *SlogBenchmarkSuite {
	suite := &SlogBenchmarkSuite{
		Creator:        creator,
		WarningManager: NewWarningManager(creator.Name()),
	}
	suite.WarnOnly(WarnNoHandlerCreation)
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
				b.Fatalf("No results returned from benchmark")
			}
			benchmark, ok := results[0].Interface().(Benchmark)
			if !ok {
				b.Fatalf("Could not convert benchmark result %v", results[0].Interface())
			}
			if benchmark.HandlerFn() != nil && !suite.CanMakeHandler() {
				// This test requires the handler to be adjusted before creating the logger
				// but the Creator object doesn't provide a handler so skip the test.
				test.Debugf(2, ">>>     Skip:   %s\n", method.Name)
				suite.AddWarningFn(WarnNoHandlerCreation, method.Name, "")
				continue
			}
			test.Debugf(2, ">>>     Method: %s\n", method.Name)
			// TODO: If I could call the following I could haz results now?
			//       testing.Benchmark(func(b *testing.B) {
			b.Run(method.Name, func(b *testing.B) {
				var count test.CountWriter
				function := benchmark.Function()
				var logger *slog.Logger
				if benchmark.HandlerFn() != nil && suite.CanMakeHandler() {
					logger = slog.New(benchmark.HandlerFn()(
						suite.NewHandler(&count, benchmark.Options())))
				} else {
					logger = suite.NewLogger(&count, benchmark.Options())
				}
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
				if !benchmark.DontCount() && b.N != int(count.Written()) {
					b.Fatalf("Mismatch in log write count. Expected: %d, Actual: %d",
						b.N, count.Written())
				}
			})
		}
	}
}

// -----------------------------------------------------------------------------

type HandlerFn func(handler slog.Handler) slog.Handler
type BenchmarkFn func(logger *slog.Logger)

type Benchmark interface {
	Options() *slog.HandlerOptions
	HandlerFn() HandlerFn
	Function() BenchmarkFn
	DontCount() bool
	SetDontCount(bool)
}

var _ Benchmark = &benchmark{}

type benchmark struct {
	options     *slog.HandlerOptions
	handlerFn   HandlerFn
	benchmarkFn BenchmarkFn
	dontCount   bool
}

func NewBenchmark(options *slog.HandlerOptions, fn BenchmarkFn, handlerFn HandlerFn) Benchmark {
	return &benchmark{
		options:     options,
		benchmarkFn: fn,
		handlerFn:   handlerFn,
	}
}

func (b *benchmark) DontCount() bool {
	return b.dontCount
}

func (b *benchmark) SetDontCount(dontCount bool) {
	b.dontCount = dontCount
}

func (b *benchmark) Options() *slog.HandlerOptions {
	return b.options
}

func (b *benchmark) Function() BenchmarkFn {
	return b.benchmarkFn
}

func (b *benchmark) HandlerFn() HandlerFn {
	return b.handlerFn
}
