![Fastbelt](./fastbelt-full-logo.png)

Fastbelt is a high-performance DSL toolkit for Go with a parser generator and Language Server Protocol (LSP) support.

It is designed for language tooling that needs low latency and good throughput on large workspaces.

For background and benchmarks, see the [Fastbelt introduction](https://www.typefox.io/blog/fastbelt-introduction/) blog post.

## Installation

### Adding to a Module
Fastbelt ships as a Go module:

```sh
go get typefox.dev/fastbelt@latest
go get -tool typefox.dev/fastbelt/cmd/fastbelt@latest
```

### Global Install

You can also globally install the fastbelt CLI:

```sh
go install typefox.dev/fastbelt/cmd/fastbelt@latest
```

## Quick start

The first step is to write a grammar definition file, e.g. `grammar.fb`, which has a similar format as the grammar language of Langium.

To run the code generator for your grammar definition:

```sh
# globally installed
fastbelt generate ./grammar.fb -o ./
```

```sh
# on demand install
go run typefox.dev/fastbelt/cmd/fastbelt@latest generate ./grammar.fb -o ./
```

This writes generated Go files for services such as lexer, parser, linker, and type definitions.

Typically you will want to run generation using `go generate`.
Add a directive to some file in your module (assumes install with `go tool`):

```go
//go:generate go tool typefox.dev/fastbelt/cmd/fastbelt generate ./grammar.fb -o ./
```

## Scaffolding

To bootstrap a **new Go module** for a language, run:

```sh
fastbelt scaffold --module example.com/you/mylang --language "MyLanguage" --vscode
```

That creates a directory named after the last segment of `--module` (here `./mylang`) in the current working directory, runs `go mod init`, pulls in fastbelt as a library and tool dependency, and runs `go generate` and `go mod tidy` so generated parser, lexer, linker, and supporting files are ready to use.

In the generated package directory you get:

- `gen.go` with `go:generate` directives for grammar-based code generation
- `services.go` with customization points for generated services
- `<language-id>.fb` as the initial grammar definition
- `cmd/<language-id>-lsp/main.go` as the LSP server entrypoint

Pass `--vscode` to also add a `vscode-extension` subfolder for editor integration.

To add a language package to an **existing** module, omit `--module` and use `--package` (or `-p`) instead; see `fastbelt scaffold -h` for full usage.

## Examples

A minimal state machine example is available in `examples/statemachine`.

For editor integration, see the VS Code extension in `internal/vscode-extensions/statemachine`.

## Contributing

See [CONTRIBUTING.md](./CONTRIBUTING.md) for contribution guidelines.

## License

Fastbelt is licensed under the [MIT License](./LICENSE).
