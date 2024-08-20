package axis

import (
	"bytes"
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/madkins23/go-slog/internal/data"
	"github.com/madkins23/go-slog/internal/scoring/axis/bench"
	"github.com/madkins23/go-slog/internal/scoring/score"
)

var benchScores = map[data.HandlerTag]score.Value{
	"ChanchalZap":   90.40508042083333,
	"PhsymZerolog":  99.09420653858332,
	"SamberLogrus":  8.550007116333333,
	"SamberZap":     56.24375427375,
	"SamberZerolog": 59.9383113595,
	"SlogJSON":      97.29024068808333,
}

//go:embed testdata/bench.txt
var benchTxt string

// TestSetup is intended to verify that the data parsing/weighting algorithms don't drift.
func TestSetup(t *testing.T) {
	dbm := data.NewBenchmarks()
	require.NoError(t, dbm.ParseBenchmarkData(bytes.NewBuffer([]byte(benchTxt))))
	sbm := NewBenchmarks(defaultBenchmarkScoreWeight, "<p>Test!!!</p>", nil)
	require.NoError(t, sbm.Setup(dbm, nil))
	for _, hdlr := range dbm.HandlerTags() {
		assert.Equal(t, benchScores[hdlr], sbm.ScoreFor(hdlr), "Handler: "+hdlr)
	}
}

var defaultBenchmarkScoreWeight = map[bench.Weight]uint{
	bench.Allocations: 1,
	bench.AllocBytes:  2,
	bench.Nanoseconds: 3,
}
