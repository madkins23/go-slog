package axis

import (
	"bytes"
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/madkins23/go-slog/internal/data"
	"github.com/madkins23/go-slog/internal/scoring/axis/bench"
	"github.com/madkins23/go-slog/internal/scoring/axis/common"
	"github.com/madkins23/go-slog/internal/scoring/score"
)

var (
	scores = map[data.HandlerTag]score.Value{
		"ChanchalZap":   92.44847896956676,
		"PhsymZerolog":  99.09773449406401,
		"SamberLogrus":  6.50826292891853,
		"SamberZap":     54.92942528243242,
		"SamberZerolog": 58.77332449577939,
		"SlogJSON":      97.21913472265209,
	}
	scoresOriginal = map[data.HandlerTag]score.Value{
		"ChanchalZap":   90.40508042083333,
		"PhsymZerolog":  99.09420653858332,
		"SamberLogrus":  8.550007116333333,
		"SamberZap":     56.24375427375,
		"SamberZerolog": 59.9383113595,
		"SlogJSON":      97.29024068808333,
	}
	scoresByData = map[data.HandlerTag]score.Value{
		"ChanchalZap":   94.49187751833354,
		"PhsymZerolog":  99.1012624495175,
		"SamberLogrus":  4.4665187415015675,
		"SamberZap":     53.61509629114256,
		"SamberZerolog": 57.608337632072725,
		"SlogJSON":      97.14802875723512,
	}
)

//go:embed testdata/bench.txt
var benchTxt string

// TestSetup is intended to verify that the data parsing/weighting algorithms don't drift.
func TestSetup(t *testing.T) {
	dbm := data.NewBenchmarks()
	require.NoError(t, dbm.ParseBenchmarkData(bytes.NewBuffer([]byte(benchTxt))))
	sbm := NewBenchmarks(defaultBenchmarkScoreWeight, "<p>Test!!!</p>", nil)
	require.NoError(t, sbm.Setup(dbm, nil))
	for _, hdlr := range dbm.HandlerTags() {
		assert.True(t, common.FuzzyEqual(scoresOriginal[hdlr], sbm.ScoreForType(hdlr, score.Original)))
		assert.True(t, common.FuzzyEqual(scoresByData[hdlr], sbm.ScoreForType(hdlr, score.ByData)))
		assert.True(t, common.FuzzyEqual(scores[hdlr], sbm.ScoreFor(hdlr)))
	}
}

var defaultBenchmarkScoreWeight = map[bench.Weight]uint{
	bench.Allocations: 1,
	bench.AllocBytes:  2,
	bench.Nanoseconds: 3,
}
