package automatons

import (
	"iter"
	"sort"
)

type RuneTargetsSection struct {
	Range   *RuneRange
	Targets []int
}

type NFATargets struct {
	Epsilon []int
	Ranges  []RuneTargetsSection
}

func NewNFATargets() *NFATargets {
	return &NFATargets{
		Epsilon: make([]int, 0),
		Ranges: append(make([]RuneTargetsSection, 0), RuneTargetsSection{
			Range:   NewRuneRange(MinRune, MaxRune, false),
			Targets: make([]int, 0),
		}),
	}
}

func (t *NFATargets) Contains(c rune) bool {
	startFromIndex := sort.Search(len(t.Ranges), func(i int) bool {
		return t.Ranges[i].Range.Start >= c
	}) - 1
	if startFromIndex > -1 {
		return t.Ranges[startFromIndex].Range.Contains(c) && t.Ranges[startFromIndex].Range.Includes
	}
	return false
}

func (t *NFATargets) ContainsEpsilon() bool {
	return len(t.Epsilon) > 0
}

func (t *NFATargets) GetEpsilonTargets() []int {
	return append([]int{}, t.Epsilon...)
}

func (t *NFATargets) GetRuneTargets(c rune) []int {
	startFromIndex := sort.Search(len(t.Ranges), func(i int) bool {
		return t.Ranges[i].Range.Start >= c
	}) - 1
	if startFromIndex > -1 {
		if t.Ranges[startFromIndex].Range.Contains(c) && t.Ranges[startFromIndex].Range.Includes {
			return append([]int{}, t.Ranges[startFromIndex].Targets...)
		}
	}
	return []int{}
}

func (t *NFATargets) All() iter.Seq[RuneTargetsSection] {
	return func(yield func(RuneTargetsSection) bool) {
		if t.ContainsEpsilon() {
			if !yield(RuneTargetsSection{
				Range:   nil,
				Targets: append([]int{}, t.Epsilon...),
			}) {
				return
			}
		}
		for _, section := range t.Ranges {
			if !yield(section) {
				return
			}
		}
	}
}

func (t *NFATargets) AddEpsilonTargets(targets ...int) {
	t.Epsilon = append(t.Epsilon, targets...)
}

func (t *NFATargets) AddRuneTargets(r rune, targets ...int) {
	t.AddRuneRangeTargets(r, r, targets...)
}

func (t *NFATargets) AddRuneRangeTargets(start rune, end rune, targets ...int) {
	indexStart := sort.Search(len(t.Ranges), func(i int) bool {
		return t.Ranges[i].Range.Start > start
	}) - 1
	indexEnd := sort.Search(len(t.Ranges), func(i int) bool {
		return t.Ranges[i].Range.End < end
	}) - 1
	newRanges := make([]RuneTargetsSection, 0)
	if indexStart > -1 && indexStart >= 0 {
		newRanges = append(newRanges, t.Ranges[:indexStart]...)
	}
	if indexEnd == -1 {
		indexEnd = len(t.Ranges) - 1
	}
	currentRange := NewRuneRange(start, end, true)
	for i := indexStart; i <= indexEnd && currentRange != nil; i++ {
		section := &t.Ranges[i]
		switch {
		//   CCCCCCC...
		//SSSSSSS...
		case currentRange.Start > section.Range.Start:
			newRanges = append(newRanges, RuneTargetsSection{
				Range:   NewRuneRange(section.Range.Start, rune(currentRange.Start-1), section.Range.Includes),
				Targets: append([]int{}, section.Targets...),
			})
			switch {
			//   CCCCC
			//NNNSSSSSSSS
			case currentRange.End < section.Range.End:
				newRanges = append(newRanges, RuneTargetsSection{
					Range:   NewRuneRange(currentRange.Start, currentRange.End, currentRange.Includes || section.Range.Includes),
					Targets: append(append([]int{}, section.Targets...), targets...),
				}, RuneTargetsSection{
					Range:   NewRuneRange(rune(currentRange.End+1), section.Range.End, section.Range.Includes),
					Targets: append([]int{}, section.Targets...),
				})
				currentRange = nil
			//   CCCCCCCC
			//NNNSSSSSSSS
			case currentRange.End == section.Range.End:
				newRanges = append(newRanges, RuneTargetsSection{
					Range:   NewRuneRange(currentRange.Start, currentRange.End, currentRange.Includes || section.Range.Includes),
					Targets: append(append([]int{}, section.Targets...), targets...),
				})
				currentRange = nil
			//   CCCCCCCCCCC
			//NNNSSSSSSSSS
			case currentRange.End < section.Range.End:
				newRanges = append(newRanges, RuneTargetsSection{
					Range:   NewRuneRange(currentRange.Start, section.Range.End, currentRange.Includes || section.Range.Includes),
					Targets: append(append([]int{}, section.Targets...), targets...),
				})
				currentRange = NewRuneRange(rune(section.Range.End+1), currentRange.End, currentRange.Includes)
			}
		//CCCCCC...
		//SSSSSS...
		case currentRange.Start == section.Range.Start:
			switch {
			//CCCCCCCC
			//SSSS
			case currentRange.End > section.Range.End:
				newRanges = append(newRanges, RuneTargetsSection{
					Range:   NewRuneRange(section.Range.Start, section.Range.End, currentRange.Includes || section.Range.Includes),
					Targets: append(append([]int{}, section.Targets...), targets...),
				})
				currentRange = NewRuneRange(rune(section.Range.End+1), currentRange.End, currentRange.Includes)
			//CCCCCCCCCC
			//SSSSSSSSSS
			case currentRange.End == section.Range.End:
				newRanges = append(newRanges, RuneTargetsSection{
					Range:   NewRuneRange(currentRange.Start, currentRange.End, currentRange.Includes || section.Range.Includes),
					Targets: append(append([]int{}, section.Targets...), targets...),
				})
				currentRange = nil
			//CCCCC
			//SSSSSSSSS
			case currentRange.End < section.Range.End:
				newRanges = append(newRanges, RuneTargetsSection{
					Range:   NewRuneRange(currentRange.Start, currentRange.End, currentRange.Includes || section.Range.Includes),
					Targets: append(append([]int{}, section.Targets...), targets...),
				}, RuneTargetsSection{
					Range:   NewRuneRange(rune(currentRange.End+1), section.Range.End, section.Range.Includes),
					Targets: append([]int{}, section.Targets...),
				})
				currentRange = nil
			}
		//CCCCCCCC...
		//  SSSSSS...
		case currentRange.Start < section.Range.Start:
			newRanges = append(newRanges, RuneTargetsSection{
				Range:   NewRuneRange(currentRange.Start, rune(section.Range.Start-1), currentRange.Includes),
				Targets: append([]int{}, targets...),
			})
			switch {
			//NNNCCCCC
			//   SSSSSSSSS
			case currentRange.End < section.Range.End:
				newRanges = append(newRanges, RuneTargetsSection{
					Range:   NewRuneRange(section.Range.Start, currentRange.End, currentRange.Includes || section.Range.Includes),
					Targets: append(append([]int{}, section.Targets...), targets...),
				}, RuneTargetsSection{
					Range:   NewRuneRange(rune(currentRange.End+1), section.Range.End, section.Range.Includes),
					Targets: append([]int{}, section.Targets...),
				})
				currentRange = nil
			//NNNCCCCCCCC
			//   SSSSSSSS
			case currentRange.End == section.Range.End:
				newRanges = append(newRanges, RuneTargetsSection{
					Range:   NewRuneRange(section.Range.Start, currentRange.End, currentRange.Includes || section.Range.Includes),
					Targets: append(append([]int{}, section.Targets...), targets...),
				})
				currentRange = nil
			//NNNCCCCCCCCCCCC
			//   SSSSSSSS
			case currentRange.End > section.Range.End:
				newRanges = append(newRanges, RuneTargetsSection{
					Range:   NewRuneRange(section.Range.Start, section.Range.End, currentRange.Includes || section.Range.Includes),
					Targets: append(append([]int{}, section.Targets...), targets...),
				})
				currentRange = NewRuneRange(rune(section.Range.End+1), currentRange.End, currentRange.Includes)
			}
		}
		if currentRange == nil {
			indexEnd = i
		}
	}
	if currentRange != nil {
		newRanges = append(newRanges, RuneTargetsSection{
			Range:   currentRange,
			Targets: append([]int{}, targets...),
		})
	}
	if indexEnd+1 < len(t.Ranges) {
		newRanges = append(newRanges, t.Ranges[indexEnd+1:]...)
	}
	t.Ranges = newRanges
}

func (t *NFATargets) MergeNonEpsilonInto(target *NFATargets) {
	for _, section := range t.Ranges {
		if section.Range.Includes && len(section.Targets) > 0 {
			target.AddRuneRangeTargets(section.Range.Start, section.Range.End, section.Targets...)
		}
	}
}
