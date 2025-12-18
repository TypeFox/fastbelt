package automatons

import (
	"fmt"
	"iter"
)

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

func (t Targets) String() string {
	return fmt.Sprintf("%v", []int(t))
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
	return MappingBaseContains[Targets](&t.RuneRangeMappingBase, c)
}

func (t *RuneRangeTargetsMapping) ContainsEpsilon() bool {
	return MappingBaseContainsEpsilon(&t.RuneRangeMappingBase)
}

func (t *RuneRangeTargetsMapping) GetEpsilonValues() *Targets {
	return MappingBaseGetEpsilonValues[Targets](&t.RuneRangeMappingBase)
}

func (t *RuneRangeTargetsMapping) GetRuneValues(c rune) *Targets {
	return MappingBaseGetRuneValues[Targets](&t.RuneRangeMappingBase, c)
}

func (t *RuneRangeTargetsMapping) All() iter.Seq[RuneRangeMappingSection[Targets]] {
	return MappingBaseAll(&t.RuneRangeMappingBase)
}

func (t *RuneRangeTargetsMapping) AddEpsilonValues(values Targets) {
	MappingBaseAddEpsilonValues(&t.RuneRangeMappingBase, values)
}

func (t *RuneRangeTargetsMapping) AddRuneValues(r rune, targets Targets) {
	MappingBaseAddRuneRangeValues(&t.RuneRangeMappingBase, r, r, targets)
}

func (t *RuneRangeTargetsMapping) AddRuneRangeValues(start, end rune, targets Targets) {
	MappingBaseAddRuneRangeValues(&t.RuneRangeMappingBase, start, end, targets)
}

func (t *RuneRangeTargetsMapping) MergeNonEpsilonInto(target *RuneRangeMappingBase[Targets]) {
	MappingBaseMergeNonEpsilonInto(&t.RuneRangeMappingBase, target)
}
