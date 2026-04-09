# Concepts

This guide explains how the `typefox.dev/fastbelt` module hangs together:
documents, ASTs, the service-container wiring, codegen, the build pipeline, and where you extend behavior.
It assumes you already know what lexing and parsing mean in general; it focuses on fastbelt-specific structure.

For grammar syntax, see the [grammar reference](../references/grammar.md).
For bootstrapping a new module or package and for day-to-day codegen, see [Scaffolding](../guides/scaffolding.md).
For validation details and patterns, see [Validation](../guides/validation.md).
For integrating generated code into tools or servers, see [Consumption](../guides/consumption.md).
For a minimal LSP process (stdio, diagnostics, defaults), see [Language server](../guides/language-server.md).

## Documents and text

A **document** is a `fastbelt.Document` bound to a **text handle** (`textdoc.Handle`).
The handle exposes URI, language id, version, and text (full buffer or subranges via LSP-style ranges).

`NewDocument` attaches the handle and initializes empty slices for tokens, errors, references, and diagnostics.
`Root` may be nil until parsing has run.
The document embeds `sync.RWMutex`; callers synchronize access to its fields.

Adopters may stash arbitrary per-document data in `Document.Data` (`sync.Map`).
The builder does not clear `Data` during builds; you own lifecycle and consistency with edits.

**Document state** is a bitmask (`DocumentState`) recording which build phases completed, for example:
parsed, exported symbols, imported symbols, local symbols, linked, references resolved, validated.
The default document updater uses partial resets so some phases can be skipped on incremental updates
while still re-running import, link, reference metadata, and validation when text changes.

## AST roots and nodes

The parser stores its result in `Document.Root` as an `AstNode`.
Concrete node types are generated from your `.fb` grammar; the grammar’s entry rule defines the shape of the root.

`AstNode` is the interface all nodes implement: document and container pointers, tokens, text span, and tree walks
(`ForEachNode`, `ForEachReference`).
Embedding `AstNodeBase` (via generated code) provides the usual defaults; named constructs often implement `NamedNode`.

Cross-file and in-file name resolution builds on `UntypedReference` and typed `Reference[T]` values collected during linking.
Generated code constructs these references; linking services resolve them against symbol tables.

## Service container pattern

Fastbelt avoids a global runtime.
Instead you define a **service struct** that embeds small **`*SrvContBlock` structs** from `textdoc`, `workspace`, `linking`,
and the **generated** linking block for your language.

Each block exposes a method like `Workspace() *WorkspaceSrv` or `Generated() *GeneratedSrv`.
A **CreateServices** function allocates the struct, optionally overrides specific fields, then calls `CreateDefaultServices` helpers
in dependency order (`textdoc` → `workspace` → `linking` → your language’s `CreateDefaultServices` from `services_gen.go`).

`WorkspaceSrvCont` requires:

- `textdoc.TextdocSrvCont` (document store),
- `linking.LinkingSrvCont` (generic linking services),
- `workspace.GeneratedSrvCont` (`Lexer` and `Parser` for the language),
- `workspace.Workspace()` (document manager, updater, lock, builder, parser, validator, initializer).

The generated `*GeneratedSrvCont` interface extends `workspace.GeneratedSrvCont` with your language-specific linking surface
(scope provider, reference linker, references constructor) so generated parser and linker code can reach those implementations.

## Generated code versus hand-written code

Running the generator emits several `*_gen.go` files into the output directory (default: current package):

- `lexer_gen.go` — lexer implementation
- `parser_gen.go` — parser and AST node types usage
- `types_gen.go` — AST structs and accessors
- `linker_gen.go` — scope, reference linking, and reference construction
- `services_gen.go` — `CreateDefaultServices` for generated lexer/parser and language linking services

You keep **hand-written** files in the same package, typically:

- `services.go` — your `CreateServices`, wiring, `LanguageID`, `FileExtensions`, and optional overrides
- `validation.go` (or similar) — `Validator` implementations on AST nodes
- Any custom linking or scope behavior if you replace generated defaults

The `internal/grammar` package is the canonical bootstrap: it uses `grammar.fb` and checked-in generated files to parse `.fb` grammars.
The [statemachine example](../../examples/statemachine/) shows the same pattern for a sample language.

## Code generation CLI and `go:generate`

The `fastbelt` command lives in [`cmd/fastbelt`](../../cmd/fastbelt). It has two main roles:

1. **Default mode (no subcommand)** — compile a **grammar** `.fb` file into the five `*_gen.go` outputs.
2. **`fastbelt scaffold`** — lay down a new Go package or module with a starter grammar, `go:generate`, LSP `main`, and optional VS Code
   extension layout, then run `go generate` and `go mod tidy` (see [Scaffolding](../guides/scaffolding.md)).

Generate mode **only** targets the fastbelt grammar language: it does not start a language server and does not parse documents in *your*
DSL. Subcommands recognized before flag parsing are `help` and `scaffold`; anything else is treated as **generate** flags.

**Generate** flags:

- `-g` — path to the grammar file (default `./grammar.fb`)
- `-o` — output directory for generated Go files (default `./`)
- `-p` — Go package name (default: last segment of the output path)
- `-v` — print each file path as it is written

Before writing files, the tool loads the grammar text, uses `internal/grammar.CreateServices`, registers a `fastbelt.Document`, and runs
the workspace **Builder** on that document. Grammar **diagnostics** (lexer, parser, linker, and AST validators for `.fb`) are printed to
stdout; if any **error**-severity diagnostic is present, generation **stops** and no `*_gen.go` files are written. When the build is clean,
`document.Root` must be a `grammar.Grammar`, and the five files are emitted via `internal/generator`.

Typical regeneration from a package inside this repository:

```go
//go:generate go run ../../cmd/fastbelt -g ./statemachine.fb -v
```

In a consumer module, prefer `go tool` after `go get -tool typefox.dev/fastbelt/cmd/fastbelt@latest` (see the root `README.md`):

```go
//go:generate go tool typefox.dev/fastbelt/cmd/fastbelt -g ./grammar.fb -o . -p mypkg -v
```

Paths in `//go:generate` are relative to the file that contains the directive; adjust `-g` and the tool invocation to match your layout.

## Builder pipeline: parse, symbols, link, validate

`workspace.DefaultBuilder.Build` runs three phases (with cancellation checks between steps).

**Phase 1 (per document, parallel):**

1. **Parse** — `DocumentParser.Parse` runs `Generated().Lexer.Lex` then `Generated().Parser.Parse`, filling tokens, lexer/parser errors,
   and `Root`.
2. **Exported symbols** — `ExportedSymbolsProvider` records what each document exports for cross-document resolution.

**Phase 2 (per document, parallel, after all exports exist):**

1. **Imported symbols** — resolves imports against other documents in the manager.
2. **Local symbols** — intra-document symbol tables.
3. **Link** — `Linker` resolves references and populates `Document.References`.
4. **Reference descriptions** — human-readable or IDE-oriented metadata for references.

Then the workspace lock **downgrades** from exclusive to readable so clients can read while validation runs.

**Phase 3 (per document, parallel):**

1. **Validate** — `DocumentValidator.Validate` runs; results are stored in `Document.Diagnostics`.

You can register **build step listeners** on the builder to run hooks after specific state bits are reached
(listeners receive context and the document; errors are logged and do not fail the build).

## The `linking` and `parser` / `lexer` packages

- **`lexer.Lexer`** — `Lex(text string) *LexerResult` (tokens, errors, optional groups).
- **`parser.Parser`** — `Parse(document *Document) *ParseResult` (root node and parser errors).

**`linking.LinkingSrv`** holds pluggable providers: exported/imported/local symbols, reference descriptions, naming, and the linker.
`linking.CreateDefaultServices` installs default implementations if fields are nil.

Generated linker code depends on your language’s generated linking service struct for scope and per-reference link functions.

## Extension points

**Workspace / document level**

- **`DocumentValidator`** — replace or wrap `workspace.NewDefaultDocumentValidator()`.
  The default traverses the AST and calls `Validate` on any node that implements `fastbelt.Validator`.
- **`DocumentParser`** — uncommon; swap if you need a custom parse pipeline (defaults use generated lexer/parser).
- **`Builder`** — replace for alternative orchestration; must honor the `downgrade` contract described on `Builder.Build`.
- **`DocumentUpdater`** — coordinates builds after text changes; the default cancels in-flight work and applies partial `Reset`.
- **`Initializer`** — workspace folder discovery (used when the LSP client sends `initialized`).
- **`BuildStepListener`** — observe phase completion without replacing the builder.

**Per-node validation**

- Implement `Validator` on AST `Impl` types: `Validate(ctx context.Context, level string, accept ValidationAcceptor)`.
  The `level` argument identifies when validation runs (the default builder passes `"on-save"`).
  Use `NewDiagnostic` and options like `WithToken` / `WithRange` to attach messages to source.

**Generated linking services**

- Override `ScopeProvider`, `ReferenceLinker`, or `ReferencesConstructor` on your language linking struct before or after
  `CreateDefaultServices` if you need custom scoping or resolution (see grammar’s scope provider override for a precedent).

**Language server**

- **`server.ServerSrv`** — `DocumentSyncher`, `DefinitionProvider`, `ReferencesProvider`, logging handler, JSON-RPC connection hooks.
  `CreateDefaultServices` fills defaults, including a syncher that forwards LSP notifications to `DocumentUpdater`.

Capabilities advertised in `Initialize` include incremental text sync, completion (stub empty list), and definition/references
when the corresponding providers are non-nil.

## LSP integration at a high level

A minimal server composes your language service container with `server.ServerSrvContBlock`, calls `server.CreateDefaultServices`,
then `server.StartLanguageServer`.
`DefaultDocumentSyncher` handles `didOpen` / `didChange` / `didClose` (and save-related hooks), updates the `textdoc` store,
and invokes `DocumentUpdater.Update`, which eventually calls `Builder.Build` under the workspace lock.

Definition and references requests use the built documents and linking results; exact behavior depends on the default providers
and your symbol/reference setup.

The example at `examples/statemachine/server/main.go` wires `StatemachineSrv` with server defaults and starts stdio JSON-RPC.
See the [Language server guide](../guides/language-server.md) for a step-by-step reading of that pattern.

## Mental model for the statemachine example

- `statemachine.fb` is the language definition; `//go:generate go run ../../cmd/fastbelt …` keeps `*_gen.go` in sync.
- `CreateServices` builds the embedded container, sets `LanguageID` and `FileExtensions`, and may assign `DocumentValidator`.
- `validation.go` implements `Validator` on the root AST implementation to enforce domain rules (unique names, valid transition targets).
- `server/main.go` adds LSP services and starts the server.

Together, this mirrors how `internal/grammar` bootstraps fastbelt itself, with a different grammar name and optional validator overrides.
