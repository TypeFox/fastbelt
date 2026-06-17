// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

// Package linking resolves cross-references in parsed documents and builds the
// symbol tables and scopes that reference resolution depends on. The
// [typefox.dev/fastbelt/workspace] builder orchestrates these steps during
// document processing.
//
// # Services
//
// [SetupDefaultServices] registers the framework services:
//
//   - [SymbolExporter] — symbols this document exposes to other documents
//   - [SymbolImporter] — symbols visible from other documents
//   - [LocalSymbolsProvider] — symbols declared within this document
//   - [Linker] — resolves cross-references in the AST
//   - [ReferenceDescriptionsProvider] — metadata for resolved references
//   - [ReferenceDescriber] — describes individual references for the provider
//
// Code generation emits additional linking services in linker_gen.go:
//
//   - <Grammar>ReferencesConstructor — creates Reference values while parsing
//   - <Grammar>ReferenceLinker — resolves a single cross-reference
//   - <Grammar>ScopeProvider — supplies the candidate scope for each
//     cross-reference field
//
// # Build process
//
// The workspace builder runs linking-related work across two phases.
// The parse phase (build phase 1) runs two steps per document:
//
//   - The parser constructs the AST. For each cross-reference assignment, the
//     parser calls the generated ReferencesConstructor, which attaches a typed
//     Reference to the owner node. Each reference stores a resolve function that
//     delegates to the generated ReferenceLinker; resolution is deferred until
//     the link phase.
//   - [SymbolExporter] traverses the AST and fills document.ExportedSymbols
//     for each workspace document.
//
// Every document in the workspace must complete phase 1 before the next phase
// starts. The link phase (build phase 2) runs four steps per document:
//
//   - [SymbolImporter] merges exported symbols from workspace documents
//     into document.ImportedSymbols.
//   - [LocalSymbolsProvider] records per-container local symbols in
//     document.LocalSymbols.
//   - [Linker] traverses the AST and calls Resolve on every reference. Each
//     call reaches the ReferenceLinker's Link* method, which asks the
//     ScopeProvider for a scope and takes the first candidate.
//   - [ReferenceDescriptionsProvider] indexes resolved references on the
//     document for faster lookup of symbol usages.
//
// By default, each Scope* method on Default<Grammar>ScopeProvider returns
// [DefaultScopeOfType] for the reference's target type, combining imported
// symbols ([GlobalScopeOfType]) with symbols visible from enclosing containers
// ([LocalScopeOfType]).
//
// # Customization
//
// Linking can be customized at two levels. To adjust how individual AST node
// types participate in linking, implement these interfaces on the node's
// generated Impl struct:
//
//   - [Denominator] — custom naming logic; used by [Name]
//   - [ExportedSymbolDescriber] — custom export description for a node; used by
//     [DescribeExport]
//   - [LocalSymbolDescriber] — custom local symbol description and placement;
//     used by [DescribeLocal]
//
// To change a linking strategy across your whole language, replace one of the
// services:
//
//   - [SymbolExporter] — selects which symbols a document exposes to other
//     documents. The default exports the root node and its directly named
//     children; override it to honor an explicit export keyword, for example.
//   - [SymbolImporter] — selects which documents contribute imported symbols.
//     The default merges exports from every workspace document; override it to
//     follow explicit import statements or package boundaries.
//   - [LocalSymbolsProvider] — defines intra-document visibility. By default a
//     symbol is visible within its enclosing container, so visibility follows
//     the AST structure; override it for custom visibility rules.
//   - [Linker] — drives the resolution pass over all references. The default
//     resolves each reference to the first matching candidate in its scope;
//     override it for custom candidate selection, such as function overloading.
//   - <Grammar>ScopeProvider — computes the candidate scope for a
//     cross-reference field. Override Scope* methods on an embedded
//     Default<Grammar>ScopeProvider for custom scope logic, for example to
//     resolve struct or class members after a "." operator.
package linking
