# Consumption

This guide explains how to drive the generated lexer, parser, linker, and validators from Go using the same
workspace APIs a language server uses. It focuses on ideas, not a full compiler architecture.

The walkthrough follows the small CLI in
[`examples/statemachine/cmd/statemachine`](../../examples/statemachine/cmd/statemachine). Sample models live next
to the language, for example
[`examples/statemachine/traffic_light.statemachine`](../../examples/statemachine/traffic_light.statemachine) and
[`examples/statemachine/elevator.statemachine`](../../examples/statemachine/elevator.statemachine).

The command package is split to mirror how a real tool might separate concerns. All of these are methods on
[`Runner`](../../examples/statemachine/cmd/statemachine/runner.go) (same `main` package, different files for
readability):

- **`parser.go`** — `ParseArgs`, `LoadSource`, and `ParseStatemachine` (workspace pipeline: lex, parse, link,
  validate) plus unexported diagnostic helpers. This is the “compiler front end” slice: text in, document and
  typed root stored on the runner, or errors.
- **`interpreter.go`** — `PrintModelSummary` and `Interpret`: walk the linked AST, match event lines from
  `EventInput` to transitions, and print demo `emit command "…"` lines for `actions` that reference commands.
  No separate runtime exists in fastbelt; this is illustrative execution logic only.
- **`runner.go`** — `Runner` fields (`Stdout`, `Stderr`, `EventInput`, …) plus `ParseAndValidate` (wraps
  `ParseStatemachine`) and `Run` (diagnostics, summary, then `Interpret`).
- **`main.go`** — constructs a `Runner`, calls `ParseArgs`, `LoadSource`, `ParseAndValidate`, and `Run`, and exits on failure.

End to end: load a `.statemachine` file from disk, run a full build cycle, print diagnostics, read the root
AST, then step the machine by reading **event names** from stdin. The grammar only attaches names to commands
in an `actions` block; the grammar does not define real side effects. The example CLI lists declared commands in
the summary; for the elevator model it also prints `emit command "…"` lines after each transition into a state
whose `actions` block references commands (for example `bell` when returning to `Waiting`).

## The statemachine language (AST sketch)

The grammar lives in [`examples/statemachine/statemachine.fb`](../../examples/statemachine/statemachine.fb).
Generated types are in `types_gen.go`: root interface `Statemachine` with `Name`, `Events`, `Commands`,
`Init` (reference to a `State`), and `States`. Each `State` has optional `actions` (`[]*Reference[Command]`)
and `Transitions` (`Event` ref `=>` `State` ref). After linking, use `Reference.Ref(ctx)` to obtain the
target node interface.

Custom validation is implemented on `StatemachineImpl` in
[`examples/statemachine/validation.go`](../../examples/statemachine/validation.go) (unique event/state names)

## Build and run the example CLI

From the repository root:

```bash
go build -o statemachine ./examples/statemachine/cmd/statemachine/
./statemachine examples/statemachine/traffic_light.statemachine
```

The program requires exactly one argument, the path to the model file. After printing the summary it prompts on
stderr and reads **one event name per line** from stdin. It walks transitions by matching the text of the event
reference on each outgoing transition from the current state. That is the only behavior the AST guarantees;
there is no guard/action language to interpret.

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
file on disk, `core.FileURI(absPath)` then `.DocumentURI()` matches what the LSP stack and the example CLI use.

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
The example interpreter (`(*Runner).Interpret` in `interpreter.go`) resolves `Init` to the first `State`, then
for each line matches `Event().Text()` against the user’s input and follows `State().Ref(ctx)`.

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

## Related pieces

- LSP server entrypoint: [`examples/statemachine/server/main.go`](../../examples/statemachine/server/main.go)
  embeds the same `StatemachineSrv` and adds `server` services. The [Language server](language-server.md) guide explains that wiring.
- Code generator driver (grammar → Go): [`cmd/fastbelt/main.go`](../../cmd/fastbelt/main.go) runs the same `Builder.Build` loop for
  `.fb` grammar files before emitting `*_gen.go`.
