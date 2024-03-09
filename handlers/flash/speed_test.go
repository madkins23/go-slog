package flash

import (
	"bytes"
	"log/slog"
	"runtime"
	"testing"

	"github.com/madkins23/go-slog/handlers/sloggy"
	"github.com/madkins23/go-slog/handlers/sloggy/test"
)

// -----------------------------------------------------------------------------
// Compare sloggy.composer (using bytes.Buffer) with flash.composer (using byte array appends).

// BenchmarkComposeArray composes attributes by appending to byte arrays.
func BenchmarkComposeArray(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c := newComposer([]byte{}, false, nil, nil)
			if err := c.addAttributes(test.Attributes); err != nil {
				b.Errorf("add attributes: %s", err.Error())
			}
			buffer := c.getBytes()
			if len(buffer) < 1 {
				b.Error("Empty buffer")
			}
		}
	})
}

// BenchmarkComposeBuffer composes attributes using bytes.Buffer.
func BenchmarkComposeBuffer(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var buffer bytes.Buffer
			if err := sloggy.ComposeAttributes(&buffer, test.Attributes); err != nil {
				b.Errorf("Composition failed: %s", err.Error())
			}
			if buffer.Len() < 1 {
				b.Error("Empty buffer")
			}
		}
	})
}

// -----------------------------------------------------------------------------
// Compare memory allocation versus memory buffer pools.

var memTestLen uint = 1024

// BenchmarkMemoryAllocation allocates memory and leaves it to the garbage collector.
func BenchmarkMemoryAllocation(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			x := make([]byte, 0, memTestLen)
			x = append(x, 123)
			// Let garbage collection do its thing.
		}
	})
}

// BenchmarkMemoryPools acquires memory from a pool and returns it for reuse.
func BenchmarkMemoryPools(b *testing.B) {
	var logPool = newArrayPool[byte](memTestLen)
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			x := logPool.get()
			x = append(x, 123)
			logPool.put(x)
		}
	})
}

// -----------------------------------------------------------------------------
// Compare adding attributes to a composer one at a time with addAttribute vs.
// putting them in an array and adding them all at once with addAttributes.
// In the latter case the arrays are managed and reused via a sync.Pool.

// BenchmarkBasicManual adds attributes to the composer one at a time using addAttribute.
func BenchmarkBasicManual(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c := newComposer([]byte{}, false, nil, nil)
			if err := c.addAttribute(slog.Time(slog.TimeKey, test.Now)); err != nil {
				b.Errorf("add time: %s", err.Error())
			}
			if err := c.addAttribute(slog.String(slog.LevelKey, test.Level.String())); err != nil {
				b.Errorf("add level: %s", err.Error())
			}
			if err := c.addAttribute(slog.String(slog.MessageKey, test.Message)); err != nil {
				b.Errorf("add message: %s", err.Error())
			}
			reuseComposer(c)
		}
	})
}

var basicPool = newArrayPool[slog.Attr](lenBasic)

// BenchmarkBasicMultiple adds attributes to an array and then uses
// addAttributes to send them to the composer all at once.
func BenchmarkBasicMultiple(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			basic := basicPool.get()[:0]
			c := newComposer([]byte{}, false, nil, nil)
			basic = append(basic, slog.Time(slog.TimeKey, test.Now))
			basic = append(basic, slog.String(slog.LevelKey, test.Level.String()))
			basic = append(basic, slog.String(slog.MessageKey, test.Message))
			if err := c.addAttributes(basic); err != nil {
				b.Errorf("add attributes: %s", err.Error())
			}
			basicPool.put(basic)
		}
	})
}

// -----------------------------------------------------------------------------
// Compare loadSource() (with variable on stack) to newSource()/reuseSource().

// BenchmarkSourceLoad uses a local variable (presumably on the stack) and
// passes a pointer to that variable for loadSource to fill its fields.
func BenchmarkSourceLoad(b *testing.B) {
	pc, _, _, _ := runtime.Caller(0)
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var src source
			loadSource(pc, &src)
		}
	})
}

// BenchmarkSourceNewReuse acquires a pointer to a new, properly configured,
// source object using newSource and returns it afterward using reuseSource.
func BenchmarkSourceNewReuse(b *testing.B) {
	pc, _, _, _ := runtime.Caller(0)
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var src = newSource(pc)
			reuseSource(src)
		}
	})
}

// -----------------------------------------------------------------------------
// Compare using a cut-out before calling for attribute resolution.

// BenchmarkResolveAlways always calls slog.Attr.Value.Resolve.
func BenchmarkResolveAlways(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for _, attr := range test.Attributes {
				attr.Value.Resolve()
			}
		}
	})
}

// BenchmarkResolveConditional only calls attr.Value.Resolve if attr is a LogValuer.
func BenchmarkResolveConditional(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for _, attr := range test.Attributes {
				if attr.Value.Kind() == slog.KindLogValuer {
					attr.Value.Resolve()
				}
			}
		}
	})
}
