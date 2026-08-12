// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package generator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCommentStartEndParts(t *testing.T) {
	require.Equal(t, GetTerminalParts(`//[^\n\r]*`), []Part{
		{start: "//", end: ""},
	})
}

func TestJSStyleMultilineComment(t *testing.T) {
	require.Equal(t, GetTerminalParts(`/\*[\s\S]*?\*/`), []Part{
		{start: "/\\*", end: "\\*/"},
	})
}

func TestJSStyleCombinedComment(t *testing.T) {
	require.Equal(t, GetTerminalParts(`/\*[\s\S]*?\*/|//[^\n\r]*`), []Part{
		{start: "//", end: ""},
		{start: "/\\*", end: "\\*/"},
	})
}

func TestShellStyleSingleLineComment(t *testing.T) {
	require.Equal(t, GetTerminalParts(`#[^\n\r]*`), []Part{
		{start: "#", end: ""},
	})
}
