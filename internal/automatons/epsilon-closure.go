package automatons

func GetEpsilonClosure(nfa NFA, states ...int) BitMask {
	closure := NewBitMask_Empty(nfa.GetStateCount())
	queue := make([]int, 0, len(states))
	queue = append(queue, states...)
	for len(queue) > 0 {
		source := queue[0]
		queue = queue[1:]
		if closure.IsSet(source) {
			continue
		}
		closure.Set(source)
		targets := nfa.GetTransitionsBySource()[source]
		if targets == nil {
			continue
		}
		for _, target := range targets.GetEpsilonTargets() {
			if !closure.IsSet(target) {
				queue = append(queue, target)
			}
		}
	}

	return closure
}
