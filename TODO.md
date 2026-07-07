# TODO

Always write tests!!!

## Tasks

* [ ] collect keywords from entire grammar
* [ ] handle implicit default token mode
  * [ ] local TokenDecls
  * [ ] TokenUsages / from global TokenDecls
  * [ ] Keywords
  * [ ] KeywordSelectors

## Validations

* [ ] no duplicate token mode names
  * [ ] maximal one default token mode
* [ ] if there is a default token mode
  * [ ] ... and a KEYWORD is not listed in one of the token modes
    * [ ] then it is an error
