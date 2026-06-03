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
//     prefix [core.Matcher] functions and [core.TokenKindKeyword].
//   - Named `token` rules become [core.TokenType] values whose matchers are
//     compiled from the rule's regular expression, with [core.StartChars] taken
//     from the expression's possible first runes (the StartChars field).
//
// The generated `NewLexer` function returns a [DefaultLexer] constructed via
// [NewDefaultLexer] with every keyword and token type for that language.
// [typefox.dev/fastbelt/workspace.DefaultDocumentParser] obtains a [Lexer] from
// the service container, calls [Lexer.Lex], and stores [LexerResult.Tokens],
// [LexerResult.Comments], and [LexerResult.Errors] on the document before
// parsing.
//
// # Lexing model
//
// Fastbelt uses a table-driven scanner with the following approach:
//
//  1. At each input offset, decode the current rune and look up candidate token
//     types in a fixed-size map keyed by rune % 256 (each [core.TokenType]'s
//     StartChars slice).
//  2. Run each candidate's [core.Matcher] and keep the longest match (maximal
//     munch). Keywords and regex token rules compete on equal footing.
//  3. Route the match by [core.TokenType.Group]: default tokens go to
//     [LexerResult.Tokens], hidden tokens are dropped, comments go to
//     [LexerResult.Comments], and other groups are collected in
//     [LexerResult.Groups].
//  4. If no token type matches, emit a [core.LexerError] and advance by one
//     UTF-8 code point so lexing can continue.
//
// [DefaultLexer] tracks line and column while scanning and adapts its initial
// token-slice capacity from a running average of tokens per byte across prior
// [Lexer.Lex] calls on the same instance.
//
// Custom language projects rarely import this package directly unless they
// replace the generated lexer. Typical integration is to call the generated
// `SetupGeneratedServices` which registers [Lexer].
package lexer
