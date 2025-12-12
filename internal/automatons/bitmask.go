package automatons

import "math"

type BitMask []uint

func NewBitMask_Empty(bits int) BitMask {
	words := math.Ceil(float64(bits) / 32)
	return make(BitMask, int(words))
}

func NewBitMask_Bits(bits int, set []bool) BitMask {
	result := NewBitMask_Empty(bits)
	for i, v := range set {
		if v {
			result.Set(i)
		}
	}
	return result
}

func (bm BitMask) Set(bit int) {
	word, bitPos := bit/32, bit%32
	bm[word] |= 1 << bitPos
}

func (bm BitMask) Clear(bit int) {
	word, bitPos := bit/32, bit%32
	bm[word] &^= 1 << bitPos
}

func (bm BitMask) IsSet(bit int) bool {
	word, bitPos := bit/32, bit%32
	return (bm[word] & (1 << bitPos)) != 0
}

func (bm BitMask) Get(bit int) bool {
	word, bitPos := bit/32, bit%32
	return (bm[word] & (1 << bitPos)) != 0
}

func (bm BitMask) Equals(other BitMask) bool {
	if len(bm) != len(other) {
		return false
	}
	for i := range bm {
		if bm[i] != other[i] {
			return false
		}
	}
	return true
}

func (bm BitMask) Hash() int {
	hash := 17
	for _, word := range bm {
		hash = hash*31 + int(word)
	}
	return hash
}
