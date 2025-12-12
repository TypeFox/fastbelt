package automatons

import (
	"fmt"
)

const MaxRune = 0x10FFFF
const MinAscii = 0x20
const MaxAscii = 0x7E

func FormatRune(r rune) string {
	runeStr := ""
	if r == '\'' {
		runeStr = "'\\'"
	} else if r == '\\' {
		runeStr = "'\\\\'"
	} else if r >= MinAscii && r <= MaxAscii {
		runeStr = "'" + string(r) + "'"
	} else {
		runeStr = fmt.Sprintf("'\\u%04X'", int64(r))
	}
	return runeStr
}

// Helper function for max of integers
func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

type RuneSet struct {
	Ranges []RuneRange
}

func (set RuneSet) Length() int {
	var length int = 0
	for _, r := range set.Ranges {
		if r.Includes {
			length += int(r.End) - int(r.Start) + 1
		}
	}
	return length
}

func NewRuneSet_Rune(r rune) *RuneSet {
	middle := RuneRange{Start: r, End: r, Includes: true}
	ranges := make([]RuneRange, 1)
	if r-1 >= 0 {
		left := RuneRange{Start: 0, End: r - 1, Includes: false}
		ranges = append(ranges, left)
	}
	ranges = append(ranges, middle)
	if r+1 <= MaxRune {
		right := RuneRange{Start: r + 1, End: MaxRune, Includes: false}
		ranges = append(ranges, right)
	}
	return &RuneSet{Ranges: ranges}
}

func NewRuneSet_Range(start, end rune) *RuneSet {
	middle := RuneRange{Start: start, End: end, Includes: true}
	ranges := make([]RuneRange, 1)
	if start-1 >= 0 {
		left := RuneRange{Start: 0, End: start - 1, Includes: false}
		ranges = append(ranges, left)
	}
	ranges = append(ranges, middle)
	if end+1 <= MaxRune {
		right := RuneRange{Start: end + 1, End: MaxRune, Includes: false}
		ranges = append(ranges, right)
	}
	return &RuneSet{Ranges: ranges}
}

func NewRuneSet_Empty() *RuneSet {
	return &RuneSet{Ranges: []RuneRange{
		{Start: 0, End: MaxRune, Includes: false},
	}}
}

func NewRuneSet_OneOf(chars []rune) *RuneSet {
	runeSet := NewRuneSet_Empty()
	for _, c := range chars {
		runeSet.AddRune(c)
	}
	return runeSet
}

func NewRuneSet_Full() *RuneSet {
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
	return set.isXXcluded(true, r, r)
}

func (set RuneSet) IncludesRange(start rune, end rune) bool {
	return set.isXXcluded(true, start, end)
}

func (set RuneSet) ExcludesRune(r rune) bool {
	return set.isXXcluded(false, r, r)
}

func (set RuneSet) ExcludesRange(start rune, end rune) bool {
	return set.isXXcluded(false, start, end)
}

func (set RuneSet) isXXcluded(included bool, start rune, end rune) bool {
	var index int = 0
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

	var leftMostIndex int = 0
	for leftMostIndex < len(set.Ranges) && start > set.Ranges[leftMostIndex].End {
		leftMostIndex++
	}

	var rightMostIndex int = len(set.Ranges) - 1
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
