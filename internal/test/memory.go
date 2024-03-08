package test

import (
	"fmt"
	"runtime"
)

var seps = map[bool]string{
	false: "->",
	true:  "=>",
}

func MemoryTrap(level uint, name string, fn func()) {
	if *debug >= level {
		var memBefore, memAfter runtime.MemStats
		var started bool
		runtime.ReadMemStats(&memBefore)
		fn()
		runtime.ReadMemStats(&memAfter)
		for i := 0; i < 61; i++ {
			if memBefore.BySize[i] != memAfter.BySize[i] {
				if !started {
					fmt.Printf(">>> Memory Trap %s\n", name)
					started = true
				}
				fmt.Printf(">>>   [%02d] %5d %s %5d\n", i,
					memBefore.BySize[i], seps[memBefore.BySize[i] != memAfter.BySize[i]], memAfter.BySize[i])
			}
		}
	} else {
		fn()
	}
}

func MemoryTrapErr(level uint, name string, fn func() error) error {
	var err error
	if *debug >= level {
		var memBefore, memAfter runtime.MemStats
		var started bool
		runtime.ReadMemStats(&memBefore)
		err = fn()
		runtime.ReadMemStats(&memAfter)
		for i := 0; i < 61; i++ {
			if memBefore.BySize[i] != memAfter.BySize[i] {
				if !started {
					fmt.Printf(">>> Memory Trap %s\n", name)
					started = true
				}
				fmt.Printf(">>>   [%02d] %5d %s %5d\n", i,
					memBefore.BySize[i], seps[memBefore.BySize[i] != memAfter.BySize[i]], memAfter.BySize[i])
			}
		}
	} else {
		err = fn()
	}
	if err != nil {
		fmt.Printf(">>>   error: %s\n", err.Error())
	}
	return err
}
