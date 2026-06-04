package benchmarkGenerated

import (
	"sort"
	"unicode/utf8"
)

var URL_Lookup = [][]int64{
	{0x0000006700000000, 0x0000006800000068, 0x0010FFFF00000069},
	{},
	{0x0000007300000000, 0x0000007400000074, 0x0010FFFF00000075},
	{0x0000007300000000, 0x0000007400000074, 0x0010FFFF00000075},
	{0x0000006F00000000, 0x0000007000000070, 0x0010FFFF00000071},
	{0x0000003900000000, 0x0000003A0000003A, 0x000000720000003B, 0x0000007300000073, 0x0010FFFF00000074},
	{0x0000002E00000000, 0x0000002F0000002F, 0x0010FFFF00000030},
	{0x0000003900000000, 0x0000003A0000003A, 0x0010FFFF0000003B},
	{0x0000002E00000000, 0x0000002F0000002F, 0x0010FFFF00000030},
	{0x0000002200000000, 0x0000002300000023, 0x0000002400000024, 0x0000002500000025, 0x0000002A00000026, 0x0000002B0000002B, 0x0000002C0000002C, 0x0000002E0000002D, 0x0000002F0000002F, 0x0000003A00000030, 0x0000003C0000003B, 0x0000003D0000003D, 0x0000003F0000003E, 0x0000005A00000040, 0x0000005E0000005B, 0x0000005F0000005F, 0x0000006000000060, 0x0000007600000061, 0x0000007700000077, 0x0000007A00000078, 0x0000007D0000007B, 0x0000007E0000007E, 0x0010FFFF0000007F},
	{0x0000002200000000, 0x0000002300000023, 0x0000002400000024, 0x0000002500000025, 0x0000002A00000026, 0x0000002B0000002B, 0x0000002C0000002C, 0x0000002D0000002D, 0x0000002E0000002E, 0x0000002F0000002F, 0x0000003A00000030, 0x0000003C0000003B, 0x0000003D0000003D, 0x0000003F0000003E, 0x0000005A00000040, 0x0000005E0000005B, 0x0000005F0000005F, 0x0000006000000060, 0x0000007A00000061, 0x0000007D0000007B, 0x0000007E0000007E, 0x0010FFFF0000007F},
	{0x0000002200000000, 0x0000002300000023, 0x0000002400000024, 0x0000002500000025, 0x0000002700000026, 0x0000002900000028, 0x0000002A0000002A, 0x0000002B0000002B, 0x0000002C0000002C, 0x0000002D0000002D, 0x0000002E0000002E, 0x0000002F0000002F, 0x0000003900000030, 0x0000003A0000003A, 0x0000003C0000003B, 0x0000003D0000003D, 0x0000003F0000003E, 0x0000004000000040, 0x0000005A00000041, 0x0000005E0000005B, 0x0000005F0000005F, 0x0000006000000060, 0x0000007A00000061, 0x0000007D0000007B, 0x0000007E0000007E, 0x0010FFFF0000007F},
	{0x0000002200000000, 0x0000002300000023, 0x0000002400000024, 0x0000002600000025, 0x0000002700000027, 0x0000002900000028, 0x0000002A0000002A, 0x0000002B0000002B, 0x0000002C0000002C, 0x0000002F0000002D, 0x0000003900000030, 0x0000003A0000003A, 0x0000003C0000003B, 0x0000003D0000003D, 0x0000003E0000003E, 0x000000400000003F, 0x0000005A00000041, 0x0000005E0000005B, 0x0000005F0000005F, 0x0000006000000060, 0x0000007A00000061, 0x0000007D0000007B, 0x0000007E0000007E, 0x0010FFFF0000007F},
}
var URL_Next = [][]int{
	{1, 2, 1},
	{},
	{1, 3, 1},
	{1, 4, 1},
	{1, 5, 1},
	{1, 6, 1, 7, 1},
	{1, 8, 1},
	{1, 6, 1},
	{1, 9, 1},
	{1, 10, 1, 10, 1, 10, 1, 10, 1, 10, 1, 10, 1, 10, 1, 10, 1, 10, 10, 10, 1, 10, 1},
	{1, 10, 1, 10, 1, 10, 1, 10, 11, 1, 10, 1, 10, 1, 10, 1, 10, 1, 10, 1, 10, 1},
	{1, 10, 1, 10, 1, 12, 1, 10, 1, 10, 11, 1, 12, 10, 1, 10, 1, 10, 12, 1, 10, 1, 12, 1, 10, 1},
	{1, 12, 1, 12, 1, 12, 1, 12, 1, 12, 12, 12, 1, 12, 1, 12, 12, 1, 12, 1, 12, 1, 12, 1},
}
var URL_Accepting = [13]bool{
	12: true,
}

func URL(s string, offset int) int {
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
			searchIndex := sort.Search(len(next), func(i int) bool {
				return rune(lookup[i]&0xFFFFFFFF) > r
			}) - 1
			if searchIndex > -1 && rune(lookup[searchIndex]&0xFFFFFFFF) <= r && r <= rune(lookup[searchIndex]>>32) {
				nextState = next[searchIndex]
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
			searchIndex := sort.Search(len(next), func(i int) bool {
				return rune(lookup[i]&0xFFFFFFFF) > r
			}) - 1
			if searchIndex > -1 && rune(lookup[searchIndex]&0xFFFFFFFF) <= r && r <= rune(lookup[searchIndex]>>32) {
				nextState = next[searchIndex]
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
			searchIndex := sort.Search(len(next), func(i int) bool {
				return rune(lookup[i]&0xFFFFFFFF) > r
			}) - 1
			if searchIndex > -1 && rune(lookup[searchIndex]&0xFFFFFFFF) <= r && r <= rune(lookup[searchIndex]>>32) {
				nextState = next[searchIndex]
			}
			if nextState > -1 {
				state = nextState
			} else {
				break loop
			}
		case 12:
			nextState := -1
			next := URL_Next[12]
			lookup := URL_Lookup[12]
			searchIndex := sort.Search(len(next), func(i int) bool {
				return rune(lookup[i]&0xFFFFFFFF) > r
			}) - 1
			if searchIndex > -1 && rune(lookup[searchIndex]&0xFFFFFFFF) <= r && r <= rune(lookup[searchIndex]>>32) {
				nextState = next[searchIndex]
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
