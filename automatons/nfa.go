package automatons

import (
	"fmt"
	"iter"
)

// NFA represents a nondeterministic finite automaton
type NFA interface {
	GetStartState() int
	GetStateCount() int
	GetAcceptingStates() map[int]bool
	GetTransitionsBySource() map[int]TransitionTargets
	// Iterator over all transitions: (sourceState, targetState, charset)
	AllTransitions() iter.Seq[Transition]
}

// Transition represents a single transition in the NFA
type Transition struct {
	Source int
	Target int
	Chars  *RuneSet
}

// TransitionTargets represents all possible transitions from a given state
type TransitionTargets interface {
	// Contains checks if character c has any transitions
	Contains(c rune) bool
	// ContainsEpsilon checks if there are epsilon (empty) transitions
	ContainsEpsilon() bool
	// GetEpsilonTargets returns all target states reachable via epsilon transitions
	GetEpsilonTargets() []int
	// GetTargets returns all target states reachable by the given character set
	GetTargets(charset *RuneSet) []int
	// AllTransitions returns an iterator over all transitions with their character ranges
	AllTransitions() iter.Seq[TransitionInfo]
}

// TransitionInfo represents a transition with its character range and target states
type TransitionInfo struct {
	CharRange *RuneRange // nil for epsilon transitions
	Targets   []int
}

// TargetGroup groups targets for a specific character range
type TargetGroup struct {
	Range   RuneRange
	targets map[int]bool
}

// NewTargetGroup creates a new target group for the given range
func NewTargetGroup(runeRange RuneRange, targets ...int) *TargetGroup {
	tg := &TargetGroup{
		Range:   runeRange,
		targets: make(map[int]bool),
	}
	for _, target := range targets {
		tg.targets[target] = true
	}
	return tg
}

// Add adds a target to this group
func (tg *TargetGroup) Add(target int) {
	tg.targets[target] = true
}

// Length returns the number of targets
func (tg *TargetGroup) Length() int {
	return len(tg.targets)
}

// GetTargets returns all targets as a slice
func (tg *TargetGroup) GetTargets() []int {
	result := make([]int, 0, len(tg.targets))
	for target := range tg.targets {
		result = append(result, target)
	}
	return result
}

// AllTargets returns an iterator over all targets
func (tg *TargetGroup) AllTargets() iter.Seq[int] {
	return func(yield func(int) bool) {
		for target := range tg.targets {
			if !yield(target) {
				return
			}
		}
	}
}

// TransitionTargetsImpl implements TransitionTargets
type TransitionTargetsImpl struct {
	epsilonTargets map[int]bool
	nodes          []*TargetGroupNode
}

// TargetGroupNode represents a node in the transition structure
type TargetGroupNode struct {
	Range   RuneRange
	Targets *TargetGroup
}

// NewTransitionTargets creates a new TransitionTargetsImpl
func NewTransitionTargets() *TransitionTargetsImpl {
	return &TransitionTargetsImpl{
		epsilonTargets: make(map[int]bool),
		nodes:          make([]*TargetGroupNode, 0),
	}
}

// Contains checks if character c has any transitions
func (tt *TransitionTargetsImpl) Contains(c rune) bool {
	for _, node := range tt.nodes {
		if node.Range.Contains(c) {
			return true
		}
	}
	return false
}

// ContainsEpsilon checks if there are epsilon transitions
func (tt *TransitionTargetsImpl) ContainsEpsilon() bool {
	return len(tt.epsilonTargets) > 0
}

// GetEpsilonTargets returns all epsilon target states
func (tt *TransitionTargetsImpl) GetEpsilonTargets() []int {
	result := make([]int, 0, len(tt.epsilonTargets))
	for target := range tt.epsilonTargets {
		result = append(result, target)
	}
	return result
}

// GetTargets returns all target states reachable by the given character set
func (tt *TransitionTargetsImpl) GetTargets(charset *RuneSet) []int {
	if charset.Length() == 0 {
		return tt.GetEpsilonTargets()
	}

	result := make(map[int]bool)
	for _, charRange := range charset.Ranges {
		if charRange.Includes {
			for _, node := range tt.nodes {
				if tt.rangesIntersect(charRange, node.Range) {
					for target := range node.Targets.targets {
						result[target] = true
					}
				}
			}
		}
	}

	targets := make([]int, 0, len(result))
	for target := range result {
		targets = append(targets, target)
	}
	return targets
}

// AllTransitions returns an iterator over all transitions
func (tt *TransitionTargetsImpl) AllTransitions() iter.Seq[TransitionInfo] {
	return func(yield func(TransitionInfo) bool) {
		// Epsilon transitions
		if len(tt.epsilonTargets) > 0 {
			targets := tt.GetEpsilonTargets()
			if !yield(TransitionInfo{CharRange: nil, Targets: targets}) {
				return
			}
		}

		// Character transitions
		for _, node := range tt.nodes {
			targets := node.Targets.GetTargets()
			if !yield(TransitionInfo{CharRange: &node.Range, Targets: targets}) {
				return
			}
		}
	}
}

// rangesIntersect checks if two RuneRanges intersect
func (tt *TransitionTargetsImpl) rangesIntersect(r1, r2 RuneRange) bool {
	if !r1.Includes || !r2.Includes {
		return false
	}
	start := max(r1.Start, r2.Start)
	end := min(r1.End, r2.End)
	return start <= end
}

// rangeIntersection returns the intersection of two RuneRanges, or nil if no intersection
func (tt *TransitionTargetsImpl) rangeIntersection(r1, r2 RuneRange) *RuneRange {
	if !r1.Includes || !r2.Includes {
		return nil
	}
	start := max(r1.Start, r2.Start)
	end := min(r1.End, r2.End)
	if start <= end {
		return NewRuneRange(start, end, true)
	}
	return nil
}

// Add adds a transition for the given character set to the target state
func (tt *TransitionTargetsImpl) Add(charset *RuneSet, target int) {
	if charset.Length() == 0 {
		// Epsilon transition
		tt.epsilonTargets[target] = true
		return
	}

	for _, charRange := range charset.Ranges {
		if charRange.Includes {
			tt.addRange(charRange, target)
		}
	}
}

// addRange adds a single character range transition
func (tt *TransitionTargetsImpl) addRange(charRange RuneRange, target int) {
	nodeIndex := 0

	// Find the position to insert/modify
	for nodeIndex < len(tt.nodes) && tt.nodes[nodeIndex].Range.End < charRange.Start {
		nodeIndex++
	}

	if nodeIndex >= len(tt.nodes) {
		// Add at the end
		newNode := &TargetGroupNode{
			Range:   charRange,
			Targets: NewTargetGroup(charRange, target),
		}
		tt.nodes = append(tt.nodes, newNode)
		return
	}

	// Handle overlapping ranges by splitting and merging
	tt.insertRangeAt(nodeIndex, charRange, target)
}

// insertRangeAt handles the complex logic of inserting a range at a specific position
func (tt *TransitionTargetsImpl) insertRangeAt(nodeIndex int, newRange RuneRange, target int) {
	if nodeIndex >= len(tt.nodes) {
		newNode := &TargetGroupNode{
			Range:   newRange,
			Targets: NewTargetGroup(newRange, target),
		}
		tt.nodes = append(tt.nodes, newNode)
		return
	}

	currentNode := tt.nodes[nodeIndex]
	currentRange := currentNode.Range

	// No overlap, insert before
	if newRange.End < currentRange.Start {
		newNode := &TargetGroupNode{
			Range:   newRange,
			Targets: NewTargetGroup(newRange, target),
		}
		tt.nodes = append(tt.nodes[:nodeIndex], append([]*TargetGroupNode{newNode}, tt.nodes[nodeIndex:]...)...)
		return
	}

	// Handle overlapping ranges
	newNodes := make([]*TargetGroupNode, 0)

	// Left part (before overlap)
	if newRange.Start < currentRange.Start {
		leftRange := RuneRange{Start: newRange.Start, End: currentRange.Start - 1, Includes: true}
		newNodes = append(newNodes, &TargetGroupNode{
			Range:   leftRange,
			Targets: NewTargetGroup(leftRange, target),
		})
	} else if currentRange.Start < newRange.Start {
		leftRange := RuneRange{Start: currentRange.Start, End: newRange.Start - 1, Includes: true}
		newNodes = append(newNodes, &TargetGroupNode{
			Range:   leftRange,
			Targets: currentNode.Targets,
		})
	}

	// Overlapping part
	overlapStart := max(newRange.Start, currentRange.Start)
	overlapEnd := min(newRange.End, currentRange.End)
	if overlapStart <= overlapEnd {
		overlapRange := RuneRange{Start: overlapStart, End: overlapEnd, Includes: true}
		mergedTargets := NewTargetGroup(overlapRange)
		// Add existing targets
		for _, existingTarget := range currentNode.Targets.GetTargets() {
			mergedTargets.Add(existingTarget)
		}
		// Add new target
		mergedTargets.Add(target)
		newNodes = append(newNodes, &TargetGroupNode{
			Range:   overlapRange,
			Targets: mergedTargets,
		})
	}

	// Right part (after overlap)
	if newRange.End > currentRange.End {
		// Continue with the remaining part of newRange
		remainingRange := RuneRange{Start: currentRange.End + 1, End: newRange.End, Includes: true}
		// Replace current node with new nodes so far
		tt.nodes = append(tt.nodes[:nodeIndex], append(newNodes, tt.nodes[nodeIndex+1:]...)...)
		// Recursively handle the remaining range
		tt.insertRangeAt(nodeIndex+len(newNodes), remainingRange, target)
		return
	} else if currentRange.End > newRange.End {
		rightRange := RuneRange{Start: newRange.End + 1, End: currentRange.End, Includes: true}
		newNodes = append(newNodes, &TargetGroupNode{
			Range:   rightRange,
			Targets: currentNode.Targets,
		})
	}

	// Replace the current node with the new nodes
	tt.nodes = append(tt.nodes[:nodeIndex], append(newNodes, tt.nodes[nodeIndex+1:]...)...)
}

// NFAImpl is a basic implementation of NFA
type NFAImpl struct {
	startState          int
	stateCount          int
	acceptingStates     map[int]bool
	transitionsBySource map[int]TransitionTargets
}

// NewNFA creates a new NFA
func NewNFA(startState, stateCount int) *NFAImpl {
	return &NFAImpl{
		startState:          startState,
		stateCount:          stateCount,
		acceptingStates:     make(map[int]bool),
		transitionsBySource: make(map[int]TransitionTargets),
	}
}

// GetStartState returns the start state
func (nfa *NFAImpl) GetStartState() int {
	return nfa.startState
}

// GetStateCount returns the number of states
func (nfa *NFAImpl) GetStateCount() int {
	return nfa.stateCount
}

// GetAcceptingStates returns the accepting states
func (nfa *NFAImpl) GetAcceptingStates() map[int]bool {
	return nfa.acceptingStates
}

// GetTransitionsBySource returns transitions organized by source state
func (nfa *NFAImpl) GetTransitionsBySource() map[int]TransitionTargets {
	return nfa.transitionsBySource
}

// AllTransitions returns an iterator over all transitions
func (nfa *NFAImpl) AllTransitions() iter.Seq[Transition] {
	return func(yield func(Transition) bool) {
		for sourceState, targets := range nfa.transitionsBySource {
			for info := range targets.AllTransitions() {
				var charset *RuneSet
				if info.CharRange != nil {
					charset = NewRuneSet_Range(info.CharRange.Start, info.CharRange.End)
				} else {
					charset = NewRuneSet_Empty()
				}

				for _, target := range info.Targets {
					if !yield(Transition{Source: sourceState, Target: target, Chars: charset}) {
						return
					}
				}
			}
		}
	}
}

// AddAcceptingState marks a state as accepting
func (nfa *NFAImpl) AddAcceptingState(state int) {
	nfa.acceptingStates[state] = true
}

// AddTransition adds a transition from source to target with the given character set
func (nfa *NFAImpl) AddTransition(source int, charset *RuneSet, target int) {
	if nfa.transitionsBySource[source] == nil {
		nfa.transitionsBySource[source] = NewTransitionTargets()
	}
	nfa.transitionsBySource[source].(*TransitionTargetsImpl).Add(charset, target)
}

func (nfa NFAImpl) String() string {
	result := "NFA:\n"
	result += "Start State: " + fmt.Sprintf("%d", nfa.startState) + "\n"
	result += "Accepting States: "
	for state := range nfa.acceptingStates {
		result += fmt.Sprintf("%d ", state)
	}
	result += "\nTransitions:\n"
	for source, targets := range nfa.transitionsBySource {
		for info := range targets.AllTransitions() {
			var charset *RuneSet
			if info.CharRange != nil {
				charset = NewRuneSet_Range(info.CharRange.Start, info.CharRange.End)
			} else {
				charset = NewRuneSet_Empty()
			}
			for _, target := range info.Targets {
				result += fmt.Sprintf("  %d --%v--> %d\n", source, charset, target)
			}
		}
	}
	return result
}
