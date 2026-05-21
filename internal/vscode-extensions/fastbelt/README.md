# Fastbelt VS Code Extension

This extension contributes editor support for the Fastbelt Grammar language (`.fb` files) in Visual Studio Code and other compatible editors.

For Fastbelt concepts, APIs, generator usage, and CLI details, use the library documentation as the primary reference: https://pkg.go.dev/typefox.dev/fastbelt

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
//go:generate go tool typefox.dev/fastbelt/cmd/fastbelt -g ./grammar.fb -o ./
```

## Scaffolding

To bootstrap a **new Go module** for a language (minimal `.fb` grammar, `go:generate` using `go tool` on this CLI, LSP command, and VS Code extension layout), run:

```sh
fastbelt scaffold -module example.com/you/mylang -language "MyLanguage"
```

That creates a directory named after the last segment of `-module` (here `./mylang`) in the current working directory, runs `go mod init`, pulls in fastbelt as a library and tool dependency, lays down the files, and runs `go generate`. Use `fastbelt scaffold -h` for full usage.
