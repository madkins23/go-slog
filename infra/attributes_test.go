package infra

import (
	"log/slog"
	"testing"
)

// -----------------------------------------------------------------------------
// scripts/comp infra EmptyAttr

var emptyAttr = slog.Attr{}

// copyEmptyAttr duplicates the original EmptyAttr() function.
func copyEmptyAttr() slog.Attr {
	return emptyAttr
}

// BenchmarkEmptyAttrCopy tests performance of returning an empty attribute
// from a local variable (which copies it along the way).
func BenchmarkEmptyAttrCopy(b *testing.B) {
	var empty slog.Attr
	for i := 0; i < b.N; i++ {
		empty = copyEmptyAttr()
	}
	empty.Equal(emptyAttr)
}

// BenchmarkEmptyAttrMake tests performance of returning a newly made
// empty attribute each time.
func BenchmarkEmptyAttrMake(b *testing.B) {
	var empty slog.Attr
	for i := 0; i < b.N; i++ {
		empty = EmptyAttr()
	}
	empty.Equal(emptyAttr)
}
