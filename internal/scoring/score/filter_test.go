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
	assert.True(t, flt.Keep(hdlr1))
	assert.False(t, flt.Keep(hdlr2))
	assert.False(t, flt.Keep(hdlr3))
}

func TestExcludeFilter(t *testing.T) {
	flt := NewExcludeFilter(hdlr2)
	require.NotNil(t, flt)
	assert.True(t, flt.Keep(hdlr1))
	assert.False(t, flt.Keep(hdlr2))
	assert.True(t, flt.Keep(hdlr3))
}

func TestGroupFilter(t *testing.T) {
	fltGrp1 := NewFilterGroup(hdlr1)
	require.NotNil(t, fltGrp1)
	fltInGrp1 := NewIncludeFilter(fltGrp1, hdlr2)
	assert.True(t, fltInGrp1.Keep(hdlr1))
	assert.True(t, fltInGrp1.Keep(hdlr2))
	assert.False(t, fltInGrp1.Keep(hdlr3))
	fltExGrp1 := NewExcludeFilter(fltGrp1, hdlr2)
	assert.False(t, fltExGrp1.Keep(hdlr1))
	assert.False(t, fltExGrp1.Keep(hdlr2))
	assert.True(t, fltExGrp1.Keep(hdlr3))
}
