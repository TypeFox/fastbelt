package automatons

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContains(t *testing.T) {
	// Test case: [a-z] contains [c-f]
	setAZ := NewRuneSetRange('a', 'z')
	setCF := NewRuneSetRange('c', 'f')

	assert.Equal(t, true, setAZ.Contains(setCF))

	// Test case: [c-f] does not contain [a-z]
	assert.Equal(t, false, setCF.Contains(setAZ))
}

func TestAdd(t *testing.T) {
	set := NewRuneSetRange('a', 'c')
	other := NewRuneSetRange('x', 'z')

	set.Add(other)

	assert.Equal(t, true, set.IncludesRune('a') && set.IncludesRune('b') && set.IncludesRune('c'))
	assert.Equal(t, true, set.IncludesRune('x') && set.IncludesRune('y') && set.IncludesRune('z'))
	assert.Equal(t, false, set.IncludesRune('m'))
}

func TestRemove(t *testing.T) {
	set := NewRuneSetRange('a', 'z')
	toRemove := NewRuneSetRange('m', 'p')

	set.Remove(toRemove)

	assert.Equal(t, true, set.IncludesRune('a') && set.IncludesRune('l'))
	assert.Equal(t, true, set.IncludesRune('q') && set.IncludesRune('z'))
	assert.Equal(t, false, set.IncludesRune('m') || set.IncludesRune('n') || set.IncludesRune('o') || set.IncludesRune('p'))
}

func TestUnion(t *testing.T) {
	setAC := NewRuneSetRange('a', 'c')
	setXZ := NewRuneSetRange('x', 'z')

	result := Union(setAC, setXZ)

	assert.Equal(t, true, result.IncludesRune('a') && result.IncludesRune('b') && result.IncludesRune('c'))
	assert.Equal(t, true, result.IncludesRune('x') && result.IncludesRune('y') && result.IncludesRune('z'))
	assert.Equal(t, false, result.IncludesRune('m'))
}

func TestExcept(t *testing.T) {
	setAZ := NewRuneSetRange('a', 'z')
	setMP := NewRuneSetRange('m', 'p')

	result := Except(setAZ, setMP)

	assert.Equal(t, true, result.IncludesRune('a') && result.IncludesRune('l'))
	assert.Equal(t, true, result.IncludesRune('q') && result.IncludesRune('z'))
	assert.Equal(t, false, result.IncludesRune('m') || result.IncludesRune('n') || result.IncludesRune('o') || result.IncludesRune('p'))
}

func TestNegate(t *testing.T) {
	setAC := NewRuneSetRange('a', 'c')

	result := Negate(setAC)

	assert.Equal(t, false, result.IncludesRune('a') || result.IncludesRune('b') || result.IncludesRune('c'))
	assert.Equal(t, true, result.IncludesRune('A') && result.IncludesRune('d') && result.IncludesRune('z'))
}

func TestIntersect(t *testing.T) {
	setAM := NewRuneSetRange('a', 'm')
	setHZ := NewRuneSetRange('h', 'z')

	result := Intersect(setAM, setHZ)

	// Intersection should be [h-m]
	assert.Equal(t, true, result.IncludesRune('h') && result.IncludesRune('i') && result.IncludesRune('m'))
	assert.Equal(t, false, result.IncludesRune('g') || result.IncludesRune('n'))
}

func TestUnionMultiple(t *testing.T) {
	setAC := NewRuneSetRange('a', 'c')
	setGI := NewRuneSetRange('g', 'i')
	setXZ := NewRuneSetRange('x', 'z')

	result := Union(setAC, setGI, setXZ)

	assert.Equal(t, true, result.IncludesRune('a') && result.IncludesRune('h') && result.IncludesRune('y'))
	assert.Equal(t, false, result.IncludesRune('d') || result.IncludesRune('m'))
}
