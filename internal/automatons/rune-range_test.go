package automatons

import (
	"slices"
	"testing"
)

func TestRuneRange_IncludesSingleChar(t *testing.T) {
	r := NewRuneRange('A', 'A', true)

	Expect(r.Start).ToEqual('A')
	Expect(r.End).ToEqual('A')
	Expect(r.Includes).ToEqual(true)
	Expect(r.Contains('A')).ToEqual(true)
	Expect(r.Contains('@')).ToEqual(false)
	Expect(r.Contains('B')).ToEqual(false)

	// Test iterator
	var collected []rune
	for ch := range r.All() {
		collected = append(collected, ch)
	}
	expected := []rune{'A'}
	Expect(slices.Equal(collected, expected)).ToEqual(true)
	Expect(r.String()).ToEqual("[A-A]")
}

func TestRuneRange_IncludesRange(t *testing.T) {
	r := NewRuneRange('A', 'Z', true)

	Expect(r.Start).ToEqual('A')
	Expect(r.End).ToEqual('Z')
	Expect(r.Includes).ToEqual(true)
	Expect(r.Contains('A')).ToEqual(true)
	Expect(r.Contains('Z')).ToEqual(true)
	Expect(r.Contains('@')).ToEqual(false)
	Expect(r.Contains('[')).ToEqual(false)

	// Test iterator
	var collected []rune
	for ch := range r.All() {
		collected = append(collected, ch)
	}
	var expected []rune
	for i := 'A'; i <= 'Z'; i++ {
		expected = append(expected, i)
	}

	Expect(slices.Equal(collected, expected)).ToEqual(true)
	Expect(r.String()).ToEqual("[A-Z]")
}

func TestRuneRange_ExcludesSingleChar(t *testing.T) {
	r := NewRuneRange('A', 'A', false)

	Expect(r.Start).ToEqual('A')
	Expect(r.End).ToEqual('A')
	Expect(r.Includes).ToEqual(false)
	Expect(r.Contains('A')).ToEqual(true)
	Expect(r.Contains('@')).ToEqual(false)
	Expect(r.Contains('B')).ToEqual(false)

	// Test iterator
	var collected []rune
	for ch := range r.All() {
		collected = append(collected, ch)
	}
	expected := []rune{'A'}
	Expect(slices.Equal(collected, expected)).ToEqual(true)
	Expect(r.String()).ToEqual("[^A-A]")
}

func TestRuneRange_ExcludesRange(t *testing.T) {
	r := NewRuneRange('A', 'Z', false)

	Expect(r.Start).ToEqual('A')
	Expect(r.End).ToEqual('Z')
	Expect(r.Includes).ToEqual(false)
	Expect(r.Contains('A')).ToEqual(true)
	Expect(r.Contains('Z')).ToEqual(true)
	Expect(r.Contains('@')).ToEqual(false)
	Expect(r.Contains('[')).ToEqual(false)

	// Test iterator
	var collected []rune
	for ch := range r.All() {
		collected = append(collected, ch)
	}
	var expected []rune
	for i := 'A'; i <= 'Z'; i++ {
		expected = append(expected, i)
	}

	Expect(slices.Equal(collected, expected)).ToEqual(true)
	Expect(r.String()).ToEqual("[^A-Z]")
}

func TestRuneRange_InvalidRange(t *testing.T) {
	Expect(func() {
		NewRuneRange('Z', 'A', true)
	}).ToPanic()
}

func TestRuneRange_Equals(t *testing.T) {
	range1 := NewRuneRange('A', 'Z', true)
	range2 := NewRuneRange('A', 'Z', true)
	range3 := NewRuneRange('A', 'Z', false)
	range4 := NewRuneRange('B', 'Z', true)

	Expect(range1.Equals(*range2)).ToEqual(true)
	Expect(range1.Equals(*range3)).ToEqual(false)
	Expect(range1.Equals(*range4)).ToEqual(false)
}
