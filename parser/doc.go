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
//   - `parser_gen.go`: a Parser implementing [Parser], with one ParseX method
//     per parser rule that builds the rule's AST node.
//   - `parser_lookahead_gen.go`: the lookahead service, with one method per
//     decision point in the grammar. It is an interface with a default
//     implementation, so individual decisions can be overridden without
//     touching the parser.
//   - `atn_gen.go`: the [RuntimeATN], an augmented transition network that
//     encodes the grammar as a graph of states. It backs adaptive prediction,
//     error recovery, and completion. Passing --atn to the generate command
//     additionally emits `atn.md`, which renders the ATN as Mermaid diagrams
//     for debugging grammars.
//   - `completion_parser_gen.go`: a stripped-down variant of the parser that
//     builds no AST and is used for code completion.
//
// The generated `SetupGeneratedServices` registers the Parser, the lookahead
// service, and the default [ErrorRecoveryStrategy] and [ErrorMessageProvider]
// in the service container ([typefox.dev/fastbelt/util/service]), each unless
// a custom implementation was registered first.
//
// # Parsing model
//
// A generated parser is a recursive-descent parser: one method per rule,
// descending into sub-rules and consuming tokens left to right. Wherever the
// grammar offers a choice (between alternatives, at an optional element, or
// at a loop continuation), the parser asks its lookahead service. Decisions
// that need more than one token of lookahead use adaptive ALL(*) prediction
// over the [RuntimeATN]. Prediction runs in [PredictionModeSLL] by default;
// [ParserLookahead.SetPredictionMode] enables full-context [PredictionModeLL]
// for grammars that need it.
//
// Cross-reference assignments in the grammar do not resolve anything at parse
// time: the parser stores an unresolved reference, and the
// [typefox.dev/fastbelt/linking] phase resolves it later.
//
// # Error recovery
//
// Parse errors do not stop the parse. [Parser] implementations always return
// an AST, with nil tokens or unset fields where input was missing, plus the
// collected [core.ParserError] values, so language tooling keeps working on
// broken input. Error recovery never inserts synthetic tokens to produce a
// valid AST. The pluggable [ErrorRecoveryStrategy] controls resynchronization:
//
//   - [DefaultErrorRecovery] skips unexpected tokens and resumes parsing as
//     early as possible. This is the right choice for editor tooling.
//   - [BailErrorRecovery] stops at the first mismatch, for two-stage parsing
//     or when a broken parse result would be discarded anyway.
//
// Message texts come from the pluggable [ErrorMessageProvider]; replace it in
// the service container to change wording or to localize.
//
// # Code completion
//
// The generated completion parser and [LanguageCompletionAdapter]
// implementation tell [typefox.dev/fastbelt/server.DefaultCompletionProvider]
// which tokens and cross-references may follow at the cursor position. The
// code generator emits both together with the parsers; adopters normally
// customize completion on the server side rather than here.
package parser
