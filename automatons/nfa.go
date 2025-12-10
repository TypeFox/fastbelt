package automatons

type NFA struct {
	StartState     int
	StateCount     int
	AcceptedStates map[int]bool
	Transitions    map[int]TransitionMap
}

type TransitionRange struct {
	Range       *RuneRange // nil indicates epsilon transition
	TargetState int        //-1 indicates gap
}

type TransitionMap struct {
	Ranges []TransitionRange
}
