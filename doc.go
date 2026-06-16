// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

// Package fastbelt is a language engineering toolkit for Go.
//
// # Overview
//
// Fastbelt covers lexing, parsing, AST creation, cross-reference linking,
// workspace handling, and Language Server Protocol (LSP) support.
//
// Fastbelt is primarily inspired by [Xtext] (Java)
// and [Langium] (TypeScript).
//
// Like these frameworks, Fastbelt uses a grammar definition file as the
// entry point, based on a Grammar language that is itself implemented with
// Fastbelt. See [typefox.dev/fastbelt/grammar] for details.
//
// The [typefox.dev/fastbelt/cmd/fastbelt] command generates Go code from a
// grammar definition, which you can integrate and customize with additional
// code.
//
// # Defining the Grammar
//
// Defining a language in Fastbelt is an iterative design process. You
// co-design the concrete syntax that users write and the abstract syntax tree
// (AST) that your tooling consumes.
//
// In practice, you write interface declarations for AST node types and parser
// rules that describe how instances of those node types are created from input
// text. Each assignment in a parser rule maps matched input to a field on the
// node that the rule builds.
//
// Parser rules always consume contiguous regions of text. They match keywords
// and tokens directly and delegate to other rules for nested regions. This
// keeps grammar structure aligned with language structure and makes it clear
// how parsed text becomes AST subtrees.
//
// Cross-references provide built-in linking from symbol usage to symbol
// definition. Instead of creating a new node, a cross-reference records an
// identifier that the linker resolves to an existing AST node in scope.
//
// Because Go interfaces are generated directly from grammar interfaces, these
// interfaces should reflect how you plan to implement semantic checks,
// interpreters, or code generators. At the same time, there's a close
// correspondence between interfaces and parser rules so the parse model and
// runtime model stay easy to reason about.
//
// The grammar language itself is documented in
// [typefox.dev/fastbelt/grammar], while the code generation workflow is
// documented in [typefox.dev/fastbelt/cmd/fastbelt]. A common workflow is to
// evolve the grammar in small steps and run code generation after each change.
//
// # Customization
//
// There are two ways to add custom behavior to your language:
//
//   - Attach behavior directly to specific generated AST node `Impl` types.
//   - Implement custom services and wire them through the service container.
//
// Validations are implemented with the first approach. To add validation checks
// for a node type, implement [Validator] on that node's generated `Impl`
// struct:
//
//	type PersonImpl struct {
//	    PersonData
//	}
//
//	func (p *PersonImpl) Validate(_ context.Context, _ string, accept fastbelt.ValidationAcceptor) {
//	    if name := p.Name(); name != "" && !unicode.IsUpper([]rune(name)[0]) {
//	        accept(fastbelt.NewDiagnostic(
//	            fastbelt.SeverityError,
//	            "Name must start with an uppercase letter.",
//	            p,
//	            fastbelt.WithToken(p.NameToken()),
//	        ))
//	    }
//	}
//
// This keeps checks close to the node they validate and works naturally with
// generated AST types.
//
// For service-container based customization, define a language-specific
// `SetupServices` function and a `CreateServices` helper. The service container
// API is documented at [typefox.dev/fastbelt/util/service].
//
// By convention, `SetupServices`:
//   - registers custom values and service implementations,
//   - calls needed `SetupDefaultServices` functions from framework packages, and
//   - calls `SetupGeneratedServices` from generated code.
//
// Scaffolded language projects include a generated `services.go` that follows
// this pattern and can be used as a starting point for service-container
// customization.
//
// # To Go Further
//
// For implementing the exact behavior of your language, these subpackages are
// the most relevant:
//
//   - [typefox.dev/fastbelt/linking] for everything related to
//     cross-reference linking, including scopes and making symbols reachable
//     between files.
//   - [typefox.dev/fastbelt/workspace] for handling a collection of files in a
//     workspace folder, processing file changes by executing phases such as
//     parsing, linking, and validating.
//   - [typefox.dev/fastbelt/server] for everything related to the Language
//     Server Protocol (LSP).
//
// [Xtext]: https://eclipse.dev/Xtext/
// [Langium]: https://langium.org
package fastbelt
