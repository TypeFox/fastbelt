package automatons

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRuneRange_IncludesSingleChar(t *testing.T) {
	r := NewRuneRange('A', 'A', true)

	assert.Equal(t, 'A', r.Start)
	assert.Equal(t, 'A', r.End)
	assert.Equal(t, true, r.Includes)
	assert.Equal(t, true, r.Contains('A'))
	assert.Equal(t, false, r.Contains('@'))
	assert.Equal(t, false, r.Contains('B'))

	// Test iterator
	var collected []rune
	for ch := range r.All() {
		collected = append(collected, ch)
	}
	expected := []rune{'A'}
	assert.True(t, slices.Equal(collected, expected))
	assert.Equal(t, "[A-A]", r.String())
}

func TestRuneRange_IncludesRange(t *testing.T) {
	r := NewRuneRange('A', 'Z', true)

	assert.Equal(t, 'A', r.Start)
	assert.Equal(t, 'Z', r.End)
	assert.Equal(t, true, r.Includes)
	assert.Equal(t, true, r.Contains('A'))
	assert.Equal(t, true, r.Contains('Z'))
	assert.Equal(t, false, r.Contains('@'))
	assert.Equal(t, false, r.Contains('['))

	// Test iterator
	var collected []rune
	for ch := range r.All() {
		collected = append(collected, ch)
	}
	var expected []rune
	for i := 'A'; i <= 'Z'; i++ {
		expected = append(expected, i)
	}

	assert.True(t, slices.Equal(collected, expected))
	assert.Equal(t, "[A-Z]", r.String())
}

func TestRuneRange_ExcludesSingleChar(t *testing.T) {
	r := NewRuneRange('A', 'A', false)

	assert.Equal(t, 'A', r.Start)
	assert.Equal(t, 'A', r.End)
	assert.Equal(t, false, r.Includes)
	assert.Equal(t, true, r.Contains('A'))
	assert.Equal(t, false, r.Contains('@'))
	assert.Equal(t, false, r.Contains('B'))

	// Test iterator
	var collected []rune
	for ch := range r.All() {
		collected = append(collected, ch)
	}
	expected := []rune{'A'}
	assert.True(t, slices.Equal(collected, expected))
	assert.Equal(t, "[^A-A]", r.String())
}

func TestRuneRange_ExcludesRange(t *testing.T) {
	r := NewRuneRange('A', 'Z', false)

	assert.Equal(t, 'A', r.Start)
	assert.Equal(t, 'Z', r.End)
	assert.Equal(t, false, r.Includes)
	assert.Equal(t, true, r.Contains('A'))
	assert.Equal(t, true, r.Contains('Z'))
	assert.Equal(t, false, r.Contains('@'))
	assert.Equal(t, false, r.Contains('['))

	// Test iterator
	var collected []rune
	for ch := range r.All() {
		collected = append(collected, ch)
	}
	var expected []rune
	for i := 'A'; i <= 'Z'; i++ {
		expected = append(expected, i)
	}

	assert.True(t, slices.Equal(collected, expected))
	assert.Equal(t, "[^A-Z]", r.String())
}

func TestRuneRange_InvalidRange(t *testing.T) {
	assert.Panics(t, func() {
		NewRuneRange('Z', 'A', true)
	})
}

func TestRuneRange_Equals(t *testing.T) {
	range1 := NewRuneRange('A', 'Z', true)
	range2 := NewRuneRange('A', 'Z', true)
	range3 := NewRuneRange('A', 'Z', false)
	range4 := NewRuneRange('B', 'Z', true)

	assert.True(t, range1.Equals(*range2))
	assert.False(t, range1.Equals(*range3))
	assert.False(t, range1.Equals(*range4))
}
