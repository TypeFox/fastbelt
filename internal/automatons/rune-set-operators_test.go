package automatons

import (
	"testing"
)

func TestContains(t *testing.T) {
	// Test case: [a-z] contains [c-f]
	setAZ := NewRuneSet_Range('a', 'z')
	setCF := NewRuneSet_Range('c', 'f')

	Expect(setAZ.Contains(setCF)).ToEqual(true)

	// Test case: [c-f] does not contain [a-z]
	Expect(setCF.Contains(setAZ)).ToEqual(false)
}

func TestAdd(t *testing.T) {
	set := NewRuneSet_Range('a', 'c')
	other := NewRuneSet_Range('x', 'z')

	set.Add(other)

	Expect(set.IncludesRune('a') && set.IncludesRune('b') && set.IncludesRune('c')).ToEqual(true)
	Expect(set.IncludesRune('x') && set.IncludesRune('y') && set.IncludesRune('z')).ToEqual(true)
	Expect(set.IncludesRune('m')).ToEqual(false)
}

func TestRemove(t *testing.T) {
	set := NewRuneSet_Range('a', 'z')
	toRemove := NewRuneSet_Range('m', 'p')

	set.Remove(toRemove)

	Expect(set.IncludesRune('a') && set.IncludesRune('l')).ToEqual(true)
	Expect(set.IncludesRune('q') && set.IncludesRune('z')).ToEqual(true)
	Expect(set.IncludesRune('m') || set.IncludesRune('n') || set.IncludesRune('o') || set.IncludesRune('p')).ToEqual(false)
}

func TestUnion(t *testing.T) {
	setAC := NewRuneSet_Range('a', 'c')
	setXZ := NewRuneSet_Range('x', 'z')

	result := Union(setAC, setXZ)

	Expect(result.IncludesRune('a') && result.IncludesRune('b') && result.IncludesRune('c')).ToEqual(true)
	Expect(result.IncludesRune('x') && result.IncludesRune('y') && result.IncludesRune('z')).ToEqual(true)
	Expect(result.IncludesRune('m')).ToEqual(false)
}

func TestExcept(t *testing.T) {
	setAZ := NewRuneSet_Range('a', 'z')
	setMP := NewRuneSet_Range('m', 'p')

	result := Except(setAZ, setMP)

	Expect(result.IncludesRune('a') && result.IncludesRune('l')).ToEqual(true)
	Expect(result.IncludesRune('q') && result.IncludesRune('z')).ToEqual(true)
	Expect(result.IncludesRune('m') || result.IncludesRune('n') || result.IncludesRune('o') || result.IncludesRune('p')).ToEqual(false)
}

func TestNegate(t *testing.T) {
	setAC := NewRuneSet_Range('a', 'c')

	result := Negate(setAC)

	Expect(result.IncludesRune('a') || result.IncludesRune('b') || result.IncludesRune('c')).ToEqual(false)
	Expect(result.IncludesRune('A') && result.IncludesRune('d') && result.IncludesRune('z')).ToEqual(true)
}

func TestIntersect(t *testing.T) {
	setAM := NewRuneSet_Range('a', 'm')
	setHZ := NewRuneSet_Range('h', 'z')

	result := Intersect(setAM, setHZ)

	// Intersection should be [h-m]
	Expect(result.IncludesRune('h') && result.IncludesRune('i') && result.IncludesRune('m')).ToEqual(true)
	Expect(result.IncludesRune('g') || result.IncludesRune('n')).ToEqual(false)
}

func TestUnionMultiple(t *testing.T) {
	setAC := NewRuneSet_Range('a', 'c')
	setGI := NewRuneSet_Range('g', 'i')
	setXZ := NewRuneSet_Range('x', 'z')

	result := Union(setAC, setGI, setXZ)

	Expect(result.IncludesRune('a') && result.IncludesRune('h') && result.IncludesRune('y')).ToEqual(true)
	Expect(result.IncludesRune('d') || result.IncludesRune('m')).ToEqual(false)
}
