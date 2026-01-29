package atn

type ATN interface {
	States() []*ATNState
	AddState(state ATNState)
}

type ATNData struct {
	states []*ATNState
}

func NewATN() ATN {
	return &ATNData{
		states: []*ATNState{},
	}
}

func (a *ATNData) States() []*ATNState {
	return a.states
}

func (a *ATNData) AddState(state ATNState) {
	a.states = append(a.states, &state)
}
