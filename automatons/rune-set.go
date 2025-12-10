package automatons

const MaxRune = 0xFFFF

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

func NewRuneSet_Full() *RuneSet {
	return &RuneSet{Ranges: []RuneRange{
		{Start: 0, End: MaxRune, Includes: true},
	}}
}

func (set RuneSet) change(included bool, start rune, end rune) {
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

		var leftMostRange = set.Ranges[leftMostIndex]
		var rightMostRange = set.Ranges[rightMostIndex]
		var leftMost = leftMostRange.Start
		var rightMost = rightMostRange.End

		var leftList = set.Ranges[0:leftMostIndex]
		var rightList = set.Ranges[rightMostIndex+1:]

		ranges := make([]RuneRange, len(set.Ranges)+3)
		if leftMost < start {
			if rightMost > end {
				set.Ranges = append(ranges,
					leftList,
					NewRuneRange(leftMost, start-1, leftMostRange.Includes),
					NewRuneRange(start, end, included),
					NewRuneRange(end+1, rightMost, rightMostRange.Includes),
					rightList,
				)
			} else {
				set.Ranges = append(ranges,
					leftList,
					NewRuneRange(leftMost, start-1, leftMostRange.Includes),
					NewRuneRange(start, end, included),
					rightList,
				)
			}
		} else {
			if rightMost > end {
				set.Ranges = append(ranges,
					leftList,
					NewRuneRange(start, end, included),
					NewRuneRange(end+1, rightMost, rightMostRange.Includes),
					rightList,
				)
			} else {
				set.Ranges = append(ranges,
					leftList,
					NewRuneRange(start, end, included),
					rightList,
				)
			}
		}

		const leftIndex = math.max(0, leftMostIndex-1)
		const startRange = leftIndex
		const endRange = leftIndex + 4

		var index = startRange
		for index < endRange && index+1 < len(set.Ranges) {
			current = set.Ranges[index]
			next = set.Ranges[index+1]
			if current.mode == next.mode && current.to+1 == next.from {
				set.Ranges.splice(index, 2, NewRuneRange(current.from, next.to, current.mode))
			} else {
				index++
			}
		}
	}
}
