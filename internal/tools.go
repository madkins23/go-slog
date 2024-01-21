//go:build bench

package internal

// This file exists to keep the go:generate tool(s) in the module file.
// Since they are otherwise not referenced in an import statement,
// (just the go:generate statement) running
//	go tidy
// will remove them from go.mod.

import (
	_ "github.com/dmarkham/enumer"
)
