package benchmarkGenerated

import (
	"slices"
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
			searchIndex := slices.IndexFunc(lookup, func(lowHigh int64) bool {
				lo := rune(lowHigh & 0xFFFFFFFF)
				hi := rune(lowHigh >> 32)
				return lo <= r && r <= hi
			})
			if searchIndex > -1 {
				nextState = next[searchIndex]
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
			searchIndex := slices.IndexFunc(lookup, func(lowHigh int64) bool {
				lo := rune(lowHigh & 0xFFFFFFFF)
				hi := rune(lowHigh >> 32)
				return lo <= r && r <= hi
			})
			if searchIndex > -1 {
				nextState = next[searchIndex]
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
			searchIndex := slices.IndexFunc(lookup, func(lowHigh int64) bool {
				lo := rune(lowHigh & 0xFFFFFFFF)
				hi := rune(lowHigh >> 32)
				return lo <= r && r <= hi
			})
			if searchIndex > -1 {
				nextState = next[searchIndex]
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
			searchIndex := slices.IndexFunc(lookup, func(lowHigh int64) bool {
				lo := rune(lowHigh & 0xFFFFFFFF)
				hi := rune(lowHigh >> 32)
				return lo <= r && r <= hi
			})
			if searchIndex > -1 {
				nextState = next[searchIndex]
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
			searchIndex := slices.IndexFunc(lookup, func(lowHigh int64) bool {
				lo := rune(lowHigh & 0xFFFFFFFF)
				hi := rune(lowHigh >> 32)
				return lo <= r && r <= hi
			})
			if searchIndex > -1 {
				nextState = next[searchIndex]
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
			searchIndex := slices.IndexFunc(lookup, func(lowHigh int64) bool {
				lo := rune(lowHigh & 0xFFFFFFFF)
				hi := rune(lowHigh >> 32)
				return lo <= r && r <= hi
			})
			if searchIndex > -1 {
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
		if Email_Accepting[state] {
			acceptedIndex = index
		}
	}
	return acceptedIndex
}
