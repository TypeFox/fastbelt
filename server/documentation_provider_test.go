// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStripCommentMarkersLineComment(t *testing.T) {
	assert.Equal(t, "foo bar", stripCommentMarkers("// foo bar"))
}

func TestStripCommentMarkersHashComment(t *testing.T) {
	assert.Equal(t, "foo bar", stripCommentMarkers("# foo bar"))
}

func TestStripCommentMarkersDoubleDashComment(t *testing.T) {
	assert.Equal(t, "foo bar", stripCommentMarkers("-- foo bar"))
}

func TestStripCommentMarkersBlockComment(t *testing.T) {
	assert.Equal(t, "foo bar", stripCommentMarkers("/* foo bar */"))
}

func TestStripCommentMarkersMultilineBlockComment(t *testing.T) {
	input := "/*\n * line one\n * line two\n */"
	assert.Equal(t, "line one\nline two", stripCommentMarkers(input))
}

func TestStripCommentMarkersEmpty(t *testing.T) {
	assert.Equal(t, "", stripCommentMarkers("//"))
	assert.Equal(t, "", stripCommentMarkers("/* */"))
	assert.Equal(t, "", stripCommentMarkers("//   "))
}
