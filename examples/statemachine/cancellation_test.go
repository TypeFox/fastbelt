// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package statemachine

import (
	"context"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"typefox.dev/fastbelt"
	"typefox.dev/fastbelt/lexer"
	"typefox.dev/fastbelt/parser"
)

func TestParseCancelledContextTerminates(t *testing.T) {
	content, _ := generateStatemachineContent(0)
	srv := CreateServices()
	doc := fastbelt.NewDocumentFromString("file:///workspace/statemachine_0.statemachine", "statemachine", content)
	lexer := srv.MustGet[lexer.Lexer]()
	parser := srv.MustGet[parser.Parser]()
	require.NoError(t, lexer.Exec(t.Context(), doc))
	ctx, cancel := context.WithCancel(t.Context())
	// Cancel ahead of time
	cancel()
	assert.Error(t, parser.Parse(ctx, doc))
	assert.Nil(t, doc.Root, "a cancelled parse must leave the document untouched")
}

// hugeStateCount yields a ~20 MB document, large enough that timer and
// scheduling jitter stay small relative to a full lexer or parser run.
const hugeStateCount = 200_000

func TestLexerCancellationStopsEarly(t *testing.T) {
	if testing.Short() {
		t.Skip("skip long running lexer cancellation test")
	}
	content, _ := generateStatemachineContentWithStates(0, hugeStateCount)
	doc := fastbelt.NewDocumentFromString("file:///workspace/statemachine_0.statemachine", "statemachine", content)
	lexer := CreateServices().MustGet[lexer.Lexer]()
	assertCancellationStopsEarly(t, func(ctx context.Context) error {
		return lexer.Exec(ctx, doc)
	})
}

func TestParserCancellationStopsEarly(t *testing.T) {
	if testing.Short() {
		t.Skip("skip long running parser cancellation test")
	}
	content, _ := generateStatemachineContentWithStates(0, hugeStateCount)
	srv := CreateServices()
	doc := fastbelt.NewDocumentFromString("file:///workspace/statemachine_0.statemachine", "statemachine", content)
	lexer := srv.MustGet[lexer.Lexer]()
	require.NoError(t, lexer.Exec(t.Context(), doc))
	parser := srv.MustGet[parser.Parser]()
	assertCancellationStopsEarly(t, func(ctx context.Context) error {
		return parser.Parse(ctx, doc)
	})
}

// assertCancellationStopsEarly times a full run, then reruns with
// cancellation fired after 20% of that time and asserts the cancelled run
// finishes within 30% of it.
func assertCancellationStopsEarly(t *testing.T, run func(ctx context.Context) error) {
	t.Helper()
	err := run(t.Context()) // warm up
	require.NoError(t, err)

	runtime.GC()
	start := time.Now()
	err = run(t.Context())
	require.NoError(t, err)
	full := time.Since(start)

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	runtime.GC()
	start = time.Now()
	timer := time.AfterFunc(full/5, cancel)
	defer timer.Stop()
	err = run(ctx)
	assert.Error(t, err)
	cancelled := time.Since(start)

	t.Logf("full: %v, cancelled: %v", full, cancelled)
	assert.Less(t, cancelled, full*3/10)
}
