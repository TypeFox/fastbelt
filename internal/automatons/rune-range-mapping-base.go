package automatons

import (
	"fmt"
	"iter"
	"sort"
)

type Value[T any] interface {
	Join(other T) T
	Clone() T
	Empty() bool
	String() string
}

type RuneRangeMapping[T Value[T]] interface {
	Contains(c rune) bool
	ContainsEpsilon() bool
	GetRuneValues(c rune) *T
	GetEpsilonValues() *T

	All() iter.Seq[RuneRangeMappingSection[T]]

	AddEpsilonValues(values T)
	AddRuneValues(r rune, targets T)
	AddRuneRangeValues(start, end rune, targets T)

	MergeNonEpsilonInto(target *RuneRangeMappingBase[T])
}

type RuneRangeMappingSection[T Value[T]] struct {
	Range  *RuneRange
	Values T
}

type RuneRangeMappingBase[T Value[T]] struct {
	Epsilon T
	Ranges  []RuneRangeMappingSection[T]
}

func MappingBaseSection_String[T Value[T]](s RuneRangeMappingSection[T]) string {
	if s.Range == nil {
		return "ε -> " + fmt.Sprintf("%v", s.Values)
	}
	return fmt.Sprintf("%v -> %v", s.Range, s.Values)
}

func MappingBase_Contains[T Value[T]](t *RuneRangeMappingBase[T], c rune) bool {
	startFromIndex := sort.Search(len(t.Ranges), func(i int) bool {
		return t.Ranges[i].Range.Start > c
	}) - 1
	if startFromIndex > -1 {
		return t.Ranges[startFromIndex].Range.Contains(c) && t.Ranges[startFromIndex].Range.Includes
	}
	return false
}

func MappingBase_ContainsEpsilon[T Value[T]](t *RuneRangeMappingBase[T]) bool {
	return !t.Epsilon.Empty()
}

func MappingBase_GetEpsilonValues[T Value[T]](t *RuneRangeMappingBase[T]) *T {
	if t.Epsilon.Empty() {
		return nil
	}
	return &t.Epsilon
}

func MappingBase_GetRuneValues[T Value[T]](t *RuneRangeMappingBase[T], c rune) *T {
	startFromIndex := sort.Search(len(t.Ranges), func(i int) bool {
		return t.Ranges[i].Range.Start > c
	}) - 1
	if startFromIndex > -1 {
		if t.Ranges[startFromIndex].Range.Contains(c) && t.Ranges[startFromIndex].Range.Includes {
			return &t.Ranges[startFromIndex].Values
		}
	}
	return nil
}

func MappingBase_All[T Value[T]](t *RuneRangeMappingBase[T]) iter.Seq[RuneRangeMappingSection[T]] {
	return func(yield func(RuneRangeMappingSection[T]) bool) {
		if !t.Epsilon.Empty() {
			if !yield(RuneRangeMappingSection[T]{
				Range:  nil,
				Values: t.Epsilon,
			}) {
				return
			}
		}
		for _, section := range t.Ranges {
			if section.Range.Includes {
				if !yield(section) {
					return
				}
			}
		}
	}
}

func MappingBase_AddEpsilonValues[T Value[T]](t *RuneRangeMappingBase[T], values T) {
	t.Epsilon = t.Epsilon.Join(values)
}

func MappingBase_AddRuneRangeValues[T Value[T]](t *RuneRangeMappingBase[T], start, end rune, values T) {
	indexStart := sort.Search(len(t.Ranges), func(i int) bool {
		return t.Ranges[i].Range.Start > start
	}) - 1
	indexEnd := sort.Search(len(t.Ranges), func(i int) bool {
		return t.Ranges[i].Range.End < end
	}) - 1
	newRanges := make([]RuneRangeMappingSection[T], 0)
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
			newRanges = append(newRanges, RuneRangeMappingSection[T]{
				Range:  NewRuneRange(section.Range.Start, rune(currentRange.Start-1), section.Range.Includes),
				Values: section.Values.Clone(),
			})
			switch {
			//   CCCCC
			//NNNSSSSSSSS
			case currentRange.End < section.Range.End:
				newRanges = append(newRanges, RuneRangeMappingSection[T]{
					Range:  NewRuneRange(currentRange.Start, currentRange.End, currentRange.Includes || section.Range.Includes),
					Values: section.Values.Join(values),
				}, RuneRangeMappingSection[T]{
					Range:  NewRuneRange(rune(currentRange.End+1), section.Range.End, section.Range.Includes),
					Values: section.Values.Clone(),
				})
				currentRange = nil
			//   CCCCCCCC
			//NNNSSSSSSSS
			case currentRange.End == section.Range.End:
				newRanges = append(newRanges, RuneRangeMappingSection[T]{
					Range:  NewRuneRange(currentRange.Start, currentRange.End, currentRange.Includes || section.Range.Includes),
					Values: section.Values.Join(values),
				})
				currentRange = nil
			//   CCCCCCCCCCC
			//NNNSSSSSSSSS
			case currentRange.End < section.Range.End:
				newRanges = append(newRanges, RuneRangeMappingSection[T]{
					Range:  NewRuneRange(currentRange.Start, section.Range.End, currentRange.Includes || section.Range.Includes),
					Values: section.Values.Join(values),
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
				newRanges = append(newRanges, RuneRangeMappingSection[T]{
					Range:  NewRuneRange(section.Range.Start, section.Range.End, currentRange.Includes || section.Range.Includes),
					Values: section.Values.Join(values),
				})
				currentRange = NewRuneRange(rune(section.Range.End+1), currentRange.End, currentRange.Includes)
			//CCCCCCCCCC
			//SSSSSSSSSS
			case currentRange.End == section.Range.End:
				newRanges = append(newRanges, RuneRangeMappingSection[T]{
					Range:  NewRuneRange(currentRange.Start, currentRange.End, currentRange.Includes || section.Range.Includes),
					Values: section.Values.Join(values),
				})
				currentRange = nil
			//CCCCC
			//SSSSSSSSS
			case currentRange.End < section.Range.End:
				newRanges = append(newRanges, RuneRangeMappingSection[T]{
					Range:  NewRuneRange(currentRange.Start, currentRange.End, currentRange.Includes || section.Range.Includes),
					Values: section.Values.Join(values),
				}, RuneRangeMappingSection[T]{
					Range:  NewRuneRange(rune(currentRange.End+1), section.Range.End, section.Range.Includes),
					Values: section.Values.Clone(),
				})
				currentRange = nil
			}
		//CCCCCCCC...
		//  SSSSSS...
		case currentRange.Start < section.Range.Start:
			newRanges = append(newRanges, RuneRangeMappingSection[T]{
				Range:  NewRuneRange(currentRange.Start, rune(section.Range.Start-1), currentRange.Includes),
				Values: values.Clone(),
			})
			switch {
			//NNNCCCCC
			//   SSSSSSSSS
			case currentRange.End < section.Range.End:
				newRanges = append(newRanges, RuneRangeMappingSection[T]{
					Range:  NewRuneRange(section.Range.Start, currentRange.End, currentRange.Includes || section.Range.Includes),
					Values: section.Values.Join(values),
				}, RuneRangeMappingSection[T]{
					Range:  NewRuneRange(rune(currentRange.End+1), section.Range.End, section.Range.Includes),
					Values: section.Values.Clone(),
				})
				currentRange = nil
			//NNNCCCCCCCC
			//   SSSSSSSS
			case currentRange.End == section.Range.End:
				newRanges = append(newRanges, RuneRangeMappingSection[T]{
					Range:  NewRuneRange(section.Range.Start, currentRange.End, currentRange.Includes || section.Range.Includes),
					Values: section.Values.Join(values),
				})
				currentRange = nil
			//NNNCCCCCCCCCCCC
			//   SSSSSSSS
			case currentRange.End > section.Range.End:
				newRanges = append(newRanges, RuneRangeMappingSection[T]{
					Range:  NewRuneRange(section.Range.Start, section.Range.End, currentRange.Includes || section.Range.Includes),
					Values: section.Values.Join(values),
				})
				currentRange = NewRuneRange(rune(section.Range.End+1), currentRange.End, currentRange.Includes)
			}
		}
		if currentRange == nil {
			indexEnd = i
		}
	}
	if currentRange != nil {
		newRanges = append(newRanges, RuneRangeMappingSection[T]{
			Range:  currentRange,
			Values: values.Clone(),
		})
	}
	if indexEnd+1 < len(t.Ranges) {
		newRanges = append(newRanges, t.Ranges[indexEnd+1:]...)
	}
	t.Ranges = newRanges
}

func MappingBase_MergeNonEpsilonInto[T Value[T]](t *RuneRangeMappingBase[T], target *RuneRangeMappingBase[T]) {
	for _, section := range t.Ranges {
		if section.Range.Includes && !section.Values.Empty() {
			MappingBase_AddRuneRangeValues(target, section.Range.Start, section.Range.End, section.Values)
		}
	}
}
