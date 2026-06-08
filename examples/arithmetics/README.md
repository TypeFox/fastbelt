# arithmetics

This repository was generated with [fastbelt](https://github.com/TypeFox/fastbelt) using `fastbelt scaffold`.

It includes:

- **`arithmetics.fb`** — minimal fastbelt grammar (`ArithmeticsModel`).
- **`gen.go`** — `//go:generate` running `go tool typefox.dev/fastbelt/cmd/fastbelt generate` to emit `*_gen.go`.
- **`services.go`** — LSP service wiring (document, workspace, linking, generated lexer/parser).
- **`cmd/arithmetics-lsp`** — language server entrypoint (stdio JSON-RPC).
- **`vscode-extension/`** — VS Code language client + TextMate grammar skeleton.
- **Root `package.json`** — npm workspace so you can build and package the extension from the repo root.

## Requirements

- Go (see `go.mod`)
- Node.js 20+ and npm 10+ for the extension

## Regenerate Go sources

From the module root:

```sh
go generate ./...
```

## Language server

```sh
go run ./cmd/arithmetics-lsp
```

## VS Code extension

From the repository root:

```sh
npm install
npm run build
npm run package
```

- **`npm run build`** runs TypeScript checks, bundles the client with esbuild, and builds the Go language server into `vscode-extension/dist/server`.
- **`npm run package`** produces `vscode-extension/arithmetics-vscode.vsix` (a `.vsix` you can install with “Install from VSIX…”).
Upstream docs and sources: [https://github.com/TypeFox/fastbelt](https://github.com/TypeFox/fastbelt)
