package benchmarkGenerated

import (
	"slices"
	"unicode/utf8"
)

var IPv4_Lookup = [][]int64{
	{0x0000003100000030, 0x0000003200000032, 0x0000003900000033},
	{0x0000002E0000002E, 0x0000003900000030},
	{0x0000002E0000002E, 0x0000003400000030, 0x0000003500000035, 0x0000003900000036},
	{0x0000002E0000002E, 0x0000003900000030},
	{0x0000003100000030, 0x0000003200000032, 0x0000003900000033},
	{0x0000002E0000002E, 0x0000003500000030},
	{0x0000002E0000002E},
	{0x0000002E0000002E, 0x0000003900000030},
	{0x0000002E0000002E, 0x0000003400000030, 0x0000003500000035, 0x0000003900000036},
	{0x0000002E0000002E, 0x0000003900000030},
	{0x0000003100000030, 0x0000003200000032, 0x0000003900000033},
	{0x0000002E0000002E, 0x0000003500000030},
	{0x0000002E0000002E},
	{0x0000002E0000002E, 0x0000003900000030},
	{0x0000002E0000002E, 0x0000003400000030, 0x0000003500000035, 0x0000003900000036},
	{0x0000002E0000002E, 0x0000003900000030},
	{0x0000003100000030, 0x0000003200000032, 0x0000003900000033},
	{0x0000002E0000002E, 0x0000003500000030},
	{0x0000002E0000002E},
	{0x0000003900000030},
	{0x0000003400000030, 0x0000003500000035, 0x0000003900000036},
	{0x0000003900000030},
	{0x0000003500000030},
	{},
}
var IPv4_Next = [][]int{
	{1, 2, 3},
	{4, 3},
	{4, 3, 5, 6},
	{4, 6},
	{7, 8, 9},
	{4, 6},
	{4},
	{10, 9},
	{10, 9, 11, 12},
	{10, 12},
	{13, 14, 15},
	{10, 12},
	{10},
	{16, 15},
	{16, 15, 17, 18},
	{16, 18},
	{19, 20, 21},
	{16, 18},
	{16},
	{21},
	{21, 22, 23},
	{23},
	{23},
	{},
}

func IPv4(s string, offset int) int {
	input := s[offset:]
	length := len(input)
	accepted := map[int]bool{19: true, 20: true, 21: true, 22: true, 23: true}
	state := 0
	acceptedIndex := 0
	index := 0
loop:
	for index < length {
		r, runeSize := utf8.DecodeRuneInString(input[index:])
		switch state {
		case 0:
			nextState := -1
			next := IPv4_Next[0]
			lookup := IPv4_Lookup[0]
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
			next := IPv4_Next[1]
			lookup := IPv4_Lookup[1]
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
			next := IPv4_Next[2]
			lookup := IPv4_Lookup[2]
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
			next := IPv4_Next[3]
			lookup := IPv4_Lookup[3]
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
			next := IPv4_Next[4]
			lookup := IPv4_Lookup[4]
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
			next := IPv4_Next[5]
			lookup := IPv4_Lookup[5]
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
		case 6:
			nextState := -1
			next := IPv4_Next[6]
			lookup := IPv4_Lookup[6]
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
		case 7:
			nextState := -1
			next := IPv4_Next[7]
			lookup := IPv4_Lookup[7]
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
		case 8:
			nextState := -1
			next := IPv4_Next[8]
			lookup := IPv4_Lookup[8]
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
		case 9:
			nextState := -1
			next := IPv4_Next[9]
			lookup := IPv4_Lookup[9]
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
		case 10:
			nextState := -1
			next := IPv4_Next[10]
			lookup := IPv4_Lookup[10]
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
		case 11:
			nextState := -1
			next := IPv4_Next[11]
			lookup := IPv4_Lookup[11]
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
		case 12:
			nextState := -1
			next := IPv4_Next[12]
			lookup := IPv4_Lookup[12]
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
		case 13:
			nextState := -1
			next := IPv4_Next[13]
			lookup := IPv4_Lookup[13]
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
		case 14:
			nextState := -1
			next := IPv4_Next[14]
			lookup := IPv4_Lookup[14]
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
		case 15:
			nextState := -1
			next := IPv4_Next[15]
			lookup := IPv4_Lookup[15]
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
		case 16:
			nextState := -1
			next := IPv4_Next[16]
			lookup := IPv4_Lookup[16]
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
		case 17:
			nextState := -1
			next := IPv4_Next[17]
			lookup := IPv4_Lookup[17]
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
		case 18:
			nextState := -1
			next := IPv4_Next[18]
			lookup := IPv4_Lookup[18]
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
		case 19:
			nextState := -1
			next := IPv4_Next[19]
			lookup := IPv4_Lookup[19]
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
		case 20:
			nextState := -1
			next := IPv4_Next[20]
			lookup := IPv4_Lookup[20]
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
		case 21:
			nextState := -1
			next := IPv4_Next[21]
			lookup := IPv4_Lookup[21]
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
		case 22:
			nextState := -1
			next := IPv4_Next[22]
			lookup := IPv4_Lookup[22]
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
		case 23:
			nextState := -1
			next := IPv4_Next[23]
			lookup := IPv4_Lookup[23]
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
		if accepted[state] {
			acceptedIndex = index
		}
	}
	return acceptedIndex
}
