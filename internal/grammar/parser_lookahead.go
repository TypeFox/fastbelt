// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package grammar

import "typefox.dev/fastbelt/parser"

// fastbeltParserLookahead extends the generated default with disambiguation for
// group cardinality guards inside parser and composite rules. Optional semicolons
// between rules mean an ID can start either another element in the current rule
// body or the next parser rule; SLL cannot distinguish those cases when only
// LA(1) is considered. The next composite rule always starts with the "composite"
// keyword, so no extra disambiguation is needed for that case.
type fastbeltParserLookahead struct {
	DefaultFastbeltParserLookahead
}

func newFastbeltParserLookahead() FastbeltParserLookahead {
	return &fastbeltParserLookahead{}
}

// isNextParserRule reports whether the upcoming tokens start a new parser rule
// (Name=ID ":" or Name returns Type=ID ":").
func isNextParserRule(state *parser.ParserState) bool {
	if state.LA(1).Type != Token_ID {
		return false
	}
	switch state.LA(2).Type {
	case Keyword_Colon, Keyword_returns:
		return true
	default:
		return false
	}
}

func (l *fastbeltParserLookahead) groupGuardContinue(state *parser.ParserState, decision int) bool {
	if isNextParserRule(state) {
		return false
	}
	prediction, _ := state.AdaptivePredict(decision, l.PredictionMode())
	return prediction == 0
}

func (l *fastbeltParserLookahead) GroupOptional(state *parser.ParserState) bool {
	return l.groupGuardContinue(state, DecisionGroupOptional)
}

func (l *fastbeltParserLookahead) GroupElementsLoop(state *parser.ParserState) bool {
	return l.groupGuardContinue(state, DecisionGroupElementsLoop)
}

func (l *fastbeltParserLookahead) CompositeGroupOptional(state *parser.ParserState) bool {
	return l.groupGuardContinue(state, DecisionCompositeGroupOptional)
}

func (l *fastbeltParserLookahead) CompositeGroupElementsLoop(state *parser.ParserState) bool {
	return l.groupGuardContinue(state, DecisionCompositeGroupElementsLoop)
}
