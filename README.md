# Fastbelt

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
fastbelt -g ./grammar.fb -o ./
```

```sh
# on demand install
go run typefox.dev/fastbelt/cmd/fastbelt@latest -g ./grammar.fb -o ./
```

This writes generated Go files for services such as lexer, parser, linker, and type definitions.

Typically you will want to run generation using `go generate`.
Add a directive to some file in your module (assumes install with `go tool`):

```go
//go:generate go tool fastbelt -g ./grammar.fb -o ./
```

## Examples

A minimal state machine example is available in `examples/statemachine`.

For editor integration, see the VS Code extension in `internal/vscode-extensions/statemachine`.

## Contributing

Issues and pull requests are welcome.

## License

Fastbelt is licensed under the [MIT License](./LICENSE).
