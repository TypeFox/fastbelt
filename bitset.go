// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package fastbelt

// BitSet is a simple implementation of a bitset that supports insertion,
// deletion, and membership testing. It is designed to be memory efficient for
// sparse sets of integers, where the integers can be large but the number of
// integers in the set is relatively small.
//
// Inside of Fastbelt, BitSets are mostly used to handle sets of token types
// in the parser engine.
type BitSet struct {
	offset int
	set    []uint64
}

// Merges multiple bitsets into a single bitset that contains all the set bits from
// the input bitsets. The resulting bitset's offset is the minimum offset of the
// input bitsets, and its length is determined by the maximum integer in the
// input bitsets.
func MergeBitSets(bitsets []*BitSet) *BitSet {
	b := &BitSet{}
	first := -1
	for i, bs := range bitsets {
		if bs != nil {
			first = i
			break
		}
	}
	if first == -1 {
		return b
	}
	min := bitsets[first].offset
	max := 0
	for _, bs := range bitsets {
		if bs == nil {
			continue
		}
		if bs.offset < min {
			min = bs.offset
		}
		total := bs.offset + len(bs.set)
		if total > max {
			max = total
		}
	}
	b.offset = min
	b.set = make([]uint64, max-min)
	for i := min; i < max; i++ {
		for _, bs := range bitsets {
			if bs == nil {
				continue
			}
			delta := i - bs.offset
			if delta < 0 || delta >= len(bs.set) {
				// Outside of source bitset
				continue
			}
			// Combine existing word with new word
			b.set[i-min] |= bs.set[delta]
		}
	}
	return b
}

// NewBitset creates a new empty BitSet.
func NewBitset() *BitSet {
	return &BitSet{}
}

const wordSize = 64

func computeIndex(i int) (int, int) {
	return i / wordSize, i % wordSize
}

// Insert adds the integer i to the bitset, setting the corresponding bit to 1.
// The bitset will grow dynamically to accommodate larger integers as needed.
func (b *BitSet) Insert(i int) *BitSet {
	index, offset := computeIndex(i)
	delta := index - b.offset
	length := len(b.set)
	if delta < 0 {
		// index is less than offset, append to the front
		b.offset = index
		intermediate := make([]uint64, (-delta)+length)
		copy(intermediate[-delta:], b.set)
		b.set = intermediate
		delta = 0
	} else if delta >= length {
		// index is larger than offset+len, append to the end
		if length == 0 {
			// Special case if bitset is still empty
			b.offset = index
			b.set = make([]uint64, 1)
			delta = 0
		} else {
			intermediate := make([]uint64, delta+1)
			copy(intermediate, b.set)
			b.set = intermediate
		}
	}
	b.set[delta] |= 1 << offset
	return b
}

// Delete removes the integer i from the bitset, setting the corresponding bit to 0.
// If i is outside the current range of the bitset, Delete does nothing.
func (b *BitSet) Delete(i int) {
	index, offset := computeIndex(i)
	delta := index - b.offset
	if delta < 0 || delta >= len(b.set) {
		// Outside of bitset, nothing to delete
		return
	}
	b.set[delta] &^= 1 << offset
}

// At returns true if the integer i is in the bitset (i.e. the corresponding bit is 1),
// and false otherwise. If i is outside the current range of the bitset, At returns false.
func (b *BitSet) At(i int) bool {
	index, offset := computeIndex(i)
	diff := index - b.offset
	if diff < 0 || diff >= len(b.set) {
		return false
	}
	return (b.set[diff] & (1 << offset)) != 0
}

// Empty returns true if the bitset contains no set bits (i.e. all bits are 0), and false otherwise.
func (b *BitSet) Empty() bool {
	if len(b.set) == 0 {
		return true
	}
	for _, word := range b.set {
		if word != 0 {
			return false
		}
	}
	return true
}
