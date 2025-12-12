package main

func IP_FindStringIndex(input string) (loc []int) {
	length := len(input)
	accepted := map[int]bool{14: true, 15: true}
	state := 0
	acceptedIndex := -1
	if accepted[state] {
		acceptedIndex = 0
	}
	index := 0
loop:
	for index < length {
		r := rune(input[index])
		switch state {
		case 0:
			if r >= 48 && r <= 57 { // 0..9,
				state = 1
			} else {
				break loop
			}
		case 1:
			if r >= 48 && r <= 57 { // 0..9,
				state = 2
			} else {
				break loop
			}
		case 2:
			if r == 46 { // .,
				state = 3
			} else if r >= 48 && r <= 57 { // 0..9,
				state = 4
			} else {
				break loop
			}
		case 3:
			if r >= 48 && r <= 57 { // 0..9,
				state = 5
			} else {
				break loop
			}
		case 4:
			if r == 46 { // .,
				state = 3
			} else {
				break loop
			}
		case 5:
			if r >= 48 && r <= 57 { // 0..9,
				state = 6
			} else {
				break loop
			}
		case 6:
			if r == 46 { // .,
				state = 7
			} else if r >= 48 && r <= 57 { // 0..9,
				state = 8
			} else {
				break loop
			}
		case 7:
			if r >= 48 && r <= 57 { // 0..9,
				state = 9
			} else {
				break loop
			}
		case 8:
			if r == 46 { // .,
				state = 7
			} else {
				break loop
			}
		case 9:
			if r >= 48 && r <= 57 { // 0..9,
				state = 10
			} else {
				break loop
			}
		case 10:
			if r == 46 { // .,
				state = 11
			} else if r >= 48 && r <= 57 { // 0..9,
				state = 12
			} else {
				break loop
			}
		case 11:
			if r >= 48 && r <= 57 { // 0..9,
				state = 13
			} else {
				break loop
			}
		case 12:
			if r == 46 { // .,
				state = 11
			} else {
				break loop
			}
		case 13:
			if r >= 48 && r <= 57 { // 0..9,
				state = 14
			} else {
				break loop
			}
		case 14:
			if r >= 48 && r <= 57 { // 0..9,
				state = 15
			} else {
				break loop
			}
		}
		if accepted[state] {
			acceptedIndex = index
		}
		index++
	}
	if acceptedIndex == -1 {
		return nil
	}
	return []int{0, acceptedIndex}
}
