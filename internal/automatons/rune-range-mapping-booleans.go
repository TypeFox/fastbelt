package automatons

import "iter"

type Booleans bool

func (t Booleans) NewEmpty() Booleans {
	return false
}

func (t Booleans) Join(other Booleans) Booleans {
	return t || other
}

func (t Booleans) Clone() Booleans {
	return t
}

func (t Booleans) Empty() bool {
	return bool(t) == false
}

type RuneRangeBooleansMapping struct {
	RuneRangeMappingBase[Booleans]
}

func NewRuneRangeBooleans() *RuneRangeBooleansMapping {
	return &RuneRangeBooleansMapping{
		RuneRangeMappingBase: RuneRangeMappingBase[Booleans]{
			Epsilon: false,
			Ranges: append(make([]RuneRangeMappingSection[Booleans], 0), RuneRangeMappingSection[Booleans]{
				Range:  NewRuneRange(MinRune, MaxRune, false),
				Values: false,
			}),
		},
	}
}

func (t *RuneRangeBooleansMapping) Contains(c rune) bool {
	return MappingBase_Contains[Booleans](&t.RuneRangeMappingBase, c)
}

func (t *RuneRangeBooleansMapping) ContainsEpsilon() bool {
	return MappingBase_ContainsEpsilon(&t.RuneRangeMappingBase)
}

func (t *RuneRangeBooleansMapping) GetEpsilonValues() *Booleans {
	return MappingBase_GetEpsilonValues[Booleans](&t.RuneRangeMappingBase)
}

func (t *RuneRangeBooleansMapping) GetRuneValues(c rune) *Booleans {
	return MappingBase_GetRuneValues[Booleans](&t.RuneRangeMappingBase, c)
}

func (t *RuneRangeBooleansMapping) All() iter.Seq[RuneRangeMappingSection[Booleans]] {
	return MappingBase_All(&t.RuneRangeMappingBase)
}

func (t *RuneRangeBooleansMapping) AddEpsilonValues(values Booleans) {
	MappingBase_AddEpsilonValues(&t.RuneRangeMappingBase, values)
}

func (t *RuneRangeBooleansMapping) AddRuneValues(r rune, targets Booleans) {
	MappingBase_AddRuneRangeValues(&t.RuneRangeMappingBase, r, r, targets)
}

func (t *RuneRangeBooleansMapping) AddRuneRangeValues(start, end rune, targets Booleans) {
	MappingBase_AddRuneRangeValues(&t.RuneRangeMappingBase, start, end, targets)
}

func (t *RuneRangeBooleansMapping) MergeNonEpsilonInto(target *RuneRangeMappingBase[Booleans]) {
	MappingBase_MergeNonEpsilonInto(&t.RuneRangeMappingBase, target)
}
