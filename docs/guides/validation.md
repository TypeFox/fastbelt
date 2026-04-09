# Validation

Lexing and parsing build an AST, but most languages need extra rules (uniqueness, consistency, style) before the model is trustworthy.
In fastbelt, node-level checks are expressed by implementing `fastbelt.Validator` on your generated `*Impl` types.
Document-wide orchestration is done with `workspace.DocumentValidator`, which the default workspace builder runs after linking.

## `Validator` and diagnostics

The core contract lives in the root `fastbelt` package:

```go
type ValidationAcceptor func(diagnostic *Diagnostic)

type Validator interface {
	Validate(ctx context.Context, level string, accept ValidationAcceptor)
}
```

Implementations report issues by calling `accept` with a `*Diagnostic` built via `NewDiagnostic` (see `validation.go`).
Severity constants mirror LSP (`SeverityError`, `SeverityWarning`, `SeverityInfo`, `SeverityHint`).
Optional fields are set with `WithToken`, `WithRange`, `WithCode`, `WithTags`, and `WithData`.
The diagnostic range is chosen with priority: `WithRange`, then `WithToken`, then the node’s text segment.

`ValidationAcceptor` is a simple sink: the default document validator appends each diagnostic to a slice and returns it.
Your `Validate` method should not assume anything about ordering beyond the traversal order described below.

## Validation levels (`level` string)

The `level` argument is meant to distinguish *when* validation runs (the doc comment suggests values such as `"on-type"` and
`"on-save"`). The framework does not define a fixed enumeration; it is a convention between your `Validator` implementation
and whoever calls `DocumentValidator.Validate`.

**Important:** `workspace.DefaultBuilder` is the built-in driver for a full build. In phase 3 it calls:

```go
diagnostics := validator.Validate(ctx, doc, "on-save")
```

So today every validator invoked through the default builder always receives `level == "on-save"`.
If you branch on `level` inside `Validate`, that branch will not run for other levels until something else calls
`DocumentValidator.Validate` with a different string (for example a custom command, test, or future fastbelt API).

## Document-level validation

`workspace.DocumentValidator` validates an entire document:

```go
type DocumentValidator interface {
	Validate(ctx context.Context, doc *fastbelt.Document, level string) []*fastbelt.Diagnostic
}
```

`NewDefaultDocumentValidator` returns `DefaultDocumentValidator`, which:

- Returns `nil` if `doc.Root == nil` (nothing to traverse).
- Walks the AST with `fastbelt.TraverseNode` starting at `doc.Root`, visiting the root and then each descendant in a
  depth-first pre-order (each node, then its subtree).
- For every node that also implements `fastbelt.Validator`, calls `Validate(ctx, level, accept)` where `accept` appends to
  the result slice.
- Respects cancellation: before each node it checks `ctx.Err()` and skips further work if the context is done.

`workspace.CreateDefaultServices` sets `WorkspaceSrv.DocumentValidator` to `NewDefaultDocumentValidator()` when the field is
still nil, same as assigning it explicitly in service setup.

## How the builder runs validation

`workspace.DefaultBuilder.Build` runs three phases. Validation is **phase 3**, after parse, symbol export/import, linking, and
reference descriptions. The builder calls `downgrade()` before phase 3 so the workspace lock can move to a read-friendly mode
while validation runs in parallel per document.

For each document that does not yet have `fastbelt.DocStateValidated`, the builder:

1. Calls `DocumentValidator.Validate(ctx, doc, "on-save")`.
2. Assigns the returned slice to `doc.Diagnostics`.
3. Sets `DocStateValidated` on the document and notifies build-step listeners.

`Builder.Reset` clears `doc.Diagnostics` when the validated bit is cleared.

The LSP entrypoint in `server.StartLanguageServer` registers a listener on `DocStateValidated` that merges lexer, parser, and
linker diagnostics with `doc.Diagnostics` and publishes them to the client.

## Grammar (`.fb`) validators

The fastbelt grammar itself uses the same pattern. Package `internal/grammar` implements `Validate` on `GrammarImpl`
(uniqueness of rule and interface names across rules and terminals), `TokenImpl` (terminal regex must not match the empty
string), and `KeywordImpl` (non-empty, not whitespace-only, warning on embedded whitespace). That code is a reference
implementation alongside the public example.

## Example: statemachine

- [`examples/statemachine/validation.go`](../../examples/statemachine/validation.go) — `StatemachineImpl` implements
  `fastbelt.Validator` and checks unique event and state names and valid transition targets.
- [`examples/statemachine/services.go`](../../examples/statemachine/services.go) — `CreateServices` assigns
  `workspace.NewDefaultDocumentValidator()` to `srv.Workspace().DocumentValidator` (matching what
  `workspace.CreateDefaultServices` would do if left unset).

## Limitations (current behavior)

1. **Single built-in level from the builder** — Only `"on-save"` is passed by `DefaultBuilder`. Finer-grained levels require a
   custom caller of `DocumentValidator.Validate` or future framework support.

2. **No automatic parent/child deduplication** — Every AST node that implements `Validator` is visited once. If both a parent
   and a child implement `Validator`, both run. Avoid duplicate work by implementing `Validate` only where it makes sense.

3. **Per-document scope in the default pipeline** — `Validate` receives `ctx`, `level`, and `accept` only. Nodes can still reach
   the enclosing `*fastbelt.Document` via `AstNode.Document()` and ancestors via `Container()` (and helpers like
   `ContainerOfType`), but the default validator does not pass other documents; cross-file semantic checks need data already
   reflected on the node or on `doc` (for example after linking).

4. **Empty root** — If parsing failed or the root was never set, `DefaultDocumentValidator` returns no diagnostics from
   traversal (parse/link errors are surfaced separately).

5. **Context** — Traversal stops checking new validators after cancel, but diagnostics collected so far are still returned.
