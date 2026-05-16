package statemachine

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"typefox.dev/fastbelt/test"
)

// TestRecovery_MissingInitialState verifies that when the mandatory
// "initialState" keyword is absent the parser reports an error via
// single-token insertion and still returns a usable statemachine node.
// Single-token deletion for the following Token_ID consume may absorb the
// "state" keyword when a state name follows, so this test only asserts that
// errors are reported and the root node is present — not that state blocks
// survive (see TestRecovery_MissingEnd for state-survival coverage).
func TestRecovery_MissingInitialState(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		statemachine NoInit
		events flick
		state off
		  flick => off
		end
	`)

	assert.NotEmpty(t, doc.Document.ParserErrors, "missing initialState should produce at least one parse error")
	sm, ok := test.FindNode[Statemachine](doc)
	require.True(t, ok, "statemachine root node should be present despite parse errors")
	assert.Equal(t, "NoInit", sm.Name())
}

// TestRecovery_StrayToken verifies that a single stray identifier before the
// states section is discarded by Sync so the state block is parsed normally.
func TestRecovery_StrayToken(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		statemachine S
		events flick
		initialState off
		stray_token
		state off
		  flick => off
		end
	`)

	assert.NotEmpty(t, doc.Document.ParserErrors, "stray token should produce at least one parse error")
	sm, ok := test.FindNode[Statemachine](doc)
	require.True(t, ok, "statemachine node should be present")
	assert.Equal(t, "S", sm.Name())
	states := test.FindAll[State](doc)
	assert.NotEmpty(t, states, "state 'off' should survive after the stray token is discarded")
}

// TestRecovery_MissingEnd verifies that a missing "end" keyword inside a state
// block is handled by single-token insertion and parsing continues into the
// next state block.
func TestRecovery_MissingEnd(t *testing.T) {
	f := test.New(t, CreateServices())
	// "state off" is missing its "end"; the parser should insert a synthetic
	// token and continue to parse "state on".
	doc := f.Parse(`
		statemachine S
		events flick
		initialState off
		state off
		  flick => on
		state on
		  flick => off
		end
	`)

	assert.NotEmpty(t, doc.Document.ParserErrors, "missing end keyword should produce at least one parse error")
	_, ok := test.FindNode[Statemachine](doc)
	require.True(t, ok, "statemachine node should be present")
	states := test.FindAll[State](doc)
	assert.GreaterOrEqual(t, len(states), 1, "at least one state should be parsed after recovery")
}

// TestRecovery_ExtraneousTokenBeforeArrow verifies that a stray identifier
// in front of the required "=>" keyword triggers single-token deletion in
// Consume: LA(1) is wrong but LA(2) is the expected "=>", so the parser drops
// LA(1) and matches LA(2) on the next consume.
func TestRecovery_ExtraneousTokenBeforeArrow(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		statemachine S
		events flick
		initialState off
		state off
		  flick bogus => off
		end
	`)

	assert.NotEmpty(t, doc.Document.ParserErrors, "extraneous identifier should produce a parse error")
	off := test.MustFindNamedNode[State](doc, "off")
	require.Len(t, off.Transitions(), 1, "transition should still be parsed after single-token deletion")
}

// TestRecovery_GarbageBetweenStates verifies that Sync's consume-until loop
// skips several unexpected tokens in a row instead of bailing out after one.
// Two stray identifiers appear between state blocks; the second state should
// still be parsed.
func TestRecovery_GarbageBetweenStates(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		statemachine S
		events flick
		initialState off
		state off
		  flick => on
		end
		garbage1 garbage2
		state on
		  flick => off
		end
	`)

	assert.NotEmpty(t, doc.Document.ParserErrors, "garbage tokens should produce parse errors")
	states := test.FindAll[State](doc)
	require.GreaterOrEqual(t, len(states), 2, "both state blocks should survive the garbage cluster")

	names := []string{states[0].Name(), states[1].Name()}
	assert.Contains(t, names, "off")
	assert.Contains(t, names, "on")
}

// TestRecovery_MissingTransitionArrow verifies that a transition missing its
// "=>" arrow is patched by single-token insertion in Consume (the next ID is
// not the expected "=>", and the token after that is also not "=>", so the
// parser inserts a synthetic arrow rather than deleting). The state block as
// a whole must still close cleanly on its "end" keyword.
func TestRecovery_MissingTransitionArrow(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		statemachine S
		events flick
		initialState off
		state off
		  flick on
		end
		state on
		  flick => off
		end
	`)

	assert.NotEmpty(t, doc.Document.ParserErrors, "missing '=>' should produce a parse error")
	states := test.FindAll[State](doc)
	require.Len(t, states, 2, "both state blocks should still parse")
	// Recovery must not bleed into the second state.
	on := test.MustFindNamedNode[State](doc, "on")
	require.Len(t, on.Transitions(), 1)
}

// TestRecovery_MultipleErrorsContinueParsing verifies that multiple independent
// errors in different rules all get reported and the surviving structure is
// still complete. This is the analogue of the broken-grammar case that drove
// the recovery fix: the parser must not unwind on first error.
func TestRecovery_MultipleErrorsContinueParsing(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		statemachine S
		events flick
		initialState off
		state off
		  flick on
		end
		state mid
		  flick => off
		state on
		  flick => off
		end
	`)

	// First state: missing "=>". Second state: missing "end".
	require.GreaterOrEqual(t, len(doc.Document.ParserErrors), 2,
		"each independent error should be reported separately")

	states := test.FindAll[State](doc)
	assert.GreaterOrEqual(t, len(states), 2, "states surrounding the broken one should still parse")
}

// TestRecovery_ErrorRecoveryModeSuppressesCascades verifies that a single
// mistake does not generate a cascade of duplicate errors: while in
// error-recovery mode, additional reportError calls are suppressed until a
// real token is successfully matched.
func TestRecovery_ErrorRecoveryModeSuppressesCascades(t *testing.T) {
	f := test.New(t, CreateServices())
	// A run of three stray identifiers between state and the rest of the
	// statemachine. Without dedup we would emit one error per skipped token;
	// errorRecoveryMode should collapse them into a single message.
	doc := f.Parse(`
		statemachine S
		events flick
		initialState off
		junk junk junk
		state off
		  flick => off
		end
	`)

	require.NotEmpty(t, doc.Document.ParserErrors)
	assert.LessOrEqual(t, len(doc.Document.ParserErrors), 2,
		"a single garbage cluster should not produce one error per token")
	assert.NotEmpty(t, test.FindAll[State](doc), "state block should still be parsed")
}
