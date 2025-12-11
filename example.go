package main

func APlusB_FindStringIndex(input string) (loc []int) {
    length := len(input)
    accepted := map[int]bool{14: true, 15: true, }
    state := 0
    acceptedIndex := -1
    if accepted[state] {
        acceptedIndex = 0
    }
    index := 0
    halted := false
    for !halted {
        if index >= length {
            halted = true
            continue
        } else {
            r := rune(input[index])
            switch state {
            case 14:
                if r >= '0' && r <= '9' {
                    state = 15
                } else {
                    halted = true
                }
            case 6:
                if r == '.' {
                    state = 7
                } else if r >= '0' && r <= '9' {
                    state = 8
                } else {
                    halted = true
                }
            case 2:
                if r == '.' {
                    state = 3
                } else if r >= '0' && r <= '9' {
                    state = 4
                } else {
                    halted = true
                }
            case 4:
                if r == '.' {
                    state = 3
                } else {
                    halted = true
                }
            case 9:
                if r >= '0' && r <= '9' {
                    state = 10
                } else {
                    halted = true
                }
            case 10:
                if r == '.' {
                    state = 11
                } else if r >= '0' && r <= '9' {
                    state = 12
                } else {
                    halted = true
                }
            case 5:
                if r >= '0' && r <= '9' {
                    state = 6
                } else {
                    halted = true
                }
            case 12:
                if r == '.' {
                    state = 11
                } else {
                    halted = true
                }
            case 13:
                if r >= '0' && r <= '9' {
                    state = 14
                } else {
                    halted = true
                }
            case 7:
                if r >= '0' && r <= '1' {
                    state = 9
                } else if r == '2' {
                    state = 9
                } else if r >= '3' && r <= '9' {
                    state = 9
                } else {
                    halted = true
                }
            case 8:
                if r == '.' {
                    state = 7
                } else {
                    halted = true
                }
            case 1:
                if r >= '0' && r <= '9' {
                    state = 2
                } else {
                    halted = true
                }
            case 3:
                if r >= '0' && r <= '1' {
                    state = 5
                } else if r == '2' {
                    state = 5
                } else if r >= '3' && r <= '9' {
                    state = 5
                } else {
                    halted = true
                }
            case 11:
                if r >= '0' && r <= '1' {
                    state = 13
                } else if r == '2' {
                    state = 13
                } else if r >= '3' && r <= '9' {
                    state = 13
                } else {
                    halted = true
                }
            case 0:
                if r >= '0' && r <= '1' {
                    state = 1
                } else if r == '2' {
                    state = 1
                } else if r >= '3' && r <= '9' {
                    state = 1
                } else {
                    halted = true
                }
            }
        }
        if !halted && accepted[state] {
            acceptedIndex = index
        }
        index++
    }
    return []int{0, acceptedIndex}
}
