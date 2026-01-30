package atn

type Transition interface {
	Target() ATNState
	IsEpsilon() bool
}

type AbstractTransitionData struct {
	target ATNState
}

func NewTransitionData(target ATNState) AbstractTransitionData {
	return AbstractTransitionData{
		target: target,
	}
}
func (t *AbstractTransitionData) Target() ATNState {
	return t.target
}

func (t *AbstractTransitionData) IsEpsilon() bool {
	return false
}
