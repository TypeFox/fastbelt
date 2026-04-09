# Language server (LSP)

This guide shows how to expose your language through a **Language Server Protocol** server using the same service container and
workspace types as a CLI or tests. The stock setup speaks JSON-RPC over **stdio**, which is what editors like VS Code expect from a
language client.

## What you implement versus what fastbelt provides

You already have:

- A grammar, generated `*_gen.go` files, and a hand-written `CreateServices()` that sets `LanguageID` and `FileExtensions`
  (see [Scaffolding](scaffolding.md) and [Concepts](../explanations/concepts.md)).

The `server` package adds:

- **`ServerSrv`** — slog handler, LSP server implementation, document syncher, definition/references providers, and stdio transport
  hooks (`ConnectionDialer`, `ConnectionBinder`).
- **`CreateDefaultServices`** — fills in `NewDefaultLanguageServer`, `NewDefaultDocumentSyncher`, default definition and references
  providers, and a stdio dialer unless you override them.

`StartLanguageServer` dials stdio, binds your `LanguageServer` to the connection, registers a **build-step listener** on
`DocStateValidated`, merges lexer/parser/linker diagnostics with `doc.Diagnostics`, and publishes them via LSP.

## Minimal `main` (statemachine pattern)

The example at [`examples/statemachine/server/main.go`](../../examples/statemachine/server/main.go) is the canonical shape:

1. Construct your language service struct (`statemachine.CreateServices()`).
2. Embed `server.ServerSrvContBlock` on a wrapper type that also embeds or holds your language `*Srv`.
3. Call `server.CreateDefaultServices(srv)` so syncher, providers, and transport defaults are installed.
4. Call `server.StartLanguageServer(ctx, srv)`.

The wrapper must satisfy `server.ServerSrvCont`, which requires `textdoc.TextdocSrvCont` and `workspace.WorkspaceSrvCont` in
addition to `Server() *ServerSrv`. Embedding `*statemachine.StatemachineSrv` provides the textdoc and workspace sides; embedding
`ServerSrvContBlock` provides `Server()`.

## Capabilities and behavior

`Initialize` advertises incremental text sync, an (empty) completion provider, and definition/references when the corresponding
providers are non-nil. The default definition and references implementations use the built documents and linking data from your
workspace.

Document lifecycle (`didOpen` / `didChange` / `didClose`) goes through `DefaultDocumentSyncher`, which updates `textdoc` and drives
`DocumentUpdater` → `Builder.Build` like the [consumption](consumption.md) guide describes.

## Customization

- Replace **`DefinitionProvider`** or **`ReferencesProvider`** on `ServerSrv` before starting if the defaults are not enough.
- Swap **`DocumentSyncher`** only if you need different LSP buffering or batching (uncommon).
- Point **`SlogHandler`** at a custom handler to route logs to the client when `Initialize` runs.
- For non-stdio transports, assign a different **`ConnectionDialer`** (and binder if needed) on `ServerSrv` before
  `StartLanguageServer`.

## Editor integration

This repository includes a VS Code extension under `internal/vscode-extensions/` (see the root `README.md` for pointers). That
code is separate from the Go `server` package: the extension typically launches your `main` as a child process and talks LSP over
its stdin/stdout.

## Related reading

- [Concepts — LSP integration](../explanations/concepts.md#lsp-integration-at-a-high-level) for how syncher and build fit together.
- [Validation](validation.md) for where `doc.Diagnostics` come from before publish.
