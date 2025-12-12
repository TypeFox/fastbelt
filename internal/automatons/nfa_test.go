package automatons

import (
	"reflect"
	"sort"
	"testing"
)

func TestTransitionTargets_AddAndGet(t *testing.T) {
	// A
	// B
	//     C
	//       D
	targets := NewTransitionTargets()

	// Add transitions
	targets.Add(NewRuneSet_Rune(1), 2)
	targets.Add(NewRuneSet_Rune(1), 3)
	targets.Add(NewRuneSet_Rune(5), 4)
	targets.Add(NewRuneSet_Rune(7), 5)

	// Test getTargets
	testCases := []struct {
		char     rune
		expected []int
	}{
		{0, []int{}},
		{1, []int{2, 3}},
		{2, []int{}},
		{5, []int{4}},
		{6, []int{}},
		{7, []int{5}},
		{8, []int{}},
	}

	for _, tc := range testCases {
		actual := targets.GetTargets(NewRuneSet_Rune(tc.char))
		sort.Ints(actual)
		sort.Ints(tc.expected)
		if !reflect.DeepEqual(actual, tc.expected) {
			t.Errorf("GetTargets(%d): expected %v, got %v", tc.char, tc.expected, actual)
		}
	}

	// Test contains
	containsTests := []struct {
		char     rune
		expected bool
	}{
		{0, false},
		{1, true},
		{2, false},
		{5, true},
		{7, true},
		{8, false},
	}

	for _, tc := range containsTests {
		actual := targets.Contains(tc.char)
		if actual != tc.expected {
			t.Errorf("Contains(%d): expected %t, got %t", tc.char, tc.expected, actual)
		}
	}
}

func TestTransitionTargets_AddOverlappingRanges(t *testing.T) {
	// AAA
	//     BBB
	//  CCCCC
	targets := NewTransitionTargets()

	targets.Add(NewRuneSet_Range(1, 3), 2)
	targets.Add(NewRuneSet_Range(5, 7), 4)
	targets.Add(NewRuneSet_Range(2, 6), 3)

	testCases := []struct {
		char     rune
		expected []int
	}{
		{0, []int{}},
		{1, []int{2}},
		{2, []int{2, 3}},
		{3, []int{2, 3}},
		{4, []int{3}},
		{5, []int{3, 4}},
		{6, []int{3, 4}},
		{7, []int{4}},
		{8, []int{}},
	}

	for _, tc := range testCases {
		actual := targets.GetTargets(NewRuneSet_Rune(tc.char))
		sort.Ints(actual)
		sort.Ints(tc.expected)
		if !reflect.DeepEqual(actual, tc.expected) {
			t.Errorf("GetTargets(%d): expected %v, got %v", tc.char, tc.expected, actual)
		}
	}

	// Test contains
	containsTests := []struct {
		char     rune
		expected bool
	}{
		{0, false},
		{1, true},
		{2, true},
		{3, true},
		{4, true},
		{5, true},
		{6, true},
		{7, true},
		{8, false},
	}

	for _, tc := range containsTests {
		actual := targets.Contains(tc.char)
		if actual != tc.expected {
			t.Errorf("Contains(%d): expected %t, got %t", tc.char, tc.expected, actual)
		}
	}
}

func TestTransitionTargets_EpsilonTransitions(t *testing.T) {
	targets := NewTransitionTargets()

	// Add epsilon transitions (empty charset)
	targets.Add(NewRuneSet_Empty(), 2)
	targets.Add(NewRuneSet_Empty(), 3)

	epsilonTargets := targets.GetEpsilonTargets()
	sort.Ints(epsilonTargets)
	expected := []int{2, 3}
	sort.Ints(expected)

	if !reflect.DeepEqual(epsilonTargets, expected) {
		t.Errorf("GetEpsilonTargets(): expected %v, got %v", expected, epsilonTargets)
	}

	if !targets.ContainsEpsilon() {
		t.Error("ContainsEpsilon() should return true")
	}
}

func TestTransitionTargets_AddOverlappingRangesDifferentOrder(t *testing.T) {
	// AAA
	//  BBBBB
	//     CCC
	targets := NewTransitionTargets()

	targets.Add(NewRuneSet_Range(1, 3), 2)
	targets.Add(NewRuneSet_Range(2, 6), 3)
	targets.Add(NewRuneSet_Range(5, 7), 4)

	testCases := []struct {
		char     rune
		expected []int
	}{
		{0, []int{}},
		{1, []int{2}},
		{2, []int{2, 3}},
		{3, []int{2, 3}},
		{4, []int{3}},
		{5, []int{3, 4}},
		{6, []int{3, 4}},
		{7, []int{4}},
		{8, []int{}},
	}

	for _, tc := range testCases {
		actual := targets.GetTargets(NewRuneSet_Rune(tc.char))
		sort.Ints(actual)
		sort.Ints(tc.expected)
		if !reflect.DeepEqual(actual, tc.expected) {
			t.Errorf("GetTargets(%d): expected %v, got %v", tc.char, tc.expected, actual)
		}
	}

	containsTests := []struct {
		char     rune
		expected bool
	}{
		{0, false},
		{1, true},
		{2, true},
		{3, true},
		{4, true},
		{5, true},
		{6, true},
		{7, true},
		{8, false},
	}

	for _, tc := range containsTests {
		actual := targets.Contains(tc.char)
		if actual != tc.expected {
			t.Errorf("Contains(%d): expected %t, got %t", tc.char, tc.expected, actual)
		}
	}
}

func TestTransitionTargets_IterateAll(t *testing.T) {
	targets := NewTransitionTargets()

	targets.Add(NewRuneSet_Range(1, 3), 1)
	targets.Add(NewRuneSet_Range(3, 6), 2)
	targets.Add(NewRuneSet_Range(7, 9), 3)

	testCases := []struct {
		char     rune
		expected []int
	}{
		{1, []int{1}},
		{2, []int{1}},
		{3, []int{1, 2}},
		{4, []int{2}},
		{5, []int{2}},
		{6, []int{2}},
		{7, []int{3}},
		{8, []int{3}},
		{9, []int{3}},
	}

	for _, tc := range testCases {
		actual := targets.GetTargets(NewRuneSet_Rune(tc.char))
		sort.Ints(actual)
		sort.Ints(tc.expected)
		if !reflect.DeepEqual(actual, tc.expected) {
			t.Errorf("GetTargets(%d): expected %v, got %v", tc.char, tc.expected, actual)
		}
	}
}

func TestNFA_BasicFunctionality(t *testing.T) {
	nfa := NewNFA(0, 3)

	// Add some accepting states
	nfa.AddAcceptingState(2)

	// Add some transitions
	nfa.AddTransition(0, NewRuneSet_Rune('a'), 1)
	nfa.AddTransition(1, NewRuneSet_Rune('b'), 2)
	nfa.AddTransition(0, NewRuneSet_Empty(), 2) // epsilon transition

	// Test basic properties
	if nfa.GetStartState() != 0 {
		t.Errorf("Expected start state 0, got %d", nfa.GetStartState())
	}

	if nfa.GetStateCount() != 3 {
		t.Errorf("Expected state count 3, got %d", nfa.GetStateCount())
	}

	acceptingStates := nfa.GetAcceptingStates()
	if !acceptingStates[2] {
		t.Error("State 2 should be accepting")
	}

	// Test transitions
	transitions := nfa.GetTransitionsBySource()
	if len(transitions) != 2 {
		t.Errorf("Expected 2 source states with transitions, got %d", len(transitions))
	}

	// Test state 0 transitions
	state0Transitions := transitions[0]
	if state0Transitions == nil {
		t.Fatal("State 0 should have transitions")
	}

	// Test 'a' transition
	targetsA := state0Transitions.GetTargets(NewRuneSet_Rune('a'))
	expectedA := []int{1}
	if !reflect.DeepEqual(targetsA, expectedA) {
		t.Errorf("State 0 'a' transition: expected %v, got %v", expectedA, targetsA)
	}

	// Test epsilon transition
	epsilonTargets := state0Transitions.GetEpsilonTargets()
	expectedEpsilon := []int{2}
	if !reflect.DeepEqual(epsilonTargets, expectedEpsilon) {
		t.Errorf("State 0 epsilon transitions: expected %v, got %v", expectedEpsilon, epsilonTargets)
	}

	// Test state 1 transitions
	state1Transitions := transitions[1]
	if state1Transitions == nil {
		t.Fatal("State 1 should have transitions")
	}

	targetsB := state1Transitions.GetTargets(NewRuneSet_Rune('b'))
	expectedB := []int{2}
	if !reflect.DeepEqual(targetsB, expectedB) {
		t.Errorf("State 1 'b' transition: expected %v, got %v", expectedB, targetsB)
	}
}

func TestTargetGroup(t *testing.T) {
	// Test NewTargetGroup
	runeRange := RuneRange{Start: 'a', End: 'z', Includes: true}
	group := NewTargetGroup(runeRange, 1, 2, 3)

	if group.Length() != 3 {
		t.Errorf("Expected length 3, got %d", group.Length())
	}

	// Test Add
	group.Add(4)
	if group.Length() != 4 {
		t.Errorf("Expected length 4 after adding target, got %d", group.Length())
	}

	// Test GetTargets
	targets := group.GetTargets()
	sort.Ints(targets)
	expected := []int{1, 2, 3, 4}
	if !reflect.DeepEqual(targets, expected) {
		t.Errorf("GetTargets(): expected %v, got %v", expected, targets)
	}

	// Test iterator
	var iteratedTargets []int
	for target := range group.AllTargets() {
		iteratedTargets = append(iteratedTargets, target)
	}
	sort.Ints(iteratedTargets)
	if !reflect.DeepEqual(iteratedTargets, expected) {
		t.Errorf("AllTargets() iterator: expected %v, got %v", expected, iteratedTargets)
	}
}

func TestNFA_AllTransitionsIterator(t *testing.T) {
	nfa := NewNFA(0, 2)

	nfa.AddTransition(0, NewRuneSet_Rune('a'), 1)
	nfa.AddTransition(0, NewRuneSet_Empty(), 1) // epsilon
	nfa.AddTransition(1, NewRuneSet_Range('b', 'd'), 0)

	var transitions []Transition
	for transition := range nfa.AllTransitions() {
		transitions = append(transitions, transition)
	}

	// Should have at least 3 transitions: 'a'->1, epsilon->1, and one for the range 'b'-'d'->0
	if len(transitions) < 3 {
		t.Errorf("Expected at least 3 transitions, got %d", len(transitions))
	}

	// Check that we have the expected transitions
	hasATransition := false
	hasEpsilonTransition := false
	hasBRangeTransition := false

	for _, trans := range transitions {
		if trans.Source == 0 && trans.Target == 1 {
			if trans.Chars.IncludesRune('a') && trans.Chars.Length() == 1 {
				hasATransition = true
			} else if trans.Chars.Length() == 0 {
				hasEpsilonTransition = true
			}
		} else if trans.Source == 1 && trans.Target == 0 {
			if trans.Chars.IncludesRune('b') && trans.Chars.IncludesRune('c') && trans.Chars.IncludesRune('d') {
				hasBRangeTransition = true
			}
		}
	}

	if !hasATransition {
		t.Error("Missing 'a' transition from state 0 to state 1")
	}
	if !hasEpsilonTransition {
		t.Error("Missing epsilon transition from state 0 to state 1")
	}
	if !hasBRangeTransition {
		t.Error("Missing 'b'-'d' range transition from state 1 to state 0")
	}
}
