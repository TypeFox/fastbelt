package benchmarkGenerated

import (
	"unicode/utf8"
)

var Email_Lookup = [][]int64{
	{0x0000002A00000000, 0x0000002B0000002B, 0x0000002C0000002C, 0x0000002E0000002D, 0x0000002F0000002F, 0x0000003900000030, 0x000000400000003A, 0x0000005A00000041, 0x0000005E0000005B, 0x0000005F0000005F, 0x0000006000000060, 0x0000007A00000061, 0x0010FFFF0000007B},
	{0x0010FFFF00000000},
	{0x0000002A00000000, 0x0000002B0000002B, 0x0000002C0000002C, 0x0000002E0000002D, 0x0000002F0000002F, 0x0000003900000030, 0x0000003F0000003A, 0x0000004000000040, 0x0000005A00000041, 0x0000005E0000005B, 0x0000005F0000005F, 0x0000006000000060, 0x0000007A00000061, 0x0010FFFF0000007B},
	{0x0000002C00000000, 0x0000002E0000002D, 0x0000002F0000002F, 0x0000003900000030, 0x000000400000003A, 0x0000005A00000041, 0x0000005E0000005B, 0x0000005F0000005F, 0x0000006000000060, 0x0000007A00000061, 0x0010FFFF0000007B},
	{0x0000002C00000000, 0x0000002D0000002D, 0x0000002E0000002E, 0x0000002F0000002F, 0x0000003900000030, 0x000000400000003A, 0x0000005A00000041, 0x0000005E0000005B, 0x0000005F0000005F, 0x0000006000000060, 0x0000007A00000061, 0x0010FFFF0000007B},
	{0x0000002C00000000, 0x0000002D0000002D, 0x0000002E0000002E, 0x0000002F0000002F, 0x0000003900000030, 0x000000400000003A, 0x0000005A00000041, 0x0000005E0000005B, 0x0000005F0000005F, 0x0000006000000060, 0x0000007A00000061, 0x0010FFFF0000007B},
	{0x0000002C00000000, 0x0000002D0000002D, 0x0000002E0000002E, 0x0000002F0000002F, 0x0000003900000030, 0x000000400000003A, 0x0000005A00000041, 0x0000005E0000005B, 0x0000005F0000005F, 0x0000006000000060, 0x0000007A00000061, 0x0010FFFF0000007B},
}
var Email_Next = [][]int{
	{1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1},
	{1},
	{1, 2, 1, 2, 1, 2, 1, 3, 2, 1, 2, 1, 2, 1},
	{1, 4, 1, 4, 1, 4, 1, 4, 1, 4, 1},
	{1, 4, 5, 1, 4, 1, 4, 1, 4, 1, 4, 1},
	{1, 6, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1},
	{1, 6, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1},
}
var Email_Accepting = [7]bool{
	6: true,
}

func Email(s string, offset int) int {
	input := s[offset:]
	length := len(input)
	state := 0
	acceptedIndex := 0
	index := 0
loop:
	for index < length {
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
		case 6:
			nextState := -1
			next := Email_Next[6]
			lookup := Email_Lookup[6]
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
