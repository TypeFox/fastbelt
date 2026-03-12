# PLAN-INTEGRATION.md — Pluggable Lookahead Strategy

> **Goal:** Allow the writer of a generated DSL parser to choose LL(k) or ALL(*) at construction
> time. LL(k) is the default. All existing tests continue to pass.

---

## Design constraints

| Constraint | Consequence |
| --- | --- |
| Strategy is construction-time only | No setter exposed after the parser is built. |
| Strategy is never nil | `NewParser` always installs an `LLkStrategy`; `WithLookaheadStrategy` replaces it. |
| Strategy owns its state | `LookaheadStrategy` has no LL(k)-specific parameters. The `LLkStrategy` holds the tables map; the LL(k) fallback is not visible at the interface level. |
| `ParserState` stays clean | No strategy field, no extra methods — `ParserState` is only a token-stream cursor. |

---

## 1. `parser/parser.go` — interface + LL(k) implementation

Add the interface and the default LL(k) implementation.
**`ParserState.Lookahead` is deleted** — its logic moves into `LLkStrategy.Predict`/`PredictOpt`.
`LLkLookahead`, `LookaheadPath`, `LookaheadOption` types stay (used as table values).
`LA`/`LAId` stay on `ParserState` (also used by `allstar` via `TokenSource`).

```go
// LookaheadStrategy abstracts OR-decision prediction for generated parsers.
// key is the ATN decision key ("RuleName_ProdType_N", 1-based).
// Predict returns the chosen alternative index (0-based), or -1.
// PredictOpt returns true when the optional / loop body should be entered.
type LookaheadStrategy interface {
    Predict(src *ParserState, key string) int
    PredictOpt(src *ParserState, key string) bool
}

// LLkStrategy implements LookaheadStrategy using pre-built LL(k) tables.
// tables maps ATN decision key → LLkLookahead table.
// It is the default strategy installed by generated NewParser constructors.
type LLkStrategy struct {
    tables map[string]LLkLookahead
}

func NewLLkStrategy(tables map[string]LLkLookahead) *LLkStrategy {
    return &LLkStrategy{tables: tables}
}

func (s *LLkStrategy) Predict(src *ParserState, key string) int {
    t, ok := s.tables[key]
    if !ok {
        return -1
    }
    for i, option := range t {
    outer:
        for _, path := range option {
            for j, tokenType := range path {
                if src.LAId(j+1) != tokenType {
                    continue outer
                }
            }
            return i
        }
    }
    return -1
}

func (s *LLkStrategy) PredictOpt(src *ParserState, key string) bool {
    return s.Predict(src, key) == 0
}
```

`ParserState.Lookahead` is **removed** — all former call sites become `p.predict("key")`.

---

## 2. `parser/allstar/parser_bridge.go` — adapter (new file)

`allstar` imports `parser` (no cycle: `parser` does not import `allstar`).

```go
// parserBridge wraps LLStarLookahead to satisfy parser.LookaheadStrategy.
// *parser.ParserState satisfies allstar.TokenSource via its LA method.
type parserBridge struct{ *LLStarLookahead }

func (b *parserBridge) Predict(src *parser.ParserState, key string) int {
    return b.PredictAlternation(src, key)
}
func (b *parserBridge) PredictOpt(src *parser.ParserState, key string) bool {
    return b.PredictOptional(src, key)
}

// AsParserStrategy returns a parser.LookaheadStrategy backed by this ALL(*) engine.
func (s *LLStarLookahead) AsParserStrategy() parser.LookaheadStrategy {
    return &parserBridge{s}
}
```

---

## 3. `parser/allstar/convert.go` — pre-order Alternation Idx

`convertAlternatives` assigns `Alternation.Idx` **before** recursing into children so the
occurrence index matches the generator's pre-order traversal (done already).

---

## 4. Code generator — construction-time strategy on generated `Parser`

### 4.1 Generated `Parser` struct and constructor

`NewParser` builds an `LLkStrategy` with every lookahead table keyed by its ATN key.
`strategy` is never nil.

```go
type Parser struct {
    state    *parser.ParserState
    srv      GrammarGeneratedSrvCont
    strategy parser.LookaheadStrategy
}

func NewParser(srv GrammarGeneratedSrvCont) *Parser {
    return &Parser{
        srv: srv,
        strategy: parser.NewLLkStrategy(map[string]parser.LLkLookahead{
            "Rule_Alternation_1": rule_alternation_1_lookahead,
            "Rule_Option_1":      rule_option_1_lookahead,
            // … one entry per generated LookaheadValue with its atnKey
        }),
    }
}

// WithLookaheadStrategy replaces the strategy and returns the parser for
// construction-time chaining. Call before the first Parse invocation.
func (p *Parser) WithLookaheadStrategy(s parser.LookaheadStrategy) *Parser {
    p.strategy = s
    return p
}
```

### 4.2 Private dispatch helpers (generated, not part of any public interface)

No `llk` parameter — the strategy already owns all tables.

```go
func (p *Parser) predict(key string) int {
    return p.strategy.Predict(p.state, key)
}

func (p *Parser) predictOpt(key string) bool {
    return p.strategy.PredictOpt(p.state, key)
}
```

### 4.3 `Parse` — propagate strategy to per-call copy

```go
cp := &Parser{srv: p.srv, strategy: p.strategy, state: parser.NewParserState(document.Tokens)}
```

### 4.4 OR-switch calls use `predict`; optional/loop guards use `predictOpt`

| Site | Before | After |
| --- | --- | --- |
| `generateAlternativesParser` switch | `p.state.Lookahead(TableVar)` | `p.predict("Key")` |
| `generateAssignableAlternatives` switch | `p.state.Lookahead(TableVar)` | `p.predict("Key")` |
| `generateGroupParser` guard | `p.state.Lookahead(TableVar) == 0` | `p.predictOpt("Key")` |
| `generateRuleCallParser` guard | `p.state.Lookahead(TableVar) == 0` | `p.predictOpt("Key")` |
| `generateKeywordParser` guard | inline LA check | unchanged (direct token check) |

### 4.5 ATN decision keys

Keys for OR decisions: `"RuleName_Alternation_N"` (1-based, pre-order).
Keys for optional/loop: `"RuleName_Option_N"`, `"RuleName_Repetition_N"`, or `"RuleName_RepetitionMandatory_N"`.

Generator additions:
- `LookaheadValue.atnKey string` field.
- `ParserGeneratorContext.currentRule string` — set per rule in `populateContext`.
- `ParserGeneratorContext.alternationCounters map[string]int` — pre-order OR counter per rule.
- `ParserGeneratorContext.optionCounters map[string]int` — per rule per cardinality type.
- In `populateContextWithNode`: compute and store `atnKey` for each `orLookaheads` and
  `lookaheads` entry.
- `populateContext` collects all `LookaheadValue` entries (both slices) to emit the
  `NewLLkStrategy` map in the constructor.

---

## 5. `internal/grammar/parser_gen.go` — regenerate

Run the generator after step 4. All parser tests pass because `LLkStrategy.Predict` inlines the
same table-walk logic as the removed `ParserState.Lookahead`.

---

## 6. End-user example

```go
// LL(k) — default, unchanged behaviour
p := grammar.NewParser(srv)

// ALL(*) — opt-in at construction time
rules, _ := allstar.FromParserRules(grammarRules, tokenTypes)
p := grammar.NewParser(srv).WithLookaheadStrategy(
    allstar.NewLLStarLookahead(rules, nil).AsParserStrategy(),
)
```
