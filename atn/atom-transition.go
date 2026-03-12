package atn

type AtomTransition interface {
	Transition
	Atom() int
}

type AtomTransitionData struct {
	AbstractTransitionData
	atom int
}

func NewAtomTransitionData(target ATNState, atom int) *AtomTransitionData {
	return &AtomTransitionData{
		AbstractTransitionData: NewTransitionData(target),
		atom:                   atom,
	}
}

func (a *AtomTransitionData) Atom() int {
	return a.atom
}
