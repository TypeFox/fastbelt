package benchmarkGenerated

import (
	"sort"
	"unicode/utf8"
)

var URL_Lookup = [][]int64{
{0x0000006800000068, },
{0x0000007400000074, },
{0x0000007400000074, },
{0x0000007000000070, },
{0x0000003A0000003A, 0x0000007300000073, },
{0x0000002F0000002F, },
{0x0000003A0000003A, },
{0x0000002F0000002F, },
{0x0000002300000023, 0x0000002500000025, 0x0000002B0000002B, 0x0000002E0000002D, 0x0000003A00000030, 0x0000003D0000003D, 0x0000005A00000040, 0x0000005F0000005F, 0x0000007600000061, 0x0000007700000077, 0x0000007A00000078, 0x0000007E0000007E, },
{0x0000002300000023, 0x0000002500000025, 0x0000002B0000002B, 0x0000002D0000002D, 0x0000002E0000002E, 0x0000003A00000030, 0x0000003D0000003D, 0x0000005A00000040, 0x0000005F0000005F, 0x0000007A00000061, 0x0000007E0000007E, },
{0x0000002300000023, 0x0000002500000025, 0x0000002900000028, 0x0000002B0000002B, 0x0000002D0000002D, 0x0000002E0000002E, 0x0000003900000030, 0x0000003A0000003A, 0x0000003D0000003D, 0x0000004000000040, 0x0000005A00000041, 0x0000005F0000005F, 0x0000007A00000061, 0x0000007E0000007E, },
{0x0000002300000023, 0x0000002600000025, 0x0000002900000028, 0x0000002B0000002B, 0x0000002F0000002D, 0x0000003900000030, 0x0000003A0000003A, 0x0000003D0000003D, 0x000000400000003F, 0x0000005A00000041, 0x0000005F0000005F, 0x0000007A00000061, 0x0000007E0000007E, },
}
var URL_Next = [][]int{
{1, },
{2, },
{3, },
{4, },
{5, 6, },
{7, },
{5, },
{8, },
{9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, },
{9, 9, 9, 9, 10, 9, 9, 9, 9, 9, 9, },
{9, 9, 11, 9, 9, 10, 11, 9, 9, 9, 11, 9, 11, 9, },
{11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11, },
}
func URL(s string, offset int) int {
    input := s[offset:]
    length := len(input)
    accepted := map[int]bool{11: true, }
    state := 0
    acceptedIndex := 0
    index := 0
    loop: for index < length {
        r, runeSize := utf8.DecodeRuneInString(input[index:])
        switch state {
        case 0:
nextState := -1
next := URL_Next[0]
lookup := URL_Lookup[0]
searchIndex := sort.Search(len(next), func(i int) bool {
    return rune(lookup[i] & 0xFFFFFFFF) > r
}) - 1
if searchIndex > -1 && rune(lookup[searchIndex] & 0xFFFFFFFF) <= r && r <= rune(lookup[searchIndex] >> 32) {
    nextState = next[searchIndex]
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
searchIndex := sort.Search(len(next), func(i int) bool {
    return rune(lookup[i] & 0xFFFFFFFF) > r
}) - 1
if searchIndex > -1 && rune(lookup[searchIndex] & 0xFFFFFFFF) <= r && r <= rune(lookup[searchIndex] >> 32) {
    nextState = next[searchIndex]
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
searchIndex := sort.Search(len(next), func(i int) bool {
    return rune(lookup[i] & 0xFFFFFFFF) > r
}) - 1
if searchIndex > -1 && rune(lookup[searchIndex] & 0xFFFFFFFF) <= r && r <= rune(lookup[searchIndex] >> 32) {
    nextState = next[searchIndex]
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
searchIndex := sort.Search(len(next), func(i int) bool {
    return rune(lookup[i] & 0xFFFFFFFF) > r
}) - 1
if searchIndex > -1 && rune(lookup[searchIndex] & 0xFFFFFFFF) <= r && r <= rune(lookup[searchIndex] >> 32) {
    nextState = next[searchIndex]
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
searchIndex := sort.Search(len(next), func(i int) bool {
    return rune(lookup[i] & 0xFFFFFFFF) > r
}) - 1
if searchIndex > -1 && rune(lookup[searchIndex] & 0xFFFFFFFF) <= r && r <= rune(lookup[searchIndex] >> 32) {
    nextState = next[searchIndex]
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
searchIndex := sort.Search(len(next), func(i int) bool {
    return rune(lookup[i] & 0xFFFFFFFF) > r
}) - 1
if searchIndex > -1 && rune(lookup[searchIndex] & 0xFFFFFFFF) <= r && r <= rune(lookup[searchIndex] >> 32) {
    nextState = next[searchIndex]
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
searchIndex := sort.Search(len(next), func(i int) bool {
    return rune(lookup[i] & 0xFFFFFFFF) > r
}) - 1
if searchIndex > -1 && rune(lookup[searchIndex] & 0xFFFFFFFF) <= r && r <= rune(lookup[searchIndex] >> 32) {
    nextState = next[searchIndex]
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
searchIndex := sort.Search(len(next), func(i int) bool {
    return rune(lookup[i] & 0xFFFFFFFF) > r
}) - 1
if searchIndex > -1 && rune(lookup[searchIndex] & 0xFFFFFFFF) <= r && r <= rune(lookup[searchIndex] >> 32) {
    nextState = next[searchIndex]
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
searchIndex := sort.Search(len(next), func(i int) bool {
    return rune(lookup[i] & 0xFFFFFFFF) > r
}) - 1
if searchIndex > -1 && rune(lookup[searchIndex] & 0xFFFFFFFF) <= r && r <= rune(lookup[searchIndex] >> 32) {
    nextState = next[searchIndex]
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
    return rune(lookup[i] & 0xFFFFFFFF) > r
}) - 1
if searchIndex > -1 && rune(lookup[searchIndex] & 0xFFFFFFFF) <= r && r <= rune(lookup[searchIndex] >> 32) {
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
    return rune(lookup[i] & 0xFFFFFFFF) > r
}) - 1
if searchIndex > -1 && rune(lookup[searchIndex] & 0xFFFFFFFF) <= r && r <= rune(lookup[searchIndex] >> 32) {
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
    return rune(lookup[i] & 0xFFFFFFFF) > r
}) - 1
if searchIndex > -1 && rune(lookup[searchIndex] & 0xFFFFFFFF) <= r && r <= rune(lookup[searchIndex] >> 32) {
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