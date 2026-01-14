package benchmarkGenerated

import (
	"unicode/utf8"
	"sort"
)

var Email_Lookup = map[int][]rune{
0: {0x2B, 0x2B, 0x2D, 0x2E, 0x30, 0x39, 0x41, 0x5A, 0x5F, 0x5F, 0x61, 0x7A, },
1: {0x2B, 0x2B, 0x2D, 0x2E, 0x30, 0x39, 0x40, 0x40, 0x41, 0x5A, 0x5F, 0x5F, 0x61, 0x7A, },
2: {0x2D, 0x2E, 0x30, 0x39, 0x41, 0x5A, 0x5F, 0x5F, 0x61, 0x7A, },
3: {0x2D, 0x2D, 0x2E, 0x2E, 0x30, 0x39, 0x41, 0x5A, 0x5F, 0x5F, 0x61, 0x7A, },
4: {0x2D, 0x2D, 0x2E, 0x2E, 0x30, 0x39, 0x41, 0x5A, 0x5F, 0x5F, 0x61, 0x7A, },
5: {0x2D, 0x2D, 0x2E, 0x2E, 0x30, 0x39, 0x41, 0x5A, 0x5F, 0x5F, 0x61, 0x7A, },
}
var Email_Next = map[int][]int{
0: {1, 1, 1, 1, 1, 1, },
1: {1, 1, 1, 2, 1, 1, 1, },
2: {3, 3, 3, 3, 3, },
3: {3, 4, 3, 3, 3, 3, },
4: {5, 5, 5, 5, 5, 5, },
5: {5, 5, 5, 5, 5, 5, },
}
func Email(s string, offset int) int {
    input := s[offset:]
    length := len(input)
    accepted := map[int]bool{5: true, }
    state := 0
    acceptedIndex := 0
    index := 0
    loop: for index < length {
        r, runeSize := utf8.DecodeRuneInString(input[index:])
        switch state {
        case 0:
nextState := -1
next := Email_Next[0]
lookup := Email_Lookup[0]
searchIndex := sort.Search(len(next), func(i int) bool {
    return lookup[i*2] > r
}) - 1
if searchIndex > -1 && lookup[searchIndex*2] <= r && r <= lookup[searchIndex*2+1] {
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
searchIndex := sort.Search(len(next), func(i int) bool {
    return lookup[i*2] > r
}) - 1
if searchIndex > -1 && lookup[searchIndex*2] <= r && r <= lookup[searchIndex*2+1] {
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
searchIndex := sort.Search(len(next), func(i int) bool {
    return lookup[i*2] > r
}) - 1
if searchIndex > -1 && lookup[searchIndex*2] <= r && r <= lookup[searchIndex*2+1] {
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
searchIndex := sort.Search(len(next), func(i int) bool {
    return lookup[i*2] > r
}) - 1
if searchIndex > -1 && lookup[searchIndex*2] <= r && r <= lookup[searchIndex*2+1] {
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
searchIndex := sort.Search(len(next), func(i int) bool {
    return lookup[i*2] > r
}) - 1
if searchIndex > -1 && lookup[searchIndex*2] <= r && r <= lookup[searchIndex*2+1] {
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
searchIndex := sort.Search(len(next), func(i int) bool {
    return lookup[i*2] > r
}) - 1
if searchIndex > -1 && lookup[searchIndex*2] <= r && r <= lookup[searchIndex*2+1] {
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