# TODO

Always write tests!!!

## Tasks

### Current

* pick a language with token modes
* inline token declarations that have keywords only
* list the keywords in the token mode declarations
* test feature
* implement feature

### Backlog

* [ ] token groups with group and commands
* [x] collect keywords from entire grammar
* [ ] handle implicit default token mode
  * [ ] local TokenDecls
  * [x] TokenUsages / from global TokenDecls
  * [x] Keywords
  * [ ] KeywordSelectors
* [ ] respect this one: https://typefox.slack.com/archives/C0956H0EAUB/p1783583717827469?thread_ts=1783582041.497139&cid=C0956H0EAUB

## Validations

* [ ] if there is a default token mode
  * [ ] ... and a KEYWORD is not listed in one of the token modes
    * [ ] then it is an error
