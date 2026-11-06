package benchmarkGenerated

import (
	"unicode/utf8"
)

var URL_Lookup = [][]int64{
	{0x0000006800000068},
	{0x0000007400000074},
	{0x0000007400000074},
	{0x0000007000000070},
	{0x0000003A0000003A, 0x0000007300000073},
	{0x0000002F0000002F},
	{0x0000003A0000003A},
	{0x0000002F0000002F},
	{0x0000002300000023, 0x0000002500000025, 0x0000002B0000002B, 0x0000002E0000002D, 0x0000003A00000030, 0x0000003D0000003D, 0x0000005A00000040, 0x0000005F0000005F, 0x0000007600000061, 0x0000007700000077, 0x0000007A00000078, 0x0000007E0000007E},
	{0x0000002300000023, 0x0000002500000025, 0x0000002B0000002B, 0x0000002D0000002D, 0x0000002E0000002E, 0x0000003A00000030, 0x0000003D0000003D, 0x0000005A00000040, 0x0000005F0000005F, 0x0000007A00000061, 0x0000007E0000007E},
	{0x0000002300000023, 0x0000002500000025, 0x0000002900000028, 0x0000002B0000002B, 0x0000002D0000002D, 0x0000002E0000002E, 0x0000003900000030, 0x0000003A0000003A, 0x0000003D0000003D, 0x0000004000000040, 0x0000005A00000041, 0x0000005F0000005F, 0x0000007A00000061, 0x0000007E0000007E},
	{0x0000002300000023, 0x0000002600000025, 0x0000002900000028, 0x0000002B0000002B, 0x0000002F0000002D, 0x0000003900000030, 0x0000003A0000003A, 0x0000003D0000003D, 0x000000400000003F, 0x0000005A00000041, 0x0000005F0000005F, 0x0000007A00000061, 0x0000007E0000007E},
}
var URL_Next = [][]int{
	{1},
	{2},
	{3},
	{4},
	{5, 6},
	{7},
	{5},
	{8},
	{9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9},
	{9, 9, 9, 9, 10, 9, 9, 9, 9, 9, 9},
	{9, 9, 11, 9, 9, 10, 11, 9, 9, 9, 11, 9, 11, 9},
	{11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11},
}
var URL_Accepting = [12]bool{
	11: true,
}
var URL_Loop9 = [128]bool{35: true, 37: true, 43: true, 45: true, 48: true, 49: true, 50: true, 51: true, 52: true, 53: true, 54: true, 55: true, 56: true, 57: true, 58: true, 61: true, 64: true, 65: true, 66: true, 67: true, 68: true, 69: true, 70: true, 71: true, 72: true, 73: true, 74: true, 75: true, 76: true, 77: true, 78: true, 79: true, 80: true, 81: true, 82: true, 83: true, 84: true, 85: true, 86: true, 87: true, 88: true, 89: true, 90: true, 95: true, 97: true, 98: true, 99: true, 100: true, 101: true, 102: true, 103: true, 104: true, 105: true, 106: true, 107: true, 108: true, 109: true, 110: true, 111: true, 112: true, 113: true, 114: true, 115: true, 116: true, 117: true, 118: true, 119: true, 120: true, 121: true, 122: true, 126: true}
var URL_Loop10 = [128]bool{46: true}
var URL_Loop11 = [128]bool{35: true, 37: true, 38: true, 40: true, 41: true, 43: true, 45: true, 46: true, 47: true, 48: true, 49: true, 50: true, 51: true, 52: true, 53: true, 54: true, 55: true, 56: true, 57: true, 58: true, 61: true, 63: true, 64: true, 65: true, 66: true, 67: true, 68: true, 69: true, 70: true, 71: true, 72: true, 73: true, 74: true, 75: true, 76: true, 77: true, 78: true, 79: true, 80: true, 81: true, 82: true, 83: true, 84: true, 85: true, 86: true, 87: true, 88: true, 89: true, 90: true, 95: true, 97: true, 98: true, 99: true, 100: true, 101: true, 102: true, 103: true, 104: true, 105: true, 106: true, 107: true, 108: true, 109: true, 110: true, 111: true, 112: true, 113: true, 114: true, 115: true, 116: true, 117: true, 118: true, 119: true, 120: true, 121: true, 122: true, 126: true}

func URL(s string, offset int) int {
	input := s[offset:]
	length := len(input)
	state := 0
	acceptedIndex := 0
	index := 0
loop:
	for index < length {
		switch state {
		case 9:
			for index < length {
				c := input[index]
				if c >= utf8.RuneSelf || !URL_Loop9[c] {
					break
				}
				index++
			}
			if index == length {
				break loop
			}
		case 10:
			for index < length {
				c := input[index]
				if c >= utf8.RuneSelf || !URL_Loop10[c] {
					break
				}
				index++
			}
			if index == length {
				break loop
			}
		case 11:
			start := index
			for index < length {
				c := input[index]
				if c >= utf8.RuneSelf || !URL_Loop11[c] {
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
			next := URL_Next[0]
			lookup := URL_Lookup[0]
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
			next := URL_Next[1]
			lookup := URL_Lookup[1]
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
			next := URL_Next[2]
			lookup := URL_Lookup[2]
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
			next := URL_Next[3]
			lookup := URL_Lookup[3]
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
			next := URL_Next[4]
			lookup := URL_Lookup[4]
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
			next := URL_Next[5]
			lookup := URL_Lookup[5]
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
		case 6:
			nextState := -1
			next := URL_Next[6]
			lookup := URL_Lookup[6]
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
		case 7:
			nextState := -1
			next := URL_Next[7]
			lookup := URL_Lookup[7]
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
		case 8:
			nextState := -1
			next := URL_Next[8]
			lookup := URL_Lookup[8]
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
		case 9:
			nextState := -1
			next := URL_Next[9]
			lookup := URL_Lookup[9]
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
		case 10:
			nextState := -1
			next := URL_Next[10]
			lookup := URL_Lookup[10]
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
		case 11:
			nextState := -1
			next := URL_Next[11]
			lookup := URL_Lookup[11]
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
		if URL_Accepting[state] {
			acceptedIndex = index
		}
	}
	return acceptedIndex
}
