package automatons

type NFABuilder struct {
	stateCount   int
	startState   int
	acceptStates map[int]bool
}
