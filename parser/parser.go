// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package parser

import (
	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/util/collections"
)

// Parser defines the interface for parsing tokens (lexer output) into AST nodes.
type Parser interface {
	Parse(document *core.Document) *ParseResult
}

type ParseResult struct {
	Node   core.AstNode
	Errors []*core.ParserError
}

const (
	// Indicates that the parser is not currently in error mode.
	ErrorModeNone = iota
	// Indicates that the parser has encountered an error and is currently in error mode.
	ErrorModeFail
	// Indicates that the parser has encountered an error and is ready to attempt recovery.
	ErrorModeRecover
)

type ParserState struct {
	Tokens    []core.Token
	Length    int
	Index     int
	ErrorMode int
	// ErrorRecoveryMode is set by ReportError and cleared by ReportMatch
	// (called from a successful Consume). While set, further ReportError
	// calls are dropped so that a single underlying mistake produces a
	// single diagnostic instead of one per consume attempt during unwind.
	// Parsing continues throughout - this flag is purely about message
	// deduplication, never about halting.
	ErrorRecoveryMode bool
	errors            []*core.ParserError
	atn               *RuntimeATN
	followStates      []int // stack of atn.States array indices for follow-set computation
	recovery          ErrorRecoveryStrategy
	messages          ErrorMessageProvider
	sim               *parserATNSimulator // lazily created adaptive (ALL(*)) predictor
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
	if p.ErrorMode != ErrorModeNone {
		return
	}
	p.errors = append(p.errors, core.NewParserError(msg, token))
	p.ErrorMode = ErrorModeFail
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

// LL1Lookahead is an optimized LL(1) lookahead table. Types holds the
// discriminating token types for error reporting. Lookup is indexed by
// TokenType.Id; a value >= 1 is the 1-based alternative to take, 0 means
// no alternative expects that token.
type LL1Lookahead struct {
	Types  []*core.TokenType
	Lookup []int
}

func NewParserState(tokens []core.Token, atn *RuntimeATN, recovery ErrorRecoveryStrategy, messages ErrorMessageProvider) *ParserState {
	if atn == nil {
		panic("atn must be provided")
	}
	return &ParserState{
		Tokens:       tokens,
		Length:       len(tokens),
		Index:        0,
		ErrorMode:    ErrorModeNone,
		errors:       []*core.ParserError{},
		atn:          atn,
		followStates: nil,
		recovery:     recovery,
		messages:     messages,
	}
}

// LA returns the token at the given lookahead offset.
// Returns a pointer to [core.EOFToken] if the offset is out of bounds.
func (p *ParserState) LA(offset int) *core.Token {
	pos := p.Index + offset - 1
	if pos < 0 || pos >= p.Length {
		return &core.EOFToken
	}
	return &p.Tokens[pos]
}

func (p *ParserState) Consume(tokenType *core.TokenType) *core.Token {
	if p.ErrorMode != ErrorModeNone {
		return nil
	}
	current := p.LA(1)
	if !tokenType.Matches(current.Type) {
		if current.Type == core.EOF {
			p.AppendError(p.messages.UnexpectedEndOfInput(tokenType), current)
			return nil
		}
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

type ParserLookahead interface {
	PredictionMode() PredictionMode
	SetPredictionMode(mode PredictionMode)
}

type DefaultParserLookahead struct {
	predictionMode PredictionMode
}

// PredictionMode returns the current prediction mode. See [PredictionMode] for more details.
func (l *DefaultParserLookahead) PredictionMode() PredictionMode {
	return l.predictionMode
}

// SetPredictionMode sets the prediction mode. See [PredictionMode] for more details.
// Setting the mode will affect all adaptive lookahead decisions performed by the parser.
// When switching to [PredictionModeLL], it might be beneficial to override the languages'
// parser lookahead service with a custom implementation that only applies this mode to certain
// decisions, to avoid the performance cost of full-context prediction where it is not needed.
func (l *DefaultParserLookahead) SetPredictionMode(mode PredictionMode) {
	l.predictionMode = mode
}

// PredictionFailure describes why an alternatives decision could not be
// resolved. It is self-contained so it can be handed straight to
// ErrorMessageProvider.NoViableAlternative without further parser access: Token
// is the divergence point (where input dead-ended, used as the error position),
// and Expected lists the token types that would have allowed prediction to
// continue. Expected may be empty (e.g. for the static Lookahead path or an
// error-mode bail), in which case the message falls back to its generic form.
type PredictionFailure struct {
	Token    *core.Token
	Expected []*core.TokenType
}

// AdaptivePredict resolves the decision at the given decision index using
// ALL(*) prediction over the ATN, returning the predicted 0-based
// alternative (usable directly as a generated switch-case label) or -1 if no
// alternative is viable. The second result is non-nil exactly when the
// alternative is -1, carrying the divergence token and expected-token set for
// the NoViableAlternative diagnostic. Generated code emits this for decisions
// that a single lookahead token cannot disambiguate; LL(1)-unique decisions keep
// the cheaper Lookahead path.
//
// Prediction is non-consuming and must not run while unwinding after an error:
// LAId returns -1 in error mode, which the predictor would otherwise mistake
// for EOF, so we bail out to the no-viable sentinel (the generated switch then
// takes its default arm, exactly as Lookahead would).
func (p *ParserState) AdaptivePredict(decision int, mode PredictionMode) (int, *PredictionFailure) {
	if p.ErrorMode != ErrorModeNone {
		return -1, &PredictionFailure{Token: p.LA(1)}
	}
	if p.sim == nil {
		p.sim = newParserATNSimulator(p.atn, p)
	}
	return p.sim.adaptivePredict(decision, mode, p.initialPredictionContext(mode))
}

// initialPredictionContext converts the parser's real call stack (followStates,
// which are atn.States array indices pushed by EnterRule) into the
// PredictionContext chain used to seed full-context LL prediction.
func (p *ParserState) initialPredictionContext(mode PredictionMode) *predictionContext {
	ctx := emptyPredictionContext()
	if mode == PredictionModeLL {
		// Only include the full context if we're in full-context LL mode
		for _, followIdx := range p.followStates {
			ctx = singletonPredictionContext(ctx, followIdx)
		}
	}
	return ctx
}

// Lookahead resolves a decision using a pre-built LL(1) lookup table in O(1),
// returning the matching 0-based alternative or -1 if none matches.
func (p *ParserState) Lookahead(value LL1Lookahead) (int, *PredictionFailure) {
	if p.ErrorMode != ErrorModeNone {
		return -1, &PredictionFailure{Token: p.LA(1)}
	}
	la := p.LA(1)
	id := la.TypeId
	if id < len(value.Lookup) {
		// Note: generated lookup tables are 1-based to simplify code generation.
		// That way, 0 (default value) can simply mean "no match"
		if alt := value.Lookup[id]; alt >= 1 {
			return alt - 1, nil
		}
	}
	return -1, &PredictionFailure{Token: la, Expected: value.Types}
}

// EnterRule pushes a follow-state index onto the stack.
func (p *ParserState) EnterRule(followStateIdx int) {
	p.followStates = append(p.followStates, followStateIdx)
}

// ExitRule pops the top follow-state from the stack and tries to recover from any errors.
func (p *ParserState) ExitRule() {
	if len(p.followStates) > 0 {
		p.followStates = p.followStates[:len(p.followStates)-1]
	}
	if p.ErrorMode == ErrorModeFail {
		// Once we exit the rule where the error was detected, we can attempt recovery.
		p.ErrorMode = ErrorModeRecover
		p.recovery.Recover(p)
	}
}

// Sync delegates to the recovery strategy to discard unexpected tokens before
// optional/loop guards.
func (p *ParserState) Sync(decisionStateIdx int) {
	p.recovery.Sync(p, decisionStateIdx)
}

// FollowSet returns the union of NextTokensAt for every frame on the follow-state stack.
func (p *ParserState) FollowSet() *collections.BitSet {
	sets := make([]*collections.BitSet, len(p.followStates))
	for i, idx := range p.followStates {
		sets[i] = p.atn.NextTokensAt(idx)
	}
	return collections.MergeBitSets(sets)
}
