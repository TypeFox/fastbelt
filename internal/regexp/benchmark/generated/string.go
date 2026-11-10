package benchmarkGenerated

import (
	"unicode/utf8"
)

var String_Lookup = [][]int64{
	{0x0000002200000022, 0x0000002700000027},
	{0x0000002100000000, 0x0000002200000022, 0x0000005B00000023, 0x0000005C0000005C, 0x0010FFFF0000005D},
	{0x0000002600000000, 0x0000002700000027, 0x0000005B00000028, 0x0000005C0000005C, 0x0010FFFF0000005D},
	{},
	{0x0000000900000000, 0x0010FFFF0000000B},
	{0x0000000900000000, 0x0010FFFF0000000B},
}
var String_Next = [][]int{
	{1, 2},
	{1, 3, 1, 4, 1},
	{2, 3, 2, 5, 2},
	{},
	{1, 1},
	{2, 2},
}
var String_Accepting = [6]bool{
	3: true,
}
var String_Loop1 = [128]bool{0: true, 1: true, 2: true, 3: true, 4: true, 5: true, 6: true, 7: true, 8: true, 9: true, 10: true, 11: true, 12: true, 13: true, 14: true, 15: true, 16: true, 17: true, 18: true, 19: true, 20: true, 21: true, 22: true, 23: true, 24: true, 25: true, 26: true, 27: true, 28: true, 29: true, 30: true, 31: true, 32: true, 33: true, 35: true, 36: true, 37: true, 38: true, 39: true, 40: true, 41: true, 42: true, 43: true, 44: true, 45: true, 46: true, 47: true, 48: true, 49: true, 50: true, 51: true, 52: true, 53: true, 54: true, 55: true, 56: true, 57: true, 58: true, 59: true, 60: true, 61: true, 62: true, 63: true, 64: true, 65: true, 66: true, 67: true, 68: true, 69: true, 70: true, 71: true, 72: true, 73: true, 74: true, 75: true, 76: true, 77: true, 78: true, 79: true, 80: true, 81: true, 82: true, 83: true, 84: true, 85: true, 86: true, 87: true, 88: true, 89: true, 90: true, 91: true, 93: true, 94: true, 95: true, 96: true, 97: true, 98: true, 99: true, 100: true, 101: true, 102: true, 103: true, 104: true, 105: true, 106: true, 107: true, 108: true, 109: true, 110: true, 111: true, 112: true, 113: true, 114: true, 115: true, 116: true, 117: true, 118: true, 119: true, 120: true, 121: true, 122: true, 123: true, 124: true, 125: true, 126: true, 127: true}
var String_Loop2 = [128]bool{0: true, 1: true, 2: true, 3: true, 4: true, 5: true, 6: true, 7: true, 8: true, 9: true, 10: true, 11: true, 12: true, 13: true, 14: true, 15: true, 16: true, 17: true, 18: true, 19: true, 20: true, 21: true, 22: true, 23: true, 24: true, 25: true, 26: true, 27: true, 28: true, 29: true, 30: true, 31: true, 32: true, 33: true, 34: true, 35: true, 36: true, 37: true, 38: true, 40: true, 41: true, 42: true, 43: true, 44: true, 45: true, 46: true, 47: true, 48: true, 49: true, 50: true, 51: true, 52: true, 53: true, 54: true, 55: true, 56: true, 57: true, 58: true, 59: true, 60: true, 61: true, 62: true, 63: true, 64: true, 65: true, 66: true, 67: true, 68: true, 69: true, 70: true, 71: true, 72: true, 73: true, 74: true, 75: true, 76: true, 77: true, 78: true, 79: true, 80: true, 81: true, 82: true, 83: true, 84: true, 85: true, 86: true, 87: true, 88: true, 89: true, 90: true, 91: true, 93: true, 94: true, 95: true, 96: true, 97: true, 98: true, 99: true, 100: true, 101: true, 102: true, 103: true, 104: true, 105: true, 106: true, 107: true, 108: true, 109: true, 110: true, 111: true, 112: true, 113: true, 114: true, 115: true, 116: true, 117: true, 118: true, 119: true, 120: true, 121: true, 122: true, 123: true, 124: true, 125: true, 126: true, 127: true}

func String(s string, offset int) int {
	input := s[offset:]
	length := len(input)
	state := 0
	acceptedIndex := 0
	index := 0
loop:
	for index < length {
		switch state {
		case 1:
			for index < length {
				c := input[index]
				if c >= utf8.RuneSelf || !String_Loop1[c] {
					break
				}
				index++
			}
			if index == length {
				break loop
			}
		case 2:
			for index < length {
				c := input[index]
				if c >= utf8.RuneSelf || !String_Loop2[c] {
					break
				}
				index++
			}
			if index == length {
				break loop
			}
		}
		r, runeSize := utf8.DecodeRuneInString(input[index:])
		switch state {
		case 0:
			nextState := -1
			next := String_Next[0]
			lookup := String_Lookup[0]
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
			next := String_Next[1]
			lookup := String_Lookup[1]
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
		case 2:
			nextState := -1
			next := String_Next[2]
			lookup := String_Lookup[2]
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
		case 3:
			nextState := -1
			next := String_Next[3]
			lookup := String_Lookup[3]
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
		case 4:
			nextState := -1
			next := String_Next[4]
			lookup := String_Lookup[4]
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
		case 5:
			nextState := -1
			next := String_Next[5]
			lookup := String_Lookup[5]
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
		if String_Accepting[state] {
			acceptedIndex = index
		}
	}
	return acceptedIndex
}
