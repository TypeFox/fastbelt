package generated

import (
	"sort"
	"unicode/utf8"
)

func Email(s string, offset int) int {
	input := s[offset:]
	length := len(input)
	accepted := map[int]bool{5: true}
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
			lookup := []rune{0x2B, 0x2B, 0x2D, 0x2E, 0x30, 0x39, 0x41, 0x5A, 0x5F, 0x5F, 0x61, 0x7A}
			next := []int{1, 1, 1, 1, 1, 1}
			nextState := -1
			searchIndex := sort.Search(len(next), func(i int) bool {
				return lookup[i*2] > r
			}) - 1
			if searchIndex > -1 && lookup[searchIndex*2] <= r && r <= lookup[searchIndex*2+1] {
				return next[searchIndex]
			}
			if nextState > -1 {
				state = nextState
			} else {
				break loop
			}
		case 1:
			lookup := []rune{0x2B, 0x2B, 0x2D, 0x2E, 0x30, 0x39, 0x40, 0x40, 0x41, 0x5A, 0x5F, 0x5F, 0x61, 0x7A}
			next := []int{1, 1, 1, 2, 1, 1, 1}
			nextState := -1
			searchIndex := sort.Search(len(next), func(i int) bool {
				return lookup[i*2] > r
			}) - 1
			if searchIndex > -1 && lookup[searchIndex*2] <= r && r <= lookup[searchIndex*2+1] {
				return next[searchIndex]
			}
			if nextState > -1 {
				state = nextState
			} else {
				break loop
			}
		case 2:
			lookup := []rune{0x2D, 0x2E, 0x30, 0x39, 0x41, 0x5A, 0x5F, 0x5F, 0x61, 0x7A}
			next := []int{3, 3, 3, 3, 3}
			nextState := -1
			searchIndex := sort.Search(len(next), func(i int) bool {
				return lookup[i*2] > r
			}) - 1
			if searchIndex > -1 && lookup[searchIndex*2] <= r && r <= lookup[searchIndex*2+1] {
				return next[searchIndex]
			}
			if nextState > -1 {
				state = nextState
			} else {
				break loop
			}
		case 3:
			lookup := []rune{0x2D, 0x2D, 0x2E, 0x2E, 0x30, 0x39, 0x41, 0x5A, 0x5F, 0x5F, 0x61, 0x7A}
			next := []int{3, 4, 3, 3, 3, 3}
			nextState := -1
			searchIndex := sort.Search(len(next), func(i int) bool {
				return lookup[i*2] > r
			}) - 1
			if searchIndex > -1 && lookup[searchIndex*2] <= r && r <= lookup[searchIndex*2+1] {
				return next[searchIndex]
			}
			if nextState > -1 {
				state = nextState
			} else {
				break loop
			}
		case 4:
			lookup := []rune{0x2D, 0x2D, 0x2E, 0x2E, 0x30, 0x39, 0x41, 0x5A, 0x5F, 0x5F, 0x61, 0x7A}
			next := []int{5, 5, 5, 5, 5, 5}
			nextState := -1
			searchIndex := sort.Search(len(next), func(i int) bool {
				return lookup[i*2] > r
			}) - 1
			if searchIndex > -1 && lookup[searchIndex*2] <= r && r <= lookup[searchIndex*2+1] {
				return next[searchIndex]
			}
			if nextState > -1 {
				state = nextState
			} else {
				break loop
			}
		case 5:
			lookup := []rune{0x2D, 0x2D, 0x2E, 0x2E, 0x30, 0x39, 0x41, 0x5A, 0x5F, 0x5F, 0x61, 0x7A}
			next := []int{5, 5, 5, 5, 5, 5}
			nextState := -1
			searchIndex := sort.Search(len(next), func(i int) bool {
				return lookup[i*2] > r
			}) - 1
			if searchIndex > -1 && lookup[searchIndex*2] <= r && r <= lookup[searchIndex*2+1] {
				return next[searchIndex]
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
