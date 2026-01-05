package benchmarkGenerated

import (
	"sort"
	"unicode/utf8"
)

var URL_Lookup = map[int][]rune{
	8:  {0x23, 0x23, 0x25, 0x25, 0x2B, 0x2B, 0x2D, 0x2E, 0x30, 0x3A, 0x3D, 0x3D, 0x40, 0x5A, 0x5F, 0x5F, 0x61, 0x76, 0x77, 0x77, 0x78, 0x7A, 0x7E, 0x7E},
	9:  {0x23, 0x23, 0x25, 0x25, 0x2B, 0x2B, 0x2D, 0x2D, 0x2E, 0x2E, 0x30, 0x3A, 0x3D, 0x3D, 0x40, 0x5A, 0x5F, 0x5F, 0x61, 0x7A, 0x7E, 0x7E},
	10: {0x23, 0x23, 0x25, 0x25, 0x28, 0x29, 0x2B, 0x2B, 0x2D, 0x2D, 0x2E, 0x2E, 0x30, 0x39, 0x3A, 0x3A, 0x3D, 0x3D, 0x40, 0x40, 0x41, 0x5A, 0x5F, 0x5F, 0x61, 0x7A, 0x7E, 0x7E},
	11: {0x23, 0x23, 0x25, 0x26, 0x28, 0x29, 0x2B, 0x2B, 0x2D, 0x2F, 0x30, 0x39, 0x3A, 0x3A, 0x3D, 0x3D, 0x3F, 0x40, 0x41, 0x5A, 0x5F, 0x5F, 0x61, 0x7A, 0x7E, 0x7E},
}
var URL_Next = map[int][]int{
	8:  {9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9},
	9:  {9, 9, 9, 9, 10, 9, 9, 9, 9, 9, 9},
	10: {9, 9, 11, 9, 9, 10, 11, 9, 9, 9, 11, 9, 11, 9},
	11: {11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11},
}

func URL(s string, offset int) int {
	input := s[offset:]
	length := len(input)
	accepted := map[int]bool{11: true}
	state := 0
	acceptedIndex := -1
	if accepted[state] {
		acceptedIndex = 0
	}
	index := 0
loop:
	for index < length {
		r, runeSize := utf8.DecodeRuneInString(input[index:])
		switch state {
		case 0:
			if r == 0x68 { // 'h',
				state = 1
			} else {
				break loop
			}
		case 1:
			if r == 0x74 { // 't',
				state = 2
			} else {
				break loop
			}
		case 2:
			if r == 0x74 { // 't',
				state = 3
			} else {
				break loop
			}
		case 3:
			if r == 0x70 { // 'p',
				state = 4
			} else {
				break loop
			}
		case 4:
			switch r {
			case 0x3A:
				state = 5
			case 0x73:
				state = 6
			default:
				break loop
			}
		case 5:
			if r == 0x2F { // '/',
				state = 7
			} else {
				break loop
			}
		case 6:
			if r == 0x3A { // ':',
				state = 5
			} else {
				break loop
			}
		case 7:
			if r == 0x2F { // '/',
				state = 8
			} else {
				break loop
			}
		case 8:
			nextState := -1
			next := URL_Next[8]
			lookup := URL_Lookup[8]
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
		case 9:
			nextState := -1
			next := URL_Next[9]
			lookup := URL_Lookup[9]
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
		case 10:
			nextState := -1
			next := URL_Next[10]
			lookup := URL_Lookup[10]
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
		case 11:
			nextState := -1
			next := URL_Next[11]
			lookup := URL_Lookup[11]
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
