package automatons

// Contains checks if the first RuneSet contains all characters in the second RuneSet
func (set RuneSet) Contains(other *RuneSet) bool {
	for _, r := range other.Ranges {
		if r.Includes {
			if !set.IncludesRange(r.Start, r.End) {
				return false
			}
		}
	}
	return true
}

// Add adds all characters from another RuneSet to this RuneSet
func (set *RuneSet) Add(other *RuneSet) {
	for _, r := range other.Ranges {
		if r.Includes {
			set.AddRange(r.Start, r.End)
		}
	}
}

// Remove removes all characters from another RuneSet from this RuneSet
func (set *RuneSet) Remove(other *RuneSet) {
	for _, r := range other.Ranges {
		if r.Includes {
			set.RemoveRange(r.Start, r.End)
		}
	}
}

// Union creates a new RuneSet containing all characters from the given RuneSets
func Union(first *RuneSet, others ...*RuneSet) *RuneSet {
	// Create a copy of the first set
	result := &RuneSet{Ranges: make([]RuneRange, len(first.Ranges))}
	copy(result.Ranges, first.Ranges)

	for _, other := range others {
		result.Add(other)
	}
	return result
}

// Except creates a new RuneSet containing characters in the first set but not in the second
func Except(first *RuneSet, second *RuneSet) *RuneSet {
	// Create a copy of the first set
	result := &RuneSet{Ranges: make([]RuneRange, len(first.Ranges))}
	copy(result.Ranges, first.Ranges)

	result.Remove(second)
	return result
}

// Negate creates a new RuneSet containing all characters not in the given set
func Negate(set *RuneSet) *RuneSet {
	full := NewRuneSetFull()
	return Except(full, set)
}

// Intersect creates a new RuneSet containing only characters present in both sets
// Using De Morgan's law: A ∩ B = ¬(¬A ∪ ¬B)
func Intersect(first *RuneSet, second *RuneSet) *RuneSet {
	notFirst := Negate(first)
	notSecond := Negate(second)
	unionNegated := Union(notFirst, notSecond)
	return Negate(unionNegated)
}
