package automatons

import (
	"slices"
	"testing"
)

func TestRuneRange_IncludesSingleChar(t *testing.T) {
	r := NewRuneRange('A', 'A', true)

	if r.Start != 'A' {
		t.Errorf("Expected Start to be %d, got %d", 'A', r.Start)
	}
	if r.End != 'A' {
		t.Errorf("Expected End to be %d, got %d", 'A', r.End)
	}
	if !r.Includes {
		t.Errorf("Expected Includes to be true, got %t", r.Includes)
	}
	if !r.Contains('A') {
		t.Errorf("Expected Contains('A') to be true")
	}
	if r.Contains('@') {
		t.Errorf("Expected Contains('@') to be false")
	}
	if r.Contains('B') {
		t.Errorf("Expected Contains('B') to be false")
	}

	// Test iterator
	var collected []rune
	for ch := range r.All() {
		collected = append(collected, ch)
	}
	expected := []rune{'A'}
	if !slices.Equal(collected, expected) {
		t.Errorf("Expected iterator to yield %v, got %v", expected, collected)
	}

	expectedStr := "[A-A]"
	if r.String() != expectedStr {
		t.Errorf("Expected String() to be %q, got %q", expectedStr, r.String())
	}
}

func TestRuneRange_IncludesRange(t *testing.T) {
	r := NewRuneRange('A', 'Z', true)

	if r.Start != 'A' {
		t.Errorf("Expected Start to be %d, got %d", 'A', r.Start)
	}
	if r.End != 'Z' {
		t.Errorf("Expected End to be %d, got %d", 'Z', r.End)
	}
	if !r.Includes {
		t.Errorf("Expected Includes to be true, got %t", r.Includes)
	}
	if !r.Contains('A') {
		t.Errorf("Expected Contains('A') to be true")
	}
	if !r.Contains('Z') {
		t.Errorf("Expected Contains('Z') to be true")
	}
	if r.Contains('@') {
		t.Errorf("Expected Contains('@') to be false")
	}
	if r.Contains('[') {
		t.Errorf("Expected Contains('[') to be false")
	}

	// Test iterator
	var collected []rune
	for ch := range r.All() {
		collected = append(collected, ch)
	}
	var expected []rune
	for i := 'A'; i <= 'Z'; i++ {
		expected = append(expected, i)
	}
	if !slices.Equal(collected, expected) {
		t.Errorf("Expected iterator to yield %d items, got %d items", len(expected), len(collected))
	}

	expectedStr := "[A-Z]"
	if r.String() != expectedStr {
		t.Errorf("Expected String() to be %q, got %q", expectedStr, r.String())
	}
}

func TestRuneRange_ExcludesSingleChar(t *testing.T) {
	r := NewRuneRange('A', 'A', false)

	if r.Start != 'A' {
		t.Errorf("Expected Start to be %d, got %d", 'A', r.Start)
	}
	if r.End != 'A' {
		t.Errorf("Expected End to be %d, got %d", 'A', r.End)
	}
	if r.Includes {
		t.Errorf("Expected Includes to be false, got %t", r.Includes)
	}
	if !r.Contains('A') {
		t.Errorf("Expected Contains('A') to be true")
	}
	if r.Contains('@') {
		t.Errorf("Expected Contains('@') to be false")
	}
	if r.Contains('B') {
		t.Errorf("Expected Contains('B') to be false")
	}

	// Test iterator
	var collected []rune
	for ch := range r.All() {
		collected = append(collected, ch)
	}
	expected := []rune{'A'}
	if !slices.Equal(collected, expected) {
		t.Errorf("Expected iterator to yield %v, got %v", expected, collected)
	}

	expectedStr := "[^A-A]"
	if r.String() != expectedStr {
		t.Errorf("Expected String() to be %q, got %q", expectedStr, r.String())
	}
}

func TestRuneRange_ExcludesRange(t *testing.T) {
	r := NewRuneRange('A', 'Z', false)

	if r.Start != 'A' {
		t.Errorf("Expected Start to be %d, got %d", 'A', r.Start)
	}
	if r.End != 'Z' {
		t.Errorf("Expected End to be %d, got %d", 'Z', r.End)
	}
	if r.Includes {
		t.Errorf("Expected Includes to be false, got %t", r.Includes)
	}
	if !r.Contains('A') {
		t.Errorf("Expected Contains('A') to be true")
	}
	if !r.Contains('Z') {
		t.Errorf("Expected Contains('Z') to be true")
	}
	if r.Contains('@') {
		t.Errorf("Expected Contains('@') to be false")
	}
	if r.Contains('[') {
		t.Errorf("Expected Contains('[') to be false")
	}

	// Test iterator
	var collected []rune
	for ch := range r.All() {
		collected = append(collected, ch)
	}
	var expected []rune
	for i := 'A'; i <= 'Z'; i++ {
		expected = append(expected, i)
	}
	if !slices.Equal(collected, expected) {
		t.Errorf("Expected iterator to yield %d items, got %d items", len(expected), len(collected))
	}

	expectedStr := "[^A-Z]"
	if r.String() != expectedStr {
		t.Errorf("Expected String() to be %q, got %q", expectedStr, r.String())
	}
}

func TestRuneRange_InvalidRange(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic when creating invalid range")
		}
	}()
	NewRuneRange('Z', 'A', true)
}

func TestRuneRange_Equals(t *testing.T) {
	range1 := NewRuneRange('A', 'Z', true)
	range2 := NewRuneRange('A', 'Z', true)
	range3 := NewRuneRange('A', 'Z', false)
	range4 := NewRuneRange('B', 'Z', true)

	if !range1.Equals(*range2) {
		t.Errorf("Expected range1.Equals(range2) to be true")
	}
	if range1.Equals(*range3) {
		t.Errorf("Expected range1.Equals(range3) to be false")
	}
	if range1.Equals(*range4) {
		t.Errorf("Expected range1.Equals(range4) to be false")
	}
}
