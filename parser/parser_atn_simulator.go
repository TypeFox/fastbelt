// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package parser

import (
	"sort"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/util/collections"
)

// parserATNSimulator implements ALL(*) adaptive prediction over a
// Fastbelt RuntimeATN.
//
// Prediction is non-consuming: the simulator reads lookahead through la
// (1-based offset from the decision point) and never advances the underlying
// token cursor. It runs SLL first; on an SLL conflict it falls back to full-context
// LL seeded with the parser's real call stack (based on the supplied prediction mode).
type parserATNSimulator struct {
	atn         *RuntimeATN
	parserState *ParserState
	mergeCache  *mergeCache
}

func newParserATNSimulator(atn *RuntimeATN, p *ParserState) *parserATNSimulator {
	return &parserATNSimulator{
		atn:         atn,
		parserState: p,
	}
}

// PredictionMode controls whether adaptive prediction runs SLL or LL.
// See [PredictionModeSLL] and [PredictionModeLL] for details.
type PredictionMode int

const (
	// PredictionModeSLL (i.e. Strong LL) ignores the current parser context when
	// making a prediction. It corresponds to the way traditional recursive descent
	// parsers operate and is potentially faster than [PredictionModeLL], but may
	// result in syntax errors for valid inputs, depending on the grammar.
	//
	// It is the default prediction mode for any language.
	// Use [PredictionModeLL] explicitly to enable full-context prediction if necessary.
	PredictionModeSLL PredictionMode = iota
	// PredictionModeLL allows the current parser context to be used when making a prediction
	// It is more expensive than [PredictionModeSLL], but is guaranteed to produce correct
	// predictions for all syntactically valid inputs.
	PredictionModeLL
)

// FullContext reports whether this prediction mode performs full-context LL
// prediction (true) or uses Strong LL (SLL, false).
func (m PredictionMode) FullContext() bool {
	return m == PredictionModeLL
}

func (m PredictionMode) String() string {
	switch m {
	case PredictionModeLL:
		return "LL"
	case PredictionModeSLL:
		return "SLL"
	default:
		return "UNKNOWN"
	}
}

// adaptivePredict returns the predicted 0-based alternative for the given
// decision, or invalidAlt (-1) if no alternative is viable. The second result
// is non-nil exactly when the alternative is invalidAlt, carrying the divergence
// offset and expected-token set for diagnostics. outerContext is the parser's
// real call stack, used only for the full-context LL fallback.
func (s *parserATNSimulator) adaptivePredict(decision int, mode PredictionMode, outerContext *predictionContext) (int, *PredictionFailure) {
	s.mergeCache = newMergeCache()
	dfa := s.atn.decisionToDFA[decision]

	start := dfa.getStart()
	if start == nil {
		startClosure := s.computeStartState(dfa.atnStartState, emptyPredictionContext(), false)
		start = dfa.addState(newDFAState(-1, startClosure))
		dfa.setStart(start)
	}
	return s.execATN(dfa, start, mode, outerContext)
}

func (s *parserATNSimulator) execATN(dfa *dfa, start *dfaState, mode PredictionMode, outerContext *predictionContext) (int, *PredictionFailure) {
	previousState := start
	offset := 1
	token := s.parserState.LA(offset)
	for {
		d := dfa.getExistingEdge(previousState, token.Type)
		if d == nil {
			d = s.computeTargetState(dfa, previousState, mode, token.Type)
		}
		if d == errorDFAState {
			// No viable target. If a path actually finished the decision rule,
			// commit to it so the eventual syntax error is reported at a more
			// localized spot; otherwise signal no-viable-alt and capture the
			// divergence point (offset) and expected tokens for diagnostics.
			alt := s.getAltThatFinishedDecisionEntryRule(previousState.configs)
			if alt == invalidAlt {
				return alt, buildFailure(token, previousState.configs)
			}
			return alt, nil
		}
		if d.requiresFullContext && mode.FullContext() {
			// Some decision paths have a conflict that depends on the full parser context
			// We essentially restart the simulator here, using the full context
			// Note that we're not preserving the DFA in this case
			startClosure := s.computeStartState(dfa.atnStartState, outerContext, true)
			return s.execATNWithFullContext(startClosure)
		}
		if d.isAcceptState {
			return d.prediction, nil
		}
		previousState = d
		if token.Type != core.EOF {
			offset++
			token = s.parserState.LA(offset)
		}
	}
}

// buildFailure collects the no-viable-alternative diagnostics from the last live
// configuration set: the expected-token set is every atom transition reachable
// from those configs (deduped by token id, sorted by name for a deterministic
// message), and offset marks where input diverged.
func buildFailure(token *core.Token, configs *atnConfigSet) *PredictionFailure {
	seen := map[int]bool{}
	var expected []*core.TokenType
	for _, c := range configs.configs {
		for _, tr := range c.state.Transitions {
			at, ok := tr.(*RuntimeAtomTransition)
			if !ok || at.TokenType == nil || seen[at.TokenType.Id] {
				continue
			}
			seen[at.TokenType.Id] = true
			expected = append(expected, at.TokenType)
		}
	}
	sort.Slice(expected, func(i, j int) bool {
		return expected[i].Name < expected[j].Name
	})
	return &PredictionFailure{Token: token, Expected: expected}
}

func (s *parserATNSimulator) computeTargetState(dfa *dfa, previousD *dfaState, mode PredictionMode, t *core.TokenType) *dfaState {
	reach := s.computeReachSet(previousD.configs, t, mode.FullContext())
	if reach == nil {
		return dfa.addEdge(previousD, t, errorDFAState)
	}
	state := newDFAState(-1, reach)
	predictedAlt := getUniqueAlt(reach)
	if predictedAlt != invalidAlt {
		state.isAcceptState = true
		reach.uniqueAlt = predictedAlt
		state.prediction = predictedAlt
	} else if hasConflictTerminatingPrediction(reach) {
		reach.conflictingAlts = s.getConflictingAlts(reach)
		state.requiresFullContext = true
		state.isAcceptState = true
		state.prediction = reach.conflictingAlts.Min()
	}
	return dfa.addEdge(previousD, t, state)
}

func (s *parserATNSimulator) execATNWithFullContext(start *atnConfigSet) (int, *PredictionFailure) {
	const fullCtx = true
	var reach *atnConfigSet
	previousState := start
	offset := 1
	token := s.parserState.LA(offset)
	predictedAlt := invalidAlt
	for {
		reach = s.computeReachSet(previousState, token.Type, fullCtx)
		if reach == nil {
			alt := s.getAltThatFinishedDecisionEntryRule(previousState)
			if alt == invalidAlt {
				return alt, buildFailure(token, previousState)
			}
			return alt, nil
		}
		altSubSets := getConflictingAltSubsets(reach)
		reach.uniqueAlt = getUniqueAlt(reach)
		if reach.uniqueAlt != invalidAlt {
			predictedAlt = reach.uniqueAlt
			break
		}
		predictedAlt = getSingleViableAlt(altSubSets)
		if predictedAlt != invalidAlt {
			break
		}
		previousState = reach
		if token.Type != core.EOF {
			offset++
			token = s.parserState.LA(offset)
		}
	}
	return predictedAlt, nil
}

func (s *parserATNSimulator) computeReachSet(closure *atnConfigSet, tokenType *core.TokenType, fullCtx bool) *atnConfigSet {
	intermediate := newATNConfigSet(fullCtx)
	var skippedStopStates []*atnConfig

	for _, c := range closure.configs {
		if isRuleStop(c) {
			if fullCtx || tokenType == core.EOF {
				// Only track stop states if we might need to dip into the outer context.
				// Will always be done in LL prediction mode, but only for EOF in SLL.
				skippedStopStates = append(skippedStopStates, c)
			}
			continue
		}
		for _, trans := range c.state.Transitions {
			target := getReachableTarget(trans, tokenType)
			if target != nil {
				intermediate.Add(newATNConfigWithState(c, target, c.context), s.mergeCache)
			}
		}
	}

	var reach *atnConfigSet
	if skippedStopStates == nil && tokenType != core.EOF {
		if len(intermediate.configs) == 1 || getUniqueAlt(intermediate) != invalidAlt {
			reach = intermediate
		}
	}
	if reach == nil {
		reach = newATNConfigSet(fullCtx)
		busy := newClosureBusy()
		for _, c := range intermediate.configs {
			s.closure(c, reach, busy, fullCtx)
		}
	}

	if tokenType == core.EOF {
		reach = s.removeAllConfigsNotInRuleStopState(reach)
	}

	if skippedStopStates != nil && (!fullCtx || !hasConfigInRuleStopState(reach)) {
		for _, c := range skippedStopStates {
			reach.Add(c, s.mergeCache)
		}
	}

	if reach.isEmpty() {
		return nil
	}
	return reach
}

func (s *parserATNSimulator) removeAllConfigsNotInRuleStopState(configs *atnConfigSet) *atnConfigSet {
	if allConfigsInRuleStopStates(configs) {
		return configs
	}
	result := newATNConfigSet(configs.fullCtx)
	for _, c := range configs.configs {
		if isRuleStop(c) {
			result.Add(c, s.mergeCache)
		}
	}
	return result
}

func (s *parserATNSimulator) computeStartState(start *RuntimeATNState, ctx *predictionContext, fullCtx bool) *atnConfigSet {
	configs := newATNConfigSet(fullCtx)
	for i, tr := range start.Transitions {
		target := tr.GetTarget()
		// 0-based alt: the i-th transition out of the decision state is alt i,
		// matching the generated switch-case numbering.
		c := newATNConfig(target, i, ctx)
		busy := newClosureBusy()
		s.closure(c, configs, busy, fullCtx)
	}
	return configs
}

func getReachableTarget(trans RuntimeTransition, tokenType *core.TokenType) *RuntimeATNState {
	if at, ok := trans.(*RuntimeAtomTransition); ok && at.TokenType != nil && at.TokenType.Matches(tokenType) {
		return at.Target
	}
	return nil
}

func (s *parserATNSimulator) closure(config *atnConfig, configs *atnConfigSet, busy *closureBusy, fullCtx bool) {
	s.closureCheckingStopState(config, configs, busy, fullCtx, 0)
}

func (s *parserATNSimulator) closureCheckingStopState(config *atnConfig, configs *atnConfigSet, busy *closureBusy, fullCtx bool, depth int) {
	if isRuleStop(config) {
		if config.context != nil && !config.context.isEmpty() {
			for i := 0; i < config.context.length(); i++ {
				if config.context.getReturnState(i) == predictionContextEmptyReturnState {
					if fullCtx {
						configs.Add(newATNConfigWithContext(config, emptyPredictionContext()), s.mergeCache)
						continue
					}
					// No context info: chase follow links.
					s.closureWork(config, configs, busy, fullCtx, depth)
					continue
				}
				returnState := s.atn.States[config.context.getReturnState(i)]
				newContext := config.context.getParent(i) // "pop" return state
				c := newATNConfigWithState(config, returnState, newContext)
				c.reachesIntoOuterContext = config.reachesIntoOuterContext
				s.closureCheckingStopState(c, configs, busy, fullCtx, depth-1)
			}
			return
		} else if fullCtx {
			configs.Add(config, s.mergeCache) // reached end of start rule
			return
		}
		// else: no context info, fall through to closureWork to chase follow links
	}
	s.closureWork(config, configs, busy, fullCtx, depth)
}

func (s *parserATNSimulator) closureWork(config *atnConfig, configs *atnConfigSet, busy *closureBusy, fullCtx bool, depth int) {
	state := config.state
	if !state.EpsilonOnlyTransitions {
		configs.Add(config, s.mergeCache)
	}
	srcIsRuleStop := isRuleStop(config)
	for _, t := range state.Transitions {
		c := s.getEpsilonTarget(config, t)
		if c == nil {
			continue
		}
		newDepth := depth
		if srcIsRuleStop {
			// Target fell off the end of the rule: track dipping into outer ctx.
			c.reachesIntoOuterContext++
			if present := busy.add(c); present {
				continue // avoid infinite recursion for right-recursive rules
			}
			configs.dipsIntoOuterContext = true
			newDepth--
		} else {
			if !t.IsEpsilon() {
				if present := busy.add(c); present {
					continue // avoid infinite recursion for EOF* / EOF+
				}
			}
			if _, ok := t.(*RuntimeRuleTransition); ok {
				if newDepth >= 0 {
					newDepth++
				}
			}
		}
		s.closureCheckingStopState(c, configs, busy, fullCtx, newDepth)
	}
}

func (s *parserATNSimulator) getEpsilonTarget(config *atnConfig, t RuntimeTransition) *atnConfig {
	switch tt := t.(type) {
	case *RuntimeRuleTransition:
		return s.ruleTransition(config, tt)
	case *RuntimeEpsilonTransition:
		return newATNConfigWithState(config, tt.Target, config.context)
	default:
		return nil
	}
}

func (s *parserATNSimulator) ruleTransition(config *atnConfig, t *RuntimeRuleTransition) *atnConfig {
	returnState := s.atn.stateIndex(t.FollowState)
	newContext := singletonPredictionContext(config.context, returnState)
	return newATNConfigWithState(config, t.Target, newContext)
}

func getUniqueAlt(configs *atnConfigSet) int {
	alt := invalidAlt
	for _, c := range configs.configs {
		if alt == invalidAlt {
			alt = c.alt
		} else if c.alt != alt {
			return invalidAlt
		}
	}
	return alt
}

func (s *parserATNSimulator) getConflictingAlts(configs *atnConfigSet) *collections.BitSet {
	return collections.MergeBitSets(getConflictingAltSubsets(configs))
}

func (s *parserATNSimulator) getAltThatFinishedDecisionEntryRule(configs *atnConfigSet) int {
	alts := collections.NewBitset()
	for _, c := range configs.configs {
		if c.reachesIntoOuterContext > 0 || (isRuleStop(c) && c.context.hasEmptyPath()) {
			alts.Insert(c.alt)
		}
	}
	if alts.Cardinality() == 0 {
		return invalidAlt
	}
	return alts.Min()
}

// closureBusy is the visited-config set that terminates closure on recursive
// grammars (right recursion, EOF loops).
type closureBusy struct {
	m *collections.BucketSet[*atnConfig]
}

func newClosureBusy() *closureBusy {
	return &closureBusy{m: collections.NewBucketSet[*atnConfig]()}
}

// add records c and returns true if it was already present.
func (cb *closureBusy) add(c *atnConfig) bool {
	return !cb.m.Add(c)
}
