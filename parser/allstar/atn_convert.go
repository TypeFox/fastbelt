// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package allstar

// BuildRuntimeATN converts the full build-time ATN to the minimal RuntimeATN.
// The build-time ATN (with grammar back-pointers and scaffolding fields) can be
// discarded after this call.
func BuildRuntimeATN(atn *ATN) *RuntimeATN {
	// 1. Allocate RuntimeATNState objects (no transitions yet so we can wire
	//    them by pointer in a second pass).
	rStates := make([]*RuntimeATNState, len(atn.States))
	stateMap := make(map[*ATNState]*RuntimeATNState, len(atn.States))
	for i, s := range atn.States {
		rs := &RuntimeATNState{
			StateNumber:            s.StateNumber,
			Type:                   s.Type,
			Decision:               s.Decision,
			EpsilonOnlyTransitions: s.EpsilonOnlyTransitions,
		}
		if s.Rule != nil {
			rs.RuleName = s.Rule.Name
		}
		if s.Production != nil {
			rs.ProdKind = s.Production.Kind()
			rs.ProdIdx = s.Production.Occurrence()
		}
		rStates[i] = rs
		stateMap[s] = rs
	}

	// 2. Wire transitions now that all states are in stateMap.
	for _, s := range atn.States {
		rs := stateMap[s]
		ts := make([]RuntimeTransition, len(s.Transitions))
		for i, t := range s.Transitions {
			switch at := t.(type) {
			case *AtomTransition:
				ts[i] = &RuntimeAtomTransition{
					Target:          stateMap[at.target],
					TokenTypeID:     at.TokenTypeID,
					CategoryMatches: at.CategoryMatches,
				}
			case *EpsilonTransition:
				ts[i] = &RuntimeEpsilonTransition{Target: stateMap[at.target]}
			case *RuleTransition:
				ts[i] = &RuntimeRuleTransition{
					Target:      stateMap[at.target],
					FollowState: stateMap[at.FollowState],
				}
			}
		}
		rs.Transitions = ts
	}

	// 3. Build ordered DecisionStates slice (indexed by Decision value).
	var decisionStates []*RuntimeATNState
	for _, rs := range rStates {
		if rs.Decision >= 0 {
			for len(decisionStates) <= rs.Decision {
				decisionStates = append(decisionStates, nil)
			}
			decisionStates[rs.Decision] = rs
		}
	}

	// 4. Mirror the DecisionMap.
	dm := make(map[string]*RuntimeATNState, len(atn.DecisionMap))
	for key, s := range atn.DecisionMap {
		dm[key] = stateMap[s]
	}

	return &RuntimeATN{
		States:         rStates,
		DecisionStates: decisionStates,
		DecisionMap:    dm,
	}
}
