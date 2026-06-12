// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrimCommentLineComment(t *testing.T) {
	assert.Equal(t, "foo bar", (&DefaultDocumentationTrimmer{}).TrimComment("// foo bar"))
}

func TestTrimCommentBlockComment(t *testing.T) {
	assert.Equal(t, "foo bar", (&DefaultDocumentationTrimmer{}).TrimComment("/* foo bar */"))
}

func TestTrimCommentMultilineBlockComment(t *testing.T) {
	input := "/*\n * line one\n * line two\n */"
	assert.Equal(t, "line one\nline two", (&DefaultDocumentationTrimmer{}).TrimComment(input))
}

func TestTrimCommentEmpty(t *testing.T) {
	trimmer := &DefaultDocumentationTrimmer{}
	assert.Equal(t, "", trimmer.TrimComment("//"))
	assert.Equal(t, "", trimmer.TrimComment("/* */"))
	assert.Equal(t, "", trimmer.TrimComment("//   "))
}

func TestTrimCommentPreservesTrailingPeriod(t *testing.T) {
	assert.Equal(t, "foo bar.", (&DefaultDocumentationTrimmer{}).TrimComment("// foo bar."))
}
