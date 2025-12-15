package automatons

import (
	"bytes"
)

type BitMask []byte

func NewBitMask_Empty(bits int) BitMask {
	return BitMask(bytes.Repeat([]byte{'0'}, bits))
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

func (bm BitMask) String() string {
	return string(bm)
}

func (bm BitMask) Set(bit int) {
	if bit >= 0 && bit < len(bm) {
		bm[bit] = '1'
	}
}

func (bm BitMask) Clear(bit int) {
	if bit >= 0 && bit < len(bm) {
		bm[bit] = '0'
	}
}

func (bm BitMask) IsSet(bit int) bool {
	if bit >= 0 && bit < len(bm) {
		return bm[bit] == '1'
	}
	return false
}

func (bm BitMask) Get(bit int) bool {
	if bit >= 0 && bit < len(bm) {
		return bm[bit] == '1'
	}
	return false
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
