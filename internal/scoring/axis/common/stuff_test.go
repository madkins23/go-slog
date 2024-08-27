package common

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/madkins23/go-slog/internal/scoring/score"
)

func TestPercentDifference(t *testing.T) {
	var expected score.Value = 79.15407854984895
	assert.Equal(t, expected, PercentDifference(100, 231))
}

func TestPercentEqual(t *testing.T) {
	assert.True(t, PercentEqual(850, 879))
}
