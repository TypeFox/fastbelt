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

func NewRuneSetBySingle(r rune) *RuneSet {
	left := RuneRange{Start: 0, End: r - 1, Includes: false}
	middle := RuneRange{Start: r, End: r, Includes: true}
	right := RuneRange{Start: r + 1, End: MaxRune, Includes: false}
	ranges := make([]RuneRange, 1)
	if left.End >= 0 {
		ranges = append(ranges, left)
	}
	ranges = append(ranges, middle)
	if right.Start <= MaxRune {
		ranges = append(ranges, right)
	}
	return &RuneSet{Ranges: ranges}
}
