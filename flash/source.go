package flash

import "runtime"

type source struct {
	Function string `json:"function"`
	File     string `json:"file"`
	Line     int    `json:"line"`
}

func newSource(pc uintptr) *source {
	fs := runtime.CallersFrames([]uintptr{pc})
	f, _ := fs.Next()
	src := sourcePool.get()
	src.File = f.File
	src.Function = f.Function
	src.Line = f.Line
	return src
}
