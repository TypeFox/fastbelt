// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

// Package lexer provides the shared lexer runtime used by Fastbelt-generated
// languages. It does not implement lexing for any particular language; each
// grammar gets its own generated lexer that wires this package together with
// language-specific [core.TokenType] values.
//
// Run [typefox.dev/fastbelt/cmd/fastbelt] generate on a .fb grammar file to
// emit `lexer_gen.go` (token type variables plus `NewLexer`) and register that
// lexer in the generated service setup. See [typefox.dev/fastbelt] for the
// overall toolchain and [typefox.dev/fastbelt/grammar] for how keywords and
// token rules are declared in the grammar language.
//
// # Generated lexers
//
// Code generation turns grammar terminals into Go values:
//
//   - Keyword literals from parser rules become [core.TokenType] values with
//     prefix [core.TokenMatcher] functions and [core.TokenKindKeyword].
//   - Named `token` rules become [core.TokenType] values whose matchers are
//     compiled from the rule's regular expression, with StartChars taken from
//     the expression's possible first runes.
//
// The generated `NewLexer` function returns a [DefaultLexer] constructed via
// [NewDefaultLexer] with one [TokenMode] per `token mode` declaration in the
// grammar. [typefox.dev/fastbelt/workspace.DefaultDocumentParser] obtains a
// [Lexer] from the service container, calls [Lexer.Exec], and stores
// [LexerResult.Tokens], [LexerResult.Comments], and [LexerResult.Errors] on the
// document before parsing.
//
// # Lexing model
//
// Fastbelt uses a table-driven scanner with the following approach:
//
//  1. At each input offset, decode the current rune and look up candidate token
//     types in the active mode's fixed-size map keyed by rune % 256 (built from
//     each [core.TokenType]'s StartChars slice).
//  2. Run each candidate's [core.TokenMatcher] and keep the longest match
//     (maximal munch). Among equal-length matches, the first registered token
//     type wins; generated lexers list keywords before regex token rules, so
//     keywords take precedence when both match the same span.
//  3. Route the match by [TokenTypeUsage.Modifier]: default tokens go to
//     [LexerResult.Tokens], hidden tokens are dropped, comments go to
//     [LexerResult.Comments], and other modifiers are collected in
//     [LexerResult.Modifiers].
//  4. Apply the match's mode command, if any (see below).
//  5. If no token type matches, emit a [core.LexerError] and advance by one
//     UTF-8 code point so lexing can continue. The active mode is unchanged.
//
// [DefaultLexer] adapts its initial token-slice capacity from a running average
// of tokens per byte across prior [Lexer.Exec] calls on the same instance.
//
// # Token modes
//
// A [TokenMode] is a named set of token types. Only the active mode's tokens are
// candidates at a given offset, which is how a language can lex, say, string
// interpolation without an ambiguous grammar. Each [Lexer.Exec] starts in the
// mode passed to [NewDefaultLexer] and tracks the active one in its own
// [TokenModeStack], so a single [DefaultLexer] is safe for concurrent use and
// input that ends inside a pushed mode cannot affect the next call.
//
// A token type registered with a mode may carry one command:
//
//   - push: make another mode active, remembering the current one.
//   - pop: return to the mode remembered by the matching push. A pop with
//     nothing to return to keeps the start mode active rather than failing.
//   - mode: replace the active mode without remembering it, so a later pop
//     returns to whatever was below it rather than to the replaced mode.
//
// Custom language projects rarely import this package directly unless they
// replace the generated lexer. Typical integration is to call the generated
// `SetupGeneratedServices` which registers [Lexer].
package lexer
