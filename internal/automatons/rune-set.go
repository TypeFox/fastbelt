package automatons

import (
	"fmt"
)

const MinRune = 0x0000
const MaxRune = 0x10FFFF
const MinAscii = 0x20
const MaxAscii = 0x7E

func FormatInt(r rune) string {
	intStr := ""
	if r >= 0 && r <= 0xff {
		intStr = fmt.Sprintf("0x%02X", r)
	} else {
		intStr = fmt.Sprintf("0x%04X", r)
	}
	return intStr
}

func FormatLowHighInts(low rune, high rune) string {
	return fmt.Sprintf("0x%08X%08X", high, low)
}

func FormatRune(r rune) string {
	runeStr := ""
	if r == '\'' {
		runeStr = "'\\''"
	} else if r == '\\' {
		runeStr = "'\\\\'"
	} else if r >= MinAscii && r <= MaxAscii {
		runeStr = "'" + string(r) + "'"
	} else if r <= 0xffff {
		runeStr = fmt.Sprintf("'\\u%04X'", int64(r))
	} else {
		runeStr = fmt.Sprintf("'\\U%08X'", int64(r))
	}
	return runeStr
}

type RuneSet struct {
	Ranges []RuneRange
}

func (set RuneSet) Length() int {
	var length = 0
	for _, r := range set.Ranges {
		if r.Includes {
			length += int(r.End) - int(r.Start) + 1
		}
	}
	return length
}

func NewRuneSetRune(r rune) *RuneSet {
	ranges := make([]RuneRange, 0)
	if r-1 >= 0 {
		left := RuneRange{Start: 0, End: r - 1, Includes: false}
		ranges = append(ranges, left)
	}
	middle := RuneRange{Start: r, End: r, Includes: true}
	ranges = append(ranges, middle)
	if r+1 <= MaxRune {
		right := RuneRange{Start: r + 1, End: MaxRune, Includes: false}
		ranges = append(ranges, right)
	}
	return &RuneSet{Ranges: ranges}
}

func NewRuneSetRange(start, end rune) *RuneSet {
	if start > end || start < 0 || end > MaxRune {
		panic("Range limit order is invalid!")
	}
	ranges := make([]RuneRange, 0)
	if start-1 >= 0 {
		left := RuneRange{Start: 0, End: start - 1, Includes: false}
		ranges = append(ranges, left)
	}
	middle := RuneRange{Start: start, End: end, Includes: true}
	ranges = append(ranges, middle)
	if end+1 <= MaxRune {
		right := RuneRange{Start: end + 1, End: MaxRune, Includes: false}
		ranges = append(ranges, right)
	}
	return &RuneSet{Ranges: ranges}
}

func NewRuneSetEmpty() *RuneSet {
	return &RuneSet{Ranges: []RuneRange{
		{Start: 0, End: MaxRune, Includes: false},
	}}
}

func NewRuneSetOneOf(chars []rune) *RuneSet {
	runeSet := NewRuneSetEmpty()
	for _, c := range chars {
		runeSet.AddRune(c)
	}
	return runeSet
}

func NewRuneSetFull() *RuneSet {
	return &RuneSet{Ranges: []RuneRange{
		{Start: 0, End: MaxRune, Includes: true},
	}}
}

func (set *RuneSet) AddRange(start rune, end rune) {
	set.change(true, start, end)
}

func (set *RuneSet) AddRune(r rune) {
	set.change(true, r, r)
}

func (set *RuneSet) RemoveRange(start rune, end rune) {
	set.change(false, start, end)
}

func (set *RuneSet) RemoveRune(r rune) {
	set.change(false, r, r)
}

func (set RuneSet) IncludesRune(r rune) bool {
	return set.isInOrExcluded(true, r, r)
}

func (set RuneSet) IncludesRange(start rune, end rune) bool {
	return set.isInOrExcluded(true, start, end)
}

func (set RuneSet) ExcludesRune(r rune) bool {
	return set.isInOrExcluded(false, r, r)
}

func (set RuneSet) ExcludesRange(start rune, end rune) bool {
	return set.isInOrExcluded(false, start, end)
}

func (set RuneSet) isInOrExcluded(included bool, start rune, end rune) bool {
	var index = 0
	for index < len(set.Ranges) && start > set.Ranges[index].End {
		index++
	}
	if index >= len(set.Ranges) {
		return false
	}
	var r = set.Ranges[index]
	return start >= r.Start && end <= r.End && r.Includes == included
}

func (set *RuneSet) change(included bool, start rune, end rune) {
	if start > end || start < 0 || end > MaxRune {
		panic("Range limit order is invalid!")
	}

	var leftMostIndex = 0
	for leftMostIndex < len(set.Ranges) && start > set.Ranges[leftMostIndex].End {
		leftMostIndex++
	}

	var rightMostIndex = len(set.Ranges) - 1
	for rightMostIndex >= 0 && end < set.Ranges[rightMostIndex].Start {
		rightMostIndex--
	}

	var leftMostRange = set.Ranges[leftMostIndex]
	var rightMostRange = set.Ranges[rightMostIndex]
	var leftMost = leftMostRange.Start
	var rightMost = rightMostRange.End

	var leftList = make([]RuneRange, leftMostIndex)
	copy(leftList, set.Ranges[0:leftMostIndex])
	var rightList = make([]RuneRange, len(set.Ranges)-(rightMostIndex+1))
	copy(rightList, set.Ranges[rightMostIndex+1:])

	var newRanges []RuneRange

	if leftMost < start {
		if rightMost > end {
			newRanges = append(newRanges, leftList...)
			newRanges = append(newRanges, *NewRuneRange(leftMost, start-1, leftMostRange.Includes))
			newRanges = append(newRanges, *NewRuneRange(start, end, included))
			newRanges = append(newRanges, *NewRuneRange(end+1, rightMost, rightMostRange.Includes))
			newRanges = append(newRanges, rightList...)
		} else {
			newRanges = append(newRanges, leftList...)
			newRanges = append(newRanges, *NewRuneRange(leftMost, start-1, leftMostRange.Includes))
			newRanges = append(newRanges, *NewRuneRange(start, end, included))
			newRanges = append(newRanges, rightList...)
		}
	} else {
		if rightMost > end {
			newRanges = append(newRanges, leftList...)
			newRanges = append(newRanges, *NewRuneRange(start, end, included))
			newRanges = append(newRanges, *NewRuneRange(end+1, rightMost, rightMostRange.Includes))
			newRanges = append(newRanges, rightList...)
		} else {
			newRanges = append(newRanges, leftList...)
			newRanges = append(newRanges, *NewRuneRange(start, end, included))
			newRanges = append(newRanges, rightList...)
		}
	}

	set.Ranges = newRanges

	leftIndex := maxInt(0, leftMostIndex-1)
	set.tryMergeRange(leftIndex, leftIndex+4)
}

func (set *RuneSet) tryMergeRange(fromRange int, toRange int) {
	var index = fromRange
	for index < toRange && index+1 < len(set.Ranges) {
		current := set.Ranges[index]
		next := set.Ranges[index+1]
		if current.Includes == next.Includes && current.End+1 == next.Start {
			// Merge the two ranges
			merged := *NewRuneRange(current.Start, next.End, current.Includes)
			// Remove the two ranges and insert the merged one
			set.Ranges = append(set.Ranges[:index], append([]RuneRange{merged}, set.Ranges[index+2:]...)...)
		} else {
			index++
		}
	}
}

// Equals checks if this RuneSet is equal to another RuneSet
func (set RuneSet) Equals(other RuneSet) bool {
	if len(set.Ranges) != len(other.Ranges) {
		return false
	}
	for i := range set.Ranges {
		if !set.Ranges[i].Equals(other.Ranges[i]) {
			return false
		}
	}
	return true
}

// String returns a string representation of the RuneSet showing only included ranges
func (set RuneSet) String() string {
	if len(set.Ranges) == 0 {
		return ""
	}

	var result string
	first := true
	for _, r := range set.Ranges {
		if r.Includes {
			if !first {
				result += ","
			}
			first = false
			if r.Start == r.End {
				result += "[" + string(rune(r.Start)) + "]"
			} else {
				result += "[" + string(rune(r.Start)) + "-" + string(rune(r.End)) + "]"
			}
		}
	}
	return result
}
