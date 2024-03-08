package flash

import "runtime"

type source struct {
	Function string `json:"function"`
	File     string `json:"file"`
	Line     int    `json:"line"`
}

func loadSource(pc uintptr, src *source) {
	fs := runtime.CallersFrames([]uintptr{pc})
	f, _ := fs.Next()
	src.File = f.File
	src.Function = f.Function
	src.Line = f.Line
}

// TODO: are newSource() and reuseSource() in use?
//       See BenchmarkLoadSource() and BenchmarkNewSource() in speed_test.go.

var sourcePool = newGenPool[source]()

func newSource(pc uintptr) *source {
	fs := runtime.CallersFrames([]uintptr{pc})
	f, _ := fs.Next()
	src := sourcePool.get()
	src.File = f.File
	src.Function = f.Function
	src.Line = f.Line
	return src
}

func reuseSource(src *source) {
	sourcePool.put(src)
}
