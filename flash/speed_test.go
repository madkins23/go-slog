package flash

import (
	"runtime"
	"testing"
)

// TODO: Test loadSource() vs. newSource() via benchmark.

var pc uintptr

func init() {
	pc, _, _, _ = runtime.Caller(0)
}

/* Compare loadSource() (with variable on stack) to newSource()/reuseSource().
 * Manual review of results says loadSource() is maybe 10% faster.
 */

func BenchmarkLoadSource(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var src source
			loadSource(pc, &src)
		}
	})
}

func BenchmarkNewSource(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var src = newSource(pc)
			reuseSource(src)
		}
	})
}
