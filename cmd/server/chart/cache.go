package chart

import "sync"

var (
	Cache      = make(map[string][]byte)
	CacheMutex sync.Mutex
)
