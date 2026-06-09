// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package collections

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testTokenBitset(t *testing.T, bitset *BitSet, expected ...int) {
	t.Helper()
	max := expected[len(expected)-1]
	assert.Equal(t, len(expected), bitset.Cardinality(), "cardinality should match number of expected bits")
	for i := 0; i <= max+1; i++ {
		found := slices.Contains(expected, i)
		if found {
			assert.True(t, bitset.At(i), "expected bit %d to be set", i)
		} else {
			assert.False(t, bitset.At(i), "expected bit %d to be unset", i)
		}
	}
}

func TestTokenBitset_NewIsEmpty(t *testing.T) {
	b := NewBitset()
	assert.False(t, b.At(0))
	assert.False(t, b.At(63))
	assert.False(t, b.At(1000))
	assert.Equal(t, 0, b.Cardinality(), "cardinality should be 0 for empty bitset")
}

func TestTokenBitset_InsertAndAt(t *testing.T) {
	tests := []struct {
		name    string
		inserts []int
	}{
		{"single bit in first word", []int{0}},
		{"single bit at word boundary", []int{63}},
		{"single bit in second word", []int{64}},
		{"bits across multiple words", []int{0, 64, 128, 200}},
		{"multiple bits in same word", []int{1, 5, 17, 42}},
		{"inserted in descending order", []int{300, 200, 100, 5}},
		{"interleaved low/high", []int{500, 1, 300, 65, 128}},
		{"large index", []int{10000}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewBitset()
			for _, i := range tt.inserts {
				b.Insert(i)
			}
			for _, i := range tt.inserts {
				assert.True(t, b.At(i), "expected bit %d to be set", i)
			}
			assert.Equal(t, len(tt.inserts), b.Cardinality(), "cardinality should match number of inserts")
		})
	}
}

func TestTokenBitset_UnsetBitsRemainFalse(t *testing.T) {
	b := NewBitset()
	b.Insert(64)
	b.Insert(200)

	testTokenBitset(t, b, 64, 200)
}

func TestTokenBitset_InsertReturnsSelf(t *testing.T) {
	b := NewBitset()
	assert.Same(t, b, b.Insert(0))
	assert.Same(t, b, b.Insert(100))
}

func TestTokenBitset_InsertIdempotent(t *testing.T) {
	b := NewBitset()
	b.Insert(42)
	b.Insert(42)
	b.Insert(42)
	testTokenBitset(t, b, 42)
}

func TestTokenBitset_PrependGrowsFront(t *testing.T) {
	b := NewBitset()
	b.Insert(500)
	b.Insert(1)
	testTokenBitset(t, b, 1, 500)
}

func TestMergeTokenBitsets_Empty(t *testing.T) {
	merged := MergeBitSets(nil)
	assert.NotNil(t, merged)
	assert.False(t, merged.At(0))
	assert.False(t, merged.At(1000))

	merged = MergeBitSets([]*BitSet{})
	assert.NotNil(t, merged)
	assert.False(t, merged.At(0))
}

func TestMergeTokenBitsets_Single(t *testing.T) {
	a := NewBitset().Insert(3).Insert(70).Insert(300)
	merged := MergeBitSets([]*BitSet{a})
	testTokenBitset(t, merged, 3, 70, 300)
}

func TestMergeTokenBitsets_Disjoint(t *testing.T) {
	a := NewBitset().Insert(1).Insert(2)
	assert.Equal(t, 1, a.Min())
	b := NewBitset().Insert(500).Insert(1000)
	assert.Equal(t, 500, b.Min())
	merged := MergeBitSets([]*BitSet{a, b})

	testTokenBitset(t, merged, 1, 2, 500, 1000)
}

func TestMergeTokenBitsets_Overlapping(t *testing.T) {
	a := NewBitset().Insert(5).Insert(70).Insert(150)
	assert.Equal(t, 5, a.Min())
	b := NewBitset().Insert(7).Insert(70).Insert(200)
	assert.Equal(t, 7, b.Min())
	merged := MergeBitSets([]*BitSet{a, b})

	testTokenBitset(t, merged, 5, 7, 70, 150, 200)
}

func TestMergeTokenBitsets_PreservesAllOriginals(t *testing.T) {
	a := NewBitset().Insert(0).Insert(63)
	b := NewBitset().Insert(64).Insert(127)
	c := NewBitset().Insert(128).Insert(255)
	merged := MergeBitSets([]*BitSet{a, b, c})

	testTokenBitset(t, merged, 0, 63, 64, 127, 128, 255)
}

func TestMergeTokenBitsets_DoesNotMutateInputs(t *testing.T) {
	a := NewBitset().Insert(10)
	b := NewBitset().Insert(20)
	_ = MergeBitSets([]*BitSet{a, b})

	assert.True(t, a.At(10))
	assert.False(t, a.At(20))
	assert.True(t, b.At(20))
	assert.False(t, b.At(10))
}

func BenchmarkTokenBitset_Insert(b *testing.B) {
	bitset := NewBitset()
	for b.Loop() {
		bitset.Insert(b.N % 1000)
	}
}

func BenchmarkTokenBitset_At(b *testing.B) {
	bitset := NewBitset()
	for i := range 1000 {
		bitset.Insert(i)
	}
	for b.Loop() {
		bitset.At(b.N % 1000)
	}
}

func BenchmarkMergeTokenBitsets(b *testing.B) {
	bitsets := []*BitSet{
		NewBitset().Insert(1).Insert(10).Insert(100),
		NewBitset().Insert(2).Insert(20).Insert(200),
		NewBitset().Insert(3).Insert(30).Insert(300),
	}
	b.ResetTimer()
	for b.Loop() {
		MergeBitSets(bitsets)
	}
}
