package automatons

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDotFile(t *testing.T) {
	kit := NewConstructionKit()

	a := kit.Consume(NewRuneSetRune('a'))
	b := kit.Consume(NewRuneSetRune('b'))
	ab := kit.Concat(a, b)
	aOrAB := kit.Alternate(a, ab)

	actual := aOrAB.DotFile()
	expected := `digraph NFA {
  rankdir=LR;
  node [shape=circle];
  start [shape=point];
  start -> 0;
  2 [shape=doublecircle];
  3 [shape=doublecircle];
  0 -> 1 [label="['\u0000'-'` + "`" + `']"];
  0 -> 1 [label="['b'-'\U0010FFFF']"];
  0 -> 2 [label="['a']"];
  1 -> 1 [label="['\u0000'-'\U0010FFFF']"];
  2 -> 1 [label="['\u0000'-'a']"];
  2 -> 1 [label="['c'-'\U0010FFFF']"];
  2 -> 3 [label="['b']"];
  3 -> 1 [label="['\u0000'-'\U0010FFFF']"];
}
`
	require.Equal(t, expected, actual)
}

func makeDFA() *NFA {
	kit := NewConstructionKit()
	a := kit.Consume(NewRuneSetRune('a'))
	b := kit.Consume(NewRuneSetRune('b'))
	ab := kit.Concat(a, b)
	aOrAB := kit.Alternate(a, ab)
	return aOrAB
}

func makeNFA() *NFA {
	builder := NewNFABuilder()

	start := builder.AddState()
	accept := builder.AddState()
	dead1 := builder.AddState()
	dead2 := builder.AddState()

	builder.SetStartState(start)
	builder.AcceptState(accept)

	builder.AddTransitionForSingleRune(start, accept, 'a')
	builder.AddTransitionForSingleRune(start, dead1, 'b')
	builder.AddTransitionForSingleRune(start, dead2, 'c')
	builder.AddTransitionForRuneSet(dead1, dead1, NewRuneSetFull())
	builder.AddTransitionForRuneSet(dead2, dead2, NewRuneSetFull())

	nfa := builder.Build()
	return nfa
}

func TestComputeAcceptanceReachability_ConstructedDFA(t *testing.T) {
	aOrAB := makeDFA()

	actual := aOrAB.ComputeAcceptanceReachability()
	expected := map[int]bool{
		0: true,
		//unreachable, there is only one because the consstruction kit makes DFA-like NFAs which have only one dead state
		1: false,
		2: true,
		3: true,
	}
	require.Equal(t, expected, actual)
}

func TestComputeAcceptanceReachability_HandmadeNFA(t *testing.T) {
	nfa := makeNFA()
	actual := nfa.ComputeAcceptanceReachability()
	expected := map[int]bool{
		0: true,
		1: true,
		2: false,
		3: false,
	}
	require.Equal(t, expected, actual)
}

func TestDeadState_ConstructedDFA(t *testing.T) {
	nfa := makeDFA()
	actual := nfa.DeadState()
	expected := 1
	require.Equal(t, expected, actual)
}

func TestDeadState_HandmadeNFA(t *testing.T) {
	nfa := makeNFA()
	//panics because this NFA has multiple dead states
	require.Panics(t, func() { nfa.DeadState() })
}
