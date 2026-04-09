# Scaffolding and code generation

This guide covers both **bootstrapping** a language project (`fastbelt scaffold`) and **regenerating** Go sources from a `.fb` grammar
(the default `fastbelt` invocation). The implementation lives under [`cmd/fastbelt`](../../cmd/fastbelt) and
[`internal/scaffold`](../../internal/scaffold).

---

## `fastbelt scaffold` — new module or new package

Use this when you want a **starting layout**: minimal grammar, `gen.go` with `//go:generate`, `services.go`, an LSP entry command, and
(by default) a VS Code extension skeleton plus root `package.json` for npm workspaces.

### Usage

```text
fastbelt scaffold -module <path> [-package <dir>] -language <name>
fastbelt scaffold [-package <dir>] -language <name>
```

- **`-language`** (required) — human-readable name; drives derived identifiers in templates (file names, types, slugs).
- **`-package`** (default `.`) — directory **relative** to the module root where templates are written.
  With `-module`, it is relative to the **new** module directory; without `-module`, relative to the **current working directory**.
- **`-module`** (optional) — Go **module path** for `go mod init`. When set, the tool creates a **new directory** named after the last
  segment of that path (for example `-module=example.com/acme/foo` → `./foo/` under the cwd), and that directory must be **missing or
  empty**. Then it runs `go mod init`, writes templates, attempts `go get` / `go get -tool` for fastbelt, runs `go generate`, and
  `go mod tidy`.
- **Without `-module`** — requires an existing `go.mod` in the cwd or a parent. Writes templates under `-package` relative to cwd, then
  `go get`, `go generate` (scoped to the package when it is not the module root), and `go mod tidy`.
- **`-no-vscode`** — skip the VS Code extension tree, root `package.json`, and related npm layout.

Run `fastbelt help` or `fastbelt scaffold -h` for the built-in usage text.

### What scaffold emits

Files are rendered from [`internal/scaffold/templates`](../../internal/scaffold/templates) into `WriteRoot` (the resolved package
directory). Always included:

- **`README.md`** — how to regenerate, run the LSP command, and (if applicable) build the extension.
- **`.gitignore`**
- **`gen.go`** — `//go:generate go tool typefox.dev/fastbelt/cmd/fastbelt -g ./<grammar>.fb -o . -p <package> -v`
- **`services.go`** — service container with `LanguageID` and `FileExtensions` filled from the template.
- **`<grammar>.fb`** — minimal starter grammar (name derived from the language label).
- **`cmd/<lsp-slug>/main.go`** — stdio LSP server wiring `CreateServices` and `server.CreateDefaultServices`.

With VS Code (default):

- Root **`package.json`** (npm workspace).
- **`vscode-extension/`** — extension manifest, `extension.ts`, TextMate grammar stub, esbuild and static assets.

Scaffolding tries `go get typefox.dev/fastbelt@latest` and `go get -tool typefox.dev/fastbelt/cmd/fastbelt@latest`; failures are warnings
so working inside the fastbelt repo itself still works when the module already satisfies those paths.

---

## Default `fastbelt` — grammar → `*_gen.go`

When you do **not** pass `scaffold` (or `help`), the binary runs **code generation** only: read one `.fb` file, validate it with the same
builder pipeline as the library, then write five generated Go files.

### Installing the tool

From the root `README.md`:

```sh
go get -tool typefox.dev/fastbelt/cmd/fastbelt@latest
# or
go install typefox.dev/fastbelt/cmd/fastbelt@latest
```

On demand without install:

```sh
go run typefox.dev/fastbelt/cmd/fastbelt@latest -g ./grammar.fb -o ./
```

### Generate flags

Behavior matches [`cmd/fastbelt/main.go`](../../cmd/fastbelt/main.go): a `flag.FlagSet` parses:

| Flag | Default        | Meaning                                                        |
| ---- | -------------- | -------------------------------------------------------------- |
| `-g` | `./grammar.fb` | Grammar file path (absolute path used internally).             |
| `-o` | `./`           | Output directory for generated files (`MkdirAll` with `0755`). |
| `-p` | *(derived)*    | Go `package` clause; if empty, `filepath.Base(outputPath)`.    |
| `-v` | `false`        | Print `Written: <path>` for each emitted file.                 |

**Grammar validation.** After the build, diagnostics are printed as `Severity - line:col message` (1-based line/column for display). If
any diagnostic has **error** severity, the process exits with an error and **does not** write `*_gen.go` files. Fix the grammar (or
workspace issues) until the build is clean.

**Package name vs grammar name.** `-p` only sets the Go package line. The `grammar Name;` header in the `.fb` file still prefixes
generated types and service identifiers (for example `StatemachineModelLinkingSrv`). Your folder or `-p` can differ, as in
`examples/statemachine/`.

### Generated files (fixed basenames)

For a successful run, exactly five files appear under `-o`, in this order:

1. `linker_gen.go`
2. `types_gen.go`
3. `parser_gen.go`
4. `lexer_gen.go`
5. `services_gen.go`

Each begins with a generated header (`// Code generated by typefox.dev/fastbelt/cmd/fastbelt. DO NOT EDIT.`). Regenerate after grammar
changes; keep
hand-written `.go` files separate.

---

## `//go:generate` patterns

**Inside the fastbelt repo** (examples and `internal/grammar`), directives typically use a relative `go run` to the command sources:

```go
//go:generate go run ../../cmd/fastbelt -g ./statemachine.fb -v
```

Omitting `-o` writes to the **current directory** during `go generate` (the package directory).

**Consumer modules** usually depend on the tool path and use `go tool` (after `go get -tool …`), as in scaffold’s `gen.go.tmpl`:

```go
//go:generate go tool typefox.dev/fastbelt/cmd/fastbelt -g ./grammar.fb -o . -p mypkg -v
```

---

## End-to-end workflows

**Greenfield language (scaffold):**

1. `fastbelt scaffold -module example.com/you/mylang -language "MyLanguage"` (or add `-package pkg/lang` to nest the package).
2. Edit the generated `.fb` grammar and hand-written code as needed.
3. `go generate ./...` from the module root when the grammar changes.

**Existing module (scaffold a package):**

1. From the module root, `fastbelt scaffold -language "MyLanguage" -package path/to/pkg` (requires `go.mod` discoverable from that tree).
2. Same edit/regenerate loop.

**Existing package (no scaffold):**

1. Add a `.fb` grammar and a `//go:generate` line.
2. Run `go generate` or invoke `fastbelt` / `go run …/cmd/fastbelt` with `-g`, `-o`, `-p`, `-v` as needed.
3. Implement `services.go`, optional `validation.go`, and any commands (see [Validation](validation.md) and [Consumption](consumption.md)).

---

## Summary

| Goal                           | Command                                                      |
| ------------------------------ | ------------------------------------------------------------ |
| New module + starter files     | `fastbelt scaffold -module <import/path> -language "<Name>"` |
| New package in existing module | `fastbelt scaffold -language "<Name>" [-package <rel dir>]`  |
| Regenerate `*_gen.go` only     | `fastbelt -g grammar.fb [-o dir] [-p pkg] [-v]`              |
| Help                           | `fastbelt help`                                              |
