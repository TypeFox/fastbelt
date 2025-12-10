package automatons

type RuneRange struct {
	Start    rune
	End      rune
	Includes bool
}

func NewRuneRange(start, end rune, includes bool) *RuneRange {
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
