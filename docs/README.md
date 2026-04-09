# Fastbelt documentation

Fastbelt is a Go toolkit for DSLs: grammars in `.fb` files, generated lexer/parser/linker code, workspace builds, and optional LSP
servers. Start here, then follow the links that match your goal.

## Where to read next

| Document                                     | Audience                                                                                                    |
| -------------------------------------------- | ----------------------------------------------------------------------------------------------------------- |
| [Concepts](explanations/concepts.md)         | Anyone adopting the library; mental model for documents, services, builder, and extension points            |
| [Grammar reference](references/grammar.md)   | Authors of `.fb` grammar files                                                                              |
| [Scaffolding](guides/scaffolding.md)         | `fastbelt scaffold` for new modules/packages; default `fastbelt` for grammar → `*_gen.go` and `go:generate` |
| [Validation](guides/validation.md)           | Semantic checks via `Validator` and `DocumentValidator`                                                     |
| [Consumption](guides/consumption.md)         | Driving the workspace from Go (CLI-style tools)                                                             |
| [Language server](guides/language-server.md) | Minimal LSP process over stdio and how it ties to the workspace                                             |

The [statemachine example](../examples/statemachine/) in the repository is the main end-to-end sample (generated code, validation,
`cmd` tool, and `server` entrypoint).
