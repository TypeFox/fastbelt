# PLAN-ALLSTAR.md — Migration of `chevrotain-allstar` to Go

> **Source:** `chevrotain-allstar-main/` (TypeScript, ~1 500 lines)
> **Target:** `parser/allstar/` package inside `typefox.dev/fastbelt` (Go)
> **Algorithm:** ALL(*) adaptive lookahead — the LL(∞) prediction algorithm from ANTLR4
> (paper: <https://www.antlr.org/papers/allstar-techreport.pdf>)

---

## 1. Background

The TypeScript library implements the ALL(*) lookahead algorithm as a plug-in strategy for the
[Chevrotain](https://chevrotain.io/) parser library. It replaces Chevrotain's fixed LL(k) lookahead
with an unbounded, DFA-backed adaptive prediction that can handle grammars that are not LL(k) for
any fixed k.

The algorithm has three components that map directly to source files:

| TypeScript file | Responsibility |
| --- | --- |
| `atn.ts` | Build an Augmented Transition Network (ATN) from parser rules |
| `dfa.ts` | DFA state machine used to memoize prediction results |
| `all-star-lookahead.ts` | Adaptive prediction algorithm + lookahead strategy |

The existing Go project (`typefox.dev/fastbelt`) already has its own parser framework. Its
`parser/parser.go` uses a static LL(k) lookahead table (`LLkLookahead`). This migration provides a
new `allstar` sub-package that can serve as the backing engine for a more powerful lookahead.

---

## 2. New Package Layout

```text
fastbelt/
└── parser/
    └── allstar/
        ├── grammar.go          # Minimal production/rule types consumed by the ATN builder
        ├── convert.go          # Conversion from internal/grammar AST → allstar production types
        ├── atn.go              # ATN data structures + construction algorithm
        ├── dfa.go              # DFA data structures + config set management
        ├── predict.go          # Adaptive prediction algorithm (closures, DFA cache, ambiguity)
        ├── strategy.go         # LLStarLookahead: top-level entry point + LL1 fast-path
        ├── convert_test.go     # Unit tests: conversion from internal/grammar
        ├── atn_test.go         # Unit tests: ATN construction
        ├── dfa_test.go         # Unit tests: DFA / ATNConfigSet
        ├── predict_test.go     # Unit tests: prediction algorithm internals
        └── integration_test.go # Integration tests: input string → parsed result
```

The package declaration is `package allstar`.
The import path is `typefox.dev/fastbelt/parser/allstar`.

---

## 3. Type Model

### 3.1 Grammar types — `grammar.go`

The project already has a grammar AST in `internal/grammar/types_gen.go`, but that package serves
a different purpose: it is the **meta-grammar AST** for `.fb` grammar files, used by the code
generator. The `allstar` package needs a simpler, **ATN-oriented production model** whose types
map 1-to-1 to ATN construction operations and carry the occurrence indices that ATN keys require.

The conversion from `internal/grammar` types to these types happens in `convert.go` (§3.5).
The `.fb` grammar language has no separator-based repetition, so `RepetitionWithSeparator` and
`RepetitionMandatoryWithSeparator` are omitted.

```go
// ProductionKind is a discriminator for the Production union.
type ProductionKind int

const (
    ProdTerminal             ProductionKind = iota
    ProdNonTerminal
    ProdAlternative          // one branch inside an Alternation; no occurrence index
    ProdAlternation          // OR(...)
    ProdOption               // OPTION(...)
    ProdRepetition           // MANY(...)
    ProdRepetitionMandatory  // AT_LEAST_ONE(...)
)

// Production is the interface for all grammar elements.
type Production interface {
    Kind() ProductionKind
    // Children returns the sub-productions, or nil for leaves.
    Children() []Production
    // Occurrence returns the 1-based index distinguishing multiple uses of the
    // same production kind within one rule (0 for Alternative which is not indexed).
    Occurrence() int
}

// Rule is a named parser rule.
type Rule struct {
    Name       string
    Definition []Production
}

// Terminal consumes a single token.
type Terminal struct {
    TokenTypeID     int
    TokenName       string
    Idx             int
    // CategoryMatches holds token type IDs matched via category inheritance.
    CategoryMatches []int
}

// NonTerminal calls another rule.
type NonTerminal struct {
    ReferencedRule *Rule
    Idx            int
}

// Alternative is one branch inside an Alternation (no own occurrence index).
type Alternative struct {
    Definition []Production
}

// Alternation is an OR decision point.
type Alternation struct {
    Alternatives []*Alternative
    Idx          int
}

// Option wraps an optional sequence.
type Option struct {
    Definition []Production
    Idx        int
}

// Repetition is a zero-or-more loop.
type Repetition struct {
    Definition []Production
    Idx        int
}

// RepetitionMandatory is a one-or-more loop.
type RepetitionMandatory struct {
    Definition []Production
    Idx        int
}
```

Each struct implements `Production`. `Terminal` and `NonTerminal` return `nil` from `Children()`.
`Alternative` and `Alternation` implement `Children()` by flattening their definitions.

The ATN key string used for decision-state lookup. Returns `("", false)` for production kinds
that have no ATN key (Terminal, NonTerminal, Alternative) rather than panicking, following
Effective Go's preference for explicit error returns over panics in library code:

```go
// ProductionTypeName returns the ATN decision-map key segment for p, and
// whether p is a decision-producing kind at all.
func ProductionTypeName(p Production) (string, bool) {
    switch p.Kind() {
    case ProdAlternation:
        return "Alternation", true
    case ProdOption:
        return "Option", true
    case ProdRepetition:
        return "Repetition", true
    case ProdRepetitionMandatory:
        return "RepetitionMandatory", true
    default:
        return "", false
    }
}
```

---

### 3.2 ATN types — `atn.go`

The TypeScript code uses a discriminated union (`type ATNState = BasicState | RuleStartState | …`).
In Go we use a single concrete struct with a `Type` field that acts as the discriminator, because
the algorithm only ever switches on the type constant anyway:

```go
type ATNStateType int

const (
    ATNInvalidType    ATNStateType = 0
    ATNBasic          ATNStateType = 1
    ATNRuleStart      ATNStateType = 2
    ATNPlusBlockStart ATNStateType = 4
    ATNStarBlockStart ATNStateType = 5
    ATNTokenStart     ATNStateType = 6
    ATNRuleStop       ATNStateType = 7
    ATNBlockEnd       ATNStateType = 8
    ATNStarLoopBack   ATNStateType = 9
    ATNStarLoopEntry  ATNStateType = 10
    ATNPlusLoopBack   ATNStateType = 11
    ATNLoopEnd        ATNStateType = 12
)

// ATNState is the single concrete ATN state type.
// Fields specific to certain state kinds are non-nil only for those kinds.
type ATNState struct {
    ATN                    *ATN
    Production             Production // nil for rule start/stop
    StateNumber            int
    Rule                   *Rule
    EpsilonOnlyTransitions bool
    Transitions            []Transition
    Type                   ATNStateType

    // Decision index; -1 if this state is not a decision state.
    Decision int

    // Populated for BlockStartState kinds:
    End      *ATNState
    Loopback *ATNState // PlusBlockStart.loopback, StarLoopEntry.loopback, LoopEnd.loopback

    // Populated for BlockEndState:
    Start *ATNState

    // Populated for RuleStartState:
    Stop *ATNState
}

type ATN struct {
    DecisionMap      map[string]*ATNState
    States           []*ATNState
    DecisionStates   []*ATNState
    RuleToStartState map[*Rule]*ATNState
    RuleToStopState  map[*Rule]*ATNState
}
```

Transitions:

```go
type Transition interface {
    Target() *ATNState
    IsEpsilon() bool
}

// AtomTransition fires on a specific token type.
// CategoryMatches holds the IDs of all token types that match via category
// inheritance; populated from the Terminal's CategoryMatches at ATN-build time.
type AtomTransition struct {
    target          *ATNState
    TokenTypeID     int
    CategoryMatches []int
}

// EpsilonTransition fires without consuming a token.
type EpsilonTransition struct {
    target *ATNState
}

// RuleTransition enters a sub-rule and returns to FollowState.
type RuleTransition struct {
    target      *ATNState // the rule's RuleStartState
    Rule        *Rule
    FollowState *ATNState
}
```

`EpsilonTransition` and `RuleTransition` return `true` from `IsEpsilon()`.

**ATN construction functions** (all unexported except `CreateATN` and `BuildATNKey`):

| TypeScript function | Go function signature |
| --- | --- |
| `buildATNKey` | `BuildATNKey(rule *Rule, prodType string, occurrence int) string` |
| `createATN` | `CreateATN(rules []*Rule) *ATN` |
| `createRuleStartAndStopATNStates` | `createRuleStartAndStopATNStates(atn *ATN, rules []*Rule)` |
| `atom` | `atom(atn *ATN, rule *Rule, prod Production) *atnHandle` |
| `repetition` | `repetition(atn *ATN, rule *Rule, prod *Repetition) *atnHandle` |
| `repetitionMandatory` | `repetitionMandatory(atn *ATN, rule *Rule, prod *RepetitionMandatory) *atnHandle` |
| `alternation` | `alternation(atn *ATN, rule *Rule, prod *Alternation) *atnHandle` |
| `option` | `option(atn *ATN, rule *Rule, prod *Option) *atnHandle` |
| `block` | `block(atn *ATN, rule *Rule, children []Production) *atnHandle` |
| `plus` | `plus(atn *ATN, rule *Rule, prod Production, handle *atnHandle) *atnHandle` |
| `star` | `star(atn *ATN, rule *Rule, prod Production, handle *atnHandle) *atnHandle` |
| `optional` | `optional(atn *ATN, rule *Rule, prod *Option, handle *atnHandle) *atnHandle` |
| `makeAlts` | `makeAlts(atn *ATN, rule *Rule, start *ATNState, prod Production, alts ...*atnHandle) *atnHandle` |
| `makeBlock` | `makeBlock(atn *ATN, alts []*atnHandle) *atnHandle` |
| `tokenRef` | `tokenRef(atn *ATN, rule *Rule, tokenTypeID int, prod Production) *atnHandle` |
| `ruleRef` | `ruleRef(atn *ATN, currentRule *Rule, nt *NonTerminal) *atnHandle` |
| `buildRuleHandle` | `buildRuleHandle(atn *ATN, rule *Rule, block *atnHandle)` |
| `epsilon` | `addEpsilon(from, to *ATNState)` |
| `newState` | `newATNState(atn *ATN, rule *Rule, prod Production, typ ATNStateType) *ATNState` |
| `addTransition` | `addTransition(state *ATNState, t Transition)` |
| `removeState` | `removeState(atn *ATN, state *ATNState)` |
| `defineDecisionState` | `defineDecisionState(atn *ATN, state *ATNState) int` |

`atnHandle` is a private struct `{ left, right *ATNState }`.

The `sep` parameter is dropped from `plus` and `star` compared to the TypeScript originals because
separator-based repetition is not needed.

---

### 3.3 DFA types — `dfa.go`

```go
// DFA is the memoization automaton for a single decision.
type DFA struct {
    Start         *DFAState
    States        map[string]*DFAState
    Decision      int
    ATNStartState *ATNState
}

// DFAState represents a set of ATN configurations.
type DFAState struct {
    Configs       *ATNConfigSet
    Edges         map[int]*DFAState
    IsAcceptState bool
    Prediction    int // valid when IsAcceptState == true
}

// DFAError is a sentinel value returned for dead-end states.
var DFAError = &DFAState{}

// ATNConfig is one thread of the ATN simulation.
type ATNConfig struct {
    State *ATNState
    Alt   int
    Stack []*ATNState
}

// ATNConfigSet holds a deduplicated set of ATN configurations.
type ATNConfigSet struct {
    configMap map[string]int
    configs   []*ATNConfig
    UniqueAlt int // -1 means "no unique alt determined yet"
}
```

Public methods on `ATNConfigSet`:

| TypeScript | Go |
| --- | --- |
| `add(config)` | `func (s *ATNConfigSet) Add(c *ATNConfig)` |
| `finalize()` | `func (s *ATNConfigSet) Finalize()` |
| `get size()` | `func (s *ATNConfigSet) Len() int` |
| `get elements()` | `func (s *ATNConfigSet) Elements() []*ATNConfig` |
| `get alts()` | `func (s *ATNConfigSet) Alts() []int` |
| `get key()` | `func (s *ATNConfigSet) Key() string` |

```go
// atnConfigKey produces the deduplication key for a config.
// When includeAlt is false the alt index is omitted (used for conflict detection).
// Unexported: only used within the package.
func atnConfigKey(c *ATNConfig, includeAlt bool) string
```

---

### 3.4 Prediction / strategy types — `predict.go`, `strategy.go`

```go
// AmbiguityReport is a callback invoked when an ambiguity is detected.
type AmbiguityReport func(message string)

// PredicateSet records which alternatives are guarded and whether their
// gate predicate currently evaluates to true.
type PredicateSet struct {
    predicates []bool
}

// Is(index) returns true when the index is out of range (unconstrained)
// or when predicates[index] == true.
func (p *PredicateSet) Is(index int) bool
func (p *PredicateSet) Set(index int, value bool)
func (p *PredicateSet) String() string

// EmptyPredicates is the zero-value PredicateSet used when there are no gates.
var EmptyPredicates = &PredicateSet{}

// dfaCache returns a DFA for the given predicate configuration.
// Different predicate sets may produce different prediction decisions.
// Unexported: used only within the package as an element of LLStarLookahead.dfas.
type dfaCache func(predicates *PredicateSet) *DFA

// LLStarLookaheadOptions configures the strategy.
type LLStarLookaheadOptions struct {
    // Logging is called whenever an ambiguity is detected.
    // Defaults to a function that prints to stdout.
    Logging AmbiguityReport
}

// LLStarLookahead is the main entry point. It holds the ATN and the DFA caches
// produced during the Initialize step.
type LLStarLookahead struct {
    atn     *ATN
    dfas    []dfaCache
    logging AmbiguityReport
}

// TokenSource is the minimal interface required by the prediction algorithm.
// *parser.ParserState satisfies this interface.
type TokenSource interface {
    LA(offset int) *core.Token
}
```

Key exported functions:

```go
// NewLLStarLookahead creates a new strategy and calls Initialize.
func NewLLStarLookahead(rules []*Rule, opts *LLStarLookaheadOptions) *LLStarLookahead

// AdaptivePredict runs the ALL(*) algorithm.
// Returns the chosen alternative index (0-based), or -1 on parse error.
func (s *LLStarLookahead) AdaptivePredict(
    src TokenSource,
    decision int,
    predicates *PredicateSet,
) int

// BuildLookaheadForAlternation returns a function that wraps AdaptivePredict
// for a specific alternation occurrence in a rule. Falls back to an LL(1)
// table when the grammar is deterministic at depth 1.
func (s *LLStarLookahead) BuildLookaheadForAlternation(
    rule *Rule,
    occurrence int,
    hasPredicates bool,
) func(src TokenSource, gates []func() bool) int

// BuildLookaheadForOptional returns a function that wraps AdaptivePredict
// for OPTION, MANY, and AT_LEAST_ONE productions.
func (s *LLStarLookahead) BuildLookaheadForOptional(
    rule *Rule,
    occurrence int,
    prodType string,
) func(src TokenSource) bool
```

Internal functions in `predict.go` (unexported):

| TypeScript | Go |
| --- | --- |
| `createDFACache` | `newDFACache(start *ATNState, decision int) dfaCache` |
| `initATNSimulator` | `initDFACaches(atn *ATN) []dfaCache` |
| `adaptivePredict` | `adaptivePredict(src TokenSource, dfas []dfaCache, decision int, preds *PredicateSet, log AmbiguityReport) (int, error)` |
| `performLookahead` | `performLookahead(src TokenSource, dfa *DFA, s0 *DFAState, preds *PredicateSet, log AmbiguityReport) (int, error)` |
| `computeLookaheadTarget` | `computeLookaheadTarget(src TokenSource, dfa *DFA, prev *DFAState, tokenTypeID int, lookahead int, preds *PredicateSet, log AmbiguityReport) *DFAState` |
| `getExistingTargetState` | `getExistingTargetState(state *DFAState, tokenTypeID int) *DFAState` |
| `computeReachSet` | `computeReachSet(configs *ATNConfigSet, tokenTypeID int, preds *PredicateSet) *ATNConfigSet` |
| `getReachableTarget` | `getReachableTarget(t Transition, tokenTypeID int) *ATNState` |
| `closure` | `closure(config *ATNConfig, configs *ATNConfigSet)` |
| `getEpsilonTarget` | `getEpsilonTarget(config *ATNConfig, t Transition) *ATNConfig` |
| `computeStartState` | `computeStartState(atnState *ATNState) *ATNConfigSet` |
| `getUniqueAlt` | `getUniqueAlt(configs *ATNConfigSet, preds *PredicateSet) (int, bool)` |
| `hasConflictTerminatingPrediction` | `hasConflictTerminatingPrediction(configs *ATNConfigSet) bool` |
| `allConfigsInRuleStopStates` | `allConfigsInRuleStopStates(configs *ATNConfigSet) bool` |
| `hasConfigInRuleStopState` | `hasConfigInRuleStopState(configs *ATNConfigSet) bool` |
| `getConflictingAltSets` | `getConflictingAltSets(configs []*ATNConfig) map[string]map[int]bool` |
| `hasConflictingAltSet` | `hasConflictingAltSet(altSets map[string]map[int]bool) bool` |
| `hasStateAssociatedWithOneAlt` | `hasStateAssociatedWithOneAlt(altSets map[string]map[int]bool) bool` |
| `newDFAState` | `newDFAState(configs *ATNConfigSet) *DFAState` |
| `addDFAEdge` | `addDFAEdge(dfa *DFA, from *DFAState, tokenTypeID int, to *DFAState) *DFAState` |
| `addDFAState` | `addDFAState(dfa *DFA, state *DFAState) *DFAState` |
| `reportLookaheadAmbiguity` | `reportLookaheadAmbiguity(src TokenSource, dfa *DFA, lookahead int, alts []int, log AmbiguityReport)` |
| `buildAmbiguityError` | `buildAmbiguityError(rule *Rule, prod Production, prefixPath []string, alts []int) string` |
| `buildAdaptivePredictError` | `newPredictError(path []int, prev *DFAState, tokenTypeID int) error` |
| `isLL1Sequence` | `isLL1Sequence(seqs [][]tokenInfo, allowEmpty bool) bool` |
| `getProductionDslName` | `productionDSLName(prod Production) string` |

The `predictError` type replaces `AdaptivePredictError`. It is unexported but implements the
standard `error` interface, making it compatible with all Go error-handling idioms:

```go
type predictError struct {
    tokenPath         []int // token type IDs of consumed lookahead
    possibleTypeIDs   []int
    actualTokenTypeID int
}

// Error implements the error interface.
func (e *predictError) Error() string
```

Callers that need to inspect the fields (e.g. for detailed diagnostics) use `errors.As`.

---

### 3.5 Conversion layer — `convert.go`

This file bridges `internal/grammar` (the code-generator's meta-grammar AST) and the ATN-oriented
types above. It is the only place in `allstar` that imports `internal/grammar`.

**Main entry point:**

```go
// FromParserRules converts a slice of grammar.ParserRule into the allstar Rule
// slice that CreateATN expects. tokenTypes maps terminal name → TokenTypeID.
func FromParserRules(
    rules []grammar.ParserRule,
    tokenTypes map[string]TokenInfo,
) ([]*Rule, error)
```

Where `TokenInfo` carries the ID and category-match IDs for a token type:

```go
type TokenInfo struct {
    ID              int
    CategoryMatches []int
}
```

**Conversion algorithm:**

Walking a `grammar.ParserRule`:

1. Create `counters` — a `map[ProductionKind]int` scoped to the current rule, reset per rule.
2. Walk `rule.Body()` recursively via `convertElement(el grammar.Element, counters map[ProductionKind]int)`.

Mapping rules for `convertElement`:

| `internal/grammar` type | Condition | `allstar` production | Counter bumped |
| --- | --- | --- | --- |
| `grammar.Alternatives` | — | `Alternation{Idx: next(ProdAlternation)}` | `ProdAlternation` |
| `grammar.Group` | `Cardinality() == "?"` | `Option{Idx: next(ProdOption)}` | `ProdOption` |
| `grammar.Group` | `Cardinality() == "*"` | `Repetition{Idx: next(ProdRepetition)}` | `ProdRepetition` |
| `grammar.Group` | `Cardinality() == "+"` | `RepetitionMandatory{Idx: next(ProdRepetitionMandatory)}` | `ProdRepetitionMandatory` |
| `grammar.Group` | no cardinality | inline sequence (returns `[]Production`) | — |
| `grammar.RuleCall` | no cardinality | `NonTerminal{Idx: next(ProdNonTerminal)}` | `ProdNonTerminal` |
| `grammar.RuleCall` | with cardinality | wrap in `Option`/`Repetition`/`RepetitionMandatory` first | both |
| `grammar.Keyword` | — | `Terminal{TokenName: value, Idx: next(ProdTerminal)}` | `ProdTerminal` |
| `grammar.Assignment` | — | transparent — recurse into `Value()` | — |
| `grammar.CrossRef` | — | `Terminal` for the cross-reference token | `ProdTerminal` |
| `grammar.Action` | — | skip (semantic action, no ATN impact) | — |

`next(kind)` returns the current counter value for `kind` and then increments it (1-based).

The token type IDs for `Terminal` nodes are resolved from the `tokenTypes` map keyed on the
terminal name. An error is returned for any unresolved terminal.

Rule references in `NonTerminal` are resolved in a second pass once all `Rule` objects are
created, to handle forward references between rules.

---

## 4. Algorithm Notes and Go-specific Adaptations

### Effective Go conventions

All implementation files must follow <https://go.dev/doc/effective_go>.

#### Naming

- All identifiers use MixedCaps / mixedCaps — no underscores anywhere, including constants
  (`ATNBasic`, `ProdAlternation`, not `ATN_BASIC` or `PROD_ALTERNATION`).
- Constructor functions use the `New` prefix (`NewLLStarLookahead`, `newPredictError`).
- The `Get` prefix on accessor functions is non-idiomatic and is omitted throughout
  (e.g. `atnConfigKey`, not `GetATNConfigKey`).

#### Export surface

Only identifiers needed by callers outside the package are exported. The table below summarises
the decision for every type and function in the package:

| Symbol | Exported? | Reason |
| --- | --- | --- |
| `Rule`, `Terminal`, `NonTerminal`, `Alternative`, `Alternation`, `Option`, `Repetition`, `RepetitionMandatory` | ✅ | Grammar model used by callers to build rules |
| `Production`, `ProductionKind`, `ProductionTypeName` | ✅ | Part of public grammar model |
| `TokenInfo`, `FromParserRules` | ✅ | Public API of `convert.go` |
| `LLStarLookahead`, `LLStarLookaheadOptions`, `AmbiguityReport` | ✅ | Public strategy API |
| `TokenSource`, `PredicateSet`, `EmptyPredicates` | ✅ | Required by callers invoking `AdaptivePredict` |
| `DFAError` | ✅ | Exported sentinel (compared by pointer in tests) |
| `ATN`, `ATNState`, `ATNStateType` constants | ✅ | Needed by white-box unit tests in `atn_test.go` |
| `DFA`, `DFAState`, `ATNConfig`, `ATNConfigSet` | ✅ | Needed by white-box unit tests |
| `Transition`, `AtomTransition`, `EpsilonTransition`, `RuleTransition` | ✅ | Needed by white-box unit tests |
| `dfaCache` | ❌ | Implementation detail; callers never reference it |
| `atnConfigKey` | ❌ | Package-internal helper |
| `predictError` | ❌ | Unexported error type; exposed only via the `error` interface |
| All ATN/DFA construction helpers | ❌ | Package-internal |
| All prediction algorithm helpers | ❌ | Package-internal |

Test files are placed in `package allstar` (not `package allstar_test`) so that they can access
unexported helpers directly without the round-trip of exporting them just for tests.

#### Error handling

- All functions that can fail return `error` as their last result.
- `predictError` is an unexported concrete type that implements `error`. Callers that need field
  access use `errors.As`.
- Internal helpers may use `panic` only for true programming errors (e.g. an impossible branch
  reached due to a bug), never for input errors. `ProductionTypeName` returns `(string, bool)`
  rather than panicking for unrecognised production kinds.

#### Interfaces

All interfaces are small (1–3 methods): `TokenSource` (1), `Transition` (2), `Production` (3).
The "accept interfaces, return concrete types" rule is followed throughout: factory functions
return `*LLStarLookahead`, `*ATN`, `*DFAState`, etc.

#### `ATNState.Decision` zero value

The zero value of `int` is `0`, which is a valid decision index. New `ATNState` values must be
initialised with `Decision: -1` to indicate "not a decision state". `newATNState` sets this.

#### Formatting

All files must be formatted with `gofmt` before commit. The project's existing files use standard
Go tab-indented style; no deviations are permitted.

---

### Token matching

TypeScript uses `tokenMatcher(token, transition.tokenType)` from Chevrotain. This checks whether a
token matches via token-category inheritance. In Go, `getReachableTarget` needs to check:

```go
func tokenMatches(tokenTypeID int, transitionTypeID int, categoryMatches []int) bool {
    if tokenTypeID == transitionTypeID {
        return true
    }
    for _, cat := range categoryMatches {
        if tokenTypeID == cat {
            return true
        }
    }
    return false
}
```

The `CategoryMatches []int` field is carried on `AtomTransition` and populated from the source
`Terminal.CategoryMatches` at ATN-build time. `core.TokenType` is not modified.

### DFA_ERROR sentinel

TypeScript uses `export const DFA_ERROR = {} as DFAState` — a singleton object identity check.
In Go this is `var DFAError = &DFAState{}`. The check `state === DFA_ERROR` becomes
`state == DFAError` (pointer equality).

### LL(1) fast path

`buildLookaheadForAlternation` in TypeScript builds a simple lookup table when `isLL1Sequence` is
true. The Go equivalent does the same: it builds `map[int]int` (tokenTypeID → alternativeIndex)
and returns a simple function that performs a single `src.LA(1)` lookup.

### Memoization and DFA caches

`DFACache` in TypeScript is a closure over a `Record<string, DFA>`. In Go it is a
`func(*PredicateSet) *DFA` backed by a `map[string]*DFA` protected by a `sync.RWMutex`.
Reads (cache hit) acquire only a read lock; writes (cache miss) upgrade to a write lock.
This allows multiple goroutines to parse concurrently once the DFA is warm.

### EOF token

`TokenSource.LA` must return a non-nil EOF sentinel token (a `*core.Token` with
`TypeId == core.EOF.Id`) when the stream is exhausted. `parser.ParserState` already does this.
The prediction algorithm never needs to handle a `nil` return from `LA`.

---

## 5. Implementation Order

Each step should be followed by its tests before moving on.

### Step 1 — Grammar model (`grammar.go`)

Implement the production types and their `Kind()`, `Children()`, `Occurrence()` methods.
Implement `ProductionTypeName(Production) string`.

No tests for this step alone — correctness is covered by the conversion tests in Step 1b.

---

### Step 1b — Conversion layer (`convert.go` + `convert_test.go`)

Implement `FromParserRules` and `convertElement` following the mapping table in §3.5.
Implement the two-pass forward-reference resolution for `NonTerminal.ReferencedRule`.

**Tests** (see §6.1) verify that a representative set of `internal/grammar` rule trees
are converted to the expected `allstar` production trees with correct occurrence indices.

---

### Step 2 — ATN data structures (`atn.go` — types only)

Implement `ATNState`, `ATN`, `Transition` types and the three `Transition` implementations
(`AtomTransition`, `EpsilonTransition`, `RuleTransition`).
Implement `BuildATNKey`.

**Tests** verify that `BuildATNKey` produces the expected format (`ruleName_prodType_idx`).

---

### Step 3 — ATN construction algorithm (`atn.go` — functions)

Implement `CreateATN` and all private helper functions in the order they are called:
`createRuleStartAndStopATNStates` → `block` → `atom` → leaf helpers
(`tokenRef`, `ruleRef`) → composite helpers (`alternation`, `option`, `repetition`,
`repetitionMandatory`) → loop helpers (`star`, `plus`, `optional`) → structural helpers
(`makeAlts`, `makeBlock`, `buildRuleHandle`) → utility helpers (`addEpsilon`, `newATNState`,
`addTransition`, `removeState`, `defineDecisionState`).

**Tests** (see §6.2) verify the produced ATN state/transition counts and key assignments for a
representative grammar.

---

### Step 4 — DFA structures and `ATNConfigSet` (`dfa.go` + `dfa_test.go`)

Implement `ATNConfigSet` with deduplication, `GetATNConfigKey`, and the `DFAState` / `DFA` types.
Implement `DFAError` sentinel.

**Tests** (see §6.3) verify deduplication, key generation, and `Alts()`.

---

### Step 5 — Prediction algorithm (`predict.go` + `predict_test.go`)

Implement in dependency order:
`computeStartState` → `closure` → `getEpsilonTarget` → `computeReachSet` →
`getReachableTarget` → `getUniqueAlt` → conflict helpers → DFA helpers
(`newDFAState`, `addDFAState`, `addDFAEdge`) → `computeLookaheadTarget` →
`performLookahead` → `adaptivePredict` → ambiguity helpers.

**Tests** (see §6.4) verify each function in isolation using hand-constructed ATN states.

---

### Step 6 — Strategy layer (`strategy.go`)

Implement `LLStarLookahead`, `NewLLStarLookahead`, `PredicateSet`, and the two
`BuildLookaheadFor…` factory methods (including the LL(1) fast path).

---

### Step 7 — Integration tests (`integration_test.go`)

See §6.5. These are the most important validation that the whole pipeline works end-to-end.

---

## 6. Test Plan

### 6.1 Conversion layer (`convert_test.go`)

These tests build small `internal/grammar` trees by hand and assert the resulting `allstar`
production tree and occurrence indices.

| Test | What it checks |
| --- | --- |
| `TestConvert_SingleKeyword` | `Keyword("a")` → `Terminal{Idx:1}` |
| `TestConvert_RuleCall` | `RuleCall(RuleB)` → `NonTerminal{Idx:1}` pointing at RuleB |
| `TestConvert_ForwardRef` | `RuleA` calls `RuleB` defined later → reference resolved in second pass |
| `TestConvert_Alternatives` | `Alternatives([A, B])` → `Alternation{Idx:1}` with 2 `Alternative` children |
| `TestConvert_MultipleAlternations` | Two `Alternatives` in one rule → `Idx` 1 and 2 |
| `TestConvert_Group_Option` | `Group(cardinality="?", [A])` → `Option{Idx:1}` |
| `TestConvert_Group_Repetition` | `Group(cardinality="*", [A])` → `Repetition{Idx:1}` |
| `TestConvert_Group_RepetitionMandatory` | `Group(cardinality="+", [A])` → `RepetitionMandatory{Idx:1}` |
| `TestConvert_Group_Sequence` | `Group(no cardinality, [A, B])` → flat `[]Production` (no wrapper) |
| `TestConvert_Assignment_Transparent` | `Assignment(property, "=", RuleCall)` → `NonTerminal` (assignment stripped) |
| `TestConvert_RuleCall_WithCardinality` | `RuleCall(cardinality="*")` → `Repetition` wrapping `NonTerminal` |
| `TestConvert_MixedCounters` | Rule with `OR`, `OPTION`, `MANY` each twice → `Idx` 1 and 2 per kind |
| `TestConvert_UnknownToken` | Terminal name not in `tokenTypes` map → returns error |

### 6.2 ATN construction (`atn_test.go`)

| Test | What it checks |
| --- | --- |
| `TestBuildATNKey` | Correct string format for each production type |
| `TestCreateATN_SingleTerminal` | `CONSUME(A)` → 2 states, 1 atom transition, no decision states |
| `TestCreateATN_Alternation` | `OR([A, B])` → decision state, 2 epsilon exits, 1 BlockEndState |
| `TestCreateATN_Option` | `OPTION(A)` → decision state, bypass epsilon to block end |
| `TestCreateATN_Repetition` | `MANY(A)` → StarLoopEntry, StarLoopBack, LoopEnd |
| `TestCreateATN_RepetitionMandatory` | `AT_LEAST_ONE(A)` → PlusBlockStart, PlusLoopBack, LoopEnd |
| `TestCreateATN_NonTerminal` | `SUBRULE(Rule2)` → RuleTransition pointing at Rule2's start state |
| `TestCreateATN_NestedRule` | Two rules: outer calls inner, follow state is correctly wired |
| `TestCreateATN_DecisionMapKeys` | All expected keys present in `atn.DecisionMap` |
| `TestMakeBlock_Optimisation` | Consecutive basic states are merged (no extra epsilon) |

### 6.3 DFA / ATNConfigSet (`dfa_test.go`)

| Test | What it checks |
| --- | --- |
| `TestATNConfigSet_Add_Dedup` | Adding identical config twice keeps size == 1 |
| `TestATNConfigSet_Add_DifferentAlt` | Two configs with same state/stack but different alt → size == 2 |
| `TestATNConfigSet_Key_Consistency` | Key is the same after adding the same elements in the same order |
| `TestATNConfigSet_Finalize` | Finalize doesn't panic and Len() is unchanged |
| `TestATNConfigSet_Alts` | Returns correct slice of alternative indices |
| `TestATNConfigKey_WithAlt` | Key contains `a<alt>` prefix |
| `TestATNConfigKey_WithoutAlt` | Key omits alt prefix |

### 6.4 Prediction algorithm (`predict_test.go`)

| Test | What it checks |
| --- | --- |
| `TestPredicateSet_IsOutOfBounds` | `Is(100)` returns `true` (unconstrained) |
| `TestPredicateSet_IsInBounds` | `Is(0)` after `Set(0, false)` returns `false` |
| `TestPredicateSet_String` | Serialised form matches expected binary string |
| `TestComputeStartState_Simple` | From a 2-alt decision state, produces 2 configs |
| `TestClosure_Epsilon` | Epsilon transition is followed transparently |
| `TestClosure_RuleTransition` | Stack is pushed/popped correctly on enter/exit |
| `TestClosure_RuleStop_EmptyStack` | Config is added to set when stack is empty |
| `TestGetReachableTarget_Match` | AtomTransition that matches returns its target |
| `TestGetReachableTarget_NoMatch` | AtomTransition that doesn't match returns nil |
| `TestComputeReachSet_SingleAlt` | Single matching config → reach set size == 1 |
| `TestComputeReachSet_SkipsGatedAlt` | Config filtered out by predicate set |
| `TestGetUniqueAlt_Unique` | All configs have same alt → returns that alt |
| `TestGetUniqueAlt_Mixed` | Configs have different alts → no unique alt |
| `TestAllConfigsInRuleStopStates` | True only when every config is at RuleStop |
| `TestHasConflictTerminatingPrediction_AllAtStop` | True when all configs at rule stop |
| `TestHasConflictTerminatingPrediction_Conflicting` | True when multiple alts share a state |
| `TestGetConflictingAltSets` | Correct grouping by `atnConfigKey(c, false)` |
| `TestAddDFAState_Dedup` | Same config-set key returns same DFAState pointer |
| `TestAddDFAEdge` | Edge is stored in `from.Edges` |

### 6.5 Integration tests (`integration_test.go`)

These tests mirror `atn.test.ts` exactly and constitute the acceptance criteria for the migration.
They build grammars directly using the `grammar.go` types (bypassing the conversion layer),
create an `LLStarLookahead`, and drive it with a mock `TokenSource` backed by a `[]int` (token
type ID sequence).

A helper `mockTokenSource(ids ...int)` wraps a slice into a `TokenSource` that returns the token
at `LA(offset)` and an EOF token when the slice is exhausted.

#### 6.5.1 LL(*) lookahead (unbounded)

Grammar:

```text
LongRule := OR(
  alt0: ε                          // empty
  alt1: AT_LEAST_ONE(A)            // one or more A
  alt2: AT_LEAST_ONE(A) CONSUME(B) // one or more A followed by B
)
```

| Test | Input token IDs | Expected result |
| --- | --- | --- |
| `TestLL_Star_LongestAlt1` | `[A, A, A]` | `1` (greedy, no terminating B) |
| `TestLL_Star_LongestAlt2` | `[A, A, B]` | `2` (has terminating B) |
| `TestLL_Star_ShortestAlt` | `[]` | `0` (empty alternative) |

#### 6.5.2 Ambiguity detection

Grammar is more complex (see TypeScript test for `AmbigiousParser`). Use a callback to collect
ambiguity report strings. A mock rule set equivalent to:

```text
OptionRule       := OPTION(AT_LEAST_ONE(A)) AT_LEAST_ONE(A)
AltRule          := OR(SUBRULE(RuleB), SUBRULE(RuleC))
RuleB            := MANY(A)
RuleC            := MANY(A) OPTION(B)
AltRuleWithEOF   := OR(SUBRULE(RuleEOF), SUBRULE(RuleEOF))
RuleEOF          := MANY(A) CONSUME(EOF)
AltRuleWithPred  := OR(GATE(pred,CONSUME(A)), GATE(!pred,CONSUME(A)), CONSUME(B))
AltWithOption    := OR(CONSUME(A), CONSUME(B)) OPTION(CONSUME(A))
```

| Test | Input | Expected alt | Expected ambiguity reports |
| --- | --- | --- | --- |
| `TestAmbig_Option` | `[A, A, A]` | truthy (option taken) | `"<0, 1> in <OPTION>"`, `"<0, 1> in <AT_LEAST_ONE1>"` |
| `TestAmbig_FirstAltOnAmbiguity` | `[A, A, A]` | `0` | `"<0, 1> in <OR>"` |
| `TestAmbig_EOFAmbiguity` | `[]` | `0` | `"<0, 1> in <OR>"` |
| `TestAmbig_LongPrefixNoAmbiguity` | `[A, A, B]` | `1` | no reports |
| `TestAmbig_PredicateOverride_Auto` | `[A]` | `0` | `"<0, 1> in <OR>"` |
| `TestAmbig_PredicateOverride_True` | `[A]`, pred=true | `0` | no reports |
| `TestAmbig_PredicateOverride_False` | `[A]`, pred=false | `1` | no reports |
| `TestAmbig_NonAmbigInPredicated` | `[B]` | `2` | no reports |
| `TestAmbig_AltFollowedByOption` | `[B, A]` | value==5 | no ambiguity on OR |

---

## 7. Integration with Existing `ParserState`

The existing `parser.ParserState.Lookahead(LLkLookahead) int` method implements static LL(k)
prediction. To use the new ALL(*) strategy alongside it, `ParserState` satisfies `TokenSource`
as-is (it already has `LA(offset int) *core.Token`).

Generated parsers that opt-in to ALL(*) lookahead will:

1. Call `allstar.NewLLStarLookahead(rules, opts)` (import `typefox.dev/fastbelt/parser/allstar`) once at parser init time.
2. Store the returned `*LLStarLookahead` on the parser struct.
3. Replace calls like `p.Lookahead(table)` with `p.allstar.AdaptivePredict(p, decisionIndex, preds)`.

No changes to the existing `parser/parser.go` are required; `LLkLookahead` remains available for
parsers that do not need unbounded lookahead.

---

## 8. Dependencies

No new external dependencies are required. The `allstar` package uses only:

- `fmt` (string formatting for keys and error messages)
- `strings` (string building for `ATNConfigSet.Key`)
- `sync` (`sync.RWMutex` for the DFA cache)
- `typefox.dev/fastbelt` (core token/document types, imported as `core`)
- `typefox.dev/fastbelt/internal/grammar` (imported only by `convert.go`)

Test files additionally use `github.com/stretchr/testify` (already in `go.mod`).

---

## 9. Resolved Design Decisions

| # | Question | Decision |
| --- | --- | --- |
| 1 | Token category matching | `CategoryMatches []int` lives on `AtomTransition`; `core.TokenType` is unchanged. |
| 2 | EOF handling | `TokenSource.LA` always returns a non-nil token; uses `core.EOF.Id` sentinel. The algorithm never handles `nil`. |
| 3 | Package placement | `parser/allstar/` — subordinate to `parser/`, import path `typefox.dev/fastbelt/parser/allstar`. |
| 4 | Concurrency | `DFACache` map protected by `sync.RWMutex`; reads share the lock, writes are exclusive. |
| 5 | Separator repetition | Not needed; `RepetitionWithSeparator` and `RepetitionMandatoryWithSeparator` are omitted. |
| 6 | Grammar type reuse | A conversion layer (`convert.go`) translates `internal/grammar` AST → `allstar` types; no direct reuse. |

---

## 10. File Summary

| File | Lines (est.) | Mirrors |
| --- | --- | --- |
| `parser/allstar/grammar.go` | ~100 | n/a (ATN-oriented production model) |
| `parser/allstar/convert.go` | ~120 | n/a (conversion from `internal/grammar`) |
| `parser/allstar/atn.go` | ~320 | `atn.ts` (643 lines, separator variants removed) |
| `parser/allstar/dfa.go` | ~80 | `dfa.ts` (79 lines) |
| `parser/allstar/predict.go` | ~350 | `all-star-lookahead.ts` (764 lines) |
| `parser/allstar/strategy.go` | ~80 | part of `all-star-lookahead.ts` |
| `parser/allstar/convert_test.go` | ~130 | — |
| `parser/allstar/atn_test.go` | ~160 | — |
| `parser/allstar/dfa_test.go` | ~80 | — |
| `parser/allstar/predict_test.go` | ~200 | — |
| `parser/allstar/integration_test.go` | ~250 | `atn.test.ts` (303 lines) |
| **Total** | **~1 870** | |
