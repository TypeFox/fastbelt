package automatons

import "iter"

type Targets []int

func (t Targets) NewEmpty() Targets {
	return Targets{}
}

func (t Targets) Join(other Targets) Targets {
	result := make([]int, 0)
	result = append(result, t...)
	result = append(result, other...)
	return Targets(result)
}

func (t Targets) Clone() Targets {
	result := make([]int, 0)
	result = append(result, t...)
	return Targets(result)
}

func (t Targets) Empty() bool {
	return len(t) == 0
}

type RuneRangeTargetsMapping struct {
	RuneRangeMappingBase[Targets]
}

func NewRuneRangeTargets() *RuneRangeTargetsMapping {
	return &RuneRangeTargetsMapping{
		RuneRangeMappingBase: RuneRangeMappingBase[Targets]{
			Epsilon: make([]int, 0),
			Ranges: append(make([]RuneRangeMappingSection[Targets], 0), RuneRangeMappingSection[Targets]{
				Range:  NewRuneRange(MinRune, MaxRune, false),
				Values: make([]int, 0),
			}),
		},
	}
}

func (t *RuneRangeTargetsMapping) Contains(c rune) bool {
	return MappingBase_Contains[Targets](&t.RuneRangeMappingBase, c)
}

func (t *RuneRangeTargetsMapping) ContainsEpsilon() bool {
	return MappingBase_ContainsEpsilon(&t.RuneRangeMappingBase)
}

func (t *RuneRangeTargetsMapping) GetEpsilonValues() *Targets {
	return MappingBase_GetEpsilonValues[Targets](&t.RuneRangeMappingBase)
}

func (t *RuneRangeTargetsMapping) GetRuneValues(c rune) *Targets {
	return MappingBase_GetRuneValues[Targets](&t.RuneRangeMappingBase, c)
}

func (t *RuneRangeTargetsMapping) All() iter.Seq[RuneRangeMappingSection[Targets]] {
	return MappingBase_All(&t.RuneRangeMappingBase)
}

func (t *RuneRangeTargetsMapping) AddEpsilonValues(values Targets) {
	MappingBase_AddEpsilonValues(&t.RuneRangeMappingBase, values)
}

func (t *RuneRangeTargetsMapping) AddRuneValues(r rune, targets Targets) {
	MappingBase_AddRuneRangeValues(&t.RuneRangeMappingBase, r, r, targets)
}

func (t *RuneRangeTargetsMapping) AddRuneRangeValues(start, end rune, targets Targets) {
	MappingBase_AddRuneRangeValues(&t.RuneRangeMappingBase, start, end, targets)
}

func (t *RuneRangeTargetsMapping) MergeNonEpsilonInto(target *RuneRangeMappingBase[Targets]) {
	MappingBase_MergeNonEpsilonInto(&t.RuneRangeMappingBase, target)
}
