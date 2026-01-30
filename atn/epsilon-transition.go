package atn

type EpsilonTransition interface {
	Transition
}

type EpsilonTransitionData struct {
	AbstractTransitionData
}

func NewEpsilonTransitionData(target ATNState) *EpsilonTransitionData {
	return &EpsilonTransitionData{
		AbstractTransitionData: NewTransitionData(target),
	}
}

func (e *EpsilonTransitionData) IsEpsilon() bool {
	return true
}
