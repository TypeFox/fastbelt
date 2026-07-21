// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

// Package parser provides the shared parser runtime used by Fastbelt-generated
// languages. It does not implement parsing for any particular language; each
// grammar gets its own generated recursive-descent parser that drives the
// parser state, adaptive prediction, error recovery, and completion machinery
// of this package.
//
// Run [typefox.dev/fastbelt/cmd/fastbelt] generate on a .fb grammar file to
// emit the files described below. See [typefox.dev/fastbelt] for the overall
// toolchain and [typefox.dev/fastbelt/grammar] for how parser rules are
// declared in the grammar language.
//
// # Generated parsers
//
// Code generation turns the parser rules of a grammar into four cooperating
// files:
//
//   - `parser_gen.go` — a Parser implementing [Parser], with one ParseX method
//     per parser rule. Each method allocates the rule's AST node, matches
//     tokens through [ParserState], assigns them to node fields, and records
//     the node's token segment for source ranges.
//   - `parser_lookahead_gen.go` — the lookahead service, with one method per
//     decision point in the grammar: which alternative to take, whether to
//     enter an optional group, whether to iterate a loop once more. It is an
//     interface with a default implementation, so individual decisions can be
//     overridden without touching the parser.
//   - `atn_gen.go` — the [RuntimeATN], an augmented transition network that
//     encodes the grammar as a graph of states with token, epsilon, and rule
//     transitions. It backs adaptive prediction, error recovery, and
//     completion.
//   - `completion_parser_gen.go` — a stripped-down twin of the parser that
//     builds no AST and instead records the parse progress snapshots that
//     code completion needs.
//
// The generated `SetupGeneratedServices` registers the Parser, the lookahead
// service, and the default [ErrorRecoveryStrategy] and [ErrorMessageProvider]
// in the service container ([typefox.dev/fastbelt/util/service]), each unless
// a custom implementation was registered first. During document processing,
// [typefox.dev/fastbelt/workspace.DefaultDocumentParser] obtains the [Parser],
// calls Parse on the tokens produced by [typefox.dev/fastbelt/lexer], and
// stores the resulting root [core.AstNode] and syntax errors ([ParseResult])
// on the document.
//
// # Parsing model
//
// A generated parser is a recursive-descent parser: one method per rule,
// descending into sub-rules and consuming tokens left to right. [ParserState]
// carries the per-parse state: token slice, cursor, error mode, and rule stack.
//
// Wherever the grammar offers a choice — between alternatives, at an optional
// element, at a loop continuation — the parser asks its lookahead service.
// Depending on the grammar, the generated decision method uses one of three
// strategies, cheapest first:
//
//  1. A direct check of the next token's type, when a single expected token
//     decides.
//  2. An [LL1Lookahead] table resolved by [ParserState.Lookahead] in O(1),
//     when one token of lookahead decides.
//  3. Adaptive ALL(*) prediction via [ParserState.AdaptivePredict], when
//     arbitrary lookahead is required. The predictor simulates the
//     [RuntimeATN] ahead of the current position and caches its decisions in
//     per-decision DFAs, so hot decisions approach table-lookup cost over
//     time. The cache is internally synchronized and shared across
//     concurrent parses.
//
// Adaptive prediction runs in [PredictionModeSLL] by default;
// [ParserLookahead.SetPredictionMode] enables full-context [PredictionModeLL]
// for grammars that need it.
//
// Cross-reference assignments in the grammar do not resolve anything at parse
// time: the parser stores an unresolved reference created by the generated
// references constructor, and the [typefox.dev/fastbelt/linking] phase
// resolves it later.
//
// # Error recovery
//
// Parse errors do not stop the parse. [Parser] implementations always return
// an AST — with nil tokens or unset fields where input was missing — plus the
// collected [core.ParserError] values, so language tooling keeps working on
// broken input. While the parser is recovering, follow-up errors are
// suppressed so one underlying mistake produces one diagnostic. The pluggable
// [ErrorRecoveryStrategy] controls resynchronization:
//
//   - [DefaultErrorRecovery] attempts single-token deletion when Consume hits
//     an unexpected token, discards tokens before optional/loop guards until
//     one fits the decision or an enclosing follow set (Sync), and skips to
//     the follow set after unwinding from a failed rule (Recover). Because
//     generated AST nodes tolerate nil tokens, missing tokens are reported
//     but never fabricated.
//   - [BailErrorRecovery] stops at the first mismatch and unwinds, for
//     two-stage parsing or when a broken parse result would be discarded
//     anyway.
//
// Message texts come from the pluggable [ErrorMessageProvider]; replace it in
// the service container to change wording or to localize.
//
// # Code completion
//
// The package also contains the parser side of code completion, driven by
// [typefox.dev/fastbelt/server.DefaultCompletionProvider]:
//
//  1. The generated CompletionParser reparses the document's tokens up to the
//     cursor. It mirrors the main parser's control flow but mutates no AST;
//     via [CompletionParserState] it records an [ATNSnapshot] — token index,
//     ATN state, rule stack — at every rule entry and sync point.
//  2. [CompletionParseResult.SimulateAt] picks a suitable snapshot and
//     advances an NFA-style live set over the ATN ([RuntimeATN.Simulate])
//     through the remaining tokens. Unlike prediction, the simulation keeps
//     every alternative alive: the user has not committed to a branch yet.
//  3. [RuntimeATN.NextCompletionsFromSet] reports what may come next as a
//     [CompletionInfo]: valid token types, plus [CompletionHint] entries that
//     mark cross-reference positions so the provider can offer resolvable
//     symbols instead of a bare identifier token.
//
// [LanguageCompletionAdapter] is the contract between the framework's
// completion provider and these generated artifacts; the code generator emits
// its implementation together with the parsers.
package parser
