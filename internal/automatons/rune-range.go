package automatons

import "iter"

type RuneRange struct {
	Start    rune
	End      rune
	Includes bool
}

func NewRuneRange(start, end rune, includes bool) *RuneRange {
	if start > end {
		panic("Start must be less or equal to End")
	}
	return &RuneRange{
		Start:    start,
		End:      end,
		Includes: includes,
	}
}

func (r RuneRange) Contains(ch rune) bool {
	return ch >= r.Start && ch <= r.End
}

func (r RuneRange) Equals(other RuneRange) bool {
	return r.Start == other.Start && r.End == other.End && r.Includes == other.Includes
}

func (r RuneRange) String() string {
	if r.Includes {
		return "[" + string(r.Start) + "-" + string(r.End) + "]"
	}
	return "[^" + string(r.Start) + "-" + string(r.End) + "]"
}

func (r RuneRange) Length() int {
	return int(r.End - r.Start + 1)
}

func (lst RuneRange) All() iter.Seq[rune] {
	return func(yield func(rune) bool) {
		for e := lst.Start; e <= lst.End; e++ {
			if !yield(e) {
				return
			}
		}
	}
}
