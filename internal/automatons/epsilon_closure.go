package automatons

func (nfa NFA) GetEpsilonClosure(states ...int) BitMask {
	closure := NewBitMaskEmpty(nfa.StateCount)
	queue := make([]int, 0, len(states))
	queue = append(queue, states...)
	for len(queue) > 0 {
		source := queue[0]
		queue = queue[1:]
		if closure.IsSet(source) {
			continue
		}
		closure.Set(source)
		targets := nfa.TransitionsBySource[source]
		if targets == nil {
			continue
		}
		epsilonTargets := targets.GetEpsilonValues()
		if epsilonTargets == nil {
			continue
		}
		for _, target := range *epsilonTargets {
			if !closure.IsSet(target) {
				queue = append(queue, target)
			}
		}
	}

	return closure
}
