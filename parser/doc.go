// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

// Package parser provides the shared parser runtime used by Fastbelt-generated
// languages. It does not implement parsing for any particular language; each
// grammar gets its own generated recursive-descent parser that drives the
// parser state, adaptive prediction, error recovery, and completion machinery
// of this package.
//
// Run [typefox.dev/fastbelt/cmd/fastbelt] generate on a grammar directory (or
// any .fb file in it) to emit the generated sources. See [typefox.dev/fastbelt]
// for the overall toolchain and [typefox.dev/fastbelt/grammar] for how parser
// rules are declared in the grammar language.
//
// # Generated parsers
//
// Among the generated files, the parser-related ones are:
//
//   - `parser_gen.go`: a `Parser` type implementing [Parser], with one ParseX
//     method per parser rule that builds the rule's AST node.
//   - `parser_lookahead_gen.go`: the lookahead service, with one method per
//     decision point in the grammar. It is an interface with a default
//     implementation, so individual decisions can be overridden without
//     touching the parser.
//   - `atn_gen.go`: the [RuntimeATN] used for prediction, error recovery, and
//     completion. Passing --atn to the generate command also emits `atn.md`,
//     which renders the ATN as Mermaid diagrams for debugging grammars.
//   - `completion_parser_gen.go`: a variant of the parser that builds no AST
//     and is used for code completion.
//   - `completion_gen.go`: the [LanguageCompletionAdapter] implementation and
//     the `<Grammar>CompletionFilter`, registered by the generated
//     `SetupGeneratedServerServices`.
//
// The generated `SetupGeneratedServices` registers, among others, the Parser,
// the lookahead service, and the default [ErrorRecoveryStrategy] and
// [ErrorMessageProvider] in the service container
// ([typefox.dev/fastbelt/util/service]), each unless a custom implementation
// was registered first.
//
// # Parsing model
//
// The generated parser asks its lookahead service wherever the grammar offers
// a choice (between alternatives, at an optional element, or at a loop
// continuation). Decisions that need more than one token of lookahead use
// adaptive ALL(*) prediction over the [RuntimeATN]. Prediction runs in
// [PredictionModeSLL] by default; [ParserLookahead.SetPredictionMode] enables
// full-context [PredictionModeLL] for grammars that need it.
//
// Cross-reference assignments in the grammar do not resolve anything at parse
// time: the parser stores an unresolved reference, and the
// [typefox.dev/fastbelt/linking] phase resolves it later.
//
// # Error recovery
//
// Generated parsers always produce an AST in [core.Document.Root], with nil
// tokens or unset fields where input was missing, and record
// [core.ParserError] values in [core.Document.ParserErrors]. With
// [DefaultErrorRecovery], parse errors do not stop the parse, so language
// tooling keeps working on broken input. Error recovery never inserts
// synthetic tokens to produce a valid AST. The pluggable
// [ErrorRecoveryStrategy] controls resynchronization:
//
//   - [DefaultErrorRecovery] skips extraneous tokens, leaves missing tokens
//     unset, and resumes where parsing can continue. This is the right choice
//     for editor tooling.
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
// generated `<Grammar>CompletionFilter` can narrow the cross-reference
// candidates; the server's [typefox.dev/fastbelt/server.CompletionContributor]
// decides which items are emitted and how they look.
package parser
