package atn

type Transition interface {
	Target() ATNState
	SetTarget(target ATNState)
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
func (tr *AbstractTransitionData) Target() ATNState {
	return tr.target
}

func (tr *AbstractTransitionData) IsEpsilon() bool {
	return false
}

func (tr *AbstractTransitionData) SetTarget(target ATNState) {
	tr.target = target
}
