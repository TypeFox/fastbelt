package automatons

import (
	"iter"
	"sort"
)

type RuneTargetsSection struct {
	Range   *RuneRange
	Targets []int
}

type NFATargets interface {
	Contains(c rune) bool
	ContainsEpsilon() bool
	GetEpsilonTargets() []int
	GetRuneTargets(c rune) []int
	All() iter.Seq[RuneTargetsSection]
}

type NFAMutableTargets interface {
	NFATargets
	AddRuneTarget(r rune, targets ...int)
	AddRuneRangeTarget(rng *RuneRange, targets ...int)
	AddEpsilonTarget(targets ...int)
}

type NFAMutableTargets_ArrayImpl struct {
	Epsilon []int
	Ranges  []RuneTargetsSection
}

func NewNFAMutableTargets() *NFAMutableTargets_ArrayImpl {
	return &NFAMutableTargets_ArrayImpl{
		Epsilon: make([]int, 1),
		Ranges: append(make([]RuneTargetsSection, 0), RuneTargetsSection{
			Range:   NewRuneRange(MinRune, MaxRune, false),
			Targets: make([]int, 0),
		}),
	}
}

func (t *NFAMutableTargets_ArrayImpl) Contains(c rune) bool {
	startFromIndex := sort.Search(len(t.Ranges), func(i int) bool {
		return t.Ranges[i].Range.Start >= c
	})
	if startFromIndex > -1 {
		return t.Ranges[startFromIndex].Range.Contains(c) && t.Ranges[startFromIndex].Range.Includes
	}
	return false
}

func (t *NFAMutableTargets_ArrayImpl) ContainsEpsilon() bool {
	return len(t.Epsilon) > 0
}

func (t *NFAMutableTargets_ArrayImpl) GetEpsilonTargets() []int {
	return append([]int{}, t.Epsilon...)
}

func (t *NFAMutableTargets_ArrayImpl) GetRuneTargets(c rune) []int {
	startFromIndex := sort.Search(len(t.Ranges), func(i int) bool {
		return t.Ranges[i].Range.Start >= c
	})
	if startFromIndex > -1 {
		if t.Ranges[startFromIndex].Range.Contains(c) && t.Ranges[startFromIndex].Range.Includes {
			return append([]int{}, t.Ranges[startFromIndex].Targets...)
		}
	}
	return []int{}
}

func (t *NFAMutableTargets_ArrayImpl) All() iter.Seq[RuneTargetsSection] {
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

func (t *NFAMutableTargets_ArrayImpl) AddRuneTarget(r rune, targets ...int) {
	t.AddRuneRangeTarget(r, r, targets...)
}

func (t *NFAMutableTargets_ArrayImpl) AddRuneRangeTarget(start rune, end rune, targets ...int) {
	indexStart := sort.Search(len(t.Ranges), func(i int) bool {
		return t.Ranges[i].Range.Start >= start
	})
	indexEnd := sort.Search(len(t.Ranges), func(i int) bool {
		return t.Ranges[i].Range.End <= end
	})
	newRanges := make([]RuneTargetsSection, 0)
	if indexStart > -1 && indexStart-1 >= 0 {
		newRanges = append(newRanges, t.Ranges[:indexStart-1]...)
	}
	if indexEnd == -1 {
		indexEnd = len(t.Ranges) - 1
	}
	currentRange := NewRuneRange(start, end, true)
	for i := indexStart; i <= indexEnd; i++ {
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
			//   CCCCCCCC
			//NNNSSSSSSSS
			case currentRange.End == section.Range.End:
				newRanges = append(newRanges, RuneTargetsSection{
					Range:   NewRuneRange(currentRange.Start, currentRange.End, currentRange.Includes || section.Range.Includes),
					Targets: append(append([]int{}, section.Targets...), targets...),
				})
			//   CCCCCCCCCCC
			//NNNSSSSSSSSS
			case currentRange.End < section.Range.End:
				newRanges = append(newRanges, RuneTargetsSection{
					Range:   NewRuneRange(currentRange.Start, section.Range.End, currentRange.Includes || section.Range.Includes),
					Targets: append(append([]int{}, section.Targets...), targets...),
				}, RuneTargetsSection{
					Range:   NewRuneRange(rune(section.Range.End+1), currentRange.End, currentRange.Includes),
					Targets: append([]int{}, section.Targets...),
				})
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
				}, RuneTargetsSection{
					Range:   NewRuneRange(rune(section.Range.End+1), currentRange.End, currentRange.Includes),
					Targets: append([]int{}, targets...),
				})
			//CCCCCCCCCC
			//SSSSSSSSSS
			case currentRange.End == section.Range.End:
				newRanges = append(newRanges, RuneTargetsSection{
					Range:   NewRuneRange(currentRange.Start, currentRange.End, currentRange.Includes || section.Range.Includes),
					Targets: append(append([]int{}, section.Targets...), targets...),
				})
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
				})
			//NNNCCCCCCCC
			//   SSSSSSSS
			case currentRange.End == section.Range.End:
				newRanges = append(newRanges, RuneTargetsSection{
					Range:   NewRuneRange(section.Range.Start, currentRange.End, currentRange.Includes || section.Range.Includes),
					Targets: append(append([]int{}, section.Targets...), targets...),
				})
			//NNNCCCCCCCCCCCC
			//   SSSSSSSS
			case currentRange.End > section.Range.End:
				newRanges = append(newRanges, RuneTargetsSection{
					Range:   NewRuneRange(section.Range.Start, section.Range.End, currentRange.Includes || section.Range.Includes),
					Targets: append(append([]int{}, section.Targets...), targets...),
				}, RuneTargetsSection{
					Range:   NewRuneRange(rune(section.Range.End+1), currentRange.End, currentRange.Includes),
					Targets: append([]int{}, targets...),
				})
			}
		}
	}
	if indexEnd+1 <= len(t.Ranges) {
		newRanges = append(newRanges, t.Ranges[indexEnd+1:]...)
	}
	t.Ranges = newRanges
}
