package benchmarkGenerated

import (
	"unicode/utf8"
)

func IPv4(s string, offset int) int {
	input := s[offset:]
	length := len(input)
	accepted := map[int]bool{20: true, 21: true, 22: true, 23: true, 19: true}
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
			switch r {
			case 0x30:
				fallthrough
			case 0x31:
				state = 1
			case 0x32:
				state = 2
			case 0x33:
				fallthrough
			case 0x34:
				fallthrough
			case 0x35:
				fallthrough
			case 0x36:
				fallthrough
			case 0x37:
				fallthrough
			case 0x38:
				fallthrough
			case 0x39:
				state = 3
			default:
				break loop
			}
		case 1:
			switch r {
			case 0x2E:
				state = 4
			case 0x30:
				fallthrough
			case 0x31:
				fallthrough
			case 0x32:
				fallthrough
			case 0x33:
				fallthrough
			case 0x34:
				fallthrough
			case 0x35:
				fallthrough
			case 0x36:
				fallthrough
			case 0x37:
				fallthrough
			case 0x38:
				fallthrough
			case 0x39:
				state = 3
			default:
				break loop
			}
		case 2:
			switch r {
			case 0x2E:
				state = 4
			case 0x30:
				fallthrough
			case 0x31:
				fallthrough
			case 0x32:
				fallthrough
			case 0x33:
				fallthrough
			case 0x34:
				state = 3
			case 0x35:
				state = 5
			case 0x36:
				fallthrough
			case 0x37:
				fallthrough
			case 0x38:
				fallthrough
			case 0x39:
				state = 6
			default:
				break loop
			}
		case 3:
			switch r {
			case 0x2E:
				state = 4
			case 0x30:
				fallthrough
			case 0x31:
				fallthrough
			case 0x32:
				fallthrough
			case 0x33:
				fallthrough
			case 0x34:
				fallthrough
			case 0x35:
				fallthrough
			case 0x36:
				fallthrough
			case 0x37:
				fallthrough
			case 0x38:
				fallthrough
			case 0x39:
				state = 6
			default:
				break loop
			}
		case 4:
			switch r {
			case 0x30:
				fallthrough
			case 0x31:
				state = 7
			case 0x32:
				state = 8
			case 0x33:
				fallthrough
			case 0x34:
				fallthrough
			case 0x35:
				fallthrough
			case 0x36:
				fallthrough
			case 0x37:
				fallthrough
			case 0x38:
				fallthrough
			case 0x39:
				state = 9
			default:
				break loop
			}
		case 5:
			switch r {
			case 0x2E:
				state = 4
			case 0x30:
				fallthrough
			case 0x31:
				fallthrough
			case 0x32:
				fallthrough
			case 0x33:
				fallthrough
			case 0x34:
				fallthrough
			case 0x35:
				state = 6
			default:
				break loop
			}
		case 6:
			if r == 0x2E { // '.',
				state = 4
			} else {
				break loop
			}
		case 7:
			switch r {
			case 0x2E:
				state = 10
			case 0x30:
				fallthrough
			case 0x31:
				fallthrough
			case 0x32:
				fallthrough
			case 0x33:
				fallthrough
			case 0x34:
				fallthrough
			case 0x35:
				fallthrough
			case 0x36:
				fallthrough
			case 0x37:
				fallthrough
			case 0x38:
				fallthrough
			case 0x39:
				state = 9
			default:
				break loop
			}
		case 8:
			switch r {
			case 0x2E:
				state = 10
			case 0x30:
				fallthrough
			case 0x31:
				fallthrough
			case 0x32:
				fallthrough
			case 0x33:
				fallthrough
			case 0x34:
				state = 9
			case 0x35:
				state = 11
			case 0x36:
				fallthrough
			case 0x37:
				fallthrough
			case 0x38:
				fallthrough
			case 0x39:
				state = 12
			default:
				break loop
			}
		case 9:
			switch r {
			case 0x2E:
				state = 10
			case 0x30:
				fallthrough
			case 0x31:
				fallthrough
			case 0x32:
				fallthrough
			case 0x33:
				fallthrough
			case 0x34:
				fallthrough
			case 0x35:
				fallthrough
			case 0x36:
				fallthrough
			case 0x37:
				fallthrough
			case 0x38:
				fallthrough
			case 0x39:
				state = 12
			default:
				break loop
			}
		case 10:
			switch r {
			case 0x30:
				fallthrough
			case 0x31:
				state = 13
			case 0x32:
				state = 14
			case 0x33:
				fallthrough
			case 0x34:
				fallthrough
			case 0x35:
				fallthrough
			case 0x36:
				fallthrough
			case 0x37:
				fallthrough
			case 0x38:
				fallthrough
			case 0x39:
				state = 15
			default:
				break loop
			}
		case 11:
			switch r {
			case 0x2E:
				state = 10
			case 0x30:
				fallthrough
			case 0x31:
				fallthrough
			case 0x32:
				fallthrough
			case 0x33:
				fallthrough
			case 0x34:
				fallthrough
			case 0x35:
				state = 12
			default:
				break loop
			}
		case 12:
			if r == 0x2E { // '.',
				state = 10
			} else {
				break loop
			}
		case 13:
			switch r {
			case 0x2E:
				state = 16
			case 0x30:
				fallthrough
			case 0x31:
				fallthrough
			case 0x32:
				fallthrough
			case 0x33:
				fallthrough
			case 0x34:
				fallthrough
			case 0x35:
				fallthrough
			case 0x36:
				fallthrough
			case 0x37:
				fallthrough
			case 0x38:
				fallthrough
			case 0x39:
				state = 15
			default:
				break loop
			}
		case 14:
			switch r {
			case 0x2E:
				state = 16
			case 0x30:
				fallthrough
			case 0x31:
				fallthrough
			case 0x32:
				fallthrough
			case 0x33:
				fallthrough
			case 0x34:
				state = 15
			case 0x35:
				state = 17
			case 0x36:
				fallthrough
			case 0x37:
				fallthrough
			case 0x38:
				fallthrough
			case 0x39:
				state = 18
			default:
				break loop
			}
		case 15:
			switch r {
			case 0x2E:
				state = 16
			case 0x30:
				fallthrough
			case 0x31:
				fallthrough
			case 0x32:
				fallthrough
			case 0x33:
				fallthrough
			case 0x34:
				fallthrough
			case 0x35:
				fallthrough
			case 0x36:
				fallthrough
			case 0x37:
				fallthrough
			case 0x38:
				fallthrough
			case 0x39:
				state = 18
			default:
				break loop
			}
		case 16:
			switch r {
			case 0x30:
				fallthrough
			case 0x31:
				state = 19
			case 0x32:
				state = 20
			case 0x33:
				fallthrough
			case 0x34:
				fallthrough
			case 0x35:
				fallthrough
			case 0x36:
				fallthrough
			case 0x37:
				fallthrough
			case 0x38:
				fallthrough
			case 0x39:
				state = 21
			default:
				break loop
			}
		case 17:
			switch r {
			case 0x2E:
				state = 16
			case 0x30:
				fallthrough
			case 0x31:
				fallthrough
			case 0x32:
				fallthrough
			case 0x33:
				fallthrough
			case 0x34:
				fallthrough
			case 0x35:
				state = 18
			default:
				break loop
			}
		case 18:
			if r == 0x2E { // '.',
				state = 16
			} else {
				break loop
			}
		case 19:
			if r >= 0x30 && r <= 0x39 { // '0'..'9',
				state = 21
			} else {
				break loop
			}
		case 20:
			switch r {
			case 0x30:
				fallthrough
			case 0x31:
				fallthrough
			case 0x32:
				fallthrough
			case 0x33:
				fallthrough
			case 0x34:
				state = 21
			case 0x35:
				state = 22
			case 0x36:
				fallthrough
			case 0x37:
				fallthrough
			case 0x38:
				fallthrough
			case 0x39:
				state = 23
			default:
				break loop
			}
		case 21:
			if r >= 0x30 && r <= 0x39 { // '0'..'9',
				state = 23
			} else {
				break loop
			}
		case 22:
			if r >= 0x30 && r <= 0x35 { // '0'..'5',
				state = 23
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
