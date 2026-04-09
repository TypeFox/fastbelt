# Consumption

This guide explains how to drive the generated lexer, parser, linker, and validators from Go using the same
workspace APIs a language server uses. It focuses on ideas, not a full compiler architecture.

The walkthrough follows the small CLI in
[`examples/statemachine/cmd/statemachine`](../../examples/statemachine/cmd/statemachine): load a
`.statemachine` document, run a full build cycle, print diagnostics, read the root AST, and (when input
comes from a file) step the machine by reading **event names** from stdin. The grammar only attaches names to
commands in an `actions` block; it does not define runtime command behavior, so the example lists commands but
does not “execute” them.

## The statemachine language (AST sketch)

The grammar lives in [`examples/statemachine/statemachine.fb`](../../examples/statemachine/statemachine.fb).
Generated types are in `types_gen.go`: root interface `Statemachine` with `Name`, `Events`, `Commands`,
`Init` (reference to a `State`), and `States`. Each `State` has optional `actions` (`[]*Reference[Command]`)
and `Transitions` (`Event` ref `=>` `State` ref). After linking, use `Reference.Ref(ctx)` to obtain the
target node interface.

Custom validation is implemented on `StatemachineImpl` in
[`examples/statemachine/validation.go`](../../examples/statemachine/validation.go) (unique event/state names,
transition targets).

## Build and run the example CLI

From the repository root:

```bash
go build -o statemachine ./examples/statemachine/cmd/statemachine/
./statemachine path/to/model.statemachine
```

With no arguments, or with `-`, the program reads the **document body** from stdin. In that mode it prints a
summary only; it does not read events from stdin (stdin was already consumed).

When the path is a file, after the summary it prompts on stderr and reads **one event name per line** from
stdin. It walks transitions by matching the text of the event reference on each outgoing transition from the
current state. That is the only behavior the AST guarantees; there is no guard/action language to interpret.

## Service container setup

[`examples/statemachine/services.go`](../../examples/statemachine/services.go) wires the stack:

- Embed the generated and workspace service blocks (`TextdocSrvContBlock`, `GeneratedSrvContBlock`, etc.).
- Call `CreateServices()` which registers defaults (`textdoc`, `workspace`, `linking`, generated services) and
  sets `LanguageID` and `FileExtensions` (here `statemachine` and `.statemachine`).
- Assign `DocumentValidator`: the example uses `workspace.NewDefaultDocumentValidator()`, which traverses the
  AST and calls `Validate` on every node that implements `fastbelt.Validator` (your `StatemachineImpl`).

Production code typically mirrors this layout: one constructor that returns a concrete struct embedding the
generated container blocks, with language metadata filled in.

## Creating a document

Documents are backed by a `textdoc.Handle`. For a file-like snapshot, `textdoc.NewFile(uri, languageId,
version, content)` builds an immutable buffer. The URI should be stable and unique within the workspace; for a
real path, `core.FileURI(absPath)` then `.DocumentURI()` matches what the LSP stack uses. For stdin-only runs,
use a synthetic path (the example uses `file:///stdin.statemachine`) so the URI is still non-empty.

Wrap the handle in `core.NewDocument(handle)`. That value is what the builder mutates (`Root`, `Tokens`,
errors, `Diagnostics`, etc.).

## Registering and building

Put the document in the manager, then run the builder under the workspace lock (same pattern as
[`examples/statemachine/benchmark_test.go`](../../examples/statemachine/benchmark_test.go) and the LSP document
updater):

```go
srv.Workspace().DocumentManager.Set(document)
srv.Workspace().Lock.Write(ctx, func(ctx context.Context, downgrade func()) {
    _ = srv.Workspace().Builder.Build(ctx, []*core.Document{document}, downgrade)
})
```

`Builder.Build` runs parse, symbol phases, link, reference descriptions, then validation. Calling `downgrade()`
inside `Build` switches from the exclusive write phase to a shared read phase before validation finishes; the
lock implementation ensures this ordering. For a single-document CLI, `Write` is still the right entry point
so you do not race the same assumptions as the server.

## Reading parse, link, and validation diagnostics

After `Build` returns, inspect the document:

- **Lexer / parser:** `doc.LexerErrors`, `doc.ParserErrors`. Helpers
  `workspace.CreateLexerDiagnostics` and `workspace.CreateParserDiagnostics` turn these into `core.Diagnostic`
  values with ranges suitable for printing or LSP.
- **Linker:** unresolved references appear as `ref.Error()` on entries in `doc.References`.
  `workspace.CreateLinkerDiagnostics` maps those to diagnostics using the reference text range.
- **Custom validation:** phase 3 fills `doc.Diagnostics` with pointers from your `Validator` implementations.

The example CLI prints all four sources with a short `lexer|parser|linker|validate` prefix and 1-based line /
column labels for readability (core stores LSP-style 0-based positions).

Treat lexer/parser failures and linker errors as blocking: the AST may be missing or references may not
resolve. Your tool can still pretty-print partial trees for debugging, but evaluation should wait until the
pipeline reports a consistent state you define (often: no lexer/parser errors, no reference errors, and no
validation errors at `SeverityError`).

## Reading the root AST

After a successful parse, `document.Root` implements `fastbelt.AstNode`. For a known grammar, type-assert to
your root interface, e.g. `document.Root.(statemachine.Statemachine)`. Generated interfaces live in
`types_gen.go`; use those for traversal instead of concrete structs when you can.

References are not usable as linked targets until after the link step. In application logic, call
`Ref(ctx)` (or `RefNode`) on `*core.Reference[T]` with a `context.Context` when you need the target node.
The example
simulator resolves `Init` to the first `State`, then for each transition matches `Event().Text()` against the
user’s line and follows `State().Ref(ctx)`.

## Where to hook evaluation or codegen

- **Per-node validation:** implement `fastbelt.Validator` on generated `*…Impl` types (see `validation.go`).
  The default document validator discovers them by traversal.
- **Cross-cutting passes:** register `workspace.Builder.AddBuildStepListener` with a mask such as
  `core.DocStateLinked` if you need to run analysis after references resolve but before or alongside
  validation.
- **Evaluation:** there is no interpreter in the framework; walk the typed AST (and resolved references) from
  `document.Root` after `Build` completes. Keep behavior aligned with what the grammar actually defines.

## Error handling patterns

- Propagate I/O and setup errors from `ReadFile`, `NewFile`, and `Builder.Build` (`context.Cancellation` shows
  up here in long-running servers).
- Separate **infrastructure** failures from **document** failures: the latter are usually a non-empty
  diagnostic list rather than a returned `error` from `Build`.
- When reporting to users, prefer diagnostic messages and ranges from the helpers above so lexer, parser,
  linker, and validator output share one shape.
- Avoid shadowing `err` in tight scopes; use a distinct name (`rerr`, `ferr`) when another error is still in
  scope, matching common style in this repository.

## Related pieces

- LSP server entrypoint: [`examples/statemachine/server/main.go`](../../examples/statemachine/server/main.go)
  embeds the same `StatemachineSrv` and adds `server` services. The [Language server](language-server.md) guide explains that wiring.
- Code generator driver (grammar → Go): [`cmd/fastbelt/main.go`](../../cmd/fastbelt/main.go) runs the same `Builder.Build` loop for
  `.fb` grammar files before emitting `*_gen.go`.
