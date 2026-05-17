// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package parser

import (
	core "typefox.dev/fastbelt"
)

// Parser defines the interface for parsing tokens (lexer output) into AST nodes.
type Parser interface {
	Parse(document *core.Document) *ParseResult
}

type ParseResult struct {
	Node   core.AstNode
	Errors []*core.ParserError
}

type ParserState struct {
	Tokens []core.Token
	Length int
	Index  int
	// InError is the hard-halt signal: while true, LA() returns nil so every
	// optional/loop guard evaluates false and the parser unwinds without
	// emitting further work.
	//
	// It is set in only two situations, both via AppendError:
	//
	//  1. BailErrorRecovery.RecoverInline returns nil. This is the bail
	//     strategy's whole point — halt immediately on the first mismatch.
	//     Without this flag, a loop like `for Lookahead(X) == 0 { Parse() }`
	//     would re-enter forever, since LA still returns the real (wrong)
	//     token and ErrorRecoveryMode only suppresses messages, not parsing.
	//
	//  2. Consume hits EOF mid-rule. Strictly an optimisation here: LA's
	//     bounds check already returns nil past Length, so the parser would
	//     unwind correctly without InError; the flag just short-circuits any
	//     remaining straight-line work after the EOF is detected.
	InError bool
	// ErrorRecoveryMode is set by ReportError and cleared by ReportMatch
	// (called from a successful Consume). While set, further ReportError
	// calls are dropped so that a single underlying mistake produces a
	// single diagnostic instead of one per consume attempt during unwind.
	// Parsing continues throughout — this flag is purely about message
	// deduplication, never about halting.
	ErrorRecoveryMode bool
	errors            []*core.ParserError
	atn               *RuntimeATN
	followStates      []int // stack of atn.States array indices for follow-set computation
	recovery          ErrorRecoveryStrategy
	messages          ErrorMessageProvider
}

func (p *ParserState) ATN() *RuntimeATN {
	return p.atn
}

// Messages returns the ErrorMessageProvider currently used to format
// diagnostic messages emitted by the parser.
func (p *ParserState) Messages() ErrorMessageProvider {
	return p.messages
}

func (p *ParserState) RecoveryStrategy() ErrorRecoveryStrategy {
	return p.recovery
}

func (p *ParserState) Errors() []*core.ParserError {
	return p.errors
}

func (p *ParserState) AppendError(msg string, token *core.Token) {
	p.errors = append(p.errors, core.NewParserError(msg, token))
	p.InError = true
}

// ReportError records a non-fatal parse error and enters error-recovery mode.
// While in error-recovery mode, subsequent ReportError calls are suppressed
// until reportMatch is called after a successful token match, so a single
// underlying mistake produces a single diagnostic rather than a cascade of
// messages as the parser tries (and fails) to consume the next several tokens.
//
// ReportError does NOT set inError. Parsing continues; if recovery cannot
// make progress, callers that synthesize/skip tokens (Consume, Sync, Recover)
// are responsible for advancing Index.
func (p *ParserState) ReportError(msg string, token *core.Token) {
	if p.ErrorRecoveryMode {
		return
	}
	p.ErrorRecoveryMode = true
	p.errors = append(p.errors, core.NewParserError(msg, token))
}

// ReportMatch exits error-recovery mode after a token has been successfully
// matched (either directly or via inline recovery).
func (p *ParserState) ReportMatch() {
	p.ErrorRecoveryMode = false
}

type LookaheadPath []int
type LookaheadOption []LookaheadPath
type LLkLookahead []LookaheadOption

func NewParserState(tokens []core.Token, atn *RuntimeATN, recovery ErrorRecoveryStrategy, messages ErrorMessageProvider) *ParserState {
	return &ParserState{
		Tokens:       tokens,
		Length:       len(tokens),
		Index:        0,
		InError:      false,
		errors:       []*core.ParserError{},
		atn:          atn,
		followStates: nil,
		recovery:     recovery,
		messages:     messages,
	}
}

// LA returns the token at the given lookahead offset.
// Returns nil if the offset is out of bounds or if the parser is currently in error mode.
func (p *ParserState) LA(offset int) *core.Token {
	// Test for inError first
	// prevents LA from returning real tokens while unwinding after an error,
	// which would cause infinite loops in guards.
	// Also, enables an optimization for the common EOF case: once LA returns nil, the parser
	// can short-circuit any remaining work in the current rule.
	// This circumvents the need for goto cleanup patterns in the generated code.
	if p.InError {
		return nil
	}
	return p.LARaw(offset)
}

// LARaw returns the token at offset without checking inError.
// Only used inside recovery strategy methods that must see real tokens after an error.
func (p *ParserState) LARaw(offset int) *core.Token {
	pos := p.Index + offset - 1
	if pos < 0 || pos >= p.Length {
		return nil
	}
	return &p.Tokens[pos]
}

func (p *ParserState) LAId(offset int) int {
	la := p.LA(offset)
	if la == nil {
		return core.EOF.Id
	}
	return la.TypeId
}

func (p *ParserState) Consume(tokenType *core.TokenType) *core.Token {
	if p.InError {
		return nil
	}
	current := p.LA(1)
	if current == nil {
		p.AppendError(p.messages.UnexpectedEndOfInput(tokenType), nil)
		return nil
	}
	if current.TypeId != tokenType.Id {
		recovered, ok := p.recovery.RecoverInline(p, tokenType)
		if ok {
			return recovered
		}
		p.AppendError(p.messages.UnexpectedToken(current), current)
		return nil
	}
	p.ReportMatch()
	p.Index++
	return current
}

func (p *ParserState) Lookahead(value LLkLookahead) int {
	for i, option := range value {
	outer:
		for _, path := range option {
			for j, tokenType := range path {
				if p.LAId(j+1) != tokenType {
					continue outer
				}
			}
			return i
		}
	}
	return -1
}

// EnterRule pushes a follow-state index onto the stack and triggers recovery
// if the parser is currently in error mode.
func (p *ParserState) EnterRule(followStateIdx int) {
	p.followStates = append(p.followStates, followStateIdx)
	if p.InError {
		p.recovery.Recover(p)
	}
}

// ExitRule pops the top follow-state from the stack.
func (p *ParserState) ExitRule() {
	if len(p.followStates) > 0 {
		p.followStates = p.followStates[:len(p.followStates)-1]
	}
}

// Sync delegates to the recovery strategy to discard unexpected tokens before
// optional/loop guards.
func (p *ParserState) Sync(decisionStateIdx int) {
	p.recovery.Sync(p, decisionStateIdx)
}

// FollowSet returns the union of NextTokensAt for every frame on the follow-state stack.
// The returned slice is indexed by TokenType.Id; out-of-range indices indicate
// "not in the follow set".
func (p *ParserState) FollowSet() []bool {
	if p.atn == nil {
		return nil
	}
	result := make([]bool, p.atn.TokenSetSize())
	for _, idx := range p.followStates {
		next := p.atn.NextTokensAt(idx)
		for i, v := range next {
			if v {
				result[i] = true
			}
		}
	}
	return result
}
