package automatons

import (
	"testing"
)

func TestContains(t *testing.T) {
	// Test case: [a-z] contains [c-f]
	setAZ := NewRuneSet_Range('a', 'z')
	setCF := NewRuneSet_Range('c', 'f')

	if !setAZ.Contains(setCF) {
		t.Error("Expected [a-z] to contain [c-f]")
	}

	// Test case: [c-f] does not contain [a-z]
	if setCF.Contains(setAZ) {
		t.Error("Expected [c-f] to not contain [a-z]")
	}
}

func TestAdd(t *testing.T) {
	set := NewRuneSet_Range('a', 'c')
	other := NewRuneSet_Range('x', 'z')

	set.Add(other)

	if !set.IncludesRune('a') || !set.IncludesRune('b') || !set.IncludesRune('c') {
		t.Error("Expected set to still include [a-c]")
	}

	if !set.IncludesRune('x') || !set.IncludesRune('y') || !set.IncludesRune('z') {
		t.Error("Expected set to include [x-z] after adding")
	}

	if set.IncludesRune('m') {
		t.Error("Expected set to not include 'm'")
	}
}

func TestRemove(t *testing.T) {
	set := NewRuneSet_Range('a', 'z')
	toRemove := NewRuneSet_Range('m', 'p')

	set.Remove(toRemove)

	if !set.IncludesRune('a') || !set.IncludesRune('l') {
		t.Error("Expected set to still include characters before removed range")
	}

	if !set.IncludesRune('q') || !set.IncludesRune('z') {
		t.Error("Expected set to still include characters after removed range")
	}

	if set.IncludesRune('m') || set.IncludesRune('n') || set.IncludesRune('o') || set.IncludesRune('p') {
		t.Error("Expected set to not include removed characters")
	}
}

func TestUnion(t *testing.T) {
	setAC := NewRuneSet_Range('a', 'c')
	setXZ := NewRuneSet_Range('x', 'z')

	result := Union(setAC, setXZ)

	if !result.IncludesRune('a') || !result.IncludesRune('b') || !result.IncludesRune('c') {
		t.Error("Expected union to include [a-c]")
	}

	if !result.IncludesRune('x') || !result.IncludesRune('y') || !result.IncludesRune('z') {
		t.Error("Expected union to include [x-z]")
	}

	if result.IncludesRune('m') {
		t.Error("Expected union to not include 'm'")
	}
}

func TestExcept(t *testing.T) {
	setAZ := NewRuneSet_Range('a', 'z')
	setMP := NewRuneSet_Range('m', 'p')

	result := Except(setAZ, setMP)

	if !result.IncludesRune('a') || !result.IncludesRune('l') {
		t.Error("Expected except result to include characters before removed range")
	}

	if !result.IncludesRune('q') || !result.IncludesRune('z') {
		t.Error("Expected except result to include characters after removed range")
	}

	if result.IncludesRune('m') || result.IncludesRune('n') || result.IncludesRune('o') || result.IncludesRune('p') {
		t.Error("Expected except result to not include excepted characters")
	}
}

func TestNegate(t *testing.T) {
	setAC := NewRuneSet_Range('a', 'c')

	result := Negate(setAC)

	if result.IncludesRune('a') || result.IncludesRune('b') || result.IncludesRune('c') {
		t.Error("Expected negated set to not include [a-c]")
	}

	if !result.IncludesRune('A') || !result.IncludesRune('d') || !result.IncludesRune('z') {
		t.Error("Expected negated set to include characters outside [a-c]")
	}
}

func TestIntersect(t *testing.T) {
	setAM := NewRuneSet_Range('a', 'm')
	setHZ := NewRuneSet_Range('h', 'z')

	result := Intersect(setAM, setHZ)

	// Intersection should be [h-m]
	if !result.IncludesRune('h') || !result.IncludesRune('i') || !result.IncludesRune('m') {
		t.Error("Expected intersection to include [h-m]")
	}

	if result.IncludesRune('g') || result.IncludesRune('n') {
		t.Error("Expected intersection to not include characters outside [h-m]")
	}
}

func TestUnionMultiple(t *testing.T) {
	setAC := NewRuneSet_Range('a', 'c')
	setGI := NewRuneSet_Range('g', 'i')
	setXZ := NewRuneSet_Range('x', 'z')

	result := Union(setAC, setGI, setXZ)

	if !result.IncludesRune('a') || !result.IncludesRune('h') || !result.IncludesRune('y') {
		t.Error("Expected union to include characters from all sets")
	}

	if result.IncludesRune('d') || result.IncludesRune('m') {
		t.Error("Expected union to not include characters outside all sets")
	}
}
