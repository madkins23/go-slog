package flash

import (
	"log/slog"
	"runtime"
	"testing"
	"time"
)

/* Compare adding attributes to a composer one at a time with addAttribute vs.
 * putting them in an array and adding them all at once with addAttributes.
 * In the latter case the arrays are managed and reused via a sync.Pool.
 * Manual review of results says adding them one at a time is maybe 30% faster.
 */

// BenchmarkBasicManual adds attributes to the composer one at a time using addAttribute.
func BenchmarkBasicManual(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		now := time.Now()
		for pb.Next() {
			c := newComposer([]byte{}, false, nil, nil)
			if err := c.addAttribute(slog.Time(slog.TimeKey, now)); err != nil {
				b.Errorf("add time: %s", err.Error())
			}
			if err := c.addAttribute(slog.String(slog.LevelKey, "INFO")); err != nil {
				b.Errorf("add level: %s", err.Error())
			}
			if err := c.addAttribute(slog.String(slog.MessageKey, "message")); err != nil {
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
	b.RunParallel(func(pb *testing.PB) {
		now := time.Now()
		for pb.Next() {
			basic := basicPool.get()[:0]
			c := newComposer([]byte{}, false, nil, nil)
			basic = append(basic, slog.Time(slog.TimeKey, now))
			basic = append(basic, slog.String(slog.LevelKey, "INFO"))
			basic = append(basic, slog.String(slog.MessageKey, "message"))
			if err := c.addAttributes(basic); err != nil {
				b.Errorf("add attributes: %s", err.Error())
			}
			basicPool.put(basic)
		}
	})
}

/* Compare loadSource() (with variable on stack) to newSource()/reuseSource().
 * Manual review of results says loadSource() is maybe 10% faster.
 */

// BenchmarkSourceLoad uses a local variable (presumably on the stack) and
// passes a pointer to that variable for loadSource to fill its fields.
func BenchmarkSourceLoad(b *testing.B) {
	pc, _, _, _ := runtime.Caller(0)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var src source
			loadSource(pc, &src)
		}
	})
}

// BenchmarkSourceNewReuse acquires a pointer to a new, properly configured,
// source object using newSource and returns it afterwards using reuseSource.
func BenchmarkSourceNewReuse(b *testing.B) {
	pc, _, _, _ := runtime.Caller(0)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var src = newSource(pc)
			reuseSource(src)
		}
	})
}
