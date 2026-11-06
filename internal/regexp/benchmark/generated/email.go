package benchmarkGenerated

import (
	"unicode/utf8"
)

var Email_Lookup = [][]int64{
	{0x0000002B0000002B, 0x0000002E0000002D, 0x0000003900000030, 0x0000005A00000041, 0x0000005F0000005F, 0x0000007A00000061},
	{0x0000002B0000002B, 0x0000002E0000002D, 0x0000003900000030, 0x0000004000000040, 0x0000005A00000041, 0x0000005F0000005F, 0x0000007A00000061},
	{0x0000002E0000002D, 0x0000003900000030, 0x0000005A00000041, 0x0000005F0000005F, 0x0000007A00000061},
	{0x0000002D0000002D, 0x0000002E0000002E, 0x0000003900000030, 0x0000005A00000041, 0x0000005F0000005F, 0x0000007A00000061},
	{0x0000002D0000002D, 0x0000002E0000002E, 0x0000003900000030, 0x0000005A00000041, 0x0000005F0000005F, 0x0000007A00000061},
	{0x0000002D0000002D, 0x0000002E0000002E, 0x0000003900000030, 0x0000005A00000041, 0x0000005F0000005F, 0x0000007A00000061},
}
var Email_Next = [][]int{
	{1, 1, 1, 1, 1, 1},
	{1, 1, 1, 2, 1, 1, 1},
	{3, 3, 3, 3, 3},
	{3, 4, 3, 3, 3, 3},
	{5, 5, 5, 5, 5, 5},
	{5, 5, 5, 5, 5, 5},
}
var Email_Accepting = [6]bool{
	5: true,
}
var Email_Loop1 = [128]bool{43: true, 45: true, 46: true, 48: true, 49: true, 50: true, 51: true, 52: true, 53: true, 54: true, 55: true, 56: true, 57: true, 65: true, 66: true, 67: true, 68: true, 69: true, 70: true, 71: true, 72: true, 73: true, 74: true, 75: true, 76: true, 77: true, 78: true, 79: true, 80: true, 81: true, 82: true, 83: true, 84: true, 85: true, 86: true, 87: true, 88: true, 89: true, 90: true, 95: true, 97: true, 98: true, 99: true, 100: true, 101: true, 102: true, 103: true, 104: true, 105: true, 106: true, 107: true, 108: true, 109: true, 110: true, 111: true, 112: true, 113: true, 114: true, 115: true, 116: true, 117: true, 118: true, 119: true, 120: true, 121: true, 122: true}
var Email_Loop3 = [128]bool{45: true, 48: true, 49: true, 50: true, 51: true, 52: true, 53: true, 54: true, 55: true, 56: true, 57: true, 65: true, 66: true, 67: true, 68: true, 69: true, 70: true, 71: true, 72: true, 73: true, 74: true, 75: true, 76: true, 77: true, 78: true, 79: true, 80: true, 81: true, 82: true, 83: true, 84: true, 85: true, 86: true, 87: true, 88: true, 89: true, 90: true, 95: true, 97: true, 98: true, 99: true, 100: true, 101: true, 102: true, 103: true, 104: true, 105: true, 106: true, 107: true, 108: true, 109: true, 110: true, 111: true, 112: true, 113: true, 114: true, 115: true, 116: true, 117: true, 118: true, 119: true, 120: true, 121: true, 122: true}
var Email_Loop5 = [128]bool{45: true, 46: true, 48: true, 49: true, 50: true, 51: true, 52: true, 53: true, 54: true, 55: true, 56: true, 57: true, 65: true, 66: true, 67: true, 68: true, 69: true, 70: true, 71: true, 72: true, 73: true, 74: true, 75: true, 76: true, 77: true, 78: true, 79: true, 80: true, 81: true, 82: true, 83: true, 84: true, 85: true, 86: true, 87: true, 88: true, 89: true, 90: true, 95: true, 97: true, 98: true, 99: true, 100: true, 101: true, 102: true, 103: true, 104: true, 105: true, 106: true, 107: true, 108: true, 109: true, 110: true, 111: true, 112: true, 113: true, 114: true, 115: true, 116: true, 117: true, 118: true, 119: true, 120: true, 121: true, 122: true}

func Email(s string, offset int) int {
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
				if c >= utf8.RuneSelf || !Email_Loop1[c] {
					break
				}
				index++
			}
			if index == length {
				break loop
			}
		case 3:
			for index < length {
				c := input[index]
				if c >= utf8.RuneSelf || !Email_Loop3[c] {
					break
				}
				index++
			}
			if index == length {
				break loop
			}
		case 5:
			start := index
			for index < length {
				c := input[index]
				if c >= utf8.RuneSelf || !Email_Loop5[c] {
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
			next := Email_Next[0]
			lookup := Email_Lookup[0]
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
			next := Email_Next[1]
			lookup := Email_Lookup[1]
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
			next := Email_Next[2]
			lookup := Email_Lookup[2]
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
			next := Email_Next[3]
			lookup := Email_Lookup[3]
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
			next := Email_Next[4]
			lookup := Email_Lookup[4]
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
			next := Email_Next[5]
			lookup := Email_Lookup[5]
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
		if Email_Accepting[state] {
			acceptedIndex = index
		}
	}
	return acceptedIndex
}
