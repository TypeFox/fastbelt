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
