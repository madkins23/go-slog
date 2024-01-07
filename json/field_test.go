package json

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/madkins23/go-slog/test"
)

type testCase struct {
	test.BasicCase
	FldList    []string
	Duplicates map[string]uint
}

func Test_FieldCounter(t *testing.T) {
	for _, file := range test.Files(t, "") {
		var tc testCase
		test.Load(t, "", file, &tc)
		t.Run(tc.Name(), func(t *testing.T) {
			test.Debugf(2, "JSON: %s", tc.Source(t))
			counter := NewFieldCounter([]byte(tc.Source(t)))
			require.NoError(t, counter.Parse())
			assert.Equal(t, uint(len(tc.FldList)), counter.NumFields())
			assert.Equal(t, tc.FldList, counter.Fields())
			assert.Equal(t, tc.Duplicates, counter.Duplicates())
			test.Debugf(1, "Test: %s %s", tc.Name(), counter)
		})
	}
}
