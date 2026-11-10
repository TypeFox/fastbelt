package benchmarkGenerated

import (
	"unicode/utf8"
)

var ID_Lookup = [][]int64{
	{0x0000005A00000041, 0x0000005F0000005F, 0x0000007A00000061},
	{0x0000003900000030, 0x0000005A00000041, 0x0000005F0000005F, 0x0000007A00000061},
}
var ID_Next = [][]int{
	{1, 1, 1},
	{1, 1, 1, 1},
}
var ID_Accepting = [2]bool{
	1: true,
}
var ID_Loop1 = [128]bool{48: true, 49: true, 50: true, 51: true, 52: true, 53: true, 54: true, 55: true, 56: true, 57: true, 65: true, 66: true, 67: true, 68: true, 69: true, 70: true, 71: true, 72: true, 73: true, 74: true, 75: true, 76: true, 77: true, 78: true, 79: true, 80: true, 81: true, 82: true, 83: true, 84: true, 85: true, 86: true, 87: true, 88: true, 89: true, 90: true, 95: true, 97: true, 98: true, 99: true, 100: true, 101: true, 102: true, 103: true, 104: true, 105: true, 106: true, 107: true, 108: true, 109: true, 110: true, 111: true, 112: true, 113: true, 114: true, 115: true, 116: true, 117: true, 118: true, 119: true, 120: true, 121: true, 122: true}

func ID(s string, offset int) int {
	input := s[offset:]
	length := len(input)
	state := 0
	acceptedIndex := 0
	index := 0
loop:
	for index < length {
		switch state {
		case 1:
			start := index
			for index < length {
				c := input[index]
				if c >= utf8.RuneSelf || !ID_Loop1[c] {
					break
				}
				index++
			}
			if index != start {
				acceptedIndex = index
			}
			if index == length {
				break loop
			}
		}
		r, runeSize := utf8.DecodeRuneInString(input[index:])
		switch state {
		case 0:
			nextState := -1
			next := ID_Next[0]
			lookup := ID_Lookup[0]
			for i, lowHigh := range lookup {
				if rune(lowHigh&0xFFFFFFFF) <= r && r <= rune(lowHigh>>32) {
					nextState = next[i]
					break
				}
			}
			if nextState > -1 {
				state = nextState
			} else {
				break loop
			}
		case 1:
			nextState := -1
			next := ID_Next[1]
			lookup := ID_Lookup[1]
			for i, lowHigh := range lookup {
				if rune(lowHigh&0xFFFFFFFF) <= r && r <= rune(lowHigh>>32) {
					nextState = next[i]
					break
				}
			}
			if nextState > -1 {
				state = nextState
			} else {
				break loop
			}
		default:
			break loop
		}
		index += runeSize
		if ID_Accepting[state] {
			acceptedIndex = index
		}
	}
	return acceptedIndex
}
