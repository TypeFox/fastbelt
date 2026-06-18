// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var trimmer = &DefaultDocumentationTrimmer{}

func TestTrimCommentLineComment(t *testing.T) {
	assert.Equal(t, "foo bar", trimmer.TrimComment("// foo bar"))
}

func TestTrimCommentBlockComment(t *testing.T) {
	assert.Equal(t, "foo bar", trimmer.TrimComment("/* foo bar */"))
}

func TestTrimCommentMultilineBlockComment(t *testing.T) {
	input := "/*\n * line one\n * line two\n */"
	assert.Equal(t, "line one\nline two", trimmer.TrimComment(input))
}

func TestTrimCommentEmpty(t *testing.T) {
	assert.Equal(t, "", trimmer.TrimComment("//"))
	assert.Equal(t, "", trimmer.TrimComment("/* */"))
	assert.Equal(t, "", trimmer.TrimComment("//   "))
}

func TestTrimCommentPreservesTrailingPeriod(t *testing.T) {
	assert.Equal(t, "foo bar.", trimmer.TrimComment("// foo bar."))
}
