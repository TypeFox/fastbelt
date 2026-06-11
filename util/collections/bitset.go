// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package collections

import "math/bits"

// BitSet is a simple implementation of a bitset that supports insertion,
// deletion, and membership testing.
//
// Inside of Fastbelt, BitSets are mostly used to handle sets of token types
// in the parser engine.
type BitSet struct {
	words []uint64
}

// Merges multiple bitsets into a single bitset that contains all the set bits from
// the input bitsets. The resulting bitset's offset is the minimum offset of the
// input bitsets, and its length is determined by the maximum integer in the
// input bitsets.
func MergeBitSets(bitsets []*BitSet) *BitSet {
	b := &BitSet{}
	max := 0
	for _, bs := range bitsets {
		if bs == nil {
			continue
		}
		total := len(bs.words)
		if total > max {
			max = total
		}
	}
	b.words = make([]uint64, max)
	for _, bs := range bitsets {
		if bs == nil {
			continue
		}
		for i := range bs.words {
			// Combine existing word with new word
			b.words[i] |= bs.words[i]
		}
	}
	return b
}

// NewBitset creates a new empty BitSet.
func NewBitset() *BitSet {
	return &BitSet{}
}

// Insert adds the integer i to the bitset, setting the corresponding bit to 1.
// The bitset will grow dynamically to accommodate larger integers as needed.
func (b *BitSet) Insert(i int) *BitSet {
	w := i >> 6
	for w >= len(b.words) {
		b.words = append(b.words, 0)
	}
	b.words[w] |= 1 << (uint(i) & 63)
	return b
}

// At returns true if the integer i is in the bitset (i.e. the corresponding bit is 1),
// and false otherwise.
func (b *BitSet) At(i int) bool {
	w := i >> 6
	if w >= len(b.words) {
		return false
	}
	return b.words[w]&(1<<(uint(i)&63)) != 0
}

func (b *BitSet) Cardinality() int {
	n := 0
	for _, w := range b.words {
		n += bits.OnesCount64(w)
	}
	return n
}

func (b *BitSet) Min() int {
	for wi, w := range b.words {
		if w != 0 {
			return wi<<6 + bits.TrailingZeros64(w)
		}
	}
	return -1
}

// Empty returns true if the bitset contains no set bits, and false otherwise.
func (b *BitSet) Empty() bool {
	return b.Cardinality() == 0
}

func (b *BitSet) Equals(other *BitSet) bool {
	if b == other {
		return true
	}
	if other == nil || b == nil || len(b.words) != len(other.words) {
		return false
	}
	length := len(b.words)
	for i := range length {
		if b.words[i] != other.words[i] {
			return false
		}
	}
	return true
}
