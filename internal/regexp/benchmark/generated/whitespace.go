package benchmarkGenerated

import (
	"unicode/utf8"
)

var Whitespace_Lookup = [][]int64{
	{0x0000000A00000009, 0x0000000D0000000D, 0x0000002000000020},
	{0x0000000A00000009, 0x0000000D0000000D, 0x0000002000000020},
}
var Whitespace_Next = [][]int{
	{1, 1, 1},
	{1, 1, 1},
}
var Whitespace_Accepting = [2]bool{
	1: true,
}
var Whitespace_Loop1 = [128]bool{9: true, 10: true, 13: true, 32: true}

func Whitespace(s string, offset int) int {
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
				if c >= utf8.RuneSelf || !Whitespace_Loop1[c] {
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
			next := Whitespace_Next[0]
			lookup := Whitespace_Lookup[0]
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
			next := Whitespace_Next[1]
			lookup := Whitespace_Lookup[1]
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
		if Whitespace_Accepting[state] {
			acceptedIndex = index
		}
	}
	return acceptedIndex
}
