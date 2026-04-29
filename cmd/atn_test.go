package main

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/internal/generator"
	"typefox.dev/fastbelt/internal/grammar"
	allstar "typefox.dev/fastbelt/parser/allstar"
	"typefox.dev/fastbelt/textdoc"
	"typefox.dev/lsp"
)

func TestATNandRuntimeATNAreEqual(t *testing.T) {
	grammarText, err := os.ReadFile("../internal/grammar/grammar.fb")
	if err != nil {
		t.Fatalf("failed to read grammar file: %v", err)
	}

	srv := grammar.CreateServices()
	file, _ := textdoc.NewFile(lsp.URIFromPath("internal/grammar/grammar.fb"), "fb", 0, string(grammarText))

	document := core.NewDocument(file)
	srv.Workspace().DocumentManager.Set(document)
	if err := srv.Workspace().Builder.Build(context.Background(), []*core.Document{document}); err != nil {
		t.Fatalf("failed to build document: %v", err)
	}

	grammr, ok := document.Root.(grammar.Grammar)
	if !ok {
		t.Fatalf("parser result is not a Grammar")
	}

	parserRules := grammr.Rules()
	rules, err := generator.FromParserRules(parserRules, generator.GetTokenTypes(grammr))
	if err != nil {
		t.Fatalf("failed to create parser rules: %v", err)
	}

	atn := allstar.CreateATN(rules)
	rtn := grammar.BuildATN()

	require.Equal(t, len(atn.States), len(rtn.States))
	require.Equal(t, len(atn.DecisionStates), len(rtn.DecisionStates))
	require.Equal(t, len(atn.DecisionMap), len(rtn.DecisionMap))

	stateMap := make(map[*allstar.ATNState]*allstar.RuntimeATNState)
	for aIndex, aState := range atn.States {
		rState := rtn.States[aIndex]
		require.Equal(t, aState.StateNumber, rState.StateNumber)
		require.Equal(t, aState.Decision, rState.Decision)
		require.Equal(t, aState.Rule.Name, rState.RuleName)
		require.Equal(t, aState.EpsilonOnlyTransitions, rState.EpsilonOnlyTransitions)
		require.Equal(t, aState.Type, rState.Type)

		require.Equal(t, len(aState.Transitions), len(rState.Transitions))
		stateMap[aState] = rState
	}
	for aIndex, aState := range atn.States {
		rState := rtn.States[aIndex]
		expectedSource := stateMap[aState]
		require.Equal(t, expectedSource, rState)
		for aTransitionIndex, aTransition := range aState.Transitions {
			rTransition := rState.Transitions[aTransitionIndex]
			require.Equal(t, aTransition.IsEpsilon(), rTransition.IsEpsilon())
			actualTarget := stateMap[aTransition.Target()]
			require.Equal(t, actualTarget, rTransition.GetTarget())
		}
	}
}
