package json

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Set -debug flag to show extra print statements.
// Command line setting:
//
//	go test ./... -args -debug
var debug = flag.Bool("debug", false, "Show debug statements")

type testCase struct {
	Name       string
	FldList    []string
	Duplicates map[string]uint
	File       string
	Source     string
}

func Test_FieldCounter(t *testing.T) {
	paths, err := filepath.Glob(filepath.Join("testdata/test", "*.json"))
	require.NoError(t, err, "Get list of test files")
	for _, path := range paths {
		bytes, err := os.ReadFile(path)
		require.NoError(t, err, "Read test case")
		var test testCase
		require.NoError(t, json.Unmarshal(bytes, &test), "Unmarshal test case")
		base := filepath.Base(path)
		name := test.Name
		if name == "" {
			name = base[:len(base)-len(filepath.Ext(base))]
		}
		source := []byte(test.Source)
		if len(source) < 1 {
			source, err = os.ReadFile("testdata/json/" + base)
			require.NoError(t, err, "Read source from file %s", base)
			require.NotEmpty(t, source, "File %s is empty", base)
		}
		counter := NewFieldCounter(source)
		require.NoError(t, counter.Parse())
		assert.Equal(t, uint(len(test.FldList)), counter.NumFields())
		assert.Equal(t, test.FldList, counter.Fields())
		assert.Equal(t, test.Duplicates, counter.Duplicates())
		dbgFmt("Test: %s %s", name, counter)
	}
}

// dbgFmt will only print the specified data if the -debug command flag is set.
func dbgFmt(format string, args ...interface{}) {
	if *debug {
		fmt.Printf(">>> "+format+"\n", args...)
	}
}
