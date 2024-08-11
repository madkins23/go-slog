package score

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	hdlr1 = "Kitten"
	hdlr2 = "Puppy"
	hdlr3 = "Fawn"
)

func TestIncludeFilter(t *testing.T) {
	flt := NewIncludeFilter(hdlr1)
	require.NotNil(t, flt)
	assert.True(t, flt.Include(hdlr1))
	assert.False(t, flt.Include(hdlr2))
	assert.False(t, flt.Include(hdlr3))
}

func TestExcludeFilter(t *testing.T) {
	flt := NewExcludeFilter(hdlr2)
	require.NotNil(t, flt)
	assert.True(t, flt.Include(hdlr1))
	assert.False(t, flt.Include(hdlr2))
	assert.True(t, flt.Include(hdlr3))
}

func TestIncludeExcludeFilter(t *testing.T) {
	fltEx1 := NewExcludeFilter(hdlr1)
	require.NotNil(t, fltEx1)
	fltIn2 := NewIncludeFilter(hdlr2)
	require.NotNil(t, fltIn2)
	fltInEx1Hdlr1 := NewIncludeFilter(fltEx1, hdlr1, hdlr2)
	assert.False(t, fltInEx1Hdlr1.Include(hdlr1))
	assert.True(t, fltInEx1Hdlr1.Include(hdlr2))
	fltInHdlr1Ex1 := NewIncludeFilter(hdlr1, fltEx1, hdlr2)
	assert.True(t, fltInHdlr1Ex1.Include(hdlr1))
	assert.True(t, fltInHdlr1Ex1.Include(hdlr2))
}

func TestExcludeIncludeFilter(t *testing.T) {
	fltIn1 := NewIncludeFilter(hdlr1)
	require.NotNil(t, fltIn1)
	fltEx2 := NewExcludeFilter(hdlr2)
	require.NotNil(t, fltEx2)
	fltExIn1Hdlr1 := NewExcludeFilter(fltIn1, hdlr1, hdlr2)
	assert.True(t, fltExIn1Hdlr1.Include(hdlr1))
	assert.False(t, fltExIn1Hdlr1.Include(hdlr2))
	assert.True(t, fltExIn1Hdlr1.Include(hdlr3))
	fltExHdlr1In1 := NewExcludeFilter(hdlr1, fltIn1, hdlr2)
	assert.False(t, fltExHdlr1In1.Include(hdlr1))
	assert.False(t, fltExHdlr1In1.Include(hdlr2))
	assert.True(t, fltExIn1Hdlr1.Include(hdlr3))
}
