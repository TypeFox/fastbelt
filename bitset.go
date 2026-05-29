// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package fastbelt

type BitSet struct {
	offset int
	set    []uint64
}

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
			diff := i - bs.offset
			if diff < 0 || diff >= len(bs.set) {
				// Outside of bitset
				continue
			}
			// Combine existing word with new word
			b.set[i-min] |= bs.set[diff]
		}
	}
	return b
}

func NewBitset() *BitSet {
	return &BitSet{}
}

const wordSize = 64

func computeIndex(i int) (int, int) {
	return i / wordSize, i % wordSize
}

func (b *BitSet) Insert(i int) *BitSet {
	index, offset := computeIndex(i)
	diff := index - b.offset
	length := len(b.set)
	if diff < 0 {
		// index is less than offset, append to the front
		b.offset = index
		intermediate := make([]uint64, (-diff)+length)
		copy(intermediate[-diff:], b.set)
		b.set = intermediate
		diff = 0
	} else if diff >= length {
		// index is larger than offset+len, append to the end
		if length == 0 {
			// Special case if bitset is still empty
			b.offset = index
			b.set = make([]uint64, 1)
			diff = 0
		} else {
			intermediate := make([]uint64, diff+1)
			copy(intermediate, b.set)
			b.set = intermediate
		}
	}
	b.set[diff] |= 1 << offset
	return b
}

func (b *BitSet) Delete(i int) {
	index, offset := computeIndex(i)
	diff := index - b.offset
	if diff < 0 || diff >= len(b.set) {
		// Outside of bitset, nothing to delete
		return
	}
	b.set[diff] &^= 1 << offset
}

func (b *BitSet) At(i int) bool {
	index, offset := computeIndex(i)
	diff := index - b.offset
	if diff < 0 || diff >= len(b.set) {
		return false
	}
	return (b.set[diff] & (1 << offset)) > 0
}

func (b *BitSet) Empty() bool {
	return len(b.set) == 0
}
