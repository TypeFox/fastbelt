# PLAN-TEST.md — LL(k) vs ALL(*) Strategy Comparison Tests

## Goal

Demonstrate that grammars requiring unbounded lookahead succeed with ALL(*) but fail
with LL(k) for any fixed k. A single `Strategy` interface lets callers choose which
prediction algorithm is used at parser-construction time.

---

## 1. New Public API — `Strategy` interface in `parser/allstar/strategy.go`

```go
// Strategy is the prediction interface shared by LL(k) and ALL(*) algorithms.
type Strategy interface {
    PredictAlternation(src TokenSource, key string) int
    PredictOptional(src TokenSource, key string) bool
}
```

Add two methods to `LLStarLookahead` so it naturally satisfies `Strategy`:

| Method | What it does |
| --- | --- |
| `PredictAlternation(src, key) int` | Looks up the ATN decision state by key; calls `AdaptivePredict` |
| `PredictOptional(src, key) bool` | Same lookup; returns `alt == 0` (enter optional block) |

---

## 2. Test-internal LL(k) implementation — `compare_test.go`

`llkStrategy` (unexported, test-only) is a static LL(1) strategy backed by
`map[string]map[int]int` — decision key → first-token-ID → alt index.

It uses the same `Strategy` interface. When a token maps to multiple alts (conflict)
it returns -1.

---

## 3. Mini Parser

`streamParser` in the test file wraps a `*parser.ParserState` and a `Strategy`.
It drives actual token consumption via `ParserState.Consume()` so the test exercises
the full parsing loop, not just prediction.

```
newStreamParser(input string, strategy Strategy) *streamParser
```

---

## 4. Test Scenarios

All four grammars share token alphabet:  `a`=1  `b`=2  `c`=3

### Scenario 1 — Greedy alternation (`OR(A+, A+·B)`)

```
LongRule := OR(
    alt 0: A+
    alt 1: A+  B
)
```

| Input | LL(1) result | ALL(*) result |
| --- | --- | --- |
| "a a b" | -1 (conflict: 'a' maps to both alts) | 1 ✓ |
| "a a a" | -1 (conflict: 'a' maps to both alts) | 0 ✓ |

### Scenario 2 — Sub-rule indirection (`OR(SUBRULE(Ra), SUBRULE(Rb))`)

```
AltRule := OR(SUBRULE(Ra), SUBRULE(Rb))
Ra      := A+
Rb      := A+  B
```

Same disambiguation challenge, but prediction must enter sub-rules first.

| Input | LL(1) result | ALL(*) result |
| --- | --- | --- |
| "a a b" | -1 | 1 ✓ |
| "a a a" | -1 | 0 ✓ |

### Scenario 3 — Zero-or-more common prefix (`OR(A*, A*·C)`)

```
ManyRule := OR(
    alt 0: A*
    alt 1: A*  C
)
```

| Input | LL(1) result | ALL(*) result |
| --- | --- | --- |
| "a a c" | -1 (conflict: 'a' maps to both alts) | 1 ✓ |
| "a a"   | -1 (conflict: 'a' maps to both alts) | 0 ✓ |

### Scenario 4 — Option greediness (`OPTION(A+) · A+`)

```
OptionRule := OPTION(A+)  A+
```

With input "a" only the option-skipping path can succeed: taking the option
consumes the single A, leaving nothing for the mandatory A+ that follows.

| Input | LL(1) result (greedy) | ALL(*) result |
| --- | --- | --- |
| "a" | 0 (takes option — wrong, leaves outer A+ with no input) | 1 (skips option ✓) |
| "a a a" | 0 (takes option — ambiguous but accepted) | 0 or 1 (ambiguous) |

In scenario 4, LL(1) returns the wrong (greedy) answer for single-token input;
ALL(*) returns the only path that can succeed.

---

## 5. File Layout

```
parser/allstar/strategy.go        # add Strategy interface + methods on LLStarLookahead
parser/allstar/compare_test.go    # llkStrategy + streamParser + 4 × test functions
```

---

## 6. Test Structure (each test)

```
// arrange
tokens := tokenize("a a b")
grammar := build<X>Grammar()                     // allstar []*Rule

llk    := newLLkFor<X>(...)                      // static LL(1) table
allStr := allstar.NewLLStarLookahead(grammar, nil) // adaptive

// act
errLLk    := newStreamParser(tokens, llk).parse<X>()
errALLStar:= newStreamParser(tokens, allStr).parse<X>()

// assert
assert.Error(t, errLLk)      // LL(k) fails
assert.NoError(t, errALLStar) // ALL(*) succeeds
```

Tests are table-driven within each scenario to cover both the "B-terminated" and
"plain-repetition" inputs with a single grammar definition.
