# TODO

Always write tests!!!

## Tasks

### Backlog

* [ ] delete this file
* [x] redo open TODOs
* [x] split functions into smaller ones
* [x] token groups with group and commands
* [x] token group definition inside a token mode
* [x] collect keywords from entire grammar
* [x] handle token mode
  * [x] local TokenDecls
  * [x] TokenUsages / from global TokenDecls
  * [x] Keywords
  * [x] KeywordSelectors
* [ ] respect this one: https://typefox.slack.com/archives/C0956H0EAUB/p1783583717827469?thread_ts=1783582041.497139&cid=C0956H0EAUB

## Validations

* [x] if there is a default token mode
  * [x] ... and a KEYWORD is not listed in one of the token modes
    * [x] then it is an error
* [x] `push`/`mode` command without a target mode is an error
* [x] `pop` command with a target mode is an error
* [x] a token mode that nothing enters is a warning
* [x] an empty token mode is a warning
* [x] a token that is not listed in any token mode is a warning, not an error
      (decided; `SomeTokenGroup` in completion.fb relies on this)

## Open decisions from the token mode test sweep

* [x] `mode(X)` replaces the active mode, so a following `pop` cannot undo it
      (decided; mirrors ANTLR's `mode` - only `push` can be undone by `pop`)
* [ ] mode-local `token`/`token group` declarations are visible to parser rules
      but not to a token usage in another mode. Should cross-mode reuse work?
* [ ] a `push(X)` command resolves modes across documents. Should token modes be
      file-local instead?
* [ ] `lexer.Mode`, `lexer.NewMode`, `lexer.NewDefaultMode` and
      `lexer.DefaultMode` are dead since `TokenMode` replaced them - remove?
* [ ] `CompositeRule` has no `SymbolKind`, so it shows up as `lsp.Field` in the
      outline (`TokenMode`/`TokenGroup` were given proper kinds).
* [x] in token_modes.fb, `STRING_CONTENT` treated a lone `\` as literal text, so
      `` `a\\` `` lost its closing backtick to the longest match. Fixed by
      excluding `\` from the literal alternative.
