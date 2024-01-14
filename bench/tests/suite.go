package tests

import (
	"log/slog"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/madkins23/go-slog/infra"
	"github.com/madkins23/go-slog/test"
)

type SlogBenchmarkSuite struct {
	creator infra.Creator

	b  *testing.B
	mu sync.RWMutex
}

func NewSlogBenchmarkSuite(creator infra.Creator) *SlogBenchmarkSuite {
	return &SlogBenchmarkSuite{creator: creator}
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

	suite.SetB(b)

	suiteType := reflect.TypeOf(suite)
	test.Debugf(1, ">>>   Suite: %s\n", suiteType.Elem().Name())
	for i := 0; i < suiteType.NumMethod(); i++ {
		method := suiteType.Method(i)
		if strings.HasPrefix(method.Name, benchmarkMethodPrefix) {
			test.Debugf(2, ">>>     Method: %s\n", method.Name)
			// TODO: If I could call the following I could haz results now?
			//       testing.Benchmark(func(b *testing.B) {
			b.Run(method.Name, func(b *testing.B) {
				var count infra.CountWriter
				results := method.Func.Call([]reflect.Value{reflect.ValueOf(suite)})
				if len(results) < 1 {
					b.Fatalf("No results returned from benchmark")
				}
				benchmark, ok := results[0].Interface().(Benchmark)
				if !ok {
					b.Fatalf("Could not convert benchmark result %v", results[0].Interface())
				}
				function := benchmark.Function()
				logger := slog.New(suite.creator.NewHandle(&count, benchmark.Options()))
				b.ReportAllocs()
				b.ResetTimer()
				b.RunParallel(func(pb *testing.PB) {
					for pb.Next() {
						function(logger)
					}
				})
				b.StopTimer()
				if b.N != int(count.Written()) {
					b.Fatalf("Mismatch in log write count. Expected: %d, Actual: %d",
						b.N, count.Written())
				}
			})
		}
	}
}

// -----------------------------------------------------------------------------

type BenchmarkFn func(logger *slog.Logger)

type Benchmark interface {
	Options() *slog.HandlerOptions
	Function() BenchmarkFn
}

var _ Benchmark = &benchmark{}

type benchmark struct {
	options *slog.HandlerOptions
	fn      BenchmarkFn
}

func (b *benchmark) Options() *slog.HandlerOptions {
	return b.options
}

func (b *benchmark) Function() BenchmarkFn {
	return b.fn
}

func NewBenchmark(options *slog.HandlerOptions, fn BenchmarkFn) Benchmark {
	return &benchmark{
		options: options,
		fn:      fn,
	}
}
