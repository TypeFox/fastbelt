// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package parser

import (
	"fmt"

	"typefox.dev/fastbelt"
)

// ATNTables describes a [RuntimeATN] as plain data. The generator emits the
// ATN of a grammar in this form, and [NewRuntimeATNFromTables] builds the
// runtime structure from it. Tables of numbers and strings are static data
// for the Go compiler, so even the ATN of a large grammar compiles quickly,
// unlike code that builds it statement by statement.
type ATNTables struct {
	// States holds one entry per ATN state, indexed by state number.
	States []ATNStateEntry
	// Transitions holds the transitions of all states in state order: first
	// the transitions of state 0, then those of state 1, and so on.
	// [ATNStateEntry.Transitions] is the number of transitions of a state.
	Transitions []ATNTransitionEntry
	// DecisionStates holds the number of the decision state of each decision,
	// indexed by decision number, or -1 for a decision without a state.
	DecisionStates []int32
	// TokenTypes holds the token types of the grammar, indexed by
	// [fastbelt.TokenType.Id]. Entries that no atom transition uses may be nil.
	TokenTypes []*fastbelt.TokenType
	// CompletionHints and RuleCalls hold the completion hints and rule call
	// information that transitions refer to.
	CompletionHints []CompletionHint
	RuleCalls       []RuleCallInfo
}

// ATNStateEntry describes a state of an [ATNTables].
type ATNStateEntry struct {
	Type ATNStateType
	// EpsilonOnly reports whether all transitions of the state are epsilon
	// transitions.
	EpsilonOnly bool
	// Decision is the decision number of a decision state, or -1.
	Decision int32
	// Transitions is the number of transitions of the state.
	Transitions int32
}

// ATNTransitionKind is the kind of an [ATNTransitionEntry].
type ATNTransitionKind uint8

const (
	ATNEpsilonTransition ATNTransitionKind = iota
	ATNAtomTransition
	ATNRuleTransition
)

// ATNTransitionEntry describes a transition of an [ATNTables].
type ATNTransitionEntry struct {
	Kind ATNTransitionKind
	// Target is the number of the target state. For a rule transition, it is
	// the start state of the called rule.
	Target int32
	// Token is the token type of an atom transition, an index into
	// [ATNTables.TokenTypes].
	Token int32
	// Follow is the number of the state that a rule transition returns to.
	Follow int32
	// Hint is the completion hint of an atom or rule transition, as an index
	// into [ATNTables.CompletionHints] plus one, or 0 for none.
	Hint int32
	// Call is the rule call information of a rule transition, as an index
	// into [ATNTables.RuleCalls] plus one, or 0 for none.
	Call int32
}

// NewRuntimeATNFromTables builds the [RuntimeATN] that t describes. It panics
// if t is inconsistent, which indicates a bug in the generator.
func NewRuntimeATNFromTables(t ATNTables) *RuntimeATN {
	states := make([]RuntimeATNState, len(t.States))
	statePtrs := make([]*RuntimeATNState, len(t.States))
	for i := range states {
		states[i] = RuntimeATNState{
			StateNumber:            i,
			Type:                   t.States[i].Type,
			Decision:               int(t.States[i].Decision),
			EpsilonOnlyTransitions: t.States[i].EpsilonOnly,
		}
		statePtrs[i] = &states[i]
	}

	var atoms, epsilons, rules int
	for _, tr := range t.Transitions {
		switch tr.Kind {
		case ATNAtomTransition:
			atoms++
		case ATNEpsilonTransition:
			epsilons++
		case ATNRuleTransition:
			rules++
		}
	}
	atomTransitions := make([]RuntimeAtomTransition, 0, atoms)
	epsilonTransitions := make([]RuntimeEpsilonTransition, 0, epsilons)
	ruleTransitions := make([]RuntimeRuleTransition, 0, rules)
	transitions := make([]RuntimeTransition, len(t.Transitions))

	// Copies, so that the ATN does not point into the (generated) tables.
	hints := append([]CompletionHint(nil), t.CompletionHints...)
	calls := append([]RuleCallInfo(nil), t.RuleCalls...)
	hint := func(i int32) *CompletionHint {
		if i == 0 {
			return nil
		}
		return &hints[i-1]
	}
	call := func(i int32) *RuleCallInfo {
		if i == 0 {
			return nil
		}
		return &calls[i-1]
	}

	for i, tr := range t.Transitions {
		target := statePtrs[tr.Target]
		switch tr.Kind {
		case ATNAtomTransition:
			atomTransitions = append(atomTransitions, RuntimeAtomTransition{
				Target:         target,
				TokenType:      t.TokenTypes[tr.Token],
				CompletionHint: hint(tr.Hint),
			})
			transitions[i] = &atomTransitions[len(atomTransitions)-1]
		case ATNEpsilonTransition:
			epsilonTransitions = append(epsilonTransitions, RuntimeEpsilonTransition{Target: target})
			transitions[i] = &epsilonTransitions[len(epsilonTransitions)-1]
		case ATNRuleTransition:
			ruleTransitions = append(ruleTransitions, RuntimeRuleTransition{
				Target:         target,
				FollowState:    statePtrs[tr.Follow],
				CompletionHint: hint(tr.Hint),
				Call:           call(tr.Call),
			})
			transitions[i] = &ruleTransitions[len(ruleTransitions)-1]
		default:
			panic(fmt.Sprintf("ATN tables: transition %d has unknown kind %d", i, tr.Kind))
		}
	}

	total := 0
	for _, s := range t.States {
		total += int(s.Transitions)
	}
	if total != len(t.Transitions) {
		panic(fmt.Sprintf("ATN tables: the states have %d transitions, the transition table %d", total, len(t.Transitions)))
	}
	next := 0
	for i := range states {
		n := int(t.States[i].Transitions)
		if n > 0 {
			// The capacity limit keeps a later AppendTransitions from
			// overwriting the transitions of the next state.
			states[i].Transitions = transitions[next : next+n : next+n]
		}
		next += n
	}

	decisionStates := make([]*RuntimeATNState, len(t.DecisionStates))
	for i, s := range t.DecisionStates {
		if s >= 0 {
			decisionStates[i] = statePtrs[s]
		}
	}
	decisionMap := make([]*RuntimeATNState, len(decisionStates))
	copy(decisionMap, decisionStates)
	return NewRuntimeATN(statePtrs, decisionStates, decisionMap)
}
