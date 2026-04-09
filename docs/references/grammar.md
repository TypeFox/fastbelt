# Grammar language

Fastbelt grammar files (conventionally `.fb`) describe the **lexical** structure (tokens) and **syntactic** structure (parser rules) of a language. The notation is intentionally close to [Langium’s grammar language](https://langium.org/docs/reference/grammar-language/), but the implementation differs: Fastbelt compiles grammars to Go lexers and parsers, and **token patterns are written as [Go regular expressions](https://pkg.go.dev/regexp/syntax)** (RE2-style), not JavaScript regular expressions.

The **authoritative syntax** of the Fastbelt grammar language itself is the bootstrapped grammar in [internal/grammar/grammar.fb](../../internal/grammar/grammar.fb). This page explains that language in prose and relates it to Langium where helpful.

---

## Table of contents

- [Document structure](#document-structure)
- [Language declaration](#language-declaration)
- [Interfaces and fields](#interfaces-and-fields)
- [Parser rules](#parser-rules)
- [Entry rule](#entry-rule)
- [EBNF-like rule bodies](#ebnf-like-rule-bodies)
- [Assignments](#assignments)
- [Cross-references](#cross-references)
- [Unassigned rule calls](#unassigned-rule-calls)
- [Actions and tree rewriting](#actions-and-tree-rewriting)
- [Token rules](#token-rules)
- [Keywords](#keywords)
- [Regular expressions (Go)](#regular-expressions-go)
- [Lexer behavior](#lexer-behavior)
- [Differences from Langium](#differences-from-langium)

---

## Document structure

A grammar file is a single **language declaration** followed by any mixture of:

- `interface` definitions (AST shapes),
- `Name returns Type: …` parser rules, and
- `token` / `hidden token` definitions.

There is **no** `import` of other grammar files in the current Fastbelt grammar language.

Typical layout (order of declarations is flexible except for the [entry rule](#entry-rule)):

```fastbelt
grammar MyLanguage;

interface Root {
    Name string
}

Root returns Root:
    Name=ID;

token ID: /[a-z]+/;
hidden token WS: /\s+/;
```

**What it means:** One header, then any mix of interfaces, parser rules, and tokens in whatever order you prefer for readability.

**As a language designer:** Group related interfaces together, put the document rule near the top if you want readers to see the big picture first, and keep token definitions in one place so lexer precedence stays easy to reason about.

---

## Language declaration

Every grammar file begins with the language name:

```fastbelt
grammar MyLanguage;
```

**What it means:** This line declares the logical language name `MyLanguage`. Code generation uses it as the prefix for **service and linking types** in Go—for example `MyLanguageGeneratedSrvCont`, `MyLanguageLinking()`, and the default scope/reference helpers wired in [internal/generator/services_generator.go](../../internal/generator/services_generator.go) via `grammr.Name()`. The generated **parser** also refers to that interface when it calls into linking.

**What it does *not* do:** It does **not** set the Go `package` line in emitted files, and the **`.fb` file name is not used** for any of this. The compiler reads the path only to load the file. The package name comes from the code generator’s `-p` flag, or—if `-p` is omitted—from the **last path segment of the output directory** (`-o`); see [`cmd/fastbelt/main.go`](../../cmd/fastbelt/main.go). A generated directory like `package mylang_gen` can therefore sit beside `grammar MyLanguage;` in one project: the two names are unrelated unless you align them on purpose.

**As a language designer:** Treat `grammar …` as the stable **language id** for generated Go **type names** that your hand-written `services.go` embeds (`…GeneratedSrvCont`, `…Linking()`). Renaming it is a breaking change for that wiring. Pick `-o` / `-p` separately so the **Go package** fits your repo layout (one package per generated directory is normal in Go).

**Is mismatch a bug?** No—it follows Go’s rule that package names come from the directory (or an explicit `-p`), while the grammar header names the language for Fastbelt-specific APIs. If you want the package to mirror the grammar name, pass `-p MyLanguage` or name the output folder accordingly.

---

## Interfaces and fields

Interfaces describe AST node shapes, similar to Langium’s `interface` types.

```fastbelt
interface Person {
    Name     string
    Verified bool
}
```

**What it means:** A `Person` AST node carries a textual `Name` and a `Verified` flag. Only `string` and `bool` are allowed as built-in primitives; everything else is another interface or a slice/reference form.

**As a language designer:** Declare one interface per syntactic or semantic concept you want in the tree (statement, declaration, expression kind). Fields are the contract between the grammar assignments and the Go types codegen will emit. Give them the names you want them in downstream validators and generators.

### Field types

Supported forms (see `FieldType` in [internal/grammar/grammar.fb](../../internal/grammar/grammar.fb)):

| Syntax                     | Meaning                                                                    |
| -------------------------- | -------------------------------------------------------------------------- |
| `Name string`, `Name bool` | Primitive fields (`string` and `bool` only).                               |
| `Name OtherInterface`      | Child node of type `OtherInterface`.                                       |
| `Name *OtherInterface`     | Optional or reference-style link to `OtherInterface` (see generated code). |
| `Name []Item`              | Slice of `Item` (array type is `[]` followed by the element `FieldType`).  |

### Inheritance

```fastbelt
interface Employee extends Person {
    Id string
}
```

**What it means:** An `Employee` node is modeled as a `Person` plus an extra `Id`; generated code treats the inheritance hierarchy like a normal Go embedding/extension story for that AST type.

**As a language designer:** Use `extends` when several node kinds share the same fields (e.g. common `Location` on statements) so you do not duplicate field lists. Prefer shallow hierarchies so parser rules and validations stay easy to follow.

`extends` lists one or more parent interfaces, separated by commas.

---

## Parser rules

Parser rules define how tokens combine into AST nodes. The return type is **explicit** (there is no Langium-style `infers` in the grammar language today):

```fastbelt
Person returns Person:
    "person" Name=ID;
```

**What it means:** Input `person alice` produces a `Person` node whose `Name` comes from the `ID` token after the keyword. The leading `"person"` is fixed text the lexer recognizes as its own token type.

**As a language designer:** Use keywords to disambiguate alternatives (two rules that only differ on `ID` are ambiguous; different leading keywords are not). One parser rule implements one construction in your surface syntax and fills one interface-shaped node.

Syntax (template for any new rule):

```fastbelt
RuleName returns InterfaceName: 
  body ;
```

**What it means:** `RuleName` is invoked from other rules or as the entry rule; `InterfaceName` is the AST type this rule builds; `body` is the EBNF fragment.

**As a language designer:** Keep `RuleName` and the return interface aligned with your mental model (`IfStmt returns IfStmt`). Split large rules into smaller ones when one alternative would otherwise dominate readability or lookahead.

- `RuleName` is both the rule’s identifier and the usual name used for `Rule=[Rule:ID]` calls unless you use a different assignment.
- `InterfaceName` must be an interface declared in the same grammar.

---

## Entry rule

In Fastbelt, the **first parser rule** in the file is the document entry point: generated `Parse` calls `Parse<FirstRuleName>()`. This differs from Langium, which marks the start rule with `entry`. 

```fastbelt
interface Library {
    Decls []Decl
}

interface Decl {
    Name string
}

Library returns Library:
    Decls+=Decl*;

Decl returns Decl:
    Name=ID;
```

**What it means:** A full document is parsed by calling the generated method for whichever parser rule appeared first; everything else is reached transitively from that rule.

**As a language designer:** Put your “file root” rule above other parser rules if you rely on default `Parse`, or you will accidentally parse only a fragment. Interfaces and tokens before that rule are fine; only the ordering among **parser rules** matters for entry selection.

---

## EBNF-like rule bodies

Rule bodies use grouping, sequencing, alternatives, repetition, and parentheses, in the same spirit as Langium’s EBNF.

### Cardinality

| Operator | Meaning      |
| -------- | ------------ |
| (none)   | Exactly one  |
| `?`      | Zero or one  |
| `*`      | Zero or more |
| `+`      | One or more  |

Cardinality applies to the immediately preceding **element** (keyword, assignment, rule call, cross-reference, action, or parenthesized alternative group).

```fastbelt
interface Machine {
    Events []Event
    States []State
}

interface Event {
    Name string
}

interface State {
    Name string
    Traced bool
}

Machine returns Machine:
    ("events" Events+=Event+)?
    States+=State*;

Event returns Event:
    Name=ID;

State returns State:
    Name=ID (Trace?="trace")?;
```

**What it means:** `Name=ID` is **exactly one** identifier on each `Event` and `State`. The whole `("events" Events+=Event+)?` group is **optional** (`?`): either the keyword and **one or more** events appear together, or the block is omitted. `Events+=Event+` applies `+` to the repeated “append an `Event`” element, so the list is never empty when the block is present. `States+=State*` is **zero or more** states at top level. `(Trace?="trace")?` is an optional modifier: at most one `trace` keyword, stored as a boolean.

**As a language designer:** Put `?` on optional clauses (`else`, `where`, whole import blocks). Use `*` for “any number including none” lists (`decl*`, `statement*`). Use `+` when the syntax requires at least one instance (`argument+`, `case+`). When several tokens must move as one unit (keyword + list), wrap them in `(` `)` and hang `?`, `*`, or `+` on the group so cardinality does not attach to the wrong symbol.

### Alternatives

Use `|` between alternatives. The grammar model represents this as nested alternatives (see `Alternatives` / `Group` in [internal/grammar/grammar.fb](../../internal/grammar/grammar.fb)).

```fastbelt
interface Stmt {}

interface IfStmt extends Stmt {
    Cond string
    ThenLabel string
}

interface ReturnStmt extends Stmt {
    Value string
}

Stmt returns Stmt:
    IfStmt | ReturnStmt;

IfStmt returns IfStmt:
    "if" Cond=ID "then" ThenLabel=ID;

ReturnStmt returns ReturnStmt:
    "return" Value=ID;
```

**What it means:** Wherever the grammar expects a `Stmt`, the parser matches **either** an `IfStmt` (`if` … `then` …) **or** a `ReturnStmt` (`return` …). Concrete types extend `Stmt`, so downstream code can treat the result as one shared interface (the common pattern is an empty or minimal base interface and `extends` on each alternative).

**As a language designer:** Use alternatives for “one of several statement / expression / declaration forms.” Ensure branches **commit** differently at parse time (typically with distinct leading keywords (`if` vs `return`) or tokens) so the parser does not face two identical-looking arms. If two alternatives share the same prefix, split further with nested rules or richer lookahead, or refactor the surface syntax.

### Sequencing

Adjacent elements without `|` are a **sequence** (must appear in order).

### Parentheses

`( … )` turn several pieces of syntax into **one** unit. That matters whenever you attach `?`, `*`, or `+`, or when you need to spell **either a long sequence or an alternative** without the wrong symbol “stealing” part of the pattern.

**Optional phrase (one `?` for the whole group)**

```fastbelt
interface Endpoint {
    Host string
    Port string
}

Endpoint returns Endpoint:
    Host=ID (":" Port=ID)?;
```

**What it means:** You get `db` *or* `db:5432`. The `?` applies to **both** the colon and the port—if you wrote `":" Port=ID?`, only the port would be optional and you could end up with a stray `:`.

**Using it:** In services, URLs, or attributes where a suffix is “all or nothing”, wrap the suffix in parentheses before `?`.

**Repeating a choice (`*` / `+` over `|`)**

```fastbelt
interface Syllables {
    Parts []string
}

Syllables returns Syllables:
    Parts+=("la" | "da")*;
```

**What it means:** Zero or more times, match `la` or `da` and append each keyword’s text to `Parts`. The same `+=` + parenthesized keyword alternatives works for any fixed vocabulary (e.g. `Ops+=("+" | "-")*`). The parentheses force `*` to repeat **the whole alternative**, not just the last keyword.

**Contrast (why grouping matters):**

- `(la | da)*` — empty, or `la`, or `da`, or `lada`, …
- `la | da*` — **either** a single `la`, **or** any number of `da` (no mixing as in the first form).

**As a language designer:** Whenever you write `|` next to `*` or `+`, stop and decide which span the repetition should cover; add parentheses until the meaning matches what you say in English.

**Parentheses around recursion**

```fastbelt
interface Expr {}

Expr returns Expr:
    ID | "(" Expr ")";
```

**What it means:** An expression is an identifier or a **nested** expression in parentheses. The outer `"("` … `")"` are literal keyword tokens; the inner `Expr` is the recursive call.

**Using it — sample text the rule accepts (with a suitable `ID` token):**

| Input   | How it parses                            |
| ------- | ---------------------------------------- |
| `x`     | Single `ID` arm.                         |
| `(x)`   | Outer parens around inner `ID`.          |
| `((x))` | Nested parens; each level is one `Expr`. |

**As a language designer:** This is the usual pattern for “parentheses override precedence” in expression languages; you still add separate rules (or actions) for `+`, `*`, etc., between the atomic and parenthesized forms.

---

## Assignments

Assignments attach parsed values to AST fields. Operators match Langium:

| Operator | Use                                                                     |
| -------- | ----------------------------------------------------------------------- |
| `=`      | Single value                                                            |
| `+=`     | Append to a **slice** field                                             |
| `?=`     | Boolean: `true` if the following syntax was consumed, otherwise `false` |

Example:

```fastbelt
Employee returns Employee:
    "employee" Name=ID (Remote?="remote")?;
```

**What it means:** Every `Employee` has a `Name`. The optional keyword `remote`, if present, sets boolean field `Remote` to true; if absent, `Remote` is false. No extra token is stored—only the flag.

**As a language designer:** Use `?=` for optional modifiers (`public`, `async`, `unsafe`) that you want as clean booleans on the AST instead of “maybe a keyword token”. Use `=` / `+=` when the value itself matters (identifiers, literals, child trees).

---

## Cross-references

Reference another node by type, optionally specifying which **token or rule** supplies the identifier text:

```fastbelt
Greeting returns Greeting:
    "hello" Target=[Person:ID] "!";
```

**What it means:** After `hello`, the parser reads an `ID` and attaches a **reference** to a `Person` with that name, not an embedded `Person` subtree. The actual link is resolved later (name table / scope), like Langium cross-references.

**As a language designer:** Use cross-references whenever one construct **names** another (`extends Foo`, `goto L`, `import "path"` pointing at a module): you keep the AST small and centralize “does this name exist?” in linking and validation. Use nested parser rules with `=` when you parse the full definition inline instead of by name.

- `[Person:ID]` — reference a `Person`; the lexer token `ID` provides the name text.
- `[Person:RuleName]` — use another parser or token rule for the segment (when the grammar names a rule after `:`).

If the grammar allows omitting `: …`, the bootstrap grammar documents the optional `Rule` part of `CrossRef` (`[Type=[Interface:ID] (":" Rule=RuleCall)?]`).

Resolution (scopes, duplicate names, etc.) is handled in generated linking code and services, analogous to Langium’s scope provider.

**Multi-target references** (`[+Person:ID]` in Langium) are **not** part of the Fastbelt grammar surface in [internal/grammar/grammar.fb](../../internal/grammar/grammar.fb).

---

## Unassigned rule calls

A rule call without `=` / `+=` / `?=` does not assign into the current node; it delegates parsing to that rule for its side effect (e.g. consuming tokens) or for nested structure, depending on how the callee is defined.

```fastbelt
Block returns Block:
    "{" Statement "}";
```

**What it means:** `Statement` is invoked to parse the middle; whatever node `Statement` returns is discarded at this level because there is no assignment—only the fact that parsing succeeded matters for `Block` (you would normally assign if you need the child on `Block`).

**As a language designer:** Use unassigned calls for inline sub-rules that exist only to share syntax (`Expr` wrappers), or temporarily while prototyping—usually you switch to `Body+=Statement*` (or `Body=Statement`) once the AST needs the child. If the callee returns a node you need, assign it.

---

## Actions and tree rewriting

Curly-brace **actions** build or reshape the AST during parsing. The grammar surface is defined in [internal/grammar/grammar.fb](../../internal/grammar/grammar.fb): `{Type}` selects a node type, and `{Type.Field = current}` / `{Type.Field += current}` **rewrites** the tree by creating a new `Type`, moving the old `current` into `Field`, then continuing with the new node as `current`. This is the way to get left-associative infix chains without left-recursive rules.

### Worked example: left-associative `+` and `-`

You want input `a + b - c` to mean `(a + b) - c`, not `a + (b - c)`. A left-recursive rule like `Expr := Expr ("+" | "-") Primary` is not usable with an LL parser; instead you parse a `Primary`, then repeat `{BinaryExpr.Left=current} Op=("+" | "-") Right=Primary`.

```fastbelt
interface Expr {}

interface NameExpr extends Expr {
    Name string
}

interface BinaryExpr extends Expr {
    Left  Expr
    Op    string
    Right Expr
}

interface NegExpr extends Expr {
    Operand Expr
}

BinaryExpr returns Expr:
    Unary ({BinaryExpr.Left=current} Op=("+" | "-") Right=Unary)*;

Unary returns Expr:
    ({NegExpr} "-" Operand=Unary) | Primary;

Primary returns NameExpr:
    Name=ID;

token ID: /[a-zA-Z_][a-zA-Z0-9_]*/;
hidden token WS: /[ \n\r\t]+/;
```

**What it means:** `BinaryExpr` is listed first so it becomes the generated **parse entry** for this snippet. It starts with `Unary` (identifier, or unary `-` nested arbitrarily). That node becomes `current`. Each time the `*` group matches, the action runs **before** consuming the operator: a new `BinaryExpr` is allocated, `Left` is set to the previous `current`, then `Op` and `Right` are filled from the next `"+"` / `"-"` and following `Unary`. The new `BinaryExpr` replaces `current`. If the group never matches, the result is the lone subtree from the first `Unary` (e.g. a `NameExpr` or a `NegExpr`).

The `Unary` rule shows a **bare `{NegExpr}`** action: on the first alternative the parser switches to a `NegExpr` node, consumes `-`, then parses the operand with another `Unary` call so `- - a` becomes nested negations.

**Walkthrough for `a + b - c` (whitespace ignored):**

1. `Unary` reads `a` via `Primary` → `current` is `NameExpr(a)`.
2. First repetition: action sets `Left` to that node; `Op` is `+`; `Right` is `Unary` → `NameExpr(b)` → `current` is `BinaryExpr(+, a, b)`.
3. Second repetition: action sets `Left` to that `BinaryExpr`; `Op` is `-`; `Right` is `Unary` → `NameExpr(c)` → `current` is `BinaryExpr(-, BinaryExpr(+,a,b), c)`.

Rough shape of the final tree:

```text
BinaryExpr (op: "-")
├── Left:  BinaryExpr (op: "+")
│          ├── Left:  NameExpr "a"
│          └── Right: NameExpr "b"
└── Right: NameExpr "c"
```

**Unary minus in the same grammar:** For `-a + b`, `Unary` takes the `{NegExpr}` branch first (`Operand=Unary` reads `a`), then the binary `*` group builds `BinaryExpr` with `Left = NegExpr(a)`, `Op = +`, `Right = NameExpr(b)`.

**As a language designer:** The `{Type.Left=current} … Right=…` idiom is the standard replacement for left recursion on binary operators. Use **one** new node per operator occurrence when you want a nested binary tree; use **slice** fields with `+=` on `Op` and `Right` in the repeated group instead if you prefer a flat “N-ary” node (`Left`, `Ops[]`, `Rights[]`)—same surface syntax idea, different AST shape.

For **unary** operators, the `{NegExpr}` prefix in `Unary` is the small **bare `{Type}`** pattern: it forces the current construction to be a `NegExpr` before reading `-` and the recursive `Unary` operand. Pair `{Type}` with the tokens and assignments that follow in that alternative; use a second alternative (`| Primary`) when no wrapper is needed.

---

## Token rules

Terminal tokens are declared with the `token` keyword and a **slash-delimited** regular expression literal:

```fastbelt
token ID: /[a-zA-Z_][a-zA-Z0-9_]*/;
```

**What it means:** Any substring matching the Go regex becomes an `ID` token: ASCII letters or underscore, then optional alphanumerics/underscores. Anything else will not match this rule at that position.

**As a language designer:** Define one `token` per lexical class you need in parser rules (`ID`, `STRING`, `NUMBER`). Keep regexes **disjoint** where possible; when they overlap, remember [lexer behavior](#lexer-behavior) (longest match, keyword vs regex ordering). Prefer simpler tokens plus parser structure over one giant regex.

### Hidden tokens (whitespace, comments)

Prefix with `hidden` so the lexer skips emitting them as real tokens (they are still consumed):

```fastbelt
hidden token WS: /[ \n\r\t]+/;
```

**What it means:** Space, tabs, and newlines are consumed by the lexer but **not** passed to the parser as tokens, so parser rules do not mention whitespace between keywords and assignments.

**As a language designer:** Always add `hidden` tokens for whitespace and comments early; without them, the lexer fails on the first space. Use a single broad `WS` or split spaces vs newlines only if you need layout-sensitive rules later (Fastbelt does not yet offer Langium-style indentation recipes—keep lexer patterns maintainable).

There is no separate `terminal` keyword; Langium’s `terminal` corresponds to Fastbelt’s `token`.

---

## Keywords

**Keywords** are string literals in parser rules, written in **double quotes**:

```fastbelt
"person" Name=ID
```

**What it means:** This fragment is part of a parser rule body: the literal keyword `person` must appear, followed by whatever your `ID` token matches; `Name=ID` stores that token on the `Name` field.

**As a language designer:** Every distinct fixed spelling (`"if"`, `"return"`, `"=>"`) becomes its own keyword token in generated lexers—good for clarity and for disambiguation. Avoid spelling the same keyword differently in two rules; use one quoted form everywhere so lookahead and keyword tables stay consistent.

The bootstrap grammar only defines double-quoted strings (`/"[^"]+"/`). Unlike Langium, there is no single-quoted keyword form in the meta-grammar.

Keywords are extracted from the grammar at codegen time and become dedicated lexer token types. See [Lexer behavior](#lexer-behavior) for ordering vs regex tokens.

---

## Regular expressions (Go)

Token bodies use the **Go regexp syntax** as parsed by `regexp/syntax` with **Perl** mode (`syntax.Perl` in the compiler). That is the same flavor used by Go’s standard `regexp` package: linear-time, RE2-style.

**Notable limitations in Fastbelt’s regex compiler today:** anchors and word-boundary ops are rejected during NFA construction — specifically `^`, `$`, `\A`, `\z`, `\b`, and `\B` (see [internal/regexp/regexp.go](../../internal/regexp/regexp.go)). Prefer patterns that match concrete character sequences without these anchors, or split concerns across tokens/parser rules.

**Escape rules in `.fb` files:** the pattern is written inside `/ … /`. The bootstrap `RegexLiteral` in [internal/grammar/grammar.fb](../../internal/grammar/grammar.fb) describes which sequences are valid inside the slashes; invalid escapes or unclosed classes will fail at grammar parse or codegen time.

When porting grammars from Langium/Chevrotain, **do not assume JavaScript-only features** (e.g. some lookahead idioms); express the intent with Go-safe patterns and, if needed, extra tokens or parser rules.

---

## Lexer behavior

Understanding disambiguation is important for grammars that mix keywords and regex tokens.

1. **Keywords first** — All quoted keywords from the grammar are turned into lexer token types and registered **before** `token` rules in the generated `NewLexer()` (see [internal/generator/lexer_generator.go](../../internal/generator/lexer_generator.go)). Keywords are ordered **alphabetically** by their text value when registered.
2. **Regex tokens** — Declared in the order they appear in the grammar file.
3. **Longest match** — At each offset, among all candidate token types that can start with the current character, the lexer picks the **longest** match (`lexer.DefaultLexer` in [lexer/lexer.go](../../lexer/lexer.go)).
4. **Ties** — If two definitions match the same length, the one that appears **earlier** in that character’s candidate list wins (depends on registration order above).

**Hidden** tokens (`hidden token`) are skipped for the parser’s token stream but still participate in matching the input.

---

## Differences from Langium

Short list of Langium features that are **not** represented in [internal/grammar/grammar.fb](../../internal/grammar/grammar.fb) / Fastbelt’s current grammar language:

| Langium                                                           | Fastbelt                                      |
| ----------------------------------------------------------------- | --------------------------------------------- |
| `entry` rule keyword                                              | First parser rule is the entry rule           |
| `import` other grammars                                           | Not supported                                 |
| `terminal` name                                                   | Use `token`                                   |
| Regex in `/…/` vs Langium’s JS regex                              | **Go** `regexp/syntax` (Perl mode), RE2-style |
| `terminal` / fragment EBNF, `hidden terminal` as separate concept | Regex-only tokens; `hidden token`             |
| `returns number` etc. on terminals                                | Not supported                                 |
| `infers` / `infer`                                                | Not supported; use `returns`                  |
| Data type rules (`returns string`)                                | Not supported as a grammar construct          |
| Rule fragments (`fragment`)                                       | Not supported                                 |
| Guard parameters (`Rule<flag>`)                                   | Not supported                                 |
| Unordered groups (`&`)                                            | Not supported                                 |
| `infix` operator sections                                         | Not supported; use normal rules and actions   |
| Multi-target cross-references `[+Type:ID]`                        | Not supported                                 |
| Single-quoted keywords                                            | Double quotes only in meta-grammar            |

For a full structural definition of what *is* supported, treat [internal/grammar/grammar.fb](../../internal/grammar/grammar.fb) as the specification.
