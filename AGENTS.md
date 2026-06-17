# Fastbelt for Agents

This is a Go project.

For project documentation, read the package documentation written in Go doc format:

- `doc.go`: General overview of the Fastbelt framework
- `grammar/doc.go`: Reference documentation for the grammar language
- `lexer/doc.go`: Shared lexer runtime used by Fastbelt-generated languages
- `textdoc/doc.go`: Text documents, overlays, and LSP position mapping
- `linking/doc.go`: Cross-reference resolution, symbol tables, and scopes
- `workspace/doc.go`: Document lifecycle, loading, edits, and the build pipeline

Utility packages:

- `util/codegen/doc.go`: Building indented multi-line source for code generators
- `util/collections/doc.go`: Generic collection data structures (e.g. MultiMap)
- `util/extiter/doc.go`: Utilities for `iter.Seq` sequences
- `util/service/doc.go`: Typed dependency injection container

## VS Code extensions

`internal/vscode-extensions` is a subfolder containing a TypeScript project with VS Code extensions for the Fastbelt grammar language and for some example languages.

This subproject is built with npm; see `internal/vscode-extensions/package.json`.
