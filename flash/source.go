package flash

import "runtime"

var sourcePool = newGenPool[source]()

type source struct {
	Function string `json:"function"`
	File     string `json:"file"`
	Line     int    `json:"line"`
}

func newSource(pc uintptr) (*source, func()) {
	fs := runtime.CallersFrames([]uintptr{pc})
	f, _ := fs.Next()
	src, rtnFn := sourcePool.borrow()
	src.File = f.File
	src.Function = f.Function
	src.Line = f.Line
	return src, rtnFn
}
