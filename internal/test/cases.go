package test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// Case defines an interface for test case definitions.
// The interface provides functions to access common data items.
// Cases are read from JSON files in 'test' subdirectory of the specified test directory.
// The JSON files must contain the following common fields:
//
//	Name                   test name
//	Source                 test source string
//
// If the Source field is empty the test file name is used to read the source
// data from a JSON file in the 'json' subdirectory of the specified test directory.
type Case interface {
	Name() string
	setDir(dir string)
	setBase(base string)
	Source(t *testing.T) string
}

// Make sure the BasicCase struct implements the Case interface.
var _ Case = &BasicCase{}

// BasicCase provides the common test definition items and methods.
// Compose the specific test case struct with this struct to inherit.
// Normally the BasicName and BasicSource fields would be private,
// but they can't be because of the JSON unmarshaler.
// Do not reference them externally.
type BasicCase struct {
	BasicName   string `json:"Name"`
	BasicSource string `json:"Source"`
	base        string
	dir         string
}

// Files returns a list of file names found in the 'test' subdirectory of the specified directory.
// If the directory is not provided (the string is empty) then 'testdata' is used.
func Files(t *testing.T, dir string) []string {
	if dir == "" {
		dir = "testdata"
	}
	paths, err := filepath.Glob(filepath.Join(dir, "test", "*.json"))
	require.NoError(t, err, "Get list of test files")
	files := make([]string, len(paths))
	for i, path := range paths {
		files[i] = filepath.Base(path)
	}
	return files
}

// Load reads a test case JSON file and unmarshals it into the provided Case instantiation.
// The pathname is constructed from the named file in the 'test' subdirectory of the specified directory:
//
//	<dir>/test/<file>
//
// If the directory is not provided (the string is empty) then 'testdata' is used.
func Load(t *testing.T, dir string, file string, tc Case) {
	if dir == "" {
		dir = "testdata"
	}
	bytes, err := os.ReadFile(filepath.Join(dir, "test", file))
	require.NoError(t, err, "Read test case")
	require.NoError(t, json.Unmarshal(bytes, &tc), "Unmarshal test case")
	tc.setBase(file)
	tc.setDir(dir)
}

// Name returns the test case name.
func (bc *BasicCase) Name() string {
	if bc.BasicName == "" {
		bc.BasicName = bc.base[:len(bc.base)-len(filepath.Ext(bc.base))]
	}
	return bc.BasicName
}

// Source returns the test case source string.
// If the Source field was empty the data will be loaded from
// The pathname constructed from the named file in the 'json' subdirectory of the specified directory:
//
//	<dir>/json/<file>
//
// where the dir was specified in the preceding Load call and saved to the basic test case struct.
func (bc *BasicCase) Source(t *testing.T) string {
	if len(bc.BasicSource) < 1 {
		src, err := os.ReadFile(filepath.Join(bc.dir, "json", bc.base))
		require.NoError(t, err, "Read Source from file %s", bc.base)
		bc.BasicSource = string(src)
		require.NotEmpty(t, bc.BasicSource, "File %s is empty", bc.base)
	}

	return bc.BasicSource
}

func (bc *BasicCase) setBase(base string) {
	bc.base = base
}

func (bc *BasicCase) setDir(dir string) {
	bc.dir = dir
}
